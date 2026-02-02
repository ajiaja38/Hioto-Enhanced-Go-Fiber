package enum

type EDeviceType string

const (
	AI                  EDeviceType = "AI"
	SENSOR              EDeviceType = "SENSOR"
	SENSOR_TEMPERATURE  EDeviceType = "SENSOR_TEMPERATURE"
	SENSOR_WATER_TANK   EDeviceType = "SENSOR_WATER_TANK"
	SENSOR_CAMERA       EDeviceType = "SENSOR_CAMERA"
	SENSOR_PARKING      EDeviceType = "SENSOR_PARKING"
	SENSOR_GAS_DETECTOR EDeviceType = "SENSOR_GAS_DETECTOR"
	AKTUATOR            EDeviceType = "AKTUATOR"
)
