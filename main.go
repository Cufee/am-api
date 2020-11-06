package main

import (
	"log"

	h "github.com/cufee/am-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	externalApp := fiber.New()
	internalApp := fiber.New()

	// WG login routes
	externalApp.Get("/newlogin", h.HandleWargamingNewLogin)
	externalApp.Get("/", func(c *fiber.Ctx) error {
		c.Redirect("http://byvko.dev")
		return nil
	})

	// Login
	internalApp.Get("/redirect/:intentID", h.HandleWargamingRedirect)
	internalApp.Get("/login/:intentID", h.HandleWargamingLogin)

	// Checks
	internalApp.Get("/users/:discordID", h.HandeleUserCheck)
	internalApp.Get("/players/:playerID", h.HandelePlayerCheck)

	// Backgrounds
	internalApp.Get("/setnewbg/:discordID", h.HandleSetNewBG)
	internalApp.Get("/removebg/:discordID", h.HandleRemoveBG)

	go log.Print(externalApp.Listen(":4000"))
}
