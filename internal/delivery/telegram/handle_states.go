package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/util"
)

// handleTitleStep ...
func (h *Handler) handleTitleStep(
	ctx context.Context,
	userID int64,
	chatID int64,
	text string,
	state domain.EventState,
) error {
	if len(text) > MaxTitleLength {
		h.sendError(chatID, MsgTitleTooLong)
		return nil
	}

	state.TempEvent.Title = text
	state.Step = domain.StepDescription

	if err := h.stateRepo.SaveState(ctx, userID, state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	// Если есть сохраненный ID сообщения, обновляем его
	if state.MessageID > 0 {
		edit := tgbotapi.NewEditMessageText(chatID, state.MessageID, "Введите описание события:")
		_, err := h.bot.Send(edit)
		if err != nil {
			h.logger.Error("failed to update message for description step",
				slog.String("error", err.Error()))
			// Если не удалось обновить, отправляем новое сообщение
			msg := tgbotapi.NewMessage(chatID, "Введите описание события:")
			h.bot.Send(msg)
		}
	} else {
		// Если нет сохраненного ID, отправляем новое сообщение
		msg := tgbotapi.NewMessage(chatID, "Введите описание события:")
		h.bot.Send(msg)
	}

	return nil
}

// handleDescriptionStep ...
func (h *Handler) handleDescriptionStep(
	ctx context.Context,
	userID int64,
	chatID int64,
	text string,
	state domain.EventState,
) error {
	if len(text) > MaxDescriptionLength {
		h.sendError(chatID, MsgDescriptionTooLong)
		return nil
	}

	state.TempEvent.Description = text
	state.Step = domain.StepDate

	if err := h.stateRepo.SaveState(ctx, userID, state); err != nil {
		h.logger.Error("failed to save state in handleDescriptionStep",
			slog.String("error", err.Error()),
			slog.Int64("userID", userID))
		h.sendError(chatID, MsgSaveError)
		return fmt.Errorf("failed to save state: %w", err)
	}

	// Получаем календарь
	calendar := domain.NewCalendar()
	calendarMarkup := generateCalendar(calendar.CurrentDate, calendar.CurrentDate, userID)

	// Если есть сохраненный ID сообщения, обновляем его
	if state.MessageID > 0 {
		edit := tgbotapi.NewEditMessageText(chatID, state.MessageID, "Выберите дату события:")
		edit.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: calendarMarkup.InlineKeyboard,
		}
		_, err := h.bot.Send(edit)
		if err != nil {
			h.logger.Error("failed to update message with calendar",
				slog.String("error", err.Error()))
			// Если не удалось обновить, отправляем новый календарь
			msgID, err := h.sendDateCalendarForUser(chatID, userID)
			if err == nil {
				// Сохраняем новый ID сообщения
				state.MessageID = msgID
				_ = h.stateRepo.SaveState(ctx, userID, state)
			}
		}
		return nil
	} else {
		// Если нет сохраненного ID, отправляем новый календарь
		msgID, err := h.sendDateCalendarForUser(chatID, userID)
		if err == nil {
			// Сохраняем ID сообщения
			state.MessageID = msgID
			_ = h.stateRepo.SaveState(ctx, userID, state)
		}
		return err
	}
}

// sendDateCalendarForUser ...
func (h *Handler) sendDateCalendarForUser(chatID int64, userID int64) (int, error) {
	calendar := domain.NewCalendar()
	msg := tgbotapi.NewMessage(chatID, "Выберите дату события:")
	msg.ReplyMarkup = generateCalendar(calendar.CurrentDate, calendar.CurrentDate, userID)
	sentMsg, err := h.bot.Send(msg)
	if err != nil {
		return 0, err
	}

	return sentMsg.MessageID, nil
}

