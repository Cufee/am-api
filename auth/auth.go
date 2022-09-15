package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
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
	_, valid := validateKey(headerKey)

	// Check if the key is enabled
	if valid {
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
