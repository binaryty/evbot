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

	// Константы для валидации
	MaxTitleLength       = 100
	MaxDescriptionLength = 500

	// Сообщения об ошибках валидации
	MsgTitleTooLong       = "Слишком длинное название (макс. 100 символов)"
	MsgDescriptionTooLong = "Слишком длинное описание (макс. 500 символов)"
	MsgIncompleteData     = "Не все данные заполнены"
	MsgAdminOnly          = "🚫 Только администраторы могут создавать события"
	MsgEventSaveError     = "Ошибка сохранения события"
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

// Создайте helper-методы для логирования
func (h *Handler) logError(op string, err error, userID int64) {
	h.logger.Error(op,
		slog.String("error", err.Error()),
		slog.Int64("userID", userID))
}

func (h *Handler) logDebug(op string, userID int64, fields ...slog.Attr) {
	logger := h.logger.With(slog.Int64("userID", userID))
	if len(fields) > 0 {
		args := make([]any, len(fields))
		for i, attr := range fields {
			args[i] = attr
		}
		logger = logger.With(args...)
	}
	logger.Debug(op)
}
