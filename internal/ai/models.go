package ai

import (
	"github.com/haunt98/colossus/pkg/status"
)

type ProcessInfo struct {
	TransID    string        `json:"trans_id"`
	StatusInfo status.Status `json:"status_info"`
	InputID    string        `json:"input_id"`
	OutputID   string        `json:"output_id"`
}
