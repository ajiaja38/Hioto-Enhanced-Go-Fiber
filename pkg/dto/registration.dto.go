package dto

import (
	"go/hioto/pkg/enum"
	"time"
)

type RegistrationDto struct {
	Guid     string           `json:"guid" validate:"required"`
	Mac      string           `json:"mac" validate:"required"`
	Type     enum.EDeviceType `json:"type" validate:"required"`
	Quantity int              `json:"quantity" validate:"required,min=1"`
	Name     string           `json:"name" validate:"required"`
	Version  string           `json:"version" validate:"required"`
	Minor    string           `json:"minor" validate:"required"`
}

type ReqCloudDeviceDto struct {
	RegistrationDto
	MacServer string `json:"mac_server" validate:"required"`
}

type ResponseDeviceDto struct {
	ID           uint             `json:"id"`
	Guid         string           `json:"guid"`
	Mac          string           `json:"mac"`
	Type         enum.EDeviceType `json:"type"`
	Quantity     int              `json:"quantity"`
	Name         string           `json:"name"`
	Version      string           `json:"version"`
	Minor        string           `json:"minor"`
	Status       string           `json:"status"`
	StatusDevice string           `json:"status_device"`
	LastSeen     time.Time        `json:"last_seen"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type ResCloudDeviceDto struct {
	ResponseDeviceDto
	MacServer string `json:"mac_server"`
}

type ReqUpdateDeviceDto struct {
	Guid     string           `json:"guid" validate:"required"`
	Mac      string           `json:"mac" validate:"required"`
	Type     enum.EDeviceType `json:"type" validate:"required"`
	Quantity int              `json:"quantity" validate:"required,min=1"`
	Name     string           `json:"name" validate:"required"`
	Version  string           `json:"version" validate:"required"`
	Minor    string           `json:"minor" validate:"required"`
}

type ReqUpdateDeviceDtoCloud struct {
	ReqUpdateDeviceDto
	MacServer string `json:"mac_server" validate:"required"`
}

type ReqDeleteDeviceFromCloudDto struct {
	Guid      string `json:"guid" validate:"required"`
	MacServer string `json:"mac_server" validate:"required"`
}

type ReqDeleteDeviceToCloudDto struct {
	Guid      string `json:"guid" validate:"required"`
	MacServer string `json:"mac_server" validate:"required"`
}
