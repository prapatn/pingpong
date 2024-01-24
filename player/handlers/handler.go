package handlers

import (
	"player/services"

	"github.com/gofiber/fiber/v2"
)

type Handeler interface {
	NewMatch(ctx *fiber.Ctx) error
	GetLastMatch(c *fiber.Ctx) error
}

type handler struct {
	service services.Services
}

func NewHandler(service services.Services) Handeler {
	return handler{service: service}
}

func (h handler) NewMatch(c *fiber.Ctx) error {
	err := h.service.InsertLog()

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h handler) GetLastMatch(c *fiber.Ctx) error {
	lastMatch, err := h.service.GetLastMatch()

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Status(fiber.StatusOK)
	return c.JSON(lastMatch)
}
