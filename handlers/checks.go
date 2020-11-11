package handlers

import (
	"log"
	"strconv"
	"time"

	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

type response struct {
	DefaultPID int `json:"player_id"`

	Premium  bool `json:"premium"`
	Verified bool `json:"verified"`

	CustomBgURL string `json:"bg_url"`

	Banned      bool   `json:"banned"`
	BanReason   string `json:"ban_reason,omitempty"`
	BanNotified bool   `json:"ban_notified,omitempty"`
}

// HandeleUserCheck - Quick user check handler
func HandeleUserCheck(c *fiber.Ctx) error {
	discordID, err := strconv.Atoi(c.Params("discordID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userData, err := db.UserByDiscordID(discordID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var resData response

	// Get ban data
	banData, err := db.BanCheck(userData.ID)
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
	}

	if banData.UserID == userData.ID {
		resData.Banned = true
		resData.BanReason = banData.Reason
		resData.BanNotified = banData.Notified
	}

	resData.DefaultPID = userData.DefaultPID
	resData.Premium = false
	if time.Now().Before(userData.PremiumExpiration) {
		resData.Premium = true
		resData.CustomBgURL = userData.CustomBgURL
	}
	resData.Verified = false
	if time.Now().Before(userData.VerifiedExpiration) {
		resData.Verified = true
	}
	return c.JSON(resData)
}

// HandelePlayerCheck - Quick user check handler
func HandelePlayerCheck(c *fiber.Ctx) error {
	// Get user data
	playerID, err := strconv.Atoi(c.Params("playerID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userData, err := db.UserByPlayerID(playerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var resData response

	// Get ban data
	banData, err := db.BanCheck(userData.ID)
	if err != nil && err.Error() != "mongo: no document in result" {
		log.Println(err)
	}

	if banData.UserID == userData.ID {
		resData.Banned = true
		resData.BanReason = banData.Reason
		resData.BanNotified = banData.Notified
	}

	resData.DefaultPID = userData.DefaultPID
	resData.Premium = false
	if time.Now().Before(userData.PremiumExpiration) {
		resData.Premium = true
		resData.CustomBgURL = userData.CustomBgURL
	}
	resData.Verified = false
	if time.Now().Before(userData.VerifiedExpiration) {
		resData.Verified = true
	}
	return c.JSON(resData)
}
