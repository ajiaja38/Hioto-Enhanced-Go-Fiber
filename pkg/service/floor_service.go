package service

import (
	"go/hioto/pkg/dto"
	"go/hioto/pkg/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type FloorService struct {
	db *gorm.DB
}

func NewFloorService(db *gorm.DB) *FloorService {
	return &FloorService{db: db}
}

func (s *FloorService) CreateFloor(createDto *dto.CreateFloorDto) (*dto.ResponseFloorDto, error) {
	floor := &model.Floor{
		Name:      createDto.Name,
		CreatedAt: time.Now().In(location),
		UpdatedAt: time.Now().In(location),
	}

	if err := s.db.Create(floor).Error; err != nil {
		log.Errorf("Error creating floor: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error creating floor")
	}

	return &dto.ResponseFloorDto{
		ID:        floor.ID,
		Name:      floor.Name,
		CreatedAt: floor.CreatedAt,
		UpdatedAt: floor.UpdatedAt,
	}, nil
}

func (s *FloorService) GetAllFloors() ([]dto.ResponseAllFloorDto, error) {
	var floors []model.Floor

	if err := s.db.Find(&floors).Error; err != nil {
		log.Errorf("Error getting all floors: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error getting all floors")
	}

	var result []dto.ResponseAllFloorDto = []dto.ResponseAllFloorDto{}

	for _, floor := range floors {
		result = append(result, dto.ResponseAllFloorDto{
			ID:        floor.ID,
			Name:      floor.Name,
			CreatedAt: floor.CreatedAt,
			UpdatedAt: floor.UpdatedAt,
		})
	}

	return result, nil
}

func (s *FloorService) GetFloorByID(id string) (*dto.ResponseFloorDto, error) {
	var floor model.Floor

	if err := s.db.Preload("Rooms").First(&floor, id).Error; err != nil {
		log.Errorf("Floor not found: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Floor not found")
	}

	var rooms []dto.ResponseRoomDto
	for _, room := range floor.Rooms {
		rooms = append(rooms, dto.ResponseRoomDto{
			ID:        room.ID,
			Name:      room.Name,
			FloorID:   room.FloorID,
			CreatedAt: room.CreatedAt,
			UpdatedAt: room.UpdatedAt,
		})
	}

	return &dto.ResponseFloorDto{
		ID:        floor.ID,
		Name:      floor.Name,
		Rooms:     rooms,
		CreatedAt: floor.CreatedAt,
		UpdatedAt: floor.UpdatedAt,
	}, nil
}

func (s *FloorService) UpdateFloor(id string, updateDto *dto.UpdateFloorDto) (*dto.ResponseFloorDto, error) {
	var floor model.Floor

	if err := s.db.First(&floor, id).Error; err != nil {
		log.Errorf("Floor not found: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Floor not found")
	}

	floor.Name = updateDto.Name
	floor.UpdatedAt = time.Now().In(location)

	if err := s.db.Save(&floor).Error; err != nil {
		log.Errorf("Error updating floor: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error updating floor")
	}

	return &dto.ResponseFloorDto{
		ID:        floor.ID,
		Name:      floor.Name,
		CreatedAt: floor.CreatedAt,
		UpdatedAt: floor.UpdatedAt,
	}, nil
}

func (s *FloorService) DeleteFloor(id string) error {
	var floor model.Floor

	if err := s.db.First(&floor, id).Error; err != nil {
		log.Errorf("Floor not found: %v 💥", err)
		return fiber.NewError(fiber.StatusNotFound, "Floor not found")
	}

	if err := s.db.Delete(&floor).Error; err != nil {
		log.Errorf("Error deleting floor: %v 💥", err)
		return fiber.NewError(fiber.StatusBadRequest, "Error deleting floor")
	}

	return nil
}
