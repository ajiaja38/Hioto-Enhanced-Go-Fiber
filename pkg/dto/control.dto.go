package dto

import "go/hioto/pkg/enum"

type ControlLocalDto struct {
	Type    enum.EDeviceType `json:"type" validate:"required"`
	Message string           `json:"message" validate:"required"`
}

type ControlDto struct {
	ControlLocalDto
	MacServer string `json:"mac_server" validate:"required"`
}
