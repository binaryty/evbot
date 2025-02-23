package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

type RegistrationRepository struct {
	db *sql.DB
}

func NewRegistrationRepository(db *sql.DB) *RegistrationRepository {
	return &RegistrationRepository{
		db: db,
	}
}

func (r *RegistrationRepository) Register(ctx context.Context, eventID int64, userID int64) error {
	const query = `
		INSERT INTO registrations(event_id, user_id, created_at)
		VALUES(?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		eventID,
		userID,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	return nil
}

func (r *RegistrationRepository) Unregister(ctx context.Context, eventID int64, userID int64) error {
	const query = `
		DELETE FROM registrations
		WHERE event_id = ? AND user_id = ?`

	res, err := r.db.ExecContext(ctx, query, eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to unregister: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrRegistrationNotFound
	}

	return nil
}
func (r *RegistrationRepository) GetParticipants(ctx context.Context, eventID int64) ([]domain.Participant, error) {
	const query = `
		SELECT u.user_id, u.first_name, u.username, r.created_at
		FROM registrations r
		JOIN users u ON r.user_id = u.user_id
		WHERE r.event_id = ?`

	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrParticipantNotFound
		}

		return nil, fmt.Errorf("failed to get participants: %w", err)
	}
	defer rows.Close()

	var participants []domain.Participant
	for rows.Next() {
		var p domain.Participant
		var createdAt time.Time

		err := rows.Scan(
			&p.ID,
			&p.FirstName,
			&p.UserName,
			&createdAt,
		)
		if err != nil {
			continue
		}

		p.RegisteredAt = createdAt
		participants = append(participants, p)
	}

	return participants, nil
}

func (r *RegistrationRepository) GetParticipantsPaginated(
	ctx context.Context,
	eventID int64,
	offset int,
	limit int) ([]domain.Participant, int, error) {
	const query = `
		SELECT u.user_id, u.first_name, u.username, r.created_at
		FROM registrations r
		JOIN users u ON r.user_id = u.user_id
		WHERE r.event_id = ?
		LIMIT ?
		OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, eventID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, domain.ErrParticipantNotFound
		}

		return nil, 0, fmt.Errorf("failed to get participants: %w", err)
	}
	defer rows.Close()

	var participants []domain.Participant
	for rows.Next() {
		var p domain.Participant
		var createdAt time.Time

		err := rows.Scan(
			&p.ID,
			&p.FirstName,
			&p.UserName,
			&createdAt,
		)
		if err != nil {
			continue
		}

		p.RegisteredAt = createdAt
		participants = append(participants, p)
	}

	var total int
	err = r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM registrations WHERE event_id = ?`,
		eventID,
	).Scan(&total)

	return participants, total, nil
}

func (r *RegistrationRepository) IsRegistered(ctx context.Context, eventID int64, userID int64) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1
			FROM registrations
			WHERE event_id = ? AND user_id = ?
	)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, eventID, userID).Scan(&exists)

	return exists, err
}
