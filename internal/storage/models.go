package storage

type FileInfo struct {
	ID          string `json:"id"`
	ContentType string `json:"content_type"`
	Extension   string `json:"extension"`
}

type Response struct {
	ReturnCode    int         `json:"return_code"`
	ReturnMessage string      `json:"return_message,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}
