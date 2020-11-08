package model

type HealthResponse struct {
	DBConnection string `json:"db_connection"`
	WorkingTime  string `json:"working_time"`
}
