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

	// Auth middleware
	authDisabled := app.Group("")
	authRequired := app.Group("", auth.Validator)

	// Referrals
	authRequired.Get("/referrals/new", h.HandleNewReferral) // Generate new referral link
	authDisabled.Get("/r/:refID", h.HandleReferralLink)     // Redirect

	// WG login routes
	authRequired.Get("/newlogin", h.HandleWargamingNewLogin)          // New login intent
	authDisabled.Get("/login/r/:intentID", h.HandleWargamingRedirect) // Redirect from WG
	authDisabled.Get("/login/:intentID", h.HandleWargamingLogin)      // Login using intentID

	// Users
	authRequired.Get("/users/id/:discordID", h.HandeleUserCheck)                       // Check
	authRequired.Post("/users/id/:discordID/ban", h.HandleNewBan)                      // Ban
	authRequired.Patch("/users/id/:discordID/newdef/:playerID", h.HandleNewDefaultPID) // New default PID

	// Players
	authRequired.Get("/players/id/:playerID", h.HandelePlayerCheckByID)     // Check by ID
	authRequired.Get("/players/name/:nickname", h.HandelePlayerCheckByName) // Check by name

	// Backgrounds
	authRequired.Patch("/background/:discordID", h.HandleSetNewBG)  // Set new
	authRequired.Delete("/background/:discordID", h.HandleRemoveBG) // Delete

	// Premium
	authRequired.Get("/premium/add", h.HandleNewPremiumIntent)              // Add premium time
	authRequired.Get("/premium/newintent", h.HandleNewPremiumIntent)        // Intent for user update
	authRequired.Get("/premium/redirect/:intentID", h.HandleUpdateRedirect) // Commit using intentID

	// Payments
	authRequired.Get("/payments/new/:discordID", paypal.HandleNewSub)                                                         // Start new payment intent
	authDisabled.Get("/payments/redirect", func(ctx *fiber.Ctx) error { return ctx.Redirect("https://aftermath.link", 301) }) // PayPal redirect
	authDisabled.Post("/payments/events", paypal.HandlePaymentEvent)

	// Root
	authDisabled.Get("/", func(ctx *fiber.Ctx) error { return ctx.Redirect("https://aftermath.link", 301) }) // Root redirect

	log.Print(app.Listen(":4000"))
}
