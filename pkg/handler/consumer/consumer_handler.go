package consumer

import (
	"encoding/json"
	"go/hioto/pkg/dto"
	"go/hioto/pkg/service"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var validate = validator.New()

type ConsumerHandler struct {
	ruleService          *service.RuleService
	deviceService        *service.DeviceService
	controlDeviceService *service.ControlDeviceService
	validator            *validator.Validate
}

func NewConsumerHandler(ruleService *service.RuleService, deviceService *service.DeviceService, controlDeviceService *service.ControlDeviceService) *ConsumerHandler {
	return &ConsumerHandler{
		ruleService:          ruleService,
		deviceService:        deviceService,
		controlDeviceService: controlDeviceService,
		validator:            validator.New(),
	}
}

func (h *ConsumerHandler) RegistrationHandler(message []byte, db *gorm.DB) {
	var registrationDto dto.RegistrationDto
	err := json.Unmarshal(message, &registrationDto)

	if err != nil {
		log.Errorf("Failed to unmarshal registration message: %v", err)
		return
	}

	h.deviceService.RegisterDeviceLocal(&registrationDto)
}

func (h *ConsumerHandler) RegistrationFromCloudHandler(message []byte, db *gorm.DB) {
	var registrationDto dto.ReqCloudDeviceDto
	err := json.Unmarshal(message, &registrationDto)

	if err != nil {
		log.Errorf("Failed to unmarshal registration message: %v", err)
		return
	}

	if registrationDto.MacServer != os.Getenv("MAC_ADDRESS") {
		log.Errorf("Invalid mac server: %v", registrationDto.MacServer)
		return
	}

	h.deviceService.RegisterDeviceCloud(&registrationDto)
}

func (h *ConsumerHandler) UpdateDeviceFromCloudHandler(message []byte, db *gorm.DB) {
	var updateDeviceFromCloudDto dto.ReqUpdateDeviceDtoCloud

	err := json.Unmarshal(message, &updateDeviceFromCloudDto)

	if err != nil {
		log.Errorf("Failed to unmarshal update device message: %v", err)
		return
	}

	if updateDeviceFromCloudDto.MacServer != os.Getenv("MAC_ADDRESS") {
		log.Errorf("Invalid mac server: %v", updateDeviceFromCloudDto.MacServer)
		return
	}

	updateRawDeviceDto := dto.ReqUpdateDeviceDto{
		Guid:     updateDeviceFromCloudDto.Guid,
		Mac:      updateDeviceFromCloudDto.Mac,
		Type:     updateDeviceFromCloudDto.Type,
		Quantity: updateDeviceFromCloudDto.Quantity,
		Name:     updateDeviceFromCloudDto.Name,
		Version:  updateDeviceFromCloudDto.Version,
		Minor:    updateDeviceFromCloudDto.Minor,
	}

	h.deviceService.UpdateDeviceRMQCloud(&updateRawDeviceDto)
}

func (h *ConsumerHandler) RulesHandler(message []byte, db *gorm.DB) {
	var createRuleDto dto.CreateRuleDto
	err := json.Unmarshal(message, &createRuleDto)

	if err != nil {
		log.Errorf("Failed to unmarshal rule message: %v", err)
		return
	}

	h.ruleService.CreateRules(&createRuleDto)
}

func (h *ConsumerHandler) ControlHandler(message []byte, db *gorm.DB) {
	var controlDeviceDto dto.ControlDto
	err := json.Unmarshal(message, &controlDeviceDto)

	if err != nil {
		log.Errorf("Failed to unmarshal control message: %v", err)
		return
	}

	err = validate.Struct(controlDeviceDto)

	if err != nil {
		log.Errorf("Validation error: %v", err)
		return
	}

	if controlDeviceDto.MacServer != os.Getenv("MAC_ADDRESS") {
		log.Errorf("Invalid mac server: %v", controlDeviceDto.MacServer)
		return
	}

	h.controlDeviceService.ControlDeviceCloud(&controlDeviceDto)
}

func (h *ConsumerHandler) DeleteDeviceFromCloudHandler(message []byte, db *gorm.DB) {
	var deleteDeviceDtoFromCloud dto.ReqDeleteDeviceFromCloudDto

	err := json.Unmarshal(message, &deleteDeviceDtoFromCloud)

	if err != nil {
		log.Errorf("Failed to unmarshal delete device message: %v", err)
		return
	}

	if deleteDeviceDtoFromCloud.MacServer != os.Getenv("MAC_ADDRESS") {
		log.Errorf("Invalid mac server: %v", deleteDeviceDtoFromCloud.MacServer)
		return
	}

	h.deviceService.DeleteDeviceRMQ(deleteDeviceDtoFromCloud.Guid)
}

func (h *ConsumerHandler) ChangeStatusDevice(message []byte, db *gorm.DB) {
	messageString := string(message)

	guid := strings.Split(messageString, "#")[0]
	status := strings.Split(messageString, "#")[1]

	h.deviceService.UpdateStatusDevice(guid, status)
}

func (h *ConsumerHandler) MonitoringDataDevice(message []byte, db *gorm.DB) {
	messageString := string(message)
	guid := strings.Split(messageString, "#")[0]
	data := strings.Split(messageString, "#")[1]

	h.deviceService.UpdateStatusAsMonitoring(guid, data)
}

func (h *ConsumerHandler) TestingConsumeAktuator(message []byte, db *gorm.DB) {
	messageString := string(message)

	log.Info(messageString)
}
