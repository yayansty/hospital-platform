package repositories

import (
	"database/sql"

	"medic-api/models"
)

type PatientRepository struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	return &PatientRepository{
		db: db,
	}
}

func (r *PatientRepository) FindAll() ([]models.Patient, error) {
	query := `
		SELECT
			id,
			medical_record_number,
			name,
			birth_date,
			gender,
			created_at,
			updated_at
		FROM patients
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []models.Patient

	for rows.Next() {
		var patient models.Patient

		err := rows.Scan(
			&patient.ID,
			&patient.MedicalRecordNumber,
			&patient.Name,
			&patient.BirthDate,
			&patient.Gender,
			&patient.CreatedAt,
			&patient.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		patients = append(patients, patient)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}
func (r *PatientRepository) FindByID(id int64) (*models.Patient, error) {
	query := `
		SELECT
			id,
			medical_record_number,
			name,
			birth_date,
			gender,
			created_at,
			updated_at
		FROM patients
		WHERE id = @p1
	`

	var patient models.Patient

	err := r.db.QueryRow(
		query,
		id,
	).Scan(
		&patient.ID,
		&patient.MedicalRecordNumber,
		&patient.Name,
		&patient.BirthDate,
		&patient.Gender,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *PatientRepository) Create(patient *models.Patient) error {
	query := `
		INSERT INTO patients (
			medical_record_number,
			name,
			birth_date,
			gender
		)
		OUTPUT
			INSERTED.id,
			INSERTED.created_at,
			INSERTED.updated_at
		VALUES (
			@p1,
			@p2,
			@p3,
			@p4
		)
	`

	return r.db.QueryRow(
		query,
		patient.MedicalRecordNumber,
		patient.Name,
		patient.BirthDate,
		patient.Gender,
	).Scan(
		&patient.ID,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)
}
func (r *PatientRepository) Update(patient *models.Patient) error {
	query := `
		UPDATE patients
		SET
			medical_record_number = @p1,
			name = @p2,
			birth_date = @p3,
			gender = @p4,
			updated_at = GETDATE()
		WHERE id = @p5
	`

	result, err := r.db.Exec(
		query,
		patient.MedicalRecordNumber,
		patient.Name,
		patient.BirthDate,
		patient.Gender,
		patient.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PatientRepository) Delete(id int64) error {
	query := `
		DELETE FROM patients
		WHERE id = @p1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
