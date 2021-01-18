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

	//
	// Auth disabled
	//

	// Generate API key - localhost only
	app.Get("/keys/new", auth.GenerateKey)

	// Root
	app.Get("/", func(ctx *fiber.Ctx) error { return ctx.Redirect("https://aftermath.link", 301) }) // Root redirect

	// WG redirect
	app.Get("/redirect/:intentID", h.HandleWargamingRedirect) // Redirect from WG

	// Payments
	app.Get("/payments/redirect", func(ctx *fiber.Ctx) error { return ctx.Redirect("https://aftermath.link", 301) }) // PayPal redirect
	app.Post("/payments/events", paypal.HandlePaymentEvent)                                                          // PayPal Events

	//
	// Auth enabled
	//

	// API key validator
	app.Use(auth.Validator)

	// Referrals
	app.Get("/referrals/new", h.HandleNewReferral) // Generate new referral link
	app.Get("/r/:refID", h.HandleReferralLink)     // Redirect

	// WG login routes
	app.Get("/newlogin", h.HandleWargamingNewLogin)     // New login intent
	app.Get("/login/:intentID", h.HandleWargamingLogin) // Login using intentID

	// Users
	app.Get("/users/id/:discordID", h.HandeleUserCheck)                    // Check
	app.Post("/users/id/:discordID/ban", h.HandleNewBan)                   // Ban
	app.Patch("/users/:discordID/newdef/:playerID", h.HandleNewDefaultPID) // New default PID

	// Players
	app.Get("/players/id/:playerID", h.HandelePlayerCheckByID)     // Check by ID
	app.Get("/players/name/:nickname", h.HandelePlayerCheckByName) // Check by name

	// Backgrounds
	app.Patch("/background/:discordID", h.HandleSetNewBG)  // Set new
	app.Delete("/background/:discordID", h.HandleRemoveBG) // Delete

	// Premium
	app.Get("/premium/add", h.HandleNewPremiumIntent)              // Add premium time
	app.Get("/premium/newintent", h.HandleNewPremiumIntent)        // Intent for user update
	app.Get("/premium/redirect/:intentID", h.HandleUpdateRedirect) // Commit using intentID

	// Payments
	app.Get("/payments/new/:discordID", paypal.HandleNewSub) // Start new payment intent

	log.Print(app.Listen(":4000"))
}
