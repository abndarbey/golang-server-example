package models

type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}
