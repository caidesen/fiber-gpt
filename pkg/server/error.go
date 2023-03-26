package server

import "github.com/gofiber/fiber/v2"

func NewInputErr(msg string) error {
	return fiber.NewError(fiber.StatusBadRequest, msg)
}
