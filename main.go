package main

import (
	"log"

	h "github.com/cufee/am-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// WG login routes
	app.Get("/login/redirect/:intentID", h.HandleWargamingRedirect)
	app.Get("/login", h.HandleWargamingLogin)
	app.Get("/ping", handlePing)

	log.Print(app.Listen(":4000"))
}

func handlePing(c *fiber.Ctx) error {
	return c.SendString("Pong")
}
