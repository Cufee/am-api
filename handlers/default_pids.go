package handlers

import (
	"strconv"
	"time"

	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

type pidChangeRes struct {
	OldPID   int  `json:"old_player_id"`
	NewPID   int  `json:"new_player_id"`
	Verified bool `json:"verified"`
}

// HandleNewDefaultPID - Set a new default player ID for a user
func HandleNewDefaultPID(c *fiber.Ctx) error {
	discordID, err := strconv.Atoi(c.Params("discordID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	newPlayerID, err := strconv.Atoi(c.Params("playerID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if both IDs are provided
	if discordID == 0 || newPlayerID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "discord_id or player_id not provided",
		})
	}

	// Get user data
	userData, err := db.UserByDiscordID(discordID)

	// Save old PID
	oldPID := userData.DefaultPID
	if userData.VerifiedID != 0 {
		oldPID = userData.VerifiedID
	}

	// Create a new user record if one does not exist
	if err != nil && err.Error() == "mongo: no documents in result" {
		userData = db.UserData{ID: discordID, DefaultPID: newPlayerID}
		err = nil
	}

	// Check for other errors
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Reset verificaton
	if userData.VerifiedID != newPlayerID {
		userData.VerifiedID = 0
		userData.VerifiedExpiration = time.Now().Add(-1 * time.Minute)
	}

	// Update DB
	err = db.UpdateUser(userData, true)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Make response
	var resData pidChangeRes

	resData.OldPID = oldPID
	resData.NewPID = newPlayerID
	resData.Verified = false
	if userData.VerifiedExpiration.After(time.Now()) {
		resData.Verified = true
	}

	return c.JSON(resData)
}
