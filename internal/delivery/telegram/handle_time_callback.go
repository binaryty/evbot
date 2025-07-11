package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"

	"github.com/binaryty/evbot/internal/delivery/telegram/timepicker"
)

func (h *Handler) handleTimeCallback(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	parts := strings.Split(query.Data, ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid time callback")
	}

	switch parts[0] {
	case "time_h":
		return h.handleHoursSelection(ctx, query, parts[1])
	case "time_m":
		return h.handleMinuteSelection(ctx, query, parts[1])
	case "time_confirm":
		return h.confirmTimeSelection(ctx, query)
	case "time_cancel":
		return h.cancelTimeSelection(ctx, query)
	}
	return nil
}

func (h *Handler) handleHoursSelection(ctx context.Context, query *tgbotapi.CallbackQuery, hour string) error {
	state, err := h.stateRepo.GetState(ctx, query.From.ID)
	if err != nil {
		return err
	}

	hh, err := time.Parse("15", hour)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	state.TimePicker.TempHours = hh.Hour()
	state.TimePicker.Step = "minutes"

	editMarkup := tgbotapi.NewEditMessageReplyMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		timepicker.GenerateTimePicker(&state.TimePicker),
	)

	h.bot.Send(editMarkup)

	return h.stateRepo.SaveState(ctx, query.From.ID, *state)
}

func (h *Handler) handleMinuteSelection(ctx context.Context, query *tgbotapi.CallbackQuery, minute string) error {
	state, err := h.stateRepo.GetState(ctx, query.From.ID)
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
	state.TimePicker.Step = "hours"

	edit := tgbotapi.NewEditMessageText(
		query.Message.Chat.ID,
		query.Message.MessageID,
		fmt.Sprintf("Выбрано время: %s", newTime.Format("15:04")),
	)
	timePicker := timepicker.GenerateTimePicker(&state.TimePicker)

	edit.ReplyMarkup = &timePicker

	h.bot.Send(edit)

	return h.stateRepo.SaveState(ctx, query.From.ID, *state)
}

func (h *Handler) confirmTimeSelection(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	state, err := h.stateRepo.GetState(ctx, query.From.ID)
	if err != nil {
		return fmt.Errorf("failed to get state: %w", err)
	}

	state.TempEvent.Date = state.TimePicker.SelectedTime
	state.Step = "confirm"

	delMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
	h.bot.Send(delMsg)

	// h.sendConfirmation(query.Message.Chat.ID, state.TempEvent)
	return nil
}

func (h *Handler) cancelTimeSelection(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	if err := h.stateRepo.DeleteState(ctx, query.From.ID); err != nil {
		return fmt.Errorf("failed to delete state: %w", err)
	}
	delMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
	h.bot.Send(delMsg)

	return h.sendDateCalendar(query.Message.Chat.ID)
}
