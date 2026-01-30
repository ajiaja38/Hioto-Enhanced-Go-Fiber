package model

import (
	"time"
)

type Floor struct {
	ID        uint      `gorm:"autoIncrement;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Rooms     []Room    `gorm:"foreignKey:FloorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"rooms"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}
