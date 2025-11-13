package service

import (
	"encoding/json"
	"fmt"
	"go/hioto/config"
	"go/hioto/pkg/dto"
	"go/hioto/pkg/enum"
	messagebroker "go/hioto/pkg/handler/message_broker"
	"go/hioto/pkg/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var locations *time.Location

func init() {
	locations = time.FixedZone("WIB", 7*60*60)
}

type RuleService struct {
	db *gorm.DB
}

func NewRuleService(db *gorm.DB) *RuleService {
	return &RuleService{
		db: db,
	}
}

func generateSensorPatterns(length int) []string {
	totalPatterns := 1 << length
	patterns := make([]string, totalPatterns)

	for i := range totalPatterns {
		binaryPattern := fmt.Sprintf("%0*b", length, i)
		patterns[i] = binaryPattern
	}

	return patterns
}

func (s *RuleService) CreateRules(createRuleDto *dto.CreateRuleDto) (responseRules []dto.ResponseRuleDto, err error) {
	if err = s.db.Where("guid = ?", createRuleDto.InputGuid).First(&model.Registration{}).Error; err != nil {
		log.Errorf("Sensor is not found: %v 💥", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "The Sensor is not found")
	}

	for _, actuator := range createRuleDto.OutputGuid {
		if err = s.db.Where("guid = ?", actuator).First(&model.Registration{}).Error; err != nil {
			log.Errorf("The actuator not found: %v 💥", err)
			return nil, fiber.NewError(fiber.StatusNotFound, "The actuator is not found")
		}
	}

	length := len(createRuleDto.OutputGuid)
	sensorPatterns := generateSensorPatterns(length)

	for _, sensor := range sensorPatterns {
		for i, actuator := range createRuleDto.OutputGuid {
			outputValue := '1'

			if sensor[i] == '1' {
				outputValue = '0'
			}

			rule := &model.RuleDevice{
				InputGuid:   createRuleDto.InputGuid,
				InputValue:  sensor,
				OutputGuid:  actuator,
				OutputValue: string(outputValue),
				CreatedAt:   time.Now().In(locations),
				UpdatedAt:   time.Now().In(locations),
			}

			if err = s.db.Create(rule).Error; err != nil {
				log.Errorf("Error creating rule: %v 💥", err)
				return nil, fiber.NewError(fiber.StatusBadRequest, "Error creating rule")
			}

			responseRules = append(responseRules, dto.ResponseRuleDto{
				MacServer:   config.MAC_ADDRESS.GetValue(),
				InputGuid:   rule.InputGuid,
				InputValue:  rule.InputValue,
				OutputGuid:  rule.OutputGuid,
				OutputValue: rule.OutputValue,
				CreatedAt:   rule.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   rule.UpdatedAt.Format(time.RFC3339),
			})
		}
	}

	responseToJson, err := json.Marshal(responseRules)

	if err != nil {
		log.Errorf("Failed to marshal response: %v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Failed to marshal response")
	}

	messagebroker.PublishToRmq(
		config.RMQ_CLOUD_INSTANCE.GetValue(),
		responseToJson,
		config.RULES_RESPONSE_QUEUE.GetValue(),
		config.EXCHANGE_DIRECT.GetValue(),
	)

	log.Info("Rule was created successfully ✅")

	return responseRules, nil
}

func (s *RuleService) GetRulesByGuid(guid string) ([]dto.ResponseRuleDto, error) {
	var rules []model.RuleDevice

	if err := s.db.Where("input_guid = ? OR output_guid = ?", guid, guid).Find(&rules).Error; err != nil {
		log.Errorf("Failed to get rules: %v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Failed to get rules")
	}

	if len(rules) == 0 {
		log.Error("Failed to get rules, rules not found")
		return nil, fiber.NewError(fiber.StatusNotFound, "Failed to get rules, rules not found")
	}

	var responseRules []dto.ResponseRuleDto
	for _, rule := range rules {
		responseRules = append(responseRules, dto.ResponseRuleDto{
			MacServer:   config.MAC_ADDRESS.GetValue(),
			InputGuid:   rule.InputGuid,
			InputValue:  rule.InputValue,
			OutputGuid:  rule.OutputGuid,
			OutputValue: rule.OutputValue,
			CreatedAt:   rule.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   rule.UpdatedAt.Format(time.RFC3339),
		})
	}

	return responseRules, nil
}

func (s *RuleService) DeleteRulesByGuidSensor(guid string) error {
	var device model.Registration

	if err := s.db.Where("guid = ?", guid).First(&model.Registration{}).Scan(&device).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Sensor not found")
	}

	if device.Type != enum.SENSOR {
		return fiber.NewError(fiber.StatusBadRequest, "Device is not a sensor")
	}

	if err := s.db.Where("input_guid = ?", guid).Delete(&model.RuleDevice{}).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to delete rules")
	}

	return nil
}
