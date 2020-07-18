package gateway

import "colossus/pkg/status"

//go:generate easyjson -all models.go

type ProcessInfo struct {
	TransID    string        `json:"trans_id"`
	StatusInfo status.Status `json:"status_info"`
	EventType  int           `json:"event_type"`
	AITransID  string        `json:"ai_trans_id"`
	AIOutputID string        `json:"ai_output_id"`
}
