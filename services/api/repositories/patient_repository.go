package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"medic-api/helpers"
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

// func (r *PatientRepository) FindAll() ([]models.Patient, error) {
// 	query := `
// 		SELECT
// 			id,
// 			medical_record_number,
// 			name,
// 			birth_date,
// 			gender,
// 			created_at,
// 			updated_at
// 		FROM patients
// 		ORDER BY id DESC
// 	`

// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var patients []models.Patient

// 	for rows.Next() {
// 		var patient models.Patient

// 		err := rows.Scan(
// 			&patient.ID,
// 			&patient.MedicalRecordNumber,
// 			&patient.Name,
// 			&patient.BirthDate,
// 			&patient.Gender,
// 			&patient.CreatedAt,
// 			&patient.UpdatedAt,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		patients = append(patients, patient)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return patients, nil
// }

func (r *PatientRepository) FindAll(
	filter models.PatientFilter,
) ([]models.Patient, models.PaginationMeta, error) {

	conditions := []string{}
	args := []interface{}{}

	addCondition := func(condition string, value interface{}) {
		args = append(args, value)
		conditions = append(
			conditions,
			fmt.Sprintf(condition, len(args)),
		)
	}

	fmt.Printf(
		"DEBUG FILTER: DateFrom=%q DateTo=%q Gender=%q\n",
		filter.DateFrom,
		filter.DateTo,
		filter.Gender,
	)

	// if filter.Search != "" {
	// 	addCondition(
	// 		"(name LIKE @p%d OR medical_record_number LIKE @p%d)",
	// 		"%"+filter.Search+"%",
	// 	)
	// 	// Placeholder untuk dua kolom harus sama.
	// 	last := len(args)
	// 	conditions[len(conditions)-1] = fmt.Sprintf(
	// 		"(name LIKE @p%d OR medical_record_number LIKE @p%d)",
	// 		last,
	// 		last,
	// 	)
	// }
	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		p := len(args)

		conditions = append(conditions, fmt.Sprintf(
			"(name LIKE @p%d OR medical_record_number LIKE @p%d)",
			p,
			p,
		))
	}

	if filter.MedicalRecordNumber != "" {
		addCondition(
			"medical_record_number = @p%d",
			filter.MedicalRecordNumber,
		)
	}

	if filter.Gender != "" {
		addCondition("gender = @p%d", filter.Gender)
	}

	if filter.DateFrom != "" {
		args = append(args, filter.DateFrom)
		conditions = append(
			conditions,
			fmt.Sprintf(
				"created_at >= CAST(@p%d AS date)",
				len(args),
			),
		)
	}

	if filter.DateTo != "" {
		args = append(args, filter.DateTo)
		conditions = append(
			conditions,
			fmt.Sprintf(
				"created_at < DATEADD(day, 1, CAST(@p%d AS date))",
				len(args),
			),
		)
	}

	whereSQL := ""
	if len(conditions) > 0 {
		whereSQL = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Hitung total data yang sesuai dengan filter.
	countQuery := "SELECT COUNT(*) FROM patients" + whereSQL

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, models.PaginationMeta{}, helpers.ErrDatabase
	}

	// Hitung offset halaman.
	offset := (filter.Page - 1) * filter.Limit

	dataArgs := append([]interface{}{}, args...)
	dataArgs = append(dataArgs, offset, filter.Limit)

	offsetParam := len(dataArgs) - 1
	limitParam := len(dataArgs)

	sortColumns := map[string]string{
		"id":                    "id",
		"name":                  "name",
		"medical_record_number": "medical_record_number",
		"birth_date":            "birth_date",
		"created_at":            "created_at",
		"updated_at":            "updated_at",
	}

	sortColumn, ok := sortColumns[filter.SortBy]
	if !ok {
		sortColumn = "id"
	}

	sortOrder := strings.ToUpper(filter.SortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	query := `
		SELECT
			id,
			medical_record_number,
			name,
			birth_date,
			gender,
			created_at,
			updated_at
		FROM patients` + whereSQL + fmt.Sprintf(`
		ORDER BY %s %s
		OFFSET @p%d ROWS
		FETCH NEXT @p%d ROWS ONLY`,
		sortColumn,
		sortOrder,
		offsetParam,
		limitParam,
	)

	rows, err := r.db.Query(query, dataArgs...)
	if err != nil {
		return nil, models.PaginationMeta{}, helpers.ErrDatabase
	}
	defer rows.Close()

	patients := []models.Patient{}

	for rows.Next() {
		var patient models.Patient

		if err := rows.Scan(
			&patient.ID,
			&patient.MedicalRecordNumber,
			&patient.Name,
			&patient.BirthDate,
			&patient.Gender,
			&patient.CreatedAt,
			&patient.UpdatedAt,
		); err != nil {
			return nil, models.PaginationMeta{}, helpers.ErrDatabase
		}

		patients = append(patients, patient)
	}

	if err := rows.Err(); err != nil {
		return nil, models.PaginationMeta{}, helpers.ErrDatabase
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + filter.Limit - 1) / filter.Limit
	}

	meta := models.PaginationMeta{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return patients, meta, nil
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
		if err == sql.ErrNoRows {
			return nil, helpers.ErrNotFound
		}

		return nil, helpers.ErrDatabase
	}

	return &patient, nil
}

// func (r *PatientRepository) Create(patient *models.Patient) error {
// 	query := `
// 		INSERT INTO patients (
// 			medical_record_number,
// 			name,
// 			birth_date,
// 			gender
// 		)
// 		OUTPUT
// 			INSERTED.id,
// 			INSERTED.created_at,
// 			INSERTED.updated_at
// 		VALUES (
// 			@p1,
// 			@p2,
// 			@p3,
// 			@p4
// 		)
// 	`

//		return r.db.QueryRow(
//			query,
//			patient.MedicalRecordNumber,
//			patient.Name,
//			patient.BirthDate,
//			patient.Gender,
//		).Scan(
//			&patient.ID,
//			&patient.CreatedAt,
//			&patient.UpdatedAt,
//		)
//	}
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

	err := r.db.QueryRow(
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

	if err != nil {
		// SQL Server duplicate key
		if strings.Contains(err.Error(), "UQ_patients_medical_record_number") ||
			strings.Contains(err.Error(), "duplicate key") {
			return helpers.ErrDuplicateMedicalRecordNumber
		}

		return err
	}

	return nil
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
		return helpers.ErrDatabase
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return helpers.ErrDatabase
	}

	if rowsAffected == 0 {
		return helpers.ErrNotFound
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
		return helpers.ErrDatabase
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return helpers.ErrDatabase
	}

	if rowsAffected == 0 {
		return helpers.ErrNotFound
	}

	return nil
}

func (r *PatientRepository) Patch(
	id int64,
	patch *models.PatientPatch,
) error {
	query := `
		UPDATE patients
		SET
			name = COALESCE(@p2, name),
			birth_date = COALESCE(@p3, birth_date),
			gender = COALESCE(@p4, gender),
			updated_at = GETDATE()
		WHERE id = @p1
	`

	result, err := r.db.Exec(
		query,
		id,
		patch.Name,
		patch.BirthDate,
		patch.Gender,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return helpers.ErrNotFound
	}

	return nil
}
