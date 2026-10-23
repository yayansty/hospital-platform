package models

import "time"

type Patient struct {
	ID                  int64      `json:"id"`
	MedicalRecordNumber string     `json:"medical_record_number"`
	Name                string     `json:"name"`
	BirthDate           *time.Time `json:"birth_date,omitempty"`
	Gender              *string    `json:"gender,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type PatientPatch struct {
	Name      *string    `json:"name"`
	BirthDate *time.Time `json:"birth_date"`
	Gender    *string    `json:"gender"`
}

//	type PatientFilter struct {
//		Page                int
//		Limit               int
//		Search              string
//		MedicalRecordNumber string
//		Gender              string
//	}
type PatientFilter struct {
	Page                int
	Limit               int
	Search              string
	MedicalRecordNumber string
	Gender              string

	SortBy    string
	SortOrder string

	DateFrom string
	DateTo   string
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
