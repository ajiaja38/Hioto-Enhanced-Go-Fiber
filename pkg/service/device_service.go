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

var location *time.Location

func init() {
	location, _ = time.LoadLocation("Asia/Jakarta")

	if location == nil {
		location = time.UTC
	}
}

type DeviceService struct {
	db *gorm.DB
}

func NewDeviceService(db *gorm.DB) *DeviceService {
	return &DeviceService{
		db: db,
	}
}

func (s *DeviceService) RegisterDeviceLocal(registrationDto *dto.RegistrationDto) (registrationResponse *dto.ResponseDeviceDto, err error) {
	var status string

	if registrationDto.Type == enum.AKTUATOR {
		status = "0"
	}

	registration := &model.Registration{
		Guid:      registrationDto.Guid,
		Mac:       registrationDto.Mac,
		Type:      registrationDto.Type,
		Name:      registrationDto.Name,
		Quantity:  registrationDto.Quantity,
		Status:    status,
		Version:   registrationDto.Version,
		Minor:     registrationDto.Minor,
		LastSeen:  time.Now().In(location),
		CreatedAt: time.Now().In(location),
		UpdatedAt: time.Now().In(location),
	}

	if err = s.db.Create(registration).Error; err != nil {
		log.Errorf("Error creating device: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error creating device")
	}

	bodyToCloud := dto.ResCloudDeviceDto{
		ResponseDeviceDto: dto.ResponseDeviceDto{
			ID:           registration.ID,
			Guid:         registration.Guid,
			Mac:          registration.Mac,
			Type:         registration.Type,
			Quantity:     registration.Quantity,
			Name:         registration.Name,
			Version:      registration.Version,
			Minor:        registration.Minor,
			Status:       registration.Status,
			StatusDevice: string(registration.StatusDevice),
			CreatedAt:    registration.CreatedAt,
			UpdatedAt:    registration.UpdatedAt,
		},
		MacServer: os.Getenv("MAC_ADDRESS"),
	}

	jsonBody, err := json.Marshal(bodyToCloud)

	if err != nil {
		log.Errorf("Error marshaling JSON: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error marshaling JSON")
	}

	messagebroker.PublishToRmq(os.Getenv("RMQ_HIOTO"), "Biznet Hioto", jsonBody, os.Getenv("REGISTER_RES_CLOUD"), "amq.direct")

	registrationResponse = &dto.ResponseDeviceDto{
		ID:           registration.ID,
		Guid:         registration.Guid,
		Mac:          registration.Mac,
		Type:         registration.Type,
		Quantity:     registration.Quantity,
		Name:         registration.Name,
		Version:      registration.Version,
		Minor:        registration.Minor,
		Status:       registration.Status,
		StatusDevice: string(registration.StatusDevice),
		LastSeen:     registration.LastSeen,
		CreatedAt:    registration.CreatedAt,
		UpdatedAt:    registration.UpdatedAt,
	}

	log.Infof("Device successfully registered from local: %s ✅", registration.Name)

	return registrationResponse, nil
}

func (s *DeviceService) RegisterDeviceCloud(registrationDto *dto.ReqCloudDeviceDto) {
	var status string

	if registrationDto.Type == enum.AKTUATOR {
		status = "0"
	}

	registration := &model.Registration{
		Guid:      registrationDto.Guid,
		Mac:       registrationDto.Mac,
		Type:      registrationDto.Type,
		Name:      registrationDto.Name,
		Quantity:  registrationDto.Quantity,
		Status:    status,
		Version:   registrationDto.Version,
		Minor:     registrationDto.Minor,
		CreatedAt: time.Now().In(location),
		UpdatedAt: time.Now().In(location),
	}

	if err := s.db.Create(registration).Error; err != nil {
		log.Errorf("Error creating device: %v 💥", err)
		return
	}

	log.Infof("Your Device successfully registered from cloud: %s ✅", registration.Name)
}

func (s *DeviceService) GetAllDevice(deviceType string) ([]dto.ResponseDeviceDto, error) {
	var devices []model.Registration

	var query *gorm.DB = s.db

	if deviceType != "" {
		query = query.Where("type = ?", deviceType)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&devices).Error; err != nil {
		log.Errorf("Error getting all device: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error when getting all device")
	}

	var result []dto.ResponseDeviceDto

	for _, device := range devices {
		result = append(result, dto.ResponseDeviceDto{
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
		})
	}

	return result, nil
}

