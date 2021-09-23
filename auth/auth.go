package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// DefaultHeaderKeyIdentifier is the default api key identifier in request headers
	DefaultHeaderKeyIdentifier string = "x-api-key"
)

// Validator - Validate API key passed in header
func Validator(c *fiber.Ctx) error {
	// Parse api key
	headerKey := c.Get(DefaultHeaderKeyIdentifier)

	// Check if API key was provided
	if headerKey == "" {
		return fiber.ErrBadRequest
	}

	// Get app data
	appData, valid := validateKey(headerKey)

	// Check if the key is enabled
	if valid {
		contextCache := (*c)
		defer func() {
			// Generate IP warning
			if contextCache.IP() != "0.0.0.0" && appData.LastIP != contextCache.IP() {
				log.Print(fmt.Sprintf("Application %s changed IP address from %s to %s", appData.AppName, appData.LastIP, contextCache.IP()))

				// Update last used IP
				go updateAppLastIP(appData.AppID, contextCache.IP())
			}

			// Log request
			go logEvent(appData, contextCache)
		}()

		// Go to next middleware:
		return c.Next()
	}
	log.Printf("Key: %s | Valid: %v", headerKey, valid)
	return fiber.ErrUnauthorized
}

// GenerateKey - Generate a new API key
func GenerateKey(c *fiber.Ctx) error {
	if c.IP() != "127.0.0.1" {
		return fiber.ErrUnauthorized
	}

	// Check appName
	appName := c.Query("app-name")
	if appName == "" {
		return fiber.ErrBadRequest
	}

	// Check if the name is already taken
	_, err := appDataName(appName)
	if err == nil {
		return fiber.ErrConflict
	}

	// New app data
	appData := appllicationData{}.newApp(appName)

	// Add new app to database
	appID, err := addAPIKey(appData)
	appData.AppID = appID
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Log request
	go logEvent(appData, *c)

	// Update last used IP
	go updateAppLastIP(appData.AppID, c.IP())

	return c.SendStatus(fiber.StatusOK)
}

// validateKey - Validate API Key
func validateKey(key string) (appData appllicationData, valid bool) {
	// Get application info
	appData, err := appDataByKey(key)

	// Log error
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Print(fmt.Errorf("appDataByKey: %s", err.Error()))
	}

	// Return
	return appData, appData.Enabled
}

// updateAppLastIP - Update last IP used for app
func updateAppLastIP(appID primitive.ObjectID, IP string) {
	var appData appllicationData
	appData.AppID = appID
	appData.LastIP = IP
	appData.LastUsed = time.Now()

	err := updateAppData(appData)
	if err != nil {
		log.Print(fmt.Errorf("updateAppData: %s", err.Error()))
	}
}

// logEvent - Log access event
func logEvent(appData appllicationData, c fiber.Ctx) {
	if c.IP == nil || c.Path == nil || c.Method == nil {
		log.Printf("bad session pointer")
		return
	}

	// Prepare log data
	logData, err := appData.prepLogData()
	if err != nil {
		log.Print(fmt.Errorf("prepLogData: %s", err.Error()))
		return
	}

	// Fill log data
	logData.RequestIP = c.IP()
	logData.RequestPath = c.Path()
	logData.RequestTime = time.Now()
	logData.RequestMethod = c.Method()

	err = addLogEntry(logData)
	if err != nil {
		log.Print(fmt.Errorf("addLogEntry: %s", err.Error()))
	}
}
