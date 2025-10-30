package temperatureapi

type TemperatureResponse struct {
	Temperature float64 `json:"temperature"`
	Location    string  `json:"location"`
	SensorID    string  `json:"sensor_id"`
}

type HealthResponse struct {
	Status string `json:"status"`
}
