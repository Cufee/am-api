package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/cufee/am-api/config"
	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

type response struct {
	DefaultPID int    `json:"player_id"`
	Locale     string `json:"locale"`

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

	// Locale
	resData.Locale = userData.Locale

	// Get ban data
	banData, err := db.BanCheck(userData.ID)
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
	}

	if banData.UserID == userData.ID && banData.Expiration.After(time.Now()) {
		resData.Banned = true
		resData.BanReason = banData.Reason
		resData.BanNotified = banData.Notified
	}

	resData.DefaultPID = userData.DefaultPID

	if config.AllUsersPremium {
		// Check if premium features are enabled for all
		resData.Premium = true
		resData.CustomBgURL = userData.CustomBgURL
	} else {
		// Check premium status
		resData.Premium = false
		if time.Now().Before(userData.PremiumExpiration) {
			resData.Premium = true
			resData.CustomBgURL = userData.CustomBgURL
		}
	}

	// Check verification status
	resData.Verified = false
	if time.Now().Before(userData.VerifiedExpiration) {
		resData.Verified = true
		resData.DefaultPID = userData.VerifiedID
	}
	return c.JSON(resData)
}

// HandelePlayerCheckByID - Quick user check by player id handler
func HandelePlayerCheckByID(c *fiber.Ctx) error {
	// Get user data
	playerID, err := strconv.Atoi(c.Params("playerID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get player profile
	resData, err := checkByPID(playerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resData)
}

// HandelePlayerCheckByName - Quick user check by player name handler
func HandelePlayerCheckByName(c *fiber.Ctx) error {
	// Get user data
	c.Params("nickname")
	// Get ID from regex match to nickname
	playerID := 0

	// Get player profile
	resData, err := checkByPID(playerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resData)
}

// HandlePublicPlayerCheckByName - Quick user check by player name handler
func HandlePublicPlayerCheckByName(c *fiber.Ctx) error {
	// Get ID from name
	playerID, err := db.PlayerIDbyName(c.Params("nickname"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var publicRes response
	publicRes.DefaultPID = playerID

	return c.JSON(publicRes)
}

func checkByPID(pid int) (resData response, err error) {
	// Get user data
	userData, err := db.UserByPlayerID(pid)
	if err != nil {
		return resData, err
	}

	// Get ban data
	banData, err := db.BanCheck(userData.ID)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			err = nil
		} else {
			log.Println(err)
		}
	}
	if banData.UserID == userData.ID {
		resData.Banned = true
		resData.BanReason = banData.Reason
		resData.BanNotified = banData.Notified
	}

	if config.AllUsersPremium {
		// Check if premium features are enabled for all
		resData.Premium = true
		resData.CustomBgURL = userData.CustomBgURL
	} else {
		// Check premium status
		resData.Premium = false
		if time.Now().Before(userData.PremiumExpiration) {
			resData.Premium = true
			resData.CustomBgURL = userData.CustomBgURL
		}
	}

	// Check verified status
	resData.Verified = false
	if time.Now().Before(userData.VerifiedExpiration) {
		resData.Verified = true
	}
	return resData, err
}
