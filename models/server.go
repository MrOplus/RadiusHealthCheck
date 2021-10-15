package models

type Server struct {
	Address      string `json:"Address"`
	Port         int    `json:"Port"`
	Secret       string `json:"Secret"`
	TriggerCount int    `json:"TriggerCount"`
}
