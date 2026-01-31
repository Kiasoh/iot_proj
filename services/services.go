package services

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"iot_proj/repository"
)

type Service struct {
	Repo *repository.Repository
	Mqtt mqtt.Client
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) SetMqttClient(mqtt mqtt.Client) {
	s.Mqtt = mqtt
}
