package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/binaryty/evbot/internal/config"
	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/repository"
)

type EventUseCase struct {
	repo    repository.EventRepository
	cfg     *config.Config
	isAdmin func(int64) bool
}

// NewEventUseCase ...
func NewEventUseCase(repo repository.EventRepository, cfg *config.Config) *EventUseCase {
	return &EventUseCase{
		repo: repo,
		cfg:  cfg,
		isAdmin: func(userID int64) bool {
			for _, adminID := range cfg.AdminIDs {
				if adminID == userID {
					return true
				}
			}
			return false
		},
	}
}

// CreateEvent ...
func (uc *EventUseCase) CreateEvent(ctx context.Context, userID int64, event domain.Event) (int64, error) {
	if !uc.isAdmin(userID) {
		return 0, domain.ErrAdminOnly
	}

	if err := event.Validate(); err != nil {
		return 0, err
	}

	event.UserID = userID
	event.CreatedAt = time.Now().UTC()

	id, err := uc.repo.Save(ctx, event)
	if err != nil {
		return 0, fmt.Errorf("failed to save event: %w", err)
	}

	return id, nil
}

// GetEventByID ...
func (uc *EventUseCase) GetEventByID(ctx context.Context, eventID int64) (*domain.Event, error) {
	return uc.repo.GetByID(ctx, eventID)
}

// ListUserEvents ...
func (uc *EventUseCase) ListUserEvents(ctx context.Context, userID int64) ([]domain.Event, error) {
	return uc.repo.GetByUserID(ctx, userID)
}

// ListEvents ...
func (uc *EventUseCase) ListEvents(ctx context.Context) ([]domain.Event, error) {
	return uc.repo.GetAll(ctx)
}

// DeleteEvent ...
func (uc *EventUseCase) DeleteEvent(ctx context.Context, eventID int64) error {
	_, err := uc.repo.GetByID(ctx, eventID)
	if errors.Is(err, domain.ErrEventNotFound) {
		return err
	} else if err != nil {
		return fmt.Errorf("fetch before delete failed: %w", err)
	}

	return uc.repo.Delete(ctx, eventID)
}
