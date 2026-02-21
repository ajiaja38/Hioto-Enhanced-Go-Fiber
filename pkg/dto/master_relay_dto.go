package dto

// Payload JSON dari MQTT Master Relay
type MasterRelayMQTTPayload struct {
	Guid       string `json:"guid"`
	Mac        string `json:"mac"`
	DeviceName string `json:"deviceName"`
	Status     int    `json:"status"`
	Condition  string `json:"condition"`
	Value      struct {
		Voltage   float64 `json:"voltage"`
		Current   float64 `json:"current"`
		Power     float64 `json:"power"`
		Energy    float64 `json:"energy"`
		Frequency float64 `json:"frequency"`
		PF        float64 `json:"pf"`
	} `json:"value"`
	Unit struct {
		Voltage   string `json:"voltage"`
		Current   string `json:"current"`
		Power     string `json:"power"`
		Energy    string `json:"energy"`
		Frequency string `json:"frequency"`
		PF        string `json:"pf"`
	} `json:"unit"`
}

// Response API untuk Chart Statistik di UI
type MasterRelayStatsResponse struct {
	Labels []string  `json:"labels"`
	Data   []float64 `json:"data"`
	Unit   string    `json:"unit"`
}