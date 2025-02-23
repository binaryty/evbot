package repository

import (
	"context"
	"errors"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

var (
	ErrStateNotFound = errors.New("state not found")
)

type EventRepository interface {
	Save(ctx context.Context, event domain.Event) (int64, error)
	GetByID(ctx context.Context, eventID int64) (*domain.Event, error)
	GetByUserID(ctx context.Context, userID int64) ([]domain.Event, error)
	GetAll(ctx context.Context) ([]domain.Event, error)
	Delete(ctx context.Context, eventID int64) error
}

type StateRepository interface {
	GetState(ctx context.Context, userID int64) (*domain.EventState, error)
	SaveState(ctx context.Context, userID int64, state domain.EventState) error
	DeleteState(ctx context.Context, userID int64) error
}

type RegistrationRepository interface {
	Register(ctx context.Context, eventID int64, userID int64) error
	Unregister(ctx context.Context, eventID int64, userID int64) error
	GetParticipants(ctx context.Context, eventID int64) ([]domain.Participant, error)
	IsRegistered(ctx context.Context, eventID int64, userID int64) (bool, error)
	GetParticipantsPaginated(ctx context.Context, eventID int64, offset int, limit int) ([]domain.Participant, int, error)
}

type UserRepository interface {
	CreateOrUpdate(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, userID int64) (*domain.User, error)
}
