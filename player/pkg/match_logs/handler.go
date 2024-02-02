package matchlogs

import (
	"player/pkg/domain"

	"github.com/gofiber/fiber/v2"
)

type MatchLogHandler struct {
	matchLogUsecase domain.MatchLogUsecase
}

func NewMatchLogHandler(matchLogRoute fiber.Router, matchLogUsecase domain.MatchLogUsecase) {

	handler := &MatchLogHandler{
		matchLogUsecase: matchLogUsecase,
	}

	matchLogRoute.Get("new-match", handler.NewMatch())
	matchLogRoute.Get("match", handler.GetLastMatch())
	matchLogRoute.Get("match/:id", handler.GetMatchById())
}

func (h *MatchLogHandler) NewMatch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log, err := h.matchLogUsecase.InsertLog()

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(log)
	}
}

func (h *MatchLogHandler) GetLastMatch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lastMatch, err := h.matchLogUsecase.GetLastMatch()

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(lastMatch)
	}
}

func (h *MatchLogHandler) GetMatchById() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		matchLog, err := h.matchLogUsecase.GetMatchById(id)

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(matchLog)
	}
}
