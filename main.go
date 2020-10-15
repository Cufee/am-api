package main

import (
	h "github.com/cufee/am-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// WG login routes
	app.Get("/login/redirect/:intentID", h.HandleWargamingRedirect)
	app.Get("/login", h.HandleWargamingLogin)

	app.Listen(":80")
}
