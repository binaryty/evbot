package usecase

import (
	"context"
	"time"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/repository"
)

type EventUseCase struct {
	repo repository.EventRepository
}

func NewEventUseCase(repo repository.EventRepository) *EventUseCase {
	return &EventUseCase{
		repo: repo,
	}
}

func (uc *EventUseCase) CreateEvent(ctx context.Context, userID int64, event domain.Event) (int64, error) {
	if event.Title == "" {
		return 0, domain.ErrInvalidEventTitle
	}

	event.UserID = userID
	event.CreatedAt = time.Now().UTC()

	return uc.repo.Save(ctx, event)
}

func (uc *EventUseCase) ListUserEvents(ctx context.Context, userID int64) ([]domain.Event, error) {
	return uc.repo.GetByUserID(ctx, userID)
}

func (uc *EventUseCase) ListEvents(ctx context.Context) ([]domain.Event, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *EventUseCase) DeleteEvent(ctx context.Context, eventID int64) error {
	return uc.repo.Delete(ctx, eventID)
}
