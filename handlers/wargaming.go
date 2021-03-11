package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cufee/am-api/config"
	"github.com/cufee/am-api/intents"
	db "github.com/cufee/am-api/mongodbapi"
	"github.com/cufee/am-api/wargaming"
	"github.com/gofiber/fiber/v2"
)

// StatsRequest - Request for stats
type StatsRequest struct {
	PlayerID int    `json:"player_id"`
	Realm    string `json:"realm"`
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
	intent, err := db.GetUserIntent(c.Params("intentID"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "unable to find a valid intent",
		})
	}

	// Parse account id and time
	accID, err := strconv.Atoi(c.Query("account_id"))
	if err != nil {
		log.Print(err.Error())
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("atoi error: %v", err.Error()),
		})
	}

	// Check if accID exists
	if accID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Wargaming did not return a valid account ID, please report this error.",
		})
	}

	i, err := strconv.ParseInt(c.Query("expires_at"), 10, 64)
	if err != nil {
		log.Print(err.Error())
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("parse int error: %v", err.Error()),
		})
	}

	intent.Data.VerifiedID = accID
	intent.Data.VerifiedExpiration = time.Unix(i, 0)
	intent.Data.DefaultPID = accID

	// Clear any existing logins for this accID
	err = db.RemoveOldLogins(accID)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			log.Print(err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("remove old login error: %v", err.Error()),
			})
		}
		err = nil
	}

	// Check if player is in db
	if !db.PlayerExistsByID(accID) {
		// Add player to DB
		err = db.AddPlayerToDB(accID, intent.Realm)
		if err != nil {
			log.Print(err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("failed to add a new player to db: %v", err.Error()),
			})
		}
	}

	// Add/Update DB record
	err = db.UpdateUser(intent.Data, true)
	if err != nil {
		log.Print(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("update user error: %v", err.Error()),
		})
	}
	return c.SendString("Login success, you can close this window.")
}

// HandleWargamingLogin -
func HandleWargamingLogin(c *fiber.Ctx) error {
	// Get
	intentData, err := db.GetLoginIntent(c.Params("intentID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Link expired.",
		})
	}
	// Get user data
	userData, err := db.UserByDiscordID(intentData.DiscordID)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			log.Print(err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		userData = db.UserData{ID: intentData.DiscordID}
	}

	// Create edit intent
	newIntent, err := intents.CreateUserIntent(userData, intentData.Realm)
	if err != nil {
		log.Print(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Get redirect URL
	reqURL, err := wgAPIurl(intentData.Realm, newIntent.IntentID)
	if err != nil {
		log.Print(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	redirectURL, err := getRedirectURL(reqURL)
	if err != nil {
		log.Print(err.Error())
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
		log.Print(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	existingID := db.GetLogin(data.DiscordID)
	intentID, err := intents.CreateLoginIntent(data)
	if err != nil {
		log.Print(err.Error())
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
	tail := config.WgAPIAppID + "&redirect_uri=" + config.WgBaseRedirectURL + intentID + "&nofollow=1"
	switch realm {
	case "NA":
		return ("https://api.worldoftanks.com/wot/auth/login/?application_id=" + tail), nil
	case "RU":
		return ("https://api.worldoftanks.ru/wot/auth/login/?application_id=" + tail), nil
	case "EU":
		return ("https://api.worldoftanks.eu/wot/auth/login/?application_id=" + tail), nil
	case "ASIA":
		return ("https://api.worldoftanks.asia/wot/auth/login/?application_id=" + tail), nil
	}
	return "", errors.New("bad realm")
}

// Wargaming API Redirect URL

type redirectRes struct {
	Data struct {
		Location string `json:"location"`
	} `json:"data"`
}

// HTTP client
var clientHTTP = &http.Client{Timeout: 10 * time.Second}

func getRedirectURL(reqURL string) (string, error) {
	var resData redirectRes
	err := wargaming.GetLimited(reqURL, &resData)
	if err != nil {
		return "", err
	}
	return resData.Data.Location, nil
}
