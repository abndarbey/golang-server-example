package faulterr

// FaultErr structure
type FaultErr struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
}
