package service

import (
	"encoding/json"
	"go/hioto/pkg/dto"
	messagebroker "go/hioto/pkg/handler/message_broker"
	"go/hioto/pkg/model"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type LogService struct {
	db *gorm.DB
}

func NewLogService(db *gorm.DB) *LogService {
	return &LogService{
		db: db,
	}
}

func (s *LogService) GetAllLogs() {
	log.Info("Getting all logs...")

	var logs []model.Log

	query := s.db.Raw(`
			SELECT
				id, 
				input_guid, 
				input_name, 
				input_value, 
				output_guid, 
				output_value, 
				time 
			FROM logs
	`).Scan(&logs)

	if query.Error != nil {
		log.Errorf("Failed to get logs: %v", query.Error)
		return
	}

	if len(logs) == 0 {
		log.Info("No logs is found")
		return
	}

	var payload []dto.LogsReponseDto

	for _, log := range logs {
		payload = append(payload, dto.LogsReponseDto{
			ID:          log.ID,
			InputName:   log.InputName,
			InputGuid:   log.InputGuid,
			InputValue:  log.InputValue,
			OutputGuid:  log.OutputGuid,
			OutputValue: log.OutputValue,
			Time:        log.Time,
			MacServer:   os.Getenv("MAC_ADDRESS"),
		})
	}

	body, err := json.Marshal(payload)

	if err != nil {
		log.Errorf("Failed to marshal logs: %v", err)
		return
	}

	messagebroker.PublishToRmq(os.Getenv("RMQ_HIOTO"), "Biznet Hioto", body, os.Getenv("LOGS_QUEUE"), "amq.direct")

	queryDeleteAll := s.db.Exec("DELETE FROM logs")

	if queryDeleteAll.Error != nil {
		log.Errorf("Failed to delete logs: %v", queryDeleteAll.Error)
	}
}

func (s *LogService) GetAllLogAktuators() {
	log.Info("Getting all log aktuators...")

	var logs []model.LogAktuator

	query := s.db.Raw(`
			SELECT * FROM log_aktuators
	`).Scan(&logs)

	if query.Error != nil {
		log.Errorf("Failed to get logs: %v", query.Error)
		return
	}

	if len(logs) == 0 {
		log.Info("No logs found")
		return
	}

	var payload []dto.LogAktuatorReponseDto

	for _, log := range logs {
		payload = append(payload, dto.LogAktuatorReponseDto{
			ID:        log.ID,
			Guid:      log.InputGuid,
			Name:      log.Name,
			Value:     log.Value,
			Time:      log.Time,
			MacServer: os.Getenv("MAC_ADDRESS"),
		})
	}

	body, err := json.Marshal(payload)

	if err != nil {
		log.Errorf("Failed to json marshal logs: %v", err)
		return
	}

	messagebroker.PublishToRmq(os.Getenv("RMQ_HIOTO"), "Biznet Hioto", body, os.Getenv("LOGS_AKTUATOR_QUEUE"), "amq.direct")

	queryDeleteAll := s.db.Exec("DELETE FROM log_aktuators")

	if queryDeleteAll.Error != nil {
		log.Errorf("Failed to delete logs: %v", queryDeleteAll.Error)
		return
	}
}

func (s *LogService) GetAllMonitoringHistory() {
	log.Info("Getting all monitoring history...")

	var logs []model.MonitoringHistory

	query := s.db.Raw(`
		SELECT * FROM monitoring_histories
	`).Scan(&logs)

	if query.Error != nil {
		log.Errorf("Failed to get monitoring history: %v", query.Error)
		return
	}

	if len(logs) == 0 {
		log.Info("No monitoring history found")
		return
	}

	var payload []dto.LogMonitoringDeviceHistoryResponseDto

	for _, log := range logs {
		payload = append(payload, dto.LogMonitoringDeviceHistoryResponseDto{
			ID:         log.ID,
			DeviceGuid: log.DeviceGuid,
			DeviceName: log.DeviceName,
			DeviceType: log.DeviceType,
			Value:      log.Value,
			Time:       log.Time,
			MacServer:  os.Getenv("MAC_ADDRESS"),
		})
	}

	body, err := json.Marshal(payload)

	if err != nil {
		log.Errorf("Failed to json marshal monitoring history: %v", err)
		return
	}

	messagebroker.PublishToRmq(os.Getenv("RMQ_HIOTO"), "Biznet Hioto", body, os.Getenv("MONITORING_RESPONSE_QUEUE"), "amq.direct")

	queryDeleteAll := s.db.Exec("DELETE FROM monitoring_histories")

	if queryDeleteAll.Error != nil {
		log.Errorf("Failed to delete monitoring history: %v", queryDeleteAll.Error)
		return
	}
}
