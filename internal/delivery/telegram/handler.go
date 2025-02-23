package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/binaryty/evbot/internal/repository"
	"github.com/binaryty/evbot/internal/usecase"
)

type Handler struct {
	bot            *tgbotapi.BotAPI
	eventUC        *usecase.EventUseCase
	registrationUC *usecase.RegistrationUseCase
	userUC         *usecase.UserUseCase
	userRepo       repository.UserRepository
	stateRepo      repository.StateRepository
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	eventUC *usecase.EventUseCase,
	registrationUC *usecase.RegistrationUseCase,
	userUC *usecase.UserUseCase,
	userRepo repository.UserRepository,
	stateRepo repository.StateRepository,
) *Handler {
	return &Handler{
		bot:            bot,
		eventUC:        eventUC,
		registrationUC: registrationUC,
		userUC:         userUC,
		userRepo:       userRepo,
		stateRepo:      stateRepo,
	}
}

func (h *Handler) sendError(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, "❌ "+text)
	msg.ParseMode = "Markdown"
	_, _ = h.bot.Send(msg)
}
