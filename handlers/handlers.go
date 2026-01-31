package handlers

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"iot_proj/services"
)

type MQTTHandler struct {
	Service *services.Service
}

func NewMQTTHandler(service *services.Service) *MQTTHandler {
	return &MQTTHandler{Service: service}
}

func (h *MQTTHandler) SubscribeToCardScans(client mqtt.Client) {
	topic := "card/scan"
	token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))
		h.Service.HandleCardScan(string(msg.Payload()))
	})
	token.Wait()
	log.Printf("Subscribed to topic %s", topic)
}
