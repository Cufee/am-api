package handlers

import (
	"log"
	"strconv"
	"time"

	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

// HandleNewPmntReq - Handle new payment request
func HandleNewPmntReq(c *fiber.Ctx) error {
	discordID, err := strconv.Atoi(c.Params("discordID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user Data
	userData, err := db.UserByDiscordID(discordID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get ban data
	banData, err := db.BanCheck(userData.ID)
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println("failed to get ban data: ", err)
	}
	// Check if user is banned
	if banData.UserID == userData.ID && banData.Expiration.After(time.Now()) {
		return c.Status(500).JSON(fiber.Map{
			"error": "user is banned",
		})
	}

	return nil
}
