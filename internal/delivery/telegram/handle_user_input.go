package telegram

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/repository"
)

// handleUserInput ...
func (h *Handler) handleUserInput(ctx context.Context, update *tgbotapi.Update, text string) error {
	state, err := h.stateRepo.GetState(ctx, update.Message.From.ID)
	if err != nil {
		if errors.Is(err, repository.ErrStateNotFound) {
			return nil
		}
		return fmt.Errorf("get state error: %w", err)
	}

	defer func() {
		if state.Step == domain.StepCompleted {
			_ = h.stateRepo.DeleteState(ctx, update.Message.From.ID)
		}
	}()

	switch state.Step {
	case domain.StepTitle:
		return h.handleTitleStep(ctx, update, text, *state)
	case domain.StepDescription:
		return h.handleDescriptionStep(ctx, update, text, *state)
	case domain.StepTime:
		return h.handleFinishEventCreation(ctx, update, text)
	default:
		return nil
	}
}
