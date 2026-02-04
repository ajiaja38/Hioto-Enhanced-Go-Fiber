package dto

import "go/hioto/pkg/enum"

type ControlLocalDto struct {
	Type    enum.EDeviceType `json:"type" validate:"required"`
	Message string           `json:"message" validate:"required"`
}

type ControlGasDetector struct {
	Guid        string `json:"guid" validate:"required"`
	DeviceName  string `json:"deviceName" validate:"required"`
	Value       int    `json:"value" validate:"required"`
	Condition   int    `json:"condition"`
	Description string `json:"description" validate:"required"`
}
