package model

import "time"

type Log struct {
	ID           uint         `gorm:"autoIncrement" json:"id"`
	InputGuid    string       `gorm:"type:varchar(255);not null" json:"inputguid"`
	InputName    string       `gorm:"type:varchar(255);not null" json:"inputname"`
	InputValue   string       `gorm:"type:varchar(8);not null" json:"inputvalue"`
	OutputGuid   string       `gorm:"type:varchar(255);not null" json:"outputguid"`
	OutputValue  string       `gorm:"type:varchar(8);not null" json:"outputvalue"`
	Time         time.Time    `gorm:"not null" json:"time"`
	InputDevice  Registration `gorm:"foreignKey:InputGuid;references:Guid" json:"input_device"`
	OutputDevice Registration `gorm:"foreignKey:OutputGuid;references:Guid" json:"output_device"`
}
