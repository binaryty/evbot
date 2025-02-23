package telegram

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"time"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/util"
)

func (h *Handler) handleTitleStep(ctx context.Context, userID int64, chatID int64, text string, state domain.EventState) error {
	if strings.TrimSpace(text) == "" {
		h.sendError(chatID, "Название не может быть пустым")
		return nil
	}
	if len(text) > 100 {
		h.sendError(chatID, "Слишком длинное название (макс. 100 символов)")
		return nil
	}

	state.TempEvent.Title = text
	state.Step = domain.StepDescription

	// TODO: версионирование состояний
	if err := h.stateRepo.SaveState(ctx, userID, state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	msg := tgbotapi.NewMessage(chatID, "Введите описание события:")
	h.bot.Send(msg)

	return nil
}

func (h *Handler) handleDescriptionStep(ctx context.Context, userID int64, chatID int64, text string, state domain.EventState) error {
	if len(text) > 500 {
		h.sendError(chatID, "Слишком длинное описание (макс. 500 символов)")
		return nil
	}

	state.TempEvent.Description = text
	state.Step = domain.StepDate

	//TODO: версионирование
	if err := h.stateRepo.SaveState(ctx, userID, state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return h.sendDateCalendar(chatID)
}

// Отправка календаря для выбора даты
func (h *Handler) sendDateCalendar(chatID int64) error {
	calendar := domain.NewCalendar()
	msg := tgbotapi.NewMessage(chatID, "Выберите дату события:")
	msg.ReplyMarkup = generateCalendar(calendar)
	_, err := h.bot.Send(msg)
	return err
}

func (h *Handler) handleDateState(ctx context.Context, userID int64, chatID int64, text string, state *domain.EventState) error {
	date, err := time.Parse("02.01.2006 15:04", text)
	if err != nil {
		h.sendError(chatID, "Неверный формат даты. попробуйте снова (ДД.ММ.ГГГГ ЧЧ:ММ):")
		return nil
	}

	if date.Before(time.Now().Add(-5 * time.Minute)) {
		h.sendError(chatID, "Дата не может быть в прошлом")
		return nil
	}

	state.TempEvent.Date = date
	state.Step = domain.StepCompleted

	if err := h.eventUC.CreateEvent(ctx, userID, state.TempEvent); err != nil {
		h.sendError(chatID, "Ошибка сохранения события")
		return fmt.Errorf("failed to create event: %w", err)
	}

	msg := tgbotapi.NewMessage(chatID, "✅ Событие успешно сохранено!")
	h.bot.Send(msg)

	return nil
}

func (h *Handler) handleFinishEventCreation(ctx context.Context, userID int64, chatID int64) error {
	state, err := h.stateRepo.GetState(ctx, userID)
	if err != nil {
		h.sendError(chatID, "Ошибка создания события")
		return fmt.Errorf("get state error: %w", err)
	}

	// Валидация данных
	if state.TempEvent.Title == "" || state.TempEvent.Date.IsZero() {
		h.sendError(chatID, "Не все данные заполнены")
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
	if err := h.eventUC.CreateEvent(ctx, userID, event); err != nil {
		h.sendError(chatID, "Ошибка сохранения события")
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
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"✏️ Редактировать",
				fmt.Sprintf("edit_event:%d", event.ID),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"👥 Участники",
				fmt.Sprintf("participants:%d", event.ID),
			),
		),
	)

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
