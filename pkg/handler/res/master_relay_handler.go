package res

import (
	"go/hioto/pkg/service"
	"go/hioto/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type MasterRelayHandler struct {
	service *service.MasterRelayService
}

func NewMasterRelayHandler(s *service.MasterRelayService) *MasterRelayHandler {
	return &MasterRelayHandler{service: s}
}

func (h *MasterRelayHandler) GetStats(c *fiber.Ctx) error {
	guid := c.Params("guid")
	period := c.Query("period", "daily") // Opsi: daily, weekly, monthly

	data, err := h.service.GetStatistics(guid, period)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success get master relay statistics", data)
}