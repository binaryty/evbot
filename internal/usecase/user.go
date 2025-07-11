package usecase

import (
	"context"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/repository"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (uc *UserUseCase) User(ctx context.Context, userID int64) (*domain.User, error) {
	return uc.repo.GetByID(ctx, userID)
}

func (uc *UserUseCase) CreateOrUpdate(ctx context.Context, user *domain.User) error {
	// TODO: validate user
	return uc.repo.CreateOrUpdate(ctx, user)
}
