package route

import "github.com/gofiber/fiber/v2"

type Route struct {
	Name        string
	Description string
	Method      string
	Path        string
	Middleware  []fiber.Handler
	HandlerFunc fiber.Handler
	Test        bool
}
