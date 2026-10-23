package repositories

import (
	"context"
	"database/sql"
	"errors"

	"medic-api/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*models.User, error) {
	const query = `
		SELECT
			id,
			username,
			password_hash,
			full_name,
			role,
			is_active,
			created_at
		FROM dbo.users
		WHERE username = @p1
	`

	var user models.User

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
