package status

const (
	SuccessfulCode = 1
	ProcessingCode = 5
	FailedCode     = -1
)

type Status struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
