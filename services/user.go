package services

import "iot_proj/models"

func (s *Service) UpdateUserKeyCard(userID int, keyCard string) error {
	return s.Repo.UpdateUserKeyCard(userID, keyCard)
}

func (s *Service) GetUsers(limit, offset int) ([]*models.User, error) {
	return s.Repo.GetUsers(limit, offset)
}

func (s *Service) UpdateUserAccessLevel(userID int, accessLevel int) error {
	return s.Repo.UpdateUserAccessLevel(userID, accessLevel)
}
