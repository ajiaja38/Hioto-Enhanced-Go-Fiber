package service

import (
	"encoding/json"
	"fmt"
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
		s.ControlSensor(value[0], value[1])
		return
	}

	location = time.FixedZone("WIB", 7*60*60)

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

	log.Info("Transaction committed successfully ✅")

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
		s.ControlSensor(value[0], value[1])
		return nil
	}

	location = time.FixedZone("WIB", 7*60*60)

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Errorf("Transaction rollback due to panic: %v 💥", r)
		} else {
			if err := tx.Commit().Error; err != nil {
				log.Errorf("Error committing the transaction: %v 💥", err)
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
		return fiber.NewError(fiber.StatusBadRequest, "Error updating registration device")
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

func (s *ControlDeviceService) ControlSensor(guid, value string) {
	var ruleDevices []model.RuleDevice

	if err := s.db.Where("input_guid = ?", guid).Where("input_value = ?", value).Find(&ruleDevices).Error; err != nil {
		log.Errorf("Failed to fetch rule devices: %v 💥", err)
	}

	for _, ruleDevice := range ruleDevices {
		var aktuator model.Registration

		location = time.FixedZone("WIB", 7*60*60)

		messageToAktuator := fmt.Sprintf("%s#%s", ruleDevice.OutputGuid, ruleDevice.OutputValue)

		if err := s.db.Where("guid = ?", ruleDevice.OutputGuid).First(&aktuator).Error; err != nil {
			log.Errorf("Failed to fetch aktuator: %v 💥", err)
			continue
		}

		aktuator.Status = ruleDevice.OutputValue
		aktuator.UpdatedAt = time.Now().In(location)

		if err := s.db.Save(&aktuator).Error; err != nil {
			log.Errorf("Failed update aktuator status: %v 💥", err)
			continue
		}

		logSensor := model.Log{
			InputGuid:   guid,
			InputName:   aktuator.Name,
			InputValue:  ruleDevice.InputGuid,
			OutputGuid:  aktuator.Guid,
			OutputValue: ruleDevice.OutputValue,
			Time:        time.Now().In(location),
		}

		if err := s.db.Create(&logSensor).Error; err != nil {
			log.Errorf("Failed to insert log: %v 💥", err)
			continue
		}

		messagebroker.PublishToRoutingKey(
			os.Getenv("RMQ_URI"),
			"RMQ Local",
			[]byte(messageToAktuator),
			os.Getenv("EXCHANGE_TOPIC"),
			os.Getenv("AKTUATOR_ROUTING_KEY"),
		)

		log.Infof("Sensor rule executed for aktuator %s with value %s ✅", aktuator.Name, ruleDevice.OutputValue)
	}
}
