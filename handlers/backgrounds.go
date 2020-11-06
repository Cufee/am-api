package handlers

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cufee/am-api/config"
	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
)

// HTTP Client for outgoing requests
var client = &http.Client{Timeout: 10 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

// NSFW API response
type nsfwRes struct {
	Output struct {
		NsfwScore float64 `json:"nsfw_score"`
	} `json:"output"`
	Err string `json:"err"`
}

type cdnRes struct {
	PublicID  string    `json:"public_id"`
	Format    string    `json:"format"`
	CreatedAt time.Time `json:"created_at"`
	URL       string    `json:"url"`
}

// HandleSetNewBG - Set a new backgorund for a player
func HandleSetNewBG(c *fiber.Ctx) error {
	// Get Player ID
	discordID, err := strconv.Atoi(c.Params("discordID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user data
	userData, err := db.UserByDiscordID(discordID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if time.Now().After(userData.VerifiedExpiration) {
		return c.Status(401).JSON(fiber.Map{
			"error": "user is not verified",
		})
	}

	// Check if image is SFW
	newURL := c.Query("bgurl", "none")
	if newURL == "none" {
		return c.Status(400).JSON(fiber.Map{
			"error": "no image url provided",
		})

	}
	nsfw, err := isNSFW(newURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
			"nsfw":  nsfw,
		})
	}

	// Upload image to CDN
	imgURL, err := uploadToCDN(newURL, userData.VerifiedID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Update user data
	userData.CustomBgURL = imgURL
	err = db.UpdateUser(userData, false)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Response
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

// HandleRemoveBG - Remove backgorund for a player
func HandleRemoveBG(c *fiber.Ctx) error {
	// Get Player ID
	discordID, err := strconv.Atoi(c.Params("discordID"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user data
	userData, err := db.UserByDiscordID(discordID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if time.Now().After(userData.VerifiedExpiration) {
		return c.Status(401).JSON(fiber.Map{
			"error": "user is not verified",
		})
	}

	// Update user data
	userData.CustomBgURL = ""
	err = db.UpdateUser(userData, false)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Response
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
func isNSFW(imageURL string) (bool, error) {
	// Make request url
	reqURL := config.NSFWAPIURL
	// Send request
	form := url.Values{}
	form.Add("image", imageURL)
	req, _ := http.NewRequest("POST", reqURL, strings.NewReader(form.Encode()))
	req.Header.Add("api-key", config.NSFWAPIKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Read response
	res, err := client.Do(req)
	if err != nil {
		return true, err
	}
	defer res.Body.Close()

	// Check NSFW score
	var data nsfwRes
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return true, err
	}

	// Valid score will never be exactly 0
	if data.Output.NsfwScore == 0.0 {
		return true, fmt.Errorf("NSFW detector failed: %s", data.Err)
	}

	// Image is safe or work
	if data.Output.NsfwScore < 0.7 {
		return false, nil
	}
	return true, fmt.Errorf("image is NSFW")
}

func uploadToCDN(imageURL string, pid int) (string, error) {
	// Make request url
	reqURLP := config.CloudinaryUploadURL

	// Make signature
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	sigStr := fmt.Sprintf("format=jpg&public_id=%v&timestamp=%s", pid, timestamp) + config.CloudinaryAPISecret

	// Encode signature
	h := sha1.New()
	h.Write([]byte(sigStr))
	signature := hex.EncodeToString(h.Sum(nil))

	// Generate form
	form := url.Values{}
	form.Add("timestamp", timestamp)
	form.Add("public_id", strconv.Itoa(pid))
	form.Add("api_key", config.CloudinaryAPIKey)
	form.Add("file", imageURL)
	form.Add("signature", signature)
	// Resize and change the format
	form.Add("format", "jpg")
	form.Add("transformations ", "w_400")

	// Send post request
	req, _ := http.NewRequest("POST", reqURLP, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf(res.Status)
	}
	defer res.Body.Close()

	// Return image URL and err
	var resData cdnRes
	err = json.NewDecoder(res.Body).Decode(&resData)
	return resData.URL, err
}
