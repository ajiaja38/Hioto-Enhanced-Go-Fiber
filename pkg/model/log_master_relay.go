package model

import (
	"time"

	"gorm.io/gorm"
)

type LogMasterRelay struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Guid      string         `gorm:"type:varchar(50);index" json:"guid"`
	Mac       string         `gorm:"type:varchar(50)" json:"mac"`
	DeviceName string        `gorm:"type:varchar(100)" json:"deviceName"`
	Status    int            `json:"status"`
	Condition string         `gorm:"type:varchar(50)" json:"condition"`
	Voltage   float64        `gorm:"type:decimal(10,2)" json:"voltage"`
	Current   float64        `gorm:"type:decimal(10,2)" json:"current"`
	Power     float64        `gorm:"type:decimal(10,2)" json:"power"`
	Energy    float64        `gorm:"type:decimal(10,2)" json:"energy"`
	Frequency float64        `gorm:"type:decimal(10,2)" json:"frequency"`
	PF        float64        `gorm:"type:decimal(10,2)" json:"pf"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LogMasterRelay) TableName() string {
	return "log_master_relays"
}