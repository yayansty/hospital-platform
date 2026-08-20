package repositories

import "database/sql"

type HealthRepository struct {
	DB *sql.DB
}

func NewHealthRepository(db *sql.DB) *HealthRepository {
	return &HealthRepository{
		DB: db,
	}
}

func (r *HealthRepository) CheckDatabase() error {
	var result int

	return r.DB.QueryRow("SELECT 1").Scan(&result)
}
