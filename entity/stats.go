package entity

type (
	StatsContextTotal struct {
		Total   int    `json:"total"`
		Context string `json:"context"`
	}
	// For Response Only
	ResponseGetMostContext struct {
		Message string              `json:"message" example:"asset fetched"`
		Status  string              `json:"status" example:"success"`
		Data    []StatsContextTotal `json:"data"`
	}
)
