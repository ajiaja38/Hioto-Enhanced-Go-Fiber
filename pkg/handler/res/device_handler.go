package res

import (
	"fmt"
	"go/hioto/pkg/dto"
	"go/hioto/pkg/service"
	"go/hioto/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type DeviceHandler struct {
	deviceService *service.DeviceService
	validator     *validator.Validate
}

func NewDeviceHandler(deviceService *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{deviceService: deviceService, validator: validator.New()}
}

func (h *DeviceHandler) RegisterDevice(c *fiber.Ctx) error {
	var registrationDeviceDto dto.RegistrationDto

	if err := utils.ValidateRequestBody(c, h.validator, &registrationDeviceDto); err != nil {
		return err
	}

	response, err := h.deviceService.RegisterDeviceLocal(&registrationDeviceDto)

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Success register device", response)
}

func (h *DeviceHandler) GetAllDeviceHandler(c *fiber.Ctx) error {
	devices, err := h.deviceService.GetAllDevice(c.Query("type"))

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success get all device", devices)
}

func (h *DeviceHandler) GetDeviceByGuidHandler(c *fiber.Ctx) error {
	device, err := h.deviceService.GetDeviceByGuid(c.Params("guid"))

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success get device by guid", device)
}

func (h *DeviceHandler) UpdateDeviceByGuidHandler(c *fiber.Ctx) error {
	var updateDeviceDto dto.ReqUpdateDeviceDto

	if err := utils.ValidateRequestBody(c, h.validator, &updateDeviceDto); err != nil {
		return err
	}

	device, err := h.deviceService.UpdateDeviceAPI(&updateDeviceDto)

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fmt.Sprintf("Success Update Device %s", device.Guid), device)
}

func (h *DeviceHandler) DeleteDeviceByGuidHandler(c *fiber.Ctx) error {
	if err := h.deviceService.DeleteDevice(c.Params("guid")); err != nil {
		return err
	}

	return utils.SuccessResponse[any](c, fiber.StatusOK, "Success delete device", nil)
}
