package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

// HandleNewBan -
func HandleNewBan(c *fiber.Ctx) error {
	// Get UserID
	discordID, err := strconv.Atoi(c.Params("discordID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var banData mongodbapi.BanData
	// Parse request
	err = c.BodyParser(&banData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check body
	if banData == (mongodbapi.BanData{}) {
		return c.Status(400).JSON(fiber.Map{
			"error": "no ban data provided",
		})
	}

	// Set UserID
	banData.UserID = discordID

	// Set timestamp
	banData.Timestamp = time.Now()

	// Check 24 hour bans
	bans, err := mongodbapi.GetBansCount(discordID, 7)
	if bans > 1 {
		// Increase ban duration
		banDuration := banData.Expiration.Sub(time.Now())
		banData.Expiration = time.Now().Add(time.Duration(bans) * banDuration)
	}

	// Add ban
	err = mongodbapi.AddBanData(banData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(http.StatusOK)
}
