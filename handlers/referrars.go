package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

type newReferralRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// HandleReferralLink - Handle referral link
func HandleReferralLink(c *fiber.Ctx) error {
	// Fill click info
	var clickData mongodbapi.ReferralClick
	userID, _ := strconv.Atoi(c.Query("user_id")) // Convert user_id to int
	clickData.URL = c.BaseURL() + c.OriginalURL() // URL
	clickData.UserID = userID                     // User ID
	clickData.MetaJSON = string(c.Body())         // JSON Meta

	// Add event to db
	err := mongodbapi.RecordReferalClick(c.Params("refID"), clickData)
	if err != nil {
		log.Print(err)
	}

	// Redirect
	return c.Redirect("http://aftermath.link/")
}

// HandleNewReferral - Handle referral link
func HandleNewReferral(c *fiber.Ctx) error {
	var reqData newReferralRequest
	// Parse title and desxcription
	err := c.BodyParser(&reqData)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("bad json body: %s", err.Error()),
		})
	}

	// Generate new referral
	refData, err := mongodbapi.GenerateNewReferalCode(reqData.Title, reqData.Description)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("bad json body: %s", err.Error()),
		})
	}

	return c.JSON(fiber.Map{
		"url": refData.URL,
	})
}