func (s *DeviceService) GetDeviceByGuid(guid string) (*dto.ResponseDeviceDto, error) {
	var device model.Registration

	deviceRaw := s.db.Raw(`SELECT * FROM registrations WHERE guid = ?`, guid).Scan(&device)

	if deviceRaw.RowsAffected == 0 {
		log.Errorf("Device not found: %v 💥", deviceRaw.Error)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Device not found")
	}

	return &dto.ResponseDeviceDto{
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
	}, nil
}

func (s *DeviceService) UpdateDeviceRMQCloud(updateDto *dto.ReqUpdateDeviceDto) {
	var device model.Registration

	deviceRaw := s.db.Raw(`SELECT * FROM registrations WHERE guid = ?`, updateDto.Guid).Scan(&device)

	if deviceRaw.RowsAffected == 0 {
		log.Errorf("Device not found: %v 💥", deviceRaw.Error)
		return
	}

	deviceUpdated, err := s.updateQuery(updateDto)

	if err != nil {
		log.Errorf("Error updating device: %v 💥", err)
		return
	}

	log.Infof("Device successfully updated: %s ✅", deviceUpdated.Name)
}

func (s *DeviceService) UpdateDeviceAPI(updateDto *dto.ReqUpdateDeviceDto) (deviceUpdated *model.Registration, err error) {
	deviceUpdated, err = s.updateQuery(updateDto)

	if err != nil {
		log.Errorf("Error updating device: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error updating device")
	}

	bodyToCloud := dto.ResCloudDeviceDto{
		ResponseDeviceDto: dto.ResponseDeviceDto{
			ID:           deviceUpdated.ID,
			Guid:         deviceUpdated.Guid,
			Mac:          deviceUpdated.Mac,
			Type:         deviceUpdated.Type,
			Quantity:     deviceUpdated.Quantity,
			Name:         deviceUpdated.Name,
			Version:      deviceUpdated.Version,
			Minor:        deviceUpdated.Minor,
			Status:       deviceUpdated.Status,
			StatusDevice: string(deviceUpdated.StatusDevice),
			LastSeen:     deviceUpdated.LastSeen,
			CreatedAt:    deviceUpdated.CreatedAt,
			UpdatedAt:    deviceUpdated.UpdatedAt,
		},
		MacServer: os.Getenv("MAC_ADDRESS"),
	}

	jsonBody, err := json.Marshal(bodyToCloud)

	if err != nil {
		log.Errorf("Error marshaling JSON: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error marshaling JSON")
	}

	messagebroker.PublishToRmq(os.Getenv("RMQ_HIOTO"), "Biznet Hioto", jsonBody, os.Getenv("UPDATE_RES_CLOUD"), "amq.direct")

	return deviceUpdated, nil
}

func (s *DeviceService) updateQuery(updateDto *dto.ReqUpdateDeviceDto) (*model.Registration, error) {
	updateQuery := s.db.Exec(`
    UPDATE registrations
    SET mac = ?, type = ?, quantity = ?, name = ?, version = ?, minor = ?, updated_at = ?
    WHERE guid = ?
	`, updateDto.Mac, updateDto.Type, updateDto.Quantity, updateDto.Name, updateDto.Version, updateDto.Minor, time.Now().In(location), updateDto.Guid)

	if updateQuery.RowsAffected == 0 {
		log.Errorf("Error updating device: %v 💥", updateQuery.Error)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error updating device")
	}

	return &model.Registration{
		Mac:       updateDto.Mac,
		Type:      updateDto.Type,
		Quantity:  updateDto.Quantity,
		Name:      updateDto.Name,
		Version:   updateDto.Version,
		Minor:     updateDto.Minor,
		UpdatedAt: time.Now().In(location),
	}, nil
}

func (s *DeviceService) DeleteDeviceRMQ(guid string) {
	err := s.DeleteDevice(guid)

	if err != nil {
		log.Errorf("Error deleting device: %v 💥", err)
		return
	}
}

