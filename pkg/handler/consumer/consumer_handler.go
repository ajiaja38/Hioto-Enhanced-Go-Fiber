package consumer

import (
	"encoding/json"
	"go/hioto/pkg/dto"
	"go/hioto/pkg/service"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
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

func (h *ConsumerHandler) RegistrationHandler(message []byte) {
	var registrationDto dto.RegistrationDto

	if err := json.Unmarshal(message, &registrationDto); err != nil {
		log.Errorf("Failed to unmarshal registration message: %v", err)
		return
	}

	h.deviceService.RegisterDeviceLocal(&registrationDto)
}

func (h *ConsumerHandler) RegistrationFromCloudHandler(message []byte) {
	var registrationDto dto.RegistrationDto

	if err := json.Unmarshal(message, &registrationDto); err != nil {
		log.Errorf("Failed to unmarshal registration message: %v", err)
		return
	}

	h.deviceService.RegisterDeviceCloud(&registrationDto)
}

func (h *ConsumerHandler) UpdateDeviceFromCloudHandler(message []byte) {
	var updateDeviceFromCloudDto dto.ReqUpdateDeviceDto

	if err := json.Unmarshal(message, &updateDeviceFromCloudDto); err != nil {
		log.Errorf("Failed to unmarshal update device message: %v", err)
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

func (h *ConsumerHandler) RulesHandler(message []byte) {
	var createRuleDto dto.CreateRuleDto

	if err := json.Unmarshal(message, &createRuleDto); err != nil {
		log.Errorf("Failed to unmarshal rule message: %v", err)
		return
	}

	h.ruleService.CreateRules(&createRuleDto)
}

func (h *ConsumerHandler) ControlHandler(message []byte) {
	var controlDeviceDto dto.ControlLocalDto

	if err := json.Unmarshal(message, &controlDeviceDto); err != nil {
		log.Errorf("Failed to unmarshal control message: %v", err)
		return
	}

	if err := validate.Struct(controlDeviceDto); err != nil {
		log.Errorf("Validation error: %v", err)
		return
	}

	h.controlDeviceService.ControlDeviceCloud(&controlDeviceDto)
}

func (h *ConsumerHandler) ControlSensorHandler(message []byte) {
	messageString := string(message)

	guid := strings.Split(messageString, "#")[0]
	value := strings.Split(messageString, "#")[1]

	h.controlDeviceService.ControlSensor(guid, value)
}

func (h *ConsumerHandler) DeleteDeviceFromCloudHandler(message []byte) {
	var deleteDeviceDtoFromCloud dto.ReqDeleteDeviceFromCloudDto

	if err := json.Unmarshal(message, &deleteDeviceDtoFromCloud); err != nil {
		log.Errorf("Failed to unmarshal delete device message: %v", err)
		return
	}

	h.deviceService.DeleteDeviceRMQ(deleteDeviceDtoFromCloud.Guid)
}

func (h *ConsumerHandler) MonitoringDataDevice(message []byte) {
	messageString := string(message)
	guid := strings.Split(messageString, "#")[0]
	data := strings.Split(messageString, "#")[1]

	h.deviceService.UpdateStatusAsMonitoring(guid, data)
}

func (h *ConsumerHandler) TestingConsumeAktuator(message []byte) {
	messageString := string(message)

	log.Info(messageString)
}
