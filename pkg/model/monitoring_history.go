package model

import (
	"go/hioto/pkg/enum"
	"time"
)

type MonitoringHistory struct {
	ID         uint             `gorm:"autoIncrement" json:"id"`
	DeviceGuid string           `gorm:"type:varchar(255);not null" json:"device_guid"`
	DeviceName string           `gorm:"type:varchar(255);not null" json:"device_name"`
	DeviceType enum.EDeviceType `gorm:"type:varchar(255);not null" json:"device_type"`
	Value      string           `gorm:"type:varchar(8);not null" json:"value"`
	Device     Registration     `gorm:"foreignKey:DeviceGuid;references:Guid" json:"device"`
	Time       time.Time        `gorm:"not null" json:"time"`
}
