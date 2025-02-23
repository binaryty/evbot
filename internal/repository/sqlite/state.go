package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/repository"
)

type StateRepository struct {
	db *sql.DB
}

func NewStateRepository(db *sql.DB) *StateRepository {
	return &StateRepository{
		db: db,
	}
}

func (r *StateRepository) GetState(ctx context.Context, userID int64) (*domain.EventState, error) {
	const query = `
		SELECT state_data, created_at 
		FROM user_states 
		WHERE user_id = ?`

	var (
		stateData []byte
		createdAt time.Time
	)

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&stateData, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrStateNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	var state domain.EventState
	if err := json.Unmarshal(stateData, &state); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &state, nil
}

func (r *StateRepository) SaveState(ctx context.Context, userID int64, state domain.EventState) error {
	const query = ` INSERT INTO user_states	(user_id, state_data, created_at)
		VALUES(?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			state_data = excluded.state_data,
			created_at = excluded.created_at`

	stateData, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		userID,
		stateData,
		time.Now().UTC(),
	)

	return err
}

func (r *StateRepository) DeleteState(ctx context.Context, userID int64) error {
	const query = `
		DELETE FROM user_states 
		WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)

	return err
}
