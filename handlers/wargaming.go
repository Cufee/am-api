package handlers

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/cufee/am-api/config"
	"github.com/cufee/am-api/intents"
	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

// HandleWargamingRedirect -
func HandleWargamingRedirect(c *fiber.Ctx) error {
	// Check reponse
	if c.Query("status") != "ok" {
		// Auth failed
		statusCode, _ := strconv.Atoi(c.Query("code"))
		return c.Status(statusCode).JSON(fiber.Map{
			"error": c.Query("message"),
		})
	}
	// Get intent
	intent, err := intents.GetUserIntent(c.Params("intentID"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "unable to find a valid intent",
		})
	}
	// Parse account id and time
	accID, err := strconv.Atoi(c.Query("account_id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	i, err := strconv.ParseInt(c.Query("expires_at"), 10, 64)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	tm := time.Unix(i, 0)
	intent.Data.Verified = true
	intent.Data.VerifiedID = accID
	intent.Data.VerifiedExpiration = tm
	intent.Data.DefaultPID = accID
	// Add DB record
	err = db.UpdateUser(intent.Data, true)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendString("Login success, you can close this window.")
}

// HandleWargamingLogin -
func HandleWargamingLogin(c *fiber.Ctx) error {
	// Get
	intentData, err := intents.GetLoginIntent(c.Params("intentID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Link expired.",
		})
	}
	// Get user data
	userData, err := db.UserByDiscordID(intentData.DiscordID)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		userData = db.UserData{ID: intentData.DiscordID}
	}
	// Create edit intent
	newIntentID, err := intents.CreateUserIntent(userData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Get redirect URL
	redirectURL, err := wgAPIurl(intentData.Realm, newIntentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Redirect(redirectURL)
}

// HandleWargamingNewLogin -
func HandleWargamingNewLogin(c *fiber.Ctx) error {
	var data db.LoginData
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	existingID := db.GetLogin(data.DiscordID)
	intentID, err := intents.CreateLoginIntent(data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Return intentID
	return c.JSON(fiber.Map{
		"intent_id":   intentID,
		"existing_id": existingID,
	})
}

func wgAPIurl(realm string, intentID string) (string, error) {
	realm = strings.ToUpper(realm)
	switch realm {
	case "NA":
		return ("https://api.worldoftanks.com/wot/auth/login/?application_id=" + config.WgAPIAppID + "&redirect_uri=" + config.WgBaseRedirectURL + intentID), nil
	case "RU":
		return ("https://api.worldoftanks.ru/wot/auth/login/?application_id=" + config.WgAPIAppID + "&redirect_uri=" + config.WgBaseRedirectURL + intentID), nil
	case "EU":
		return ("https://api.worldoftanks.eu/wot/auth/login/?application_id=" + config.WgAPIAppID + "&redirect_uri=" + config.WgBaseRedirectURL + intentID), nil
	case "ASIA":
		return ("https://api.worldoftanks.asia/wot/auth/login/?application_id=" + config.WgAPIAppID + "&redirect_uri=" + config.WgBaseRedirectURL + intentID), nil
	}
	return "", errors.New("bad realm")
}
