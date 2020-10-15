package handlers

import (
	"errors"
	"time"

	"github.com/cufee/am-api/config"
	"github.com/cufee/am-api/intents"
	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

// RedirectRequest - Request from WG auth redirect
type RedirectRequest struct {
	Status     string    `json:"status"`
	Code       int       `json:"code"`
	Message    string    `json:"message"`
	Token      string    `json:"access_token"`
	Expiration time.Time `json:"expires_at"`
	AccountID  int       `json:"account_id"`
	Nickname   string    `json:"nickname "`
}

// LoginRequest - Request to login using WG
type LoginRequest struct {
	DiscordID int    `json:"discord_user_id"`
	Realm     string `json:"realm"`
}

// HandleWargamingRedirect -
func HandleWargamingRedirect(c *fiber.Ctx) error {
	var data RedirectRequest
	// Parse body and params
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Check reponse
	if data.Status != "ok" {
		// Auth failed
		return c.Status(data.Code).JSON(fiber.Map{
			"error": data.Message,
		})
	}
	// Get intent
	intent, err := intents.GetUserIntent(c.Params("intentID"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "unable to find a valid intent",
		})
	}
	// Add DB record
	intent.Data.Verified = true
	intent.Data.VerifiedID = data.AccountID
	intent.Data.VerifiedExpiration = data.Expiration
	intent.Data.DefaultPID = data.AccountID
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
	var data LoginRequest
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Get user data
	userData, err := db.UserByDiscordID(data.DiscordID)
	switch err.Error() {
	case "":
		break
	case "mongo: no documents in result":
		// Create a new user
		userData = db.UserData{ID: data.DiscordID}
		break
	default:
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Create edit intent
	intentID, err := intents.CreateUserIntent(userData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Get redirect URL
	redirectURL, err := wgAPIurl(data.Realm, intentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Redirect(redirectURL)
}

func wgAPIurl(realm string, intentID string) (string, error) {
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
