package model

import "time"

type RuleDevice struct {
	ID           uint         `gorm:"autoIncrement" json:"id"`
	InputGuid    string       `gorm:"type:varchar(255);not null" json:"input_guid"`
	InputValue   string       `gorm:"type:varchar(8);not null" json:"input_value"`
	OutputGuid   string       `gorm:"type:varchar(255);not null" json:"output_guid"`
	OutputValue  string       `gorm:"type:varchar(8);not null" json:"output_value"`
	InputDevice  Registration `gorm:"foreignKey:InputGuid;references:Guid" json:"input_device"`
	OutputDevice Registration `gorm:"foreignKey:OutputGuid;references:Guid" json:"output_device"`
	CreatedAt    time.Time    `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"not null" json:"updated_at"`
}
