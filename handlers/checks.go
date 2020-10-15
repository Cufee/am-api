package handlers

import (
	"strconv"
	"time"

	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

type response struct {
	DefaultPID  int    `json:"player_id"`
	Premium     bool   `json:"premium"`
	Verified    bool   `json:"verified"`
	CustomBgURL string `json:"bg_url"`
}

// HandeleUserCheck - Quick user check handler
func HandeleUserCheck(c *fiber.Ctx) error {
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