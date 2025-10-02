package model

import (
	"go/hioto/pkg/enum"
	"time"
)

type Registration struct {
	ID                uint                `gorm:"autoIncrement" json:"id"`
	Guid              string              `gorm:"type:varchar(255);not null;unique" json:"guid"`
	Mac               string              `gorm:"type:varchar(255);not null" json:"mac"`
	Type              enum.EDeviceType    `gorm:"type:varchar(255);not null" json:"type"`
	Quantity          int                 `gorm:"type:int;not null" json:"quantity"`
	Name              string              `gorm:"type:varchar(255);not null" json:"name"`
	Version           string              `gorm:"type:varchar(255);not null" json:"version"`
	Minor             string              `gorm:"type:varchar(255);not null" json:"minor"`
	Status            string              `gorm:"type:varchar(255);" json:"status"`
	StatusDevice      enum.EDeviceStatus  `gorm:"type:varchar(255);default:'1" json:"status_device"`
	LastSeen          time.Time           `gorm:"" json:"last_seen"`
	CreatedAt         time.Time           `gorm:"not null" json:"created_at"`
	UpdatedAt         time.Time           `gorm:"not null" json:"updated_at"`
	RulesInput        []RuleDevice        `gorm:"foreignKey:InputGuid;references:Guid" json:"rules_input"`
	RulesOutput       []RuleDevice        `gorm:"foreignKey:OutputGuid;references:Guid" json:"rules_output"`
	LogsInput         []Log               `gorm:"foreignKey:InputGuid;references:Guid" json:"logs_input"`
	LogsOutput        []Log               `gorm:"foreignKey:OutputGuid;references:Guid" json:"logs_output"`
	LogsInputAktuator []LogAktuator       `gorm:"foreignKey:InputGuid;references:Guid" json:"logs_aktuator_input"`
	MonitoringHistory []MonitoringHistory `gorm:"foreignKey:DeviceGuid;references:Guid" json:"monitoring_history"`
}
