package services

import (
	"log"

	"iot_proj/models"
)

func (s *Service) HandleCardScan(cardUUIDStr string) {
	user, err := s.Repo.GetUser(cardUUIDStr)
	if err != nil {
		log.Printf("User not found for card %s: %v", cardUUIDStr, err)
		entryLog := &models.EntryLog{
			KeyCard: &cardUUIDStr,
			Status:  "denied",
			Message: models.StringPtr("User not found"),
		}
		if err := s.Repo.CreateEntryLog(entryLog); err != nil {
			log.Printf("Error creating entry log: %v", err)
		}
		s.PublishLockAction("denied")
		return
	}

	if user.AccessLevel > 0 {
		log.Printf("Access granted for user %d with card %s", user.ID, cardUUIDStr)
		entryLog := &models.EntryLog{
			KeyCard: &cardUUIDStr,
			Status:  "granted",
		}
		if err := s.Repo.CreateEntryLog(entryLog); err != nil {
			log.Printf("Error creating entry log: %v", err)
		}
		if err := s.Repo.UpdateUserLastAccessed(cardUUIDStr); err != nil {
			log.Printf("Error updating last accessed time: %v", err)
		}
		s.PublishLockAction("granted")
	} else {
		log.Printf("Access denied for user %d with card %s", user.ID, cardUUIDStr)
		entryLog := &models.EntryLog{
			KeyCard: &cardUUIDStr,
			Status:  "denied",
			Message: models.StringPtr("Insufficient access level"),
		}
		if err := s.Repo.CreateEntryLog(entryLog); err != nil {
			log.Printf("Error creating entry log: %v", err)
		}
		s.PublishLockAction("denied")
	}
}

func (s *Service) PublishLockAction(action string) {
	topic := "lock/action"
	token := s.Mqtt.Publish(topic, 1, false, action)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Failed to publish to topic %s: %v", topic, token.Error())
	} else {
		log.Printf("Published message to topic %s: %s", topic, action)
	}
}

func (s *Service) PublishUnregisteredCard(cardUUID string) {
	topic := "card/unreg"
	token := s.Mqtt.Publish(topic, 1, false, cardUUID)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Failed to publish to topic %s: %v", topic, token.Error())
	} else {
		log.Printf("Published message to topic %s: %s", topic, cardUUID)
	}
}
