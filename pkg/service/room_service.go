package service

import (
	"go/hioto/pkg/dto"
	"go/hioto/pkg/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type RoomService struct {
	db *gorm.DB
}

func NewRoomService(db *gorm.DB) *RoomService {
	return &RoomService{db: db}
}

func (s *RoomService) CreateRoom(createDto *dto.CreateRoomDto) (*dto.ResponseRoomDto, error) {
	room := &model.Room{
		Name:      createDto.Name,
		FloorID:   createDto.FloorID,
		CreatedAt: time.Now().In(location),
		UpdatedAt: time.Now().In(location),
	}

	if err := s.db.Create(room).Error; err != nil {
		log.Errorf("Error creating room: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error creating room")
	}

	return &dto.ResponseRoomDto{
		ID:        room.ID,
		Name:      room.Name,
		FloorID:   room.FloorID,
		CreatedAt: room.CreatedAt,
		UpdatedAt: room.UpdatedAt,
	}, nil
}

func (s *RoomService) GetAllRooms() ([]dto.ResponseRoomDto, error) {
	var rooms []model.Room

	if err := s.db.Find(&rooms).Error; err != nil {
		log.Errorf("Error getting all rooms: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error getting all rooms")
	}

	var result []dto.ResponseRoomDto

	for _, room := range rooms {
		result = append(result, dto.ResponseRoomDto{
			ID:        room.ID,
			Name:      room.Name,
			FloorID:   room.FloorID,
			CreatedAt: room.CreatedAt,
			UpdatedAt: room.UpdatedAt,
		})
	}

	return result, nil
}

func (s *RoomService) GetRoomsByFloorID(floorID string) ([]dto.ResponseRoomDto, error) {
	var rooms []model.Room = []model.Room{}

	if err := s.db.Where("floor_id = ?", floorID).Find(&rooms).Error; err != nil {
		log.Errorf("Error getting rooms by floor id: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error getting rooms by floor id")
	}

	var result []dto.ResponseRoomDto = []dto.ResponseRoomDto{}

	for _, room := range rooms {
		result = append(result, dto.ResponseRoomDto{
			ID:        room.ID,
			Name:      room.Name,
			FloorID:   room.FloorID,
			CreatedAt: room.CreatedAt,
			UpdatedAt: room.UpdatedAt,
		})
	}

	return result, nil
}

func (s *RoomService) GetRoomByID(id string) (*dto.ResponseRoomDto, error) {
	var room model.Room

	if err := s.db.First(&room, id).Error; err != nil {
		log.Errorf("Room not found: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Room not found")
	}

	return &dto.ResponseRoomDto{
		ID:        room.ID,
		Name:      room.Name,
		FloorID:   room.FloorID,
		CreatedAt: room.CreatedAt,
		UpdatedAt: room.UpdatedAt,
	}, nil
}

func (s *RoomService) UpdateRoom(id string, updateDto *dto.UpdateRoomDto) (*dto.ResponseRoomDto, error) {
	var room model.Room

	if err := s.db.First(&room, id).Error; err != nil {
		log.Errorf("Room not found: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Room not found")
	}

	room.Name = updateDto.Name
	room.FloorID = updateDto.FloorID
	room.UpdatedAt = time.Now().In(location)

	if err := s.db.Save(&room).Error; err != nil {
		log.Errorf("Error updating room: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error updating room")
	}

	return &dto.ResponseRoomDto{
		ID:        room.ID,
		Name:      room.Name,
		FloorID:   room.FloorID,
		CreatedAt: room.CreatedAt,
		UpdatedAt: room.UpdatedAt,
	}, nil
}

func (s *RoomService) DeleteRoom(id string) error {
	var room model.Room

	if err := s.db.First(&room, id).Error; err != nil {
		log.Errorf("Room not found: %v 💥", err)
		return fiber.NewError(fiber.StatusNotFound, "Room not found")
	}

	if err := s.db.Delete(&room).Error; err != nil {
		log.Errorf("Error deleting room: %v 💥", err)
		return fiber.NewError(fiber.StatusBadRequest, "Error deleting room")
	}

	return nil
}
