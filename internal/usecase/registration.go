package usecase

import (
	"context"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/repository"
)

type RegistrationUseCase struct {
	eventRepo        repository.EventRepository
	registrationRepo repository.RegistrationRepository
}

func NewRegistrationUseCase(
	eventRepo repository.EventRepository,
	registrationRepo repository.RegistrationRepository,
) *RegistrationUseCase {
	return &RegistrationUseCase{
		eventRepo:        eventRepo,
		registrationRepo: registrationRepo,
	}
}

func (uc *RegistrationUseCase) ToggleRegistration(
	ctx context.Context,
	eventID int64,
	user *domain.User,
) (bool, error) {
	_, err := uc.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return false, domain.ErrEventNotFound
	}

	isRegistered, err := uc.registrationRepo.IsRegistered(ctx, eventID, user.ID)
	if err != nil {
		return false, err
	}

	if isRegistered {
		err = uc.registrationRepo.Unregister(ctx, eventID, user.ID)
		return false, err
	}

	if err = uc.registrationRepo.Register(ctx, eventID, user.ID); err != nil {
		return false, err
	}

	return true, nil
}

func (uc *RegistrationUseCase) GetParticipants(
	ctx context.Context,
	eventID int64,
) ([]domain.Participant, error) {
	return uc.registrationRepo.GetParticipants(ctx, eventID)
}

func (uc *RegistrationUseCase) GetParticipantsPaginated(
	ctx context.Context,
	eventID int64,
	offset int,
	limit int,
) ([]domain.Participant, int, error) {
	return uc.registrationRepo.GetParticipantsPaginated(ctx, eventID, offset, limit)
}

func (uc *RegistrationUseCase) IsRegistered(ctx context.Context, eventID int64, userID int64) (bool, error) {
	return uc.registrationRepo.IsRegistered(ctx, eventID, userID)
}
