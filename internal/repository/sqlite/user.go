package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateOrUpdate(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (user_id, first_name, username)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			first_name = excluded.first_name,
			username = excluded.username,
			updated_at = CURRENT_TIMESTAMP`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.FirstName,
		user.UserName,
	)
	if err != nil {
		return fmt.Errorf("failed to create or update user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	const query = `
		SELECT user_id, first_name, username
		FROM users
		WHERE user_id = ?`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.UserName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
