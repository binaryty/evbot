package usecase

import (
	"context"
	"fmt"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/repository"
)

type UserUseCase struct {
	repo repository.UserRepository
}

// NewUserUseCase ...
func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

// GetUserByID ...
func (uc *UserUseCase) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id %d: %w", userID, err)
	}

	return user, nil
}

// CreateOrUpdate ...
func (uc *UserUseCase) CreateOrUpdate(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	return uc.repo.CreateOrUpdate(ctx, user)
}
