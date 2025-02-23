package telegram

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/repository"
)

func (h *Handler) handleUserInput(ctx context.Context, userID int64, chatID int64, text string) error {
	state, err := h.stateRepo.GetState(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrStateNotFound) {
			return nil
		}
		return fmt.Errorf("get state error: %w", err)
	}

	defer func() {
		if state.Step == domain.StepCompleted {
			_ = h.stateRepo.DeleteState(ctx, userID)
		}
	}()

	switch state.Step {
	case domain.StepTitle:
		return h.handleTitleStep(ctx, userID, chatID, text, *state)
	case domain.StepDescription:
		return h.handleDescriptionStep(ctx, userID, chatID, text, *state)
	case domain.StepDate:
		return h.handleDateState(ctx, userID, chatID, text, *state)
	case domain.StepTime:
		return h.handleFinishEventCreation(ctx, userID, chatID, text)
	default:
		return nil
	}
}
