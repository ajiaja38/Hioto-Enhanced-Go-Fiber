package service

import (
	"encoding/json"
	"go/hioto/pkg/dto"
	"go/hioto/pkg/enum"
	messagebroker "go/hioto/pkg/handler/message_broker"
	"go/hioto/pkg/model"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type ControlDeviceService struct {
	db *gorm.DB
}

func NewControlDeviceService(db *gorm.DB) *ControlDeviceService {
	return &ControlDeviceService{
		db: db,
	}
}

func (s *ControlDeviceService) ControlDeviceCloud(controlDto *dto.ControlDto) {
	var device model.Registration
	value := strings.Split(controlDto.Message, "#")

	if err := s.db.Where("guid = ?", value[0]).First(&device).Error; err != nil {
		log.Errorf("Device not found: %v 💥", err)
		return
	}

	if controlDto.Type == enum.SENSOR {
		messagebroker.PublishToRmq(
			os.Getenv("RMQ_URI"),
			"RMQ Local",
			[]byte(controlDto.Message),
			os.Getenv("SENSOR_QUEUE"),
			os.Getenv("EXCHANGE_TOPIC"),
		)
		return
	}

	location, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		log.Errorf("Failed to load location: %v 💥", err)
		return
	}

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Errorf("Transaction rollback due to panic: %v 💥", r)
		} else {
			if err := tx.Commit().Error; err != nil {
				log.Errorf("Error committing transaction: %v 💥", err)
				tx.Rollback()
			}
		}
	}()

	if tx.Error != nil {
		log.Errorf("Error starting transaction: %v 💥", tx.Error)
		return
	}

	status := value[1] == "1"

	if err := tx.Model(&device).Updates(map[string]any{
		"status":     status,
		"updated_at": time.Now().In(location),
	}).Error; err != nil {
		log.Errorf("Error updating registration: %v 💥", err)
		tx.Rollback()
		return
	}

	logEntry := model.LogAktuator{
		InputGuid: value[0],
		Name:      device.Name,
		Value:     value[1],
		Time:      time.Now().In(location),
	}

	if err := tx.Create(&logEntry).Error; err != nil {
		log.Errorf("Error inserting log: %v 💥", err)
		tx.Rollback()
		return
	}

	log.Info("Transaction was committed successfully ✅")

	messagebroker.PublishToRoutingKey(
		os.Getenv("RMQ_URI"),
		"RMQ Local",
		[]byte(controlDto.Message),
		os.Getenv("EXCHANGE_TOPIC"),
		os.Getenv("AKTUATOR_ROUTING_KEY"),
	)
}

func (s *ControlDeviceService) ControlDeviceLocal(controlDto *dto.ControlLocalDto) error {
	var device model.Registration
	value := strings.Split(controlDto.Message, "#")

	if err := s.db.Where("guid = ?", value[0]).First(&device).Error; err != nil {
		log.Errorf("Device not found: %v 💥", err)
		return fiber.NewError(fiber.StatusNotFound, "Device not found")
	}

	if controlDto.Type == enum.SENSOR {
		messagebroker.PublishToRmq(
			os.Getenv("RMQ_URI"),
			"RMQ Local",
			[]byte(controlDto.Message),
			os.Getenv("SENSOR_QUEUE"),
			os.Getenv("EXCHANGE_TOPIC"),
		)
		return nil
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Errorf("Failed to load location: %v 💥", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load location")
	}

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Errorf("Transaction rollback due to panic: %v 💥", r)
		} else {
			if err := tx.Commit().Error; err != nil {
				log.Errorf("Error committing transaction: %v 💥", err)
				tx.Rollback()
			}
		}
	}()

	if tx.Error != nil {
		log.Errorf("Error starting transaction: %v 💥", tx.Error)
		return fiber.NewError(fiber.StatusBadRequest, "Error starting transaction")
	}

	status := value[1] == "1"

	if err := tx.Model(&device).Updates(map[string]any{
		"status":     status,
		"updated_at": time.Now().In(location),
	}).Error; err != nil {
		log.Errorf("Error updating registration: %v 💥", err)
		tx.Rollback()
		return fiber.NewError(fiber.StatusBadRequest, "Error updating registration")
	}

	logEntry := model.LogAktuator{
		InputGuid: value[0],
		Name:      device.Name,
		Value:     value[1],
		Time:      time.Now().In(location),
	}

	if err := tx.Create(&logEntry).Error; err != nil {
		log.Errorf("Error inserting log: %v 💥", err)
		tx.Rollback()
		return fiber.NewError(fiber.StatusBadRequest, "Error inserting log")
	}

	log.Info("Transaction for local committed successfully ✅")

	messagebroker.PublishToRoutingKey(
		os.Getenv("RMQ_URI"),
		"RMQ Local",
		[]byte(controlDto.Message),
		os.Getenv("EXCHANGE_TOPIC"),
		os.Getenv("AKTUATOR_ROUTING_KEY"),
	)

	bodyToCloud := dto.ResCloudDeviceDto{
		ResponseDeviceDto: dto.ResponseDeviceDto{
			ID:           device.ID,
			Guid:         device.Guid,
			Mac:          device.Mac,
			Type:         device.Type,
			Quantity:     device.Quantity,
			Name:         device.Name,
			Version:      device.Version,
			Minor:        device.Minor,
			Status:       device.Status,
			StatusDevice: string(device.StatusDevice),
			LastSeen:     device.LastSeen,
			CreatedAt:    device.CreatedAt,
			UpdatedAt:    device.UpdatedAt,
		},
		MacServer: os.Getenv("MAC_ADDRESS"),
	}

	jsonBody, err := json.Marshal(bodyToCloud)

	if err != nil {
		log.Errorf("Error marshaling JSON: %v 💥", err)
		return fiber.NewError(fiber.StatusBadRequest, "Error marshaling JSON")
	}

	messagebroker.PublishToRmq(os.Getenv("RMQ_HIOTO"), "Biznet Hioto", jsonBody, os.Getenv("UPDATE_RES_CLOUD"), "amq.direct")

	return nil
}
