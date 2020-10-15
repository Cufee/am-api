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

// LoginRequest - Request to login using WG
type LoginRequest struct {
	DiscordID int    `json:"discord_user_id"`
	Realm     string `json:"realm"`
}

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
	accID, err := strconv.Atoi(c.Query("account_id "))
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
	IDasInt, err := strconv.Atoi(c.Params("discordID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Get user data
	userData, err := db.UserByDiscordID(IDasInt)
	switch err.Error() {
	case "":
		break
	case "mongo: no documents in result":
		// Create a new user
		userData = db.UserData{ID: IDasInt}
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
	redirectURL, err := wgAPIurl(c.Params("realm"), intentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Redirect(redirectURL)
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
