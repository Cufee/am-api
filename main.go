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

	// WG login routes
	app.Get("/redirect/:intentID", h.HandleWargamingRedirect)
	app.Get("/login/:intentID", h.HandleWargamingLogin)
	app.Get("/newlogin", h.HandleWargamingNewLogin)

	// Checks
	app.Get("/users/:discordID", h.HandeleUserCheck)
	app.Patch("/users/:discordID/newdef/:playerID", h.HandleNewDefaultPID)
	app.Get("/players/:playerID", h.HandelePlayerCheck)

	// Backgrounds
	app.Get("/setnewbg/:discordID", h.HandleSetNewBG)
	app.Get("/removebg/:discordID", h.HandleRemoveBG)

	// Premium
	app.Get("/premium/add", h.HandleNewPremiumIntent)
	app.Get("/premium/newintent", h.HandleNewPremiumIntent)
	app.Get("/premium/redirect/:intentID", h.HandleUpdateRedirect)

	// Payments
	app.Get("/payments/new/:discordID", paypal.HandleNewSub)
	app.Post("/payments/events", paypal.HandlePaymentEvent)

	log.Print(app.Listen(":4000"))
}
