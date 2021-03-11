package handlers

import (
	"time"

	"github.com/cufee/am-api/intents"
	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

type ptwReq struct {
	ID          int `json:"user_id"`
	PremiumDays int `json:"premium_days"`
}

// HandleNewPremiumIntent -
func HandleNewPremiumIntent(c *fiber.Ctx) error {
	var request ptwReq
	// Parse request
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate
	if request.ID == 0 || request.PremiumDays == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Bad request, user ID or days is 0.",
		})
	}

	// Find user data
	userData, err := db.UserByDiscordID(request.ID)
	// Create a new user record if one does not exist
	if err != nil && err.Error() == "mongo: no documents in result" {
		userData = db.UserData{ID: request.ID}
		err = nil
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Update premium time
	if userData.PremiumExpiration.After(time.Now()) {
		userData.PremiumExpiration = userData.PremiumExpiration.Add(time.Duration(request.PremiumDays) * time.Minute * 1440)
	} else {
		userData.PremiumExpiration = time.Now().Add(time.Duration(request.PremiumDays) * time.Minute * 1440)
	}

	// Create new update intent
	var intent db.UserDataIntent
	intent.Data = userData

	intentID, err := intents.CreateUserIntent(intent.Data, "")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return intentID
	return c.JSON(fiber.Map{
		"intent_id": intentID,
	})
}

// HandleUpdateRedirect -
func HandleUpdateRedirect(c *fiber.Ctx) error {
	// Get intent
	intent, err := db.GetUserIntent(c.Params("intentID"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "unable to find a valid intent",
		})
	}

	// Add/Update DB record
	err = db.UpdateUser(intent.Data, true)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
