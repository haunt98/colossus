package gateway

import "github.com/haunt98/colossus/pkg/status"

type ProcessInfo struct {
	TransID    string        `json:"trans_id"`
	StatusInfo status.Status `json:"status_info"`
	EventType  string        `json:"event_type"`
	AITransID  string        `json:"ai_trans_id"`
	AIOutputID string        `json:"ai_output_id"`
}
