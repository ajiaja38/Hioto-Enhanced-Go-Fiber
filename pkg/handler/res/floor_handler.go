package res

import (
	"go/hioto/pkg/dto"
	"go/hioto/pkg/service"
	"go/hioto/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type FloorHandler struct {
	floorService *service.FloorService
	validator    *validator.Validate
}

func NewFloorHandler(floorService *service.FloorService) *FloorHandler {
	return &FloorHandler{floorService: floorService, validator: validator.New()}
}

func (h *FloorHandler) CreateFloor(c *fiber.Ctx) error {
	var createDto dto.CreateFloorDto

	if err := utils.ValidateRequestBody(c, h.validator, &createDto); err != nil {
		return err
	}

	response, err := h.floorService.CreateFloor(&createDto)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Success create floors", response)
}

func (h *FloorHandler) GetAllFloors(c *fiber.Ctx) error {
	response, err := h.floorService.GetAllFloors()
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success get all floors", response)
}

func (h *FloorHandler) GetFloorByID(c *fiber.Ctx) error {
	response, err := h.floorService.GetFloorByID(c.Params("id"))
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success get floor by id", response)
}

func (h *FloorHandler) UpdateFloor(c *fiber.Ctx) error {
	var updateDto dto.UpdateFloorDto

	if err := utils.ValidateRequestBody(c, h.validator, &updateDto); err != nil {
		return err
	}

	response, err := h.floorService.UpdateFloor(c.Params("id"), &updateDto)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success update floor", response)
}

func (h *FloorHandler) DeleteFloor(c *fiber.Ctx) error {
	if err := h.floorService.DeleteFloor(c.Params("id")); err != nil {
		return err
	}

	return utils.SuccessResponse[any](c, fiber.StatusOK, "Success delete floor", nil)
}
