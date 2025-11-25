package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

// handleCalendarCallback ...
func (h *Handler) handleCalendarCallback(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	parts := strings.Split(query.Data, ":")
	if len(parts) < 3 {
		return fmt.Errorf("invalid calendar callback")
	}

	userID := query.From.ID
	chatID := query.Message.Chat.ID

	// Получаем ID владельца события из callback данных
	var ownerID int64
	if len(parts) >= 4 {
		ownerID, _ = strconv.ParseInt(parts[3], 10, 64)
		h.logger.Debug("handleCalendarCallback",
			slog.Int64("userID", userID),
			slog.Int64("chatID", chatID),
			slog.Int64("ownerID", ownerID),
			slog.String("data", query.Data))
	} else {
		h.logger.Debug("handleCalendarCallback (no owner)",
			slog.Int64("userID", userID),
			slog.Int64("chatID", chatID),
			slog.String("data", query.Data))
	}

	// Если указан владелец и текущий пользователь не владелец
	if ownerID != 0 && userID != ownerID {
		h.logger.Debug("calendar accessed by different user",
			slog.Int64("userID", userID),
			slog.Int64("ownerID", ownerID))
		h.sendCallback(query.ID, EmCross, "У вас нет доступа к календарю")
		return nil
	}

	// Используем правильный ID для получения состояния
	stateUserID := userID
	if ownerID != 0 {
		stateUserID = ownerID
	}

	state, err := h.stateRepo.GetState(ctx, stateUserID)
	if err != nil {
		h.logger.Error("failed to get state",
			slog.String("error", err.Error()),
			slog.Int64("userID", stateUserID),
			slog.Int64("chatID", chatID))
		h.sendError(chatID, MsgSessionExpired)
		return err
	}

	switch parts[1] {
	case "prev", "next":
		// Обработка навигации
		date, _ := time.Parse(dateFormat, parts[2])
		calendar := &domain.Calendar{CurrentDate: date}

		if parts[1] == "prev" {
			calendar.PrevMonth()
		} else {
			calendar.NextMonth()
		}

		state.SelectedDate = time.Time{}
		editMarkup := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			query.Message.MessageID,
			generateCalendar(calendar.CurrentDate, state.SelectedDate, stateUserID),
		)
		_, err := h.bot.Send(editMarkup)
		return err

	case "select":
		// выбор даты
		h.logger.Debug("handleCalendarCallback.select",
			slog.Int64("userID", stateUserID),
			slog.String("data", parts[2]))

		// 1. Парсим и сохраняем выбранную дату
		selectedDate, _ := time.Parse(dateFormat, parts[2])
		state.SelectedDate = selectedDate
		state.TempEvent.Date = selectedDate
		state.Step = domain.StepTime

		// 2. Обновляем календарь для отображения выбранной даты
		edit := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			query.Message.MessageID,
			generateCalendar(selectedDate, state.SelectedDate, stateUserID),
		)
		_, err := h.bot.Send(edit)
		if err != nil {
			h.logger.Error("failed to update calendar markup",
				slog.String("error", err.Error()))
		}

		// 3. Сохраняем состояние
		state.MessageID = query.Message.MessageID // Сохраняем ID сообщения
		err = h.stateRepo.SaveState(ctx, stateUserID, *state)
		if err != nil {
			h.logger.Error("failed to save state after date selection",
				slog.String("error", err.Error()),
				slog.Int64("userID", stateUserID))
			h.sendError(chatID, MsgSaveError)
			return fmt.Errorf("failed to save state: %w", err)
		}

		// 4. Отправляем выбор времени
		h.logger.Debug("proceeding to time selection", slog.Int64("userID", stateUserID))
		return h.handleTimeStep(ctx, stateUserID, chatID, query.Message.MessageID)

	case "confirm":
		// Подтвержение даты
		if state.TempEvent.Date.IsZero() {
			h.sendError(chatID, " Дата не выбрана")
			return nil
		}
		return h.handleTimeStep(ctx, stateUserID, chatID, query.Message.MessageID)
	}

	return nil
}
