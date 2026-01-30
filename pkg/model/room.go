package model

import "time"

type Room struct {
	ID        uint      `gorm:"autoIncrement;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	FloorID   uint      `gorm:"not null" json:"floor_id"`
	Floor     Floor     `gorm:"foreignKey:FloorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"floor"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}
