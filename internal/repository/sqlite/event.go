package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) Save(ctx context.Context, e domain.Event) (int64, error) {
	const query = `
		INSERT INTO events
			(user_id, title, description, date, created_at)
		VALUES (?, ?, ?, ?, ?) RETURNING id`

	res, err := r.db.ExecContext(ctx, query,
		e.UserID,
		e.Title,
		e.Description,
		e.Date.UTC(),
		e.CreatedAt.UTC(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to save event: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *EventRepository) GetByID(ctx context.Context, eventID int64) (*domain.Event, error) {
	const query = `
		SELECT id, user_id, title, description, date, created_at
		FROM events 
		WHERE id = ?`

	var event domain.Event
	var dateStr, createdAtStr string

	err := r.db.QueryRowContext(ctx, query, eventID).Scan(
		&event.ID,
		&event.UserID,
		&event.Title,
		&event.Description,
		&dateStr,
		&createdAtStr,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}

		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	event.Date, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	event.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	return &event, nil
}

func (r *EventRepository) GetByUserID(ctx context.Context, userID int64) ([]domain.Event, error) {
	const query = `
		SELECT id, user_id, title, description, date, created_at
		FROM events
		WHERE user_id = ?
		ORDER BY date DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		var dateStr, createdAtStr string

		if err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Title,
			&event.Description,
			&dateStr,
			&createdAtStr,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		event.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		event.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetAll(ctx context.Context) ([]domain.Event, error) {
	const query = `
		SELECT id, user_id, title, description, date, created_at
		FROM events
		ORDER BY date DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}

		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		var dateStr, createdAtStr string

		if err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Title,
			&event.Description,
			&dateStr,
			&createdAtStr,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		event.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		event.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) Delete(ctx context.Context, eventID int64) error {
	const query = `
		DELETE FROM events
		WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrEventNotFound
	}

	return err
}
