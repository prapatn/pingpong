package handlers

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handeler interface {
	Ping(ctx *fiber.Ctx) error
}

type handler struct {
	// service     services.CatalogServices
	// redisClient *redis.Client
}

func NewHandler() Handeler {
	return handler{}
}

func (h handler) Ping(c *fiber.Ctx) error {
	ballPower, _ := strconv.Atoi(c.Query("ball_power"))
	log.Printf("Original power : %v\n", ballPower)
	// For 70-90% of the original power
	modifiedPower := int(float64(ballPower) * (0.7 + rand.Float64()*0.2))
	log.Printf("Modified power : %v\n", modifiedPower)
	return c.SendString(strconv.Itoa(modifiedPower))
}
