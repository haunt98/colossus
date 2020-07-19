package storage

//go:generate easyjson -all models.go

type FileInfo struct {
	ID          string `json:"id"`
	ContentType string `json:"content_type"`
	Extension   string `json:"extension"`
	Size        int64  `json:"size"`
	Checksum    string `json:"checksum"`
}

type Response struct {
	ReturnCode    int         `json:"return_code"`
	ReturnMessage string      `json:"return_message,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}
