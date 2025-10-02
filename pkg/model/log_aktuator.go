package model

import "time"

type LogAktuator struct {
	ID          uint         `gorm:"autoIncrement" json:"id"`
	InputGuid   string       `gorm:"type:varchar(255);not null" json:"guid"`
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	Value       string       `gorm:"type:varchar(8);not null" json:"inputvalue"`
	Time        time.Time    `gorm:"not null" json:"time"`
	InputDevice Registration `gorm:"foreignKey:InputGuid;references:Guid" json:"input_device"`
}
