package services

import (
	"strings"

	"medic-api/helpers"
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
	if err := validatePatient(patient); err != nil {
		return err
	}
	return s.repository.Create(patient)
}
func (s *PatientService) UpdatePatient(patient *models.Patient) error {
	if err := validatePatient(patient); err != nil {
		return err
	}
	return s.repository.Update(patient)
}
func (s *PatientService) DeletePatient(id int64) error {
	return s.repository.Delete(id)
}
func validatePatient(patient *models.Patient) error {
	validationErrors := helpers.ValidationErrors{}

	patient.MedicalRecordNumber = strings.TrimSpace(patient.MedicalRecordNumber)
	patient.Name = strings.TrimSpace(patient.Name)

	if !helpers.Required(patient.MedicalRecordNumber) {
		validationErrors["medical_record_number"] = "is required"
	} else if !helpers.MaxLength(patient.MedicalRecordNumber, 50) {
		validationErrors["medical_record_number"] = "maximum 50 characters"
	}

	if !helpers.Required(patient.Name) {
		validationErrors["name"] = "is required"
	} else if !helpers.MaxLength(patient.Name, 255) {
		validationErrors["name"] = "maximum 255 characters"
	}

	if patient.Gender != nil {
		gender := strings.ToUpper(strings.TrimSpace(*patient.Gender))

		if gender != "M" && gender != "F" {
			validationErrors["gender"] = "must be M or F"
		} else {
			*patient.Gender = gender
		}
	}

	if len(validationErrors) > 0 {
		return &helpers.ValidationError{
			Errors: validationErrors,
		}
	}

	return nil
}
