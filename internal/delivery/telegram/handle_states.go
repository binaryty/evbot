package telegram

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/util"
)

// handleTitleStep ...
func (h *Handler) handleTitleStep(ctx context.Context, update *tgbotapi.Update, text string, state domain.EventState) error {
	if len(text) > 100 {
		h.sendError(update.Message.Chat.ID, "Слишком длинное название (макс. 100 символов)")
		return nil
	}

	state.TempEvent.Title = text
	state.Step = domain.StepDescription

	if err := h.stateRepo.SaveState(ctx, update.Message.From.ID, state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите описание события:")
	h.bot.Send(msg)

	return nil
}

// handleDescriptionStep ...
func (h *Handler) handleDescriptionStep(ctx context.Context, update *tgbotapi.Update, text string, state domain.EventState) error {
	if len(text) > 500 {
		h.sendError(update.Message.Chat.ID, "Слишком длинное описание (макс. 500 символов)")
		return nil
	}

	state.TempEvent.Description = text
	state.Step = domain.StepDate

	if err := h.stateRepo.SaveState(ctx, update.Message.From.ID, state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return h.sendDateCalendar(update.Message.Chat.ID)
}

// sendDateCalendar ...
func (h *Handler) sendDateCalendar(chatID int64) error {
	calendar := domain.NewCalendar()
	msg := tgbotapi.NewMessage(chatID, "Выберите дату события:")
	msg.ReplyMarkup = generateCalendar(calendar.CurrentDate, calendar.CurrentDate)
	h.bot.Send(msg)

	return nil
}

// handleFinishEventCreation ...
func (h *Handler) handleFinishEventCreation(ctx context.Context, update *tgbotapi.Update, text string) error {
	state, err := h.stateRepo.GetState(ctx, update.Message.From.ID)
	if err != nil {
		h.sendError(update.Message.Chat.ID, "Ошибка создания события")
		return fmt.Errorf("get state error: %w", err)
	}

	t, err := time.Parse("15:04", text)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	d := state.TempEvent.Date

	state.TempEvent.Date = time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)

	// Валидация данных
	if state.TempEvent.Title == "" || state.TempEvent.Date.IsZero() || state.TempEvent.Date.Hour() == 0 {
		h.sendError(update.Message.Chat.ID, "Не все данные заполнены")
		return errors.New("incomplete event data")
	}

	// создаем полный объект события
	event := domain.Event{
		UserID:      update.Message.From.ID,
		Title:       state.TempEvent.Title,
		Description: state.TempEvent.Description,
		Date:        state.TempEvent.Date,
		CreatedAt:   time.Now().UTC(),
	}

	// Сохраняем в БД
	event.ID, err = h.eventUC.CreateEvent(ctx, update.Message.From.ID, event)
	if err != nil {
		h.sendError(update.Message.Chat.ID, "Ошибка сохранения события")
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
	isAdmin := h.isAdmin(update.Message.From.ID)
	markup := createEventButtons(event.ID, false, isAdmin)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.ReplyMarkup = markup

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send confirmation: %w", err)
	}

	// Очищаем состояние
	if err := h.stateRepo.DeleteState(ctx, update.Message.From.ID); err != nil {
		log.Printf("Failed to clear user state: %v", err)
	}

	return nil
}
