package dto

import "time"

type CreateRoomDto struct {
	Name    string `json:"name" validate:"required"`
	FloorID uint   `json:"floor_id" validate:"required"`
}

type UpdateRoomDto struct {
	Name    string `json:"name" validate:"required"`
	FloorID uint   `json:"floor_id" validate:"required"`
}

type ResponseRoomSimpleDto struct {
	ID    uint                   `json:"id"`
	Name  string                 `json:"name"`
	Floor ResponseFloorSimpleDto `json:"floor"`
}

type ResponseRoomDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	FloorID   uint      `json:"floor_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
