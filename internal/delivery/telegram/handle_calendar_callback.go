package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"

	"github.com/binaryty/evbot/internal/domain/entities"
)

// handleCalendarCallback ...
func (h *Handler) handleCalendarCallback(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	parts := strings.Split(query.Data, ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid calendar callback")
	}

	userID := query.From.ID
	state, err := h.stateRepo.GetState(ctx, userID)
	if err != nil {
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
			query.Message.Chat.ID,
			query.Message.MessageID,
			generateCalendar(calendar.CurrentDate, state.SelectedDate),
		)
		_, err := h.bot.Send(editMarkup)
		return err

	case "select":
		// выбор даты
		selectedDate, _ := time.Parse(dateFormat, parts[2])
		state.SelectedDate = selectedDate
		state.TempEvent.Date = selectedDate

		state.Step = domain.StepTime
		err := h.stateRepo.SaveState(ctx, userID, *state)
		if err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
		edit := tgbotapi.NewEditMessageReplyMarkup(
			query.Message.Chat.ID,
			query.Message.MessageID,
			generateCalendar(selectedDate, state.SelectedDate),
		)
		h.bot.Send(edit)

		return nil

	case "confirm":
		// Подтвержение даты
		if state.TempEvent.Date.IsZero() {
			h.sendError(query.Message.Chat.ID, " Дата не выбрана")
			return nil
		}

		// Запросить время
		msg := tgbotapi.NewMessage(query.Message.Chat.ID,
			fmt.Sprintf("Выбрана дата: %s\nВведите время в формате ЧЧ:ММ",
				state.TempEvent.Date.Format(dateFormat)))
		h.bot.Send(msg)
	}

	return nil
}
