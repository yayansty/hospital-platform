package services

import "medic-api/repositories"

type HealthService struct {
	Repository *repositories.HealthRepository
}

func NewHealthService(repository *repositories.HealthRepository) *HealthService {
	return &HealthService{
		Repository: repository,
	}
}

func (s *HealthService) CheckDatabase() error {
	return s.Repository.CheckDatabase()
}
