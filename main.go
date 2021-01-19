package main

import (
	"log"

	"github.com/cufee/am-api/auth"
	h "github.com/cufee/am-api/handlers"
	"github.com/cufee/am-api/paypal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// Logger
	app.Use(logger.New())

	// Generate API key - localhost only
	app.Get("/keys/new", auth.GenerateKey)

	// Referrals
	app.Get("/referrals/new", auth.Validator, h.HandleNewReferral) // Generate new referral link
	app.Get("/r/:refID", h.HandleReferralLink)                     // Redirect

	// WG login routes
	app.Get("/newlogin", auth.Validator, h.HandleWargamingNewLogin) // New login intent
	app.Get("/login/r/:intentID", h.HandleWargamingRedirect)        // Redirect from WG
	app.Get("/login/:intentID", h.HandleWargamingLogin)             // Login using intentID

	// Users
	app.Get("/users/id/:discordID", auth.Validator, h.HandeleUserCheck)                       // Check
	app.Post("/users/id/:discordID/ban", auth.Validator, h.HandleNewBan)                      // Ban
	app.Patch("/users/id/:discordID/newdef/:playerID", auth.Validator, h.HandleNewDefaultPID) // New default PID

	// Players
	app.Get("/players/id/:playerID", auth.Validator, h.HandelePlayerCheckByID)     // Check by ID
	app.Get("/players/name/:nickname", auth.Validator, h.HandelePlayerCheckByName) // Check by name

	// Backgrounds
	app.Patch("/background/:discordID", auth.Validator, h.HandleSetNewBG)  // Set new
	app.Delete("/background/:discordID", auth.Validator, h.HandleRemoveBG) // Delete

	// Premium
	app.Get("/premium/add", auth.Validator, h.HandleNewPremiumIntent)              // Add premium time
	app.Get("/premium/newintent", auth.Validator, h.HandleNewPremiumIntent)        // Intent for user update
	app.Get("/premium/redirect/:intentID", auth.Validator, h.HandleUpdateRedirect) // Commit using intentID

	// Payments
	app.Get("/payments/new/:discordID", paypal.HandleNewSub)                                                         // Start new payment intent
	app.Get("/payments/redirect", func(ctx *fiber.Ctx) error { return ctx.Redirect("https://aftermath.link", 301) }) // PayPal redirect
	app.Post("/payments/events", paypal.HandlePaymentEvent)

	// Root
	app.Get("/", func(ctx *fiber.Ctx) error { return ctx.Redirect("https://aftermath.link", 301) }) // Root redirect

	log.Print(app.Listen(":4000"))
}
