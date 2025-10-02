package enum

type EDeviceType string

const (
	AI                 EDeviceType = "AI"
	SENSOR             EDeviceType = "SENSOR"
	SENSOR_TEMPERATURE EDeviceType = "SENSOR_TEMPERATURE"
	SENSOR_WATER_LEVEL EDeviceType = "SENSOR_WATER_LEVEL"
	SENSOR_CAMERA      EDeviceType = "SENSOR_CAMERA"
	SENSOR_PARKING     EDeviceType = "SENSOR_PARKING"
	AKTUATOR           EDeviceType = "AKTUATOR"
)
