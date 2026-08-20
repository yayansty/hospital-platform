package services

import (
	"medic-api/models"
	"medic-api/repositories"
)

type PatientService struct {
	repository *repositories.PatientRepository
}

func NewPatientService(repository *repositories.PatientRepository) *PatientService {
	return &PatientService{
		repository: repository,
	}
}

func (s *PatientService) GetAllPatients() ([]models.Patient, error) {
	return s.repository.FindAll()
}

func (s *PatientService) GetPatientByID(id int64) (*models.Patient, error) {
	return s.repository.FindByID(id)
}

func (s *PatientService) CreatePatient(patient *models.Patient) error {
	return s.repository.Create(patient)
}
func (s *PatientService) UpdatePatient(patient *models.Patient) error {
	return s.repository.Update(patient)
}
func (s *PatientService) DeletePatient(id int64) error {
	return s.repository.Delete(id)
}
