package res

import (
	"go/hioto/pkg/dto"
	"go/hioto/pkg/service"
	"go/hioto/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct {
	roomService *service.RoomService
	validator   *validator.Validate
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService, validator: validator.New()}
}

func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var createDto dto.CreateRoomDto

	if err := utils.ValidateRequestBody(c, h.validator, &createDto); err != nil {
		return err
	}

	response, err := h.roomService.CreateRoom(&createDto)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Success create room", response)
}

func (h *RoomHandler) GetAllRooms(c *fiber.Ctx) error {
	response, err := h.roomService.GetAllRooms()
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success get all rooms", response)
}

func (h *RoomHandler) GetRoomsByFloorID(c *fiber.Ctx) error {
	response, err := h.roomService.GetRoomsByFloorID(c.Params("floor_id"))
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success get rooms by floor id", response)
}

func (h *RoomHandler) GetRoomByID(c *fiber.Ctx) error {
	response, err := h.roomService.GetRoomByID(c.Params("id"))
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success get room by id", response)
}

func (h *RoomHandler) UpdateRoom(c *fiber.Ctx) error {
	var updateDto dto.UpdateRoomDto

	if err := utils.ValidateRequestBody(c, h.validator, &updateDto); err != nil {
		return err
	}

	response, err := h.roomService.UpdateRoom(c.Params("id"), &updateDto)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success update room", response)
}

func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
	if err := h.roomService.DeleteRoom(c.Params("id")); err != nil {
		return err
	}

	return utils.SuccessResponse[any](c, fiber.StatusOK, "Success delete room", nil)
}
