package dto

import "time"

type CreateFloorDto struct {
	Name string `json:"name" validate:"required"`
}

type UpdateFloorDto struct {
	Name string `json:"name" validate:"required"`
}

type ResponseAllFloorDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseFloorSimpleDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ResponseFloorDto struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Rooms     []ResponseRoomDto `json:"rooms"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
