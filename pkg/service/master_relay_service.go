package service

import (
	"go/hioto/pkg/dto"
	"go/hioto/pkg/model"
	"time"
	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2/log"
)

var lastSaveMap = make(map[string]time.Time)

type MasterRelayService struct{
	db *gorm.DB
}

func NewMasterRelayService(db *gorm.DB) *MasterRelayService {
	return &MasterRelayService{
		db: db,
	}
}

func (s *MasterRelayService) SaveLog(payload dto.MasterRelayMQTTPayload) error {
	now := time.Now()
	lastTime, exists := lastSaveMap[payload.Guid]

	if !exists || now.Sub(lastTime).Seconds() >= 43200 {
		logData := model.LogMasterRelay{
			Guid:       payload.Guid,
			Mac:        payload.Mac,
			DeviceName: payload.DeviceName,
			Status:     payload.Status,
			Condition:  payload.Condition,
			Voltage:    payload.Value.Voltage,
			Current:    payload.Value.Current,
			Power:      payload.Value.Power,
			Energy:     payload.Value.Energy,
			Frequency:  payload.Value.Frequency,
			PF:         payload.Value.PF,
			CreatedAt:  now,
		}

		err := s.db.Create(&logData).Error
		if err == nil {
			lastSaveMap[payload.Guid] = now
			log.Info("📊 [Saved] Master Relay %s: %v Watt", payload.Guid, payload.Value.Power)
		}
		return err
	}

	return nil
}

func (s *MasterRelayService) GetStatistics(guid string, period string) (dto.MasterRelayStatsResponse, error) {
	var resp dto.MasterRelayStatsResponse
	var results []struct {
		Label string
		Value float64
	}

	db := s.db.Model(&model.LogMasterRelay{}).Where("guid = ?", guid)

	switch period {
	case "daily":
		resp.Unit = "W"
		db.Select("strftime('%H:00', created_at, 'localtime') as label, ROUND(AVG(power), 2) as value")
		db.Where("date(created_at, 'localtime') = date('now', 'localtime')")
		db.Group("strftime('%H:00', created_at, 'localtime')")
		db.Order("label ASC")

	case "weekly":
		resp.Unit = "kWh"
		db.Select("strftime('%d-%m', created_at, 'localtime') as label, ROUND(MAX(energy) - MIN(energy), 2) as value")
		db.Where("created_at >= datetime('now', '-7 days', 'localtime')")
		db.Group("strftime('%d-%m', created_at, 'localtime')")
		db.Order("created_at ASC")

	case "monthly":
		resp.Unit = "kWh"
		db.Select("strftime('%d', created_at, 'localtime') as label, ROUND(MAX(energy) - MIN(energy), 2) as value")
		db.Where("strftime('%Y-%m', created_at, 'localtime') = strftime('%Y-%m', 'now', 'localtime')")
		db.Group("strftime('%d', created_at, 'localtime')")
		db.Order("label ASC")

	default:
		resp.Unit = "W"
		db.Select("strftime('%H:00', created_at, 'localtime') as label, ROUND(AVG(power), 2) as value")
		db.Where("date(created_at, 'localtime') = date('now', 'localtime')")
		db.Group("strftime('%H:00', created_at, 'localtime')")
		db.Order("label ASC")
	}

	if err := db.Scan(&results).Error; err != nil {
		return resp, err
	}

	resp.Labels = []string{}
	resp.Data = []float64{}

	for _, r := range results {
		resp.Labels = append(resp.Labels, r.Label)
		resp.Data = append(resp.Data, r.Value)
	}

	return resp, nil
}