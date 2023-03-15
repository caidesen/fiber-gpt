package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"server/pkg/config"
)

func Setup() *fiber.App {
	// Define a new Fiber app with config.
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{"error": true, "message": e.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
		},
	})
	return app
}

func StartServer(a *fiber.App) {
	c := config.GetConfig()
	fiberConnURL := fmt.Sprintf("0.0.0.0:%d", c.Port)
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}

func WithJSONBody[T any, R any](handle func(*fiber.Ctx, T) (R, error)) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		input := new(T)
		if err := c.BodyParser(input); err != nil {
			return err
		}
		result, err := handle(c, *input)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}