// handleTimeStep ...
func (h *Handler) handleTimeStep(ctx context.Context, userID int64, chatID int64, messageID int) error {
	h.logger.Debug("handleTimeStep - START",
		slog.Int64("userID", userID),
		slog.Int64("chatID", chatID),
		slog.Int("messageID", messageID))

	state, err := h.stateRepo.GetState(ctx, userID)
	if err != nil {
		h.logger.Error("failed to get state in handleTimeStep:",
			slog.String("error", err.Error()),
			slog.Int64("userID", userID),
			slog.Int64("chatID", chatID))
		h.sendError(chatID, MsgSessionExpired)
		return fmt.Errorf("failed to get state: %w", err)
	}

	h.logger.Debug("handleTimeStep - state found",
		slog.Int64("userID", userID),
		slog.String("step", state.Step),
		slog.Any("tempEvent", state.TempEvent))

	tp := domain.TimePicker{
		SelectedTime: state.TempEvent.Date,
		Step:         domain.StepHours,
	}

	// Если messageID = 0, используем сохраненный ID из состояния или отправляем новое сообщение
	// Иначе обновляем указанное сообщение
	if messageID == 0 && state.MessageID > 0 {
		messageID = state.MessageID
	}

	if messageID == 0 {
		msg := tgbotapi.NewMessage(chatID, "Выберите время:")
		msg.ReplyMarkup = GenerateTimePicker(&tp, userID)
		sentMsg, err := h.bot.Send(msg)
		if err != nil {
			h.logger.Error("failed to send time picker message",
				slog.String("error", err.Error()))
			return fmt.Errorf("failed to send time picker: %w", err)
		}
		h.logger.Debug("time picker message sent", slog.Any("message", sentMsg))

		// Сохраняем ID нового сообщения
		state.MessageID = sentMsg.MessageID
	} else {
		// Обновляем существующее сообщение
		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Выберите время:")
		edit.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: GenerateTimePicker(&tp, userID).InlineKeyboard,
		}
		_, err = h.bot.Send(edit)
		if err != nil {
			h.logger.Error("failed to update message with time picker",
				slog.String("error", err.Error()))
			return fmt.Errorf("failed to update message: %w", err)
		}
		h.logger.Debug("calendar message updated to time picker",
			slog.Int("messageID", messageID))

		// Сохраняем ID сообщения
		state.MessageID = messageID
	}

	state.TimePicker = tp

	err = h.stateRepo.SaveState(ctx, userID, *state)
	if err != nil {
		h.logger.Error("failed to save state with time picker",
			slog.String("error", err.Error()),
			slog.Int64("userID", userID))
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}

// handleFinishEventCreation ...
func (h *Handler) handleFinishEventCreation(
	ctx context.Context,
	userID int64,
	chatID int64,
	state *domain.EventState,
) error {
	// Валидация данных
	if state.TempEvent.Title == "" || state.TempEvent.Date.IsZero() {
		h.sendError(chatID, MsgIncompleteData)
		return errors.New("incomplete event data")
	}

	// создаем полный объект события
	event := domain.Event{
		UserID:      userID,
		Title:       state.TempEvent.Title,
		Description: state.TempEvent.Description,
		Date:        state.TempEvent.Date,
		CreatedAt:   time.Now().UTC(),
	}

	// Сохраняем в БД
	eventID, err := h.eventUC.CreateEvent(ctx, userID, event)
	if err != nil {
		if errors.Is(err, domain.ErrAdminOnly) {
			h.sendError(chatID, MsgAdminOnly)
		} else {
			h.sendError(chatID, MsgEventSaveError)
		}
		return fmt.Errorf("failed to create event: %w", err)
	}

	// Отправляем подтверждение
	msgText := fmt.Sprintf(
		"🎉 *Событие успешно создано\\!*\n\n"+
			"📌 *Название:* %s\n"+
			"📝 *Описание:* %s\n"+
			"⏰ *Дата и время:* %s",
		util.EscapeMarkdownV2(event.Title),
		util.EscapeMarkdownV2(event.Description),
		event.Date.Format("02\\.01\\.2006 15\\:04"),
	)

	// Создаем кнопки управления
	isAdmin := h.isAdmin(userID)
	markup := createEventButtons(eventID, false, isAdmin)

	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.ReplyMarkup = markup

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send confirmation: %w", err)
	}

	// Очищаем состояние
	if err := h.stateRepo.DeleteState(ctx, userID); err != nil {
		log.Printf("Failed to clear user state: %v", err)
	}

	return nil
}
