package res

import (
	"go/hioto/pkg/dto"
	"go/hioto/pkg/service"
	"go/hioto/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ControlDeviceHandler struct {
	controlDeviceService *service.ControlDeviceService
	validator            *validator.Validate
}

func NewControlDeviceHandler(controlDeviceService *service.ControlDeviceService) *ControlDeviceHandler {
	return &ControlDeviceHandler{
		controlDeviceService: controlDeviceService,
		validator:            validator.New(),
	}
}

func (h *ControlDeviceHandler) ControlDeviceHandler(c *fiber.Ctx) error {
	var controlDto dto.ControlLocalDto

	if err := utils.ValidateRequestBody(c, h.validator, &controlDto); err != nil {
		return err
	}

	if err := h.controlDeviceService.ControlDeviceLocal(&controlDto); err != nil {
		return err
	}

	return utils.SuccessResponse[any](c, fiber.StatusOK, "Success control device ✅", nil)
}
