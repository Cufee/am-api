package main

import (
	"log"

	h "github.com/cufee/am-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// WG login routes
	app.Get("/redirect/:intentID", h.HandleWargamingRedirect)
	app.Get("/login/:intentID", h.HandleWargamingLogin)
	app.Get("/newlogin", h.HandleWargamingNewLogin)
	app.Get("/", func (c *fiber.Ctx) error {
		c.Redirect("http://byvko.dev")
		return nil
	})

	// Checks
	app.Get("/users/:discordID", h.HandeleUserCheck)
	app.Get("/players/:playerID", h.HandelePlayerCheck)

	log.Print(app.Listen(":4000"))
}
