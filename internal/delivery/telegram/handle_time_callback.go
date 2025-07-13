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

// handleTimeCallback ...
func (h *Handler) handleTimeCallback(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	parts := strings.Split(query.Data, ":")
	if len(parts) < 3 {
		return fmt.Errorf("invalid time callback")
	}

	userID := query.From.ID
	chatID := query.Message.Chat.ID

	h.logger.Debug("handleTimeCallback",
		slog.Int64("userID", userID),
		slog.Int64("chatID", chatID),
		slog.String("data", query.Data))

	// Получаем ID владельца события из callback данных
	var ownerID int64
	if len(parts) >= 4 {
		ownerID, _ = strconv.ParseInt(parts[3], 10, 64)
		h.logger.Debug("time callback with owner",
			slog.Int64("ownerID", ownerID),
			slog.Int64("userID", userID))
	}

	// Если указан владелец и текущий пользователь не владелец
	if ownerID != 0 && userID != ownerID {
		h.logger.Debug("time picker accessed by different user",
			slog.Int64("userID", userID),
			slog.Int64("ownerID", ownerID))
		h.sendCallback(query.ID, EmCross, "Только создатель может управлять выбором времени")
		return nil
	}

	// Используем правильный ID для получения состояния
	stateUserID := userID
	if ownerID != 0 {
		stateUserID = ownerID
	}

	// Проверяем наличие состояния пользователя
	_, err := h.stateRepo.GetState(ctx, stateUserID)
	if err != nil {
		h.logger.Error("failed to get state in time callback",
			slog.String("error", err.Error()),
			slog.Int64("userID", stateUserID))
		h.sendError(chatID, MsgSessionExpired)
		return err
	}

	switch parts[1] {
	case "time_h_select":
		return h.handleHoursSelection(ctx, query, parts[2], stateUserID)
	case "time_m_select":
		return h.handleMinuteSelection(ctx, query, parts[2], stateUserID)
	case "time_confirm":
		return h.confirmTimeSelection(ctx, query, stateUserID)
	case "time_cancel":
		return h.cancelTimeSelection(ctx, query, stateUserID)
	}

	return nil
}

func (h *Handler) handleHoursSelection(ctx context.Context, query *tgbotapi.CallbackQuery, hour string, userID int64) error {
	state, err := h.stateRepo.GetState(ctx, userID)
	if err != nil {
		return err
	}

	hh, err := time.Parse("15", hour)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	state.TimePicker.TempHours = hh.Hour()
	state.TimePicker.Step = domain.StepMinutes

	editMarkup := tgbotapi.NewEditMessageReplyMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		GenerateTimePicker(&state.TimePicker, userID),
	)

	h.bot.Send(editMarkup)

	return h.stateRepo.SaveState(ctx, userID, *state)
}

func (h *Handler) handleMinuteSelection(ctx context.Context, query *tgbotapi.CallbackQuery, minute string, userID int64) error {
	state, err := h.stateRepo.GetState(ctx, userID)
	if err != nil {
		return err
	}

	m, err := time.Parse("04", minute)
	if err != nil {
		return err
	}
	state.TimePicker.TempMinutes = m.Minute()

	newTime := time.Date(
		state.TempEvent.Date.Year(),
		state.TempEvent.Date.Month(),
		state.TempEvent.Date.Day(),
		state.TimePicker.TempHours,
		state.TimePicker.TempMinutes,
		0, 0, time.UTC,
	)

	state.TimePicker.SelectedTime = newTime
	state.TimePicker.Step = "minutes"

	edit := tgbotapi.NewEditMessageText(
		query.Message.Chat.ID,
		query.Message.MessageID,
		fmt.Sprintf("Выбрано время: %s", newTime.Format("15:04")),
	)
	timePicker := GenerateTimePicker(&state.TimePicker, userID)

	edit.ReplyMarkup = &timePicker

	h.bot.Send(edit)

	return h.stateRepo.SaveState(ctx, userID, *state)
}

func (h *Handler) confirmTimeSelection(ctx context.Context, query *tgbotapi.CallbackQuery, userID int64) error {
	state, err := h.stateRepo.GetState(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get state: %w", err)
	}

	state.TempEvent.Date = state.TimePicker.SelectedTime
	state.Step = "confirm"
	state.MessageID = query.Message.MessageID // Сохраняем ID сообщения

	// Удаляем сообщение с выбором времени
	delMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
	h.bot.Send(delMsg)

	return h.handleFinishEventCreation(ctx, userID, query.Message.Chat.ID, state)
}

func (h *Handler) cancelTimeSelection(ctx context.Context, query *tgbotapi.CallbackQuery, userID int64) error {
	chatID := query.Message.Chat.ID

	// Получаем текущее состояние пользователя
	state, err := h.stateRepo.GetState(ctx, userID)
	if err != nil {
		h.logger.Error("failed to get state in cancelTimeSelection",
			slog.String("error", err.Error()),
			slog.Int64("userID", userID))
		return fmt.Errorf("failed to get state: %w", err)
	}

	// Возвращаемся к выбору даты
	state.Step = domain.StepDate

	// Сохраняем состояние с обновленным шагом
	if err := h.stateRepo.SaveState(ctx, userID, *state); err != nil {
		h.logger.Error("failed to save state in cancelTimeSelection",
			slog.String("error", err.Error()),
			slog.Int64("userID", userID))
		return fmt.Errorf("failed to save state: %w", err)
	}

	// Удаляем сообщение с выбором времени
	delMsg := tgbotapi.NewDeleteMessage(chatID, query.Message.MessageID)
	h.bot.Send(delMsg)

	// Отправляем календарь выбора даты
	msgID, err := h.sendDateCalendarForUser(chatID, userID)
	if err == nil && msgID > 0 {
		// Сохраняем новый ID сообщения в состоянии
		state.MessageID = msgID
		if saveErr := h.stateRepo.SaveState(ctx, userID, *state); saveErr != nil {
			h.logger.Error("failed to save state with new message ID",
				slog.String("error", saveErr.Error()),
				slog.Int64("userID", userID))
		}
	}
	return err
}
