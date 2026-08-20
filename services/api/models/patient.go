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
