package services

import (
	"errors"

	"iot_proj/models"
)

func (s *Service) GetEntryLogsForUser(userID, limit, offset int) ([]*models.EntryLog, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if user.KeyCard == nil {
		return nil, errors.New("user does not have a key card")
	}

	return s.Repo.GetEntryLogsPaginated(*user.KeyCard, limit, offset)
}
