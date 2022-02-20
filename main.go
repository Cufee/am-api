package main

import (
	"log"
	"os"

	"github.com/cufee/am-api/auth"
	h "github.com/cufee/am-api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// Logger
	app.Use(logger.New())
	// CORS
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	api := app.Group("/users/v1")

	// Public endpoints
	api.Get("/public/:realm/players/name/:nickname", h.HandlePublicPlayerCheckByName) // Check by name - Public

	// Generate API key - localhost only
	api.Get("/keys/new", auth.GenerateKey)

	// Referrals
	api.Get("/referrals/new", auth.Validator, h.HandleNewReferral) // Generate new referral link
	api.Get("/r/:refID", h.HandleReferralLink)                     // Redirect

	// WG login routes
	api.Get("/newlogin", auth.Validator, h.HandleWargamingNewLogin) // New login intent
	api.Get("/login/r/:intentID", h.HandleWargamingRedirect)        // Redirect from WG
	api.Get("/login/:intentID", h.HandleWargamingLogin)             // Login using intentID

	// Users
	api.Get("/users/id/:discordID", auth.Validator, h.HandeleUserCheck)                       // Check
	api.Post("/users/id/:discordID/ban", auth.Validator, h.HandleNewBan)                      // Ban
	api.Patch("/users/id/:discordID/newdef/:playerID", auth.Validator, h.HandleNewDefaultPID) // New default PID

	// Players
	api.Get("/players/id/:playerID", auth.Validator, h.HandelePlayerCheckByID)     // Check by ID
	api.Get("/players/name/:nickname", auth.Validator, h.HandelePlayerCheckByName) // Check by name

	// Backgrounds
	api.Patch("/background/:discordID", auth.Validator, h.HandleSetNewBG)  // Set new
	api.Delete("/background/:discordID", auth.Validator, h.HandleRemoveBG) // Delete

	// Premium
	api.Get("/premium/add", auth.Validator, h.HandleNewPremiumIntent)              // Add premium time
	api.Get("/premium/newintent", auth.Validator, h.HandleNewPremiumIntent)        // Intent for user update
	api.Get("/premium/redirect/:intentID", auth.Validator, h.HandleUpdateRedirect) // Commit using intentID

	// // Payments
	// api.Get("/payments/new/:discordID", paypal.HandleNewSub)                                                         // Start new payment intent
	// api.Get("/payments/redirect", func(ctx *fiber.Ctx) error { return ctx.Redirect("https://aftermath.link", 301) }) // PayPal redirect
	// api.Post("/payments/events", paypal.HandlePaymentEvent)

	// Root
	api.Get("/", func(ctx *fiber.Ctx) error { return ctx.Redirect("https://aftermath.link", 301) }) // Root redirect

	log.Panic(app.Listen(":" + os.Getenv("PORT")))
}
