package service

import (
	"encoding/json"
	"go/hioto/config"
	"go/hioto/pkg/dto"
	"go/hioto/pkg/enum"
	messagebroker "go/hioto/pkg/handler/message_broker"
	"go/hioto/pkg/model"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var location *time.Location

func init() {
	location = time.FixedZone("WIB", 7*60*60)
}

type DeviceService struct {
	db *gorm.DB
}

func NewDeviceService(db *gorm.DB) *DeviceService {
	return &DeviceService{
		db: db,
	}
}

func (s *DeviceService) RegisterDeviceLocal(registrationDto *dto.RegistrationDto) (registrationResponse *dto.ResponseDeviceDetailDto, err error) {
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
		RoomID:    registrationDto.RoomID,
		LastSeen:  time.Now().In(location),
		CreatedAt: time.Now().In(location),
		UpdatedAt: time.Now().In(location),
	}

	if err = s.db.Create(registration).Error; err != nil {
		log.Errorf("Error creating device: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error creating device")
	}

	deviceResponse, err := s.GetDeviceByGuid(registration.Guid)
	if err != nil {
		log.Errorf("Error fetching created device: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error fetching created device")
	}

	bodyToCloud := dto.ResCloudDeviceDto{
		ResponseDeviceDetailDto: *deviceResponse,
		MacServer:               config.MAC_ADDRESS.GetValue(),
	}

	jsonBody, err := json.Marshal(bodyToCloud)

	if err != nil {
		log.Errorf("Error marshaling JSON: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error marshaling JSON")
	}

	messagebroker.PublishToRmq(
		config.RMQ_CLOUD_INSTANCE.GetValue(),
		jsonBody,
		config.REGISTER_RES_CLOUD.GetValue(),
		config.EXCHANGE_DIRECT.GetValue(),
	)

	log.Infof("Device successfully registered from local: %s ✅", registration.Name)

	return deviceResponse, nil
}

func (s *DeviceService) RegisterDeviceCloud(registrationDto *dto.RegistrationDto) {
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

func (s *DeviceService) GetAllDevice(deviceType, floorID, roomID string) ([]dto.ResponseDeviceListDto, error) {
	var devices []model.Registration

	var query *gorm.DB = s.db.Preload("Room.Floor")

	if deviceType != "" {
		query = query.Where("type = ?", deviceType)
	}

	if floorID != "" {
		query = query.Joins("JOIN rooms ON rooms.id = registrations.room_id").
			Where("rooms.floor_id = ?", floorID)
	}

	if roomID != "" {
		query = query.Where("room_id = ?", roomID)
	}

	query = query.Order("registrations.created_at DESC")

	if err := query.Find(&devices).Error; err != nil {
		log.Errorf("Error getting all device: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error when getting all device")
	}

	var result []dto.ResponseDeviceListDto = []dto.ResponseDeviceListDto{}

	for _, device := range devices {
		var roomName *string
		if device.Room != nil {
			roomName = &device.Room.Name
		}

		result = append(result, dto.ResponseDeviceListDto{
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
			RoomName:     roomName,
		})
	}

	return result, nil
}

func (s *DeviceService) GetDeviceByGuid(guid string) (*dto.ResponseDeviceDetailDto, error) {
	var device model.Registration

	if err := s.db.Preload("Room.Floor").Where("guid = ?", guid).First(&device).Error; err != nil {
		log.Errorf("Device not found: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Device not found")
	}

	var (
		roomID    *uint
		roomName  *string
		floorID   *uint
		floorName *string
	)

	if device.Room != nil {
		roomID = &device.Room.ID
		roomName = &device.Room.Name
		floorID = &device.Room.Floor.ID
		floorName = &device.Room.Floor.Name
	}

	return &dto.ResponseDeviceDetailDto{
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
		RoomID:       roomID,
		RoomName:     roomName,
		FloorID:      floorID,
		FloorName:    floorName,
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

func (s *DeviceService) UpdateDeviceAPI(updateDto *dto.ReqUpdateDeviceDto) (*dto.ResponseDeviceDetailDto, error) {
	_, err := s.updateQuery(updateDto)

	if err != nil {
		log.Errorf("Error updating device: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error updating device")
	}

	deviceResponse, err := s.GetDeviceByGuid(updateDto.Guid)
	if err != nil {
		return nil, err
	}

	bodyToCloud := dto.ResCloudDeviceDto{
		ResponseDeviceDetailDto: *deviceResponse,
		MacServer:               config.MAC_ADDRESS.GetValue(),
	}

	jsonBody, err := json.Marshal(bodyToCloud)

	if err != nil {
		log.Errorf("Error marshaling JSON: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Error marshaling JSON")
	}

	messagebroker.PublishToRmq(
		config.RMQ_CLOUD_INSTANCE.GetValue(),
		jsonBody,
		config.UPDATE_RES_CLOUD.GetValue(),
		config.EXCHANGE_DIRECT.GetValue(),
	)

	return deviceResponse, nil
}

func (s *DeviceService) updateQuery(updateDto *dto.ReqUpdateDeviceDto) (*model.Registration, error) {
	updateQuery := s.db.Exec(`
        UPDATE registrations
        SET mac = ?,
            type = ?,
            quantity = ?,
            name = ?,
            version = ?,
            minor = ?,
            room_id = ?,
            updated_at = ?
        WHERE guid = ?
	`, updateDto.Mac, updateDto.Type, updateDto.Quantity, updateDto.Name, updateDto.Version, updateDto.Minor, updateDto.RoomID, time.Now().In(location), updateDto.Guid)

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
		RoomID:    updateDto.RoomID,
		UpdatedAt: time.Now().In(location),
	}, nil
}

func (s *DeviceService) DeleteDeviceRMQ(guid string) {
	if err := s.DeleteDevice(guid); err != nil {
		log.Errorf("Error deleting device: %v 💥", err)
		return
	}
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
		MacServer: config.MAC_ADDRESS.GetValue(),
	}

	jsonBody, err := json.Marshal(payloadToCloud)

	if err != nil {
		log.Errorf("Error marshaling JSON: %v 💥", err)
		return fiber.NewError(fiber.StatusBadRequest, "Error marshaling JSON")
	}

	messagebroker.PublishToRmq(
		config.RMQ_CLOUD_INSTANCE.GetValue(),
		jsonBody,
		config.DELETE_RES_CLOUD.GetValue(),
		config.EXCHANGE_DIRECT.GetValue(),
	)

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
