package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"

	"github.com/binaryty/evbot/internal/config"
	"github.com/binaryty/evbot/internal/repository"
	"github.com/binaryty/evbot/internal/usecase"
)

const (
	EmReg    = "🎫"
	EmCross  = "❌"
	EmOk     = "✅"
	EmPeople = "👥"
	EmList   = "📋"
	EmPin    = "📌"
	EmPrev   = "◀️"
	EmNext   = "▶️"
)

type Handler struct {
	cfg            *config.Config
	bot            *tgbotapi.BotAPI
	logger         *slog.Logger
	eventUC        *usecase.EventUseCase
	registrationUC *usecase.RegistrationUseCase
	userUC         *usecase.UserUseCase
	stateRepo      repository.StateRepository
}

func NewHandler(
	cfg *config.Config,
	bot *tgbotapi.BotAPI,
	logger *slog.Logger,
	eventUC *usecase.EventUseCase,
	registrationUC *usecase.RegistrationUseCase,
	userUC *usecase.UserUseCase,
	//userRepo repository.UserRepository,
	stateRepo repository.StateRepository,
) *Handler {
	return &Handler{
		cfg:            cfg,
		bot:            bot,
		logger:         logger,
		eventUC:        eventUC,
		registrationUC: registrationUC,
		userUC:         userUC,
		stateRepo:      stateRepo,
	}
}

func (h *Handler) sendError(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, "❌ "+text)
	msg.ParseMode = "Markdown"
	h.bot.Send(msg)
}

func (h *Handler) sendCallback(queryID string, icon string, text string) {
	callback := tgbotapi.NewCallbackWithAlert(queryID, fmt.Sprintf("%s %s", icon, text))
	h.bot.Send(callback)
}

func (h *Handler) sendMsg(chatID int64, icon string, text string) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s %s", icon, text))
	msg.ParseMode = "Markdown"
	h.bot.Send(msg)
}
