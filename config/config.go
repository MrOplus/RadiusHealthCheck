package config

import "github.com/kooroshh/RadiusHealthCheck/models"

type Config struct {
	Servers     []models.Server         `json:"Servers"`
	Credentials models.BasicCredentials `json:"Credentials"`
	Interval    int                     `json:"Interval"`
	Hook        struct {
		Enabled     bool                    `json:"Enabled"`
		Url         string                  `json:"Url"`
		Credentials models.BasicCredentials `json:"Credentials"`
	} `json:"Hook"`
	ContainerControl struct {
		Enabled       bool   `json:"Enabled"`
		ContainerName string `json:"ContainerName"`
	} `json:"ContainerControl"`
}
