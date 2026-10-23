package services

import (
	"strings"
	"time"

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

// func (s *PatientService) GetAllPatients() ([]models.Patient, error) {
// 	return s.repository.FindAll()
// }

func (s *PatientService) GetAllPatients(
	filter models.PatientFilter,
) ([]models.Patient, models.PaginationMeta, error) {

	validationErrors := helpers.ValidationErrors{}

	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	if filter.Page < 1 {
		validationErrors["page"] = "must be at least 1"
	}

	if filter.Limit < 1 || filter.Limit > 100 {
		validationErrors["limit"] = "must be between 1 and 100"
	}

	filter.Search = strings.TrimSpace(filter.Search)
	filter.MedicalRecordNumber = strings.TrimSpace(
		filter.MedicalRecordNumber,
	)
	filter.Gender = strings.ToUpper(strings.TrimSpace(filter.Gender))

	if filter.Gender != "" &&
		filter.Gender != "M" &&
		filter.Gender != "F" {
		validationErrors["gender"] = "must be M or F"
	}

	// Default pengurutan
	if filter.SortBy == "" {
		filter.SortBy = "id"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "DESC"
	}

	filter.SortBy = strings.ToLower(strings.TrimSpace(filter.SortBy))
	filter.SortOrder = strings.ToUpper(strings.TrimSpace(filter.SortOrder))

	allowedSortColumns := map[string]bool{
		"id":                    true,
		"name":                  true,
		"medical_record_number": true,
		"birth_date":            true,
		"created_at":            true,
		"updated_at":            true,
	}

	if !allowedSortColumns[filter.SortBy] {
		validationErrors["sort_by"] = "is not a supported column"
	}

	if filter.SortOrder != "ASC" && filter.SortOrder != "DESC" {
		validationErrors["sort_order"] = "must be ASC or DESC"
	}

	// Validasi format tanggal dan rentang tanggal
	var dateFrom, dateTo time.Time

	if filter.DateFrom != "" {
		parsed, err := time.Parse("2006-01-02", filter.DateFrom)
		if err != nil {
			validationErrors["date_from"] = "must use YYYY-MM-DD format"
		} else {
			dateFrom = parsed
		}
	}

	if filter.DateTo != "" {
		parsed, err := time.Parse("2006-01-02", filter.DateTo)
		if err != nil {
			validationErrors["date_to"] = "must use YYYY-MM-DD format"
		} else {
			dateTo = parsed
		}
	}

	if !dateFrom.IsZero() &&
		!dateTo.IsZero() &&
		dateFrom.After(dateTo) {
		validationErrors["date_from"] = "must not be after date_to"
	}

	if len(validationErrors) > 0 {
		return nil, models.PaginationMeta{},
			&helpers.ValidationError{Errors: validationErrors}
	}

	return s.repository.FindAll(filter)
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

	// Medical Record Number
	if !helpers.Required(patient.MedicalRecordNumber) {
		validationErrors["medical_record_number"] = "is required"
	} else if !helpers.MaxLength(patient.MedicalRecordNumber, 50) {
		validationErrors["medical_record_number"] = "maximum 50 characters"
	}

	// Name
	if !helpers.Required(patient.Name) {
		validationErrors["name"] = "is required"
	} else if !helpers.MaxLength(patient.Name, 255) {
		validationErrors["name"] = "maximum 255 characters"
	}

	// Gender
	if patient.Gender != nil {
		gender := strings.ToUpper(strings.TrimSpace(*patient.Gender))

		if gender != "M" && gender != "F" {
			validationErrors["gender"] = "must be M or F"
		} else {
			*patient.Gender = gender
		}
	}

	// Birth Date
	if patient.BirthDate != nil {
		if patient.BirthDate.After(time.Now()) {
			validationErrors["birth_date"] = "cannot be in the future"
		}
	}

	if len(validationErrors) > 0 {
		return &helpers.ValidationError{
			Errors: validationErrors,
		}
	}

	return nil
}

func (s *PatientService) PatchPatient(
	id int64,
	patch *models.PatientPatch,
) error {
	if err := validatePatientPatch(patch); err != nil {
		return err
	}

	return s.repository.Patch(id, patch)
}

func validatePatientPatch(patch *models.PatientPatch) error {
	validationErrors := helpers.ValidationErrors{}

	if patch.Name != nil {
		name := strings.TrimSpace(*patch.Name)

		if !helpers.Required(name) {
			validationErrors["name"] = "is required"
		} else if !helpers.MaxLength(name, 255) {
			validationErrors["name"] = "maximum 255 characters"
		} else {
			*patch.Name = name
		}
	}

	if patch.Gender != nil {
		gender := strings.ToUpper(strings.TrimSpace(*patch.Gender))

		if gender != "M" && gender != "F" {
			validationErrors["gender"] = "must be M or F"
		} else {
			*patch.Gender = gender
		}
	}

	if patch.BirthDate != nil {
		if patch.BirthDate.After(time.Now()) {
			validationErrors["birth_date"] = "cannot be in the future"
		}
	}

	if len(validationErrors) > 0 {
		return &helpers.ValidationError{
			Errors: validationErrors,
		}
	}

	return nil
}
