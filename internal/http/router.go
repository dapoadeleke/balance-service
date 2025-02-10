package http

import "github.com/gofiber/fiber/v2"

func BuildRoutes(app *fiber.App, h *Handler) {
	user := app.Group("/user/:userId")

	user.Post("/transaction", h.PostTransaction)
	user.Get("/balance", h.GetBalance)
}
