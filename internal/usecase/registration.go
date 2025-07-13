package usecase

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/repository"
)

type RegistrationUseCase struct {
	eventRepo        repository.EventRepository
	registrationRepo repository.RegistrationRepository
}

// NewRegistrationUseCase ...
func NewRegistrationUseCase(eventRepo repository.EventRepository, registrationRepo repository.RegistrationRepository) *RegistrationUseCase {
	return &RegistrationUseCase{
		eventRepo:        eventRepo,
		registrationRepo: registrationRepo,
	}
}

// ToggleRegistration ...
func (uc *RegistrationUseCase) ToggleRegistration(ctx context.Context, eventID int64, user *domain.User) (bool, error) {
	if _, err := uc.eventRepo.GetByID(ctx, eventID); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return false, err
		}

		return false, fmt.Errorf("failed to get event: %w", err)
	}

	isRegistered, err := uc.registrationRepo.IsRegistered(ctx, eventID, user.ID)
	if err != nil {
		return false, fmt.Errorf("failed to check registration: %w", err)
	}

	if isRegistered {
		if err := uc.registrationRepo.Unregister(ctx, eventID, user.ID); err != nil {
			return false, fmt.Errorf("failed to unregister: %w", err)
		}
		return false, nil
	}

	if err := uc.registrationRepo.Register(ctx, eventID, user.ID); err != nil {
		return false, fmt.Errorf("failed to register: %w", err)
	}

	return true, nil
}

// GetParticipants ...
func (uc *RegistrationUseCase) GetParticipants(ctx context.Context, eventID int64) ([]domain.Participant, error) {

	return uc.registrationRepo.GetParticipants(ctx, eventID)
}

// GetParticipantsPaginated ...
func (uc *RegistrationUseCase) GetParticipantsPaginated(ctx context.Context, eventID int64, offset int, limit int) ([]domain.Participant, int, error) {

	return uc.registrationRepo.GetParticipantsPaginated(ctx, eventID, offset, limit)
}

// IsRegistered ...
func (uc *RegistrationUseCase) IsRegistered(ctx context.Context, eventID int64, userID int64) (bool, error) {

	return uc.registrationRepo.IsRegistered(ctx, eventID, userID)
}
