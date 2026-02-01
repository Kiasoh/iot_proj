package services

import (

	"iot_proj/models"
)

func (s *Service) GetEntryLogsForUser(userID, limit, offset int) ([]*models.EntryLog, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.Repo.GetEntryLogsPaginated(*user.KeyCard, limit, offset)
}