func (s *DeviceService) UpdateStatusDevice(guid string, status string) {
	var device model.Registration

	err := s.db.First(&device, "guid = ?", guid).Error

	if err != nil {
		log.Errorf("Device not found: %v 💥", err)
		return
	}

	device.StatusDevice = enum.EDeviceStatus(status)
	device.LastSeen = time.Now().In(location)

	err = s.db.Save(&device).Error

	if err != nil {
		log.Errorf("Error updating status device: %v 💥", err)
		return
	}

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
		return
	}

	messagebroker.PublishToRmq(os.Getenv("RMQ_HIOTO"), "Biznet Hioto", jsonBody, os.Getenv("UPDATE_RES_CLOUD"), "amq.direct")

	log.Infof("Status device successfully updated: %s ✅", status)
}

func (s *DeviceService) DeleteDevice(guid string) error {
	var device model.Registration

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
		return fiber.NewError(fiber.StatusBadGateway, "Error starting transaction")
	}

	if err := tx.Where("guid = ?", guid).First(&device).Error; err != nil {
		log.Error("Device not found 💥")
		tx.Rollback()
		return fiber.NewError(fiber.StatusNotFound, "Device not found")
	}

	if err := tx.Delete(&device).Error; err != nil {
		log.Errorf("Error deleting device: %v 💥", err)
		tx.Rollback()
		return fiber.NewError(fiber.StatusBadRequest, "Error deleting device")
	}

	switch device.Type {
	case enum.SENSOR:
		if err := tx.Where("input_guid = ?", guid).Delete(&model.RuleDevice{}).Error; err != nil {
			log.Errorf("Error deleting rule devices: %v 💥", err)
			tx.Rollback()
			return fiber.NewError(fiber.StatusBadRequest, "Error deleting rule devices")
		}
	case enum.AKTUATOR:
		if err := tx.Where("output_guid = ?", guid).Delete(&model.RuleDevice{}).Error; err != nil {
			log.Errorf("Error deleting rule devices: %v 💥", err)
			tx.Rollback()
			return fiber.NewError(fiber.StatusBadRequest, "Error deleting rule devices")
		}
	}

	payloadToCloud := dto.ReqDeleteDeviceToCloudDto{
		Guid:      guid,
		MacServer: os.Getenv("MAC_ADDRESS"),
	}

	jsonBody, err := json.Marshal(payloadToCloud)

	if err != nil {
		log.Errorf("Error marshaling JSON: %v 💥", err)
		return fiber.NewError(fiber.StatusBadRequest, "Error marshaling JSON")
	}

	messagebroker.PublishToRmq(os.Getenv("RMQ_HIOTO"), "Biznet Hioto", jsonBody, os.Getenv("DELETE_RES_CLOUD"), "amq.direct")

	log.Infof("Device successfully deleted: %s ✅", guid)
	return nil
}

func (s *DeviceService) CheckInactiveDevice() {
	ticker := time.NewTicker(60 * time.Second)

	for {
		<-ticker.C
		treshold := time.Now().Add(-10 * time.Second)

		err := s.db.Model(&model.Registration{}).
			Where("last_seen < ?", treshold).
			Update("status_device", enum.OFF).Error

		if err != nil {
			log.Errorf("Error checking for inactive device: %v 💥", err)
		} else {
			log.Infof("Inactive devices marked as offline 🔻")
		}
	}
}

func (s *DeviceService) UpdateStatusAsMonitoring(guid, payload string) {
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

	var device model.Registration

	if err := tx.Where("guid = ?", guid).First(&device).Error; err != nil {
		log.Errorf("Device not found: %v 💥", err)
		return
	}

	device.Status = payload

	if err := tx.Save(&device).Error; err != nil {
		log.Errorf("Error updating status device: %v 💥", err)
		return
	}

	MonitoringHistories := &model.MonitoringHistory{
		DeviceGuid: device.Guid,
		DeviceName: device.Name,
		DeviceType: device.Type,
		Value:      payload,
		Time:       time.Now().In(location),
	}

	if err := tx.Create(MonitoringHistories).Error; err != nil {
		log.Errorf("Error creating monitoring history: %v 💥", err)
	}

	log.Infof("Data Monitoring device %s successfully updated: %s ✅", strings.Split(device.Name, "-")[0], payload)
}
