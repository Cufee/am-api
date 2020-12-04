package main

import (
	"log"

	h "github.com/cufee/am-api/handlers"
	"github.com/cufee/am-api/paypal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// Logger
	app.Use(logger.New())

	// Referrals
	app.Get("/referrals/new", h.HandleNewReferral)
	app.Get("/referrals/:refID", h.HandleReferralLink)

	// WG login routes
	app.Get("/redirect/:intentID", h.HandleWargamingRedirect)
	app.Get("/login/:intentID", h.HandleWargamingLogin)
	app.Get("/newlogin", h.HandleWargamingNewLogin)

	// Checks
	app.Get("/users/:discordID", h.HandeleUserCheck) // Will be dropped
	app.Get("/users/id/:discordID", h.HandeleUserCheck)

	app.Patch("/users/:discordID/newdef/:playerID", h.HandleNewDefaultPID)

	app.Get("/players/:playerID", h.HandelePlayerCheckByID) // Will be dropped
	app.Get("/players/id/:playerID", h.HandelePlayerCheckByID)
	app.Get("/players/name/:nickname", h.HandelePlayerCheckByName)

	// Backgrounds
	app.Get("/setnewbg/:discordID", h.HandleSetNewBG)
	app.Get("/removebg/:discordID", h.HandleRemoveBG)

	// Premium
	app.Get("/premium/add", h.HandleNewPremiumIntent)
	app.Get("/premium/newintent", h.HandleNewPremiumIntent)
	app.Get("/premium/redirect/:intentID", h.HandleUpdateRedirect)

	// Payments
	app.Get("/payments/new/:discordID", paypal.HandleNewSub)
	app.Get("/payments/redirect", func(ctx *fiber.Ctx) error { return ctx.Redirect("https://aftermath.link", 301) })
	app.Post("/payments/events", paypal.HandlePaymentEvent)

	// Root
	app.Get("/", func(ctx *fiber.Ctx) error { return ctx.Redirect("https://aftermath.link", 301) })

	log.Print(app.Listen(":4000"))
}
