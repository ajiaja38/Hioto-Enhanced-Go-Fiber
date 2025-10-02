package dto

import (
	"go/hioto/pkg/enum"
	"time"
)

type LogsReponseDto struct {
	ID          uint      `json:"id"`
	InputGuid   string    `json:"input_guid"`
	InputName   string    `json:"input_name"`
	InputValue  string    `json:"input_value"`
	OutputGuid  string    `json:"output_guid"`
	OutputValue string    `json:"output_value"`
	Time        time.Time `json:"time"`
	MacServer   string    `json:"mac_server"`
}

type LogAktuatorReponseDto struct {
	ID        uint      `json:"id"`
	Guid      string    `json:"guid"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	Time      time.Time `json:"time"`
	MacServer string    `json:"mac_server"`
}

type LogMonitoringDeviceHistoryResponseDto struct {
	ID         uint             `json:"id"`
	DeviceGuid string           `json:"device_guid"`
	DeviceName string           `json:"device_name"`
	DeviceType enum.EDeviceType `json:"device_type"`
	Value      string           `json:"value"`
	Time       time.Time        `json:"time"`
	MacServer  string           `json:"mac_server"`
}
