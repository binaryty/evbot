package usecase

import (
	"context"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

type UserUsecase interface {
	GetUserByID(ctx context.Context, userID int64) (*domain.User, error)
	CreateOrUpdate(ctx context.Context, user *domain.User) error
}

type EventUsecase interface {
	CreateEvent(ctx context.Context, userID int64, event domain.Event) (int64, error)
	GetEventByID(ctx context.Context, userID int64) (*domain.Event, error)
	ListEvents(ctx context.Context) ([]domain.Event, error)
	ListUserEvents(ctx context.Context, userID int64) ([]domain.Event, error)
	DeleteEvent(ctx context.Context, eventID int64) error
}

type RegistrationUsecase interface {
	ToggleRegistration(ctx context.Context, eventID int64, user *domain.User) (bool, error)
	GetParticipants(ctx context.Context, eventID int64) ([]domain.Participant, error)
	GetParticipantsPaginated(ctx context.Context, eventID int64, offset int, limit int) ([]domain.Participant, int, error)
	IsRegistered(ctx context.Context, eventID int64, userID int64) (bool, error)
}
