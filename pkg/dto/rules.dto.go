package dto

type CreateRulePatternDto struct {
	InputGuid  string `json:"input_guid" validate:"required"`
	OutputGuid string `json:"output_guid" validate:"required"`
	Pattern    int    `json:"pattern" validate:"required,min=1,max=8"`
}

type CreateRuleDto struct {
	InputGuid  string   `json:"input_guid" validate:"required"`
	OutputGuid []string `json:"output_guid" validate:"required,min=1,max=8"`
}

type ResponseRuleDto struct {
	MacServer   string `json:"mac_server"`
	InputGuid   string `json:"input_guid"`
	InputValue  string `json:"input_value"`
	OutputGuid  string `json:"output_guid"`
	OutputValue string `json:"output_value"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ResponseGetRulesDto struct {
	ID                  uint   `json:"id"`
	GuidSensor          string `json:"guid_sensor"`
	SensorName          string `json:"sensor_name"`
	SensorInputValue    string `json:"sensor_input_value"`
	GuidAktuator        string `json:"guid_aktuator"`
	AktuatorName        string `json:"aktuator_name"`
	AktuatorOutputValue string `json:"aktuator_output_value"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}
