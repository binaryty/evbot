package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

type Scanner interface {
	Scan(dest ...any) error
}

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

// Save ...
func (r *EventRepository) Save(ctx context.Context, e domain.Event) (int64, error) {
	const query = `
		INSERT INTO events
			(user_id, title, description, date, created_at)
		VALUES (?, ?, ?, ?, ?)`

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

// GetByID ...
func (r *EventRepository) GetByID(ctx context.Context, eventID int64) (*domain.Event, error) {
	const query = `
		SELECT id, user_id, title, description, date, created_at
		FROM events 
		WHERE id = ?`

	event, err := scanEvent(r.db.QueryRowContext(ctx, query, eventID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}

		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return event, nil
}

// GetByUserID ...
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
		event, err := scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		events = append(events, *event)
	}

	return events, nil
}

// GetAll ...
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
		event, err := scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		events = append(events, *event)

	}

	return events, nil
}

// Delete ...
func (r *EventRepository) Delete(ctx context.Context, eventID int64) error {
	const query = `
		DELETE FROM events
		WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected err: %w", err)
	}
	if rows == 0 {
		return domain.ErrEventNotFound
	}

	return err
}

// scanEvent ...
func scanEvent(row Scanner) (*domain.Event, error) {
	var event domain.Event

	err := row.Scan(
		&event.ID,
		&event.UserID,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan event: %w", err)
	}

	return &event, nil
}
