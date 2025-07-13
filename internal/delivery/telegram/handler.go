package telegram

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/binaryty/evbot/internal/config"
	"github.com/binaryty/evbot/internal/repository"
	"github.com/binaryty/evbot/internal/usecase"
)

const (
	EmReg             = "🎫"
	EmCross           = "❌"
	EmOk              = "✅"
	EmPeople          = "👥"
	EmList            = "📋"
	EmPin             = "📌"
	EmPrev            = "◀️"
	EmNext            = "▶️"
	MsgSessionExpired = "Ошибка: сессия создания события истекла. Пожалуйста, начните заново с команды /new_event"
	MsgSaveError      = "Ошибка сохранения данных"
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

	if _, err := h.bot.Send(msg); err != nil {
		h.logger.Error("failed to send error message",
			slog.Int64("chatID", chatID),
			slog.String("text", text),
			slog.String("error", err.Error()))
	}
}

func (h *Handler) sendCallback(queryID string, icon string, text string) {
	callback := tgbotapi.NewCallbackWithAlert(queryID, fmt.Sprintf("%s %s", icon, text))

	if _, err := h.bot.Send(callback); err != nil {
		h.logger.Error("failed to send callback",
			slog.String("queryID", queryID),
			slog.String("text", text),
			slog.String("error", err.Error()))
	}
}

func (h *Handler) sendMsg(chatID int64, icon string, text string) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s %s", icon, text))
	msg.ParseMode = "Markdown"
	if _, err := h.bot.Send(msg); err != nil {
		h.logger.Error("failed to send message",
			slog.Int64("chatID", chatID),
			slog.String("text", text),
			slog.String("error", err.Error()))
	}
}
