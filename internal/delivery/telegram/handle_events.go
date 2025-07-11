package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/util"
)

// startNewEvent ...
func (h *Handler) startNewEvent(ctx context.Context, update *tgbotapi.Update) error {
	initialState := domain.EventState{
		Step: domain.StepTitle,
		TempEvent: domain.Event{
			UserID: update.Message.Chat.ID,
		},
	}

	if err := h.stateRepo.SaveState(ctx, update.Message.From.ID, initialState); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название события:")
	_, err := h.bot.Send(msg)
	return err
}

// listEvents ...
func (h *Handler) listEvents(ctx context.Context, update *tgbotapi.Update) error {
	const op = "handler.listEvents"
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	events, err := h.eventUC.ListEvents(ctx)
	if err != nil {
		h.sendError(chatID, "Ошибка получения событий")
		return fmt.Errorf("%s:list events error: %w", op, err)
	}

	if len(events) == 0 {
		msg := tgbotapi.NewMessage(chatID, "У вас пока нет событий")
		h.bot.Send(msg)
		return nil
	}

	var messages []tgbotapi.Chattable
	isAdmin := h.isAdmin(userID)

	for _, event := range events {
		// Проверяем регистрацию пользователя
		isRegistered, err := h.registrationUC.IsRegistered(ctx, event.ID, userID)
		if err != nil {
			log.Printf("failed to check if user is registered: %v", err)
			continue
		}

		// Формируем сообщение для каждого события
		eventOwner, _ := h.userUC.User(ctx, event.UserID)

		text := fmt.Sprintf(
			"📌 %s\n"+
				"📝 %s\n"+
				"⏰ %s\n"+
				"*Автор:* %s",
			util.EscapeMarkdownV2(event.Title),
			util.EscapeMarkdownV2(event.Description),
			event.Date.Format("02\\.01\\.2006 15\\:04"),
			util.EscapeMarkdownV2(eventOwner.UserName),
		)

		buttons := createEventButtons(event.ID, isRegistered, isAdmin)

		// Создаем сообщение с кнопками
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = buttons
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		messages = append(messages, msg)
	}

	// Отправляем основное сообщение с инструкцией
	infoMsg := tgbotapi.NewMessage(chatID,
		EmList+" *Список ваших событий*\n"+
			"Используйте кнопки под каждым событием для управления:")
	infoMsg.ParseMode = "Markdown"
	messages = append([]tgbotapi.Chattable{infoMsg}, messages...)

	// Отправляем все сообщения
	for _, msg := range messages {
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("%s: failed to sending message: %v", op, err)
		}
	}

	return nil
}

// createEventButtons ...
func createEventButtons(eventID int64, isRegistered bool, isAdmin bool) tgbotapi.InlineKeyboardMarkup {
	row := []tgbotapi.InlineKeyboardButton{
		createRegButton(eventID, isRegistered),
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", EmPeople, "Участники"),
			fmt.Sprintf("participants:%d", eventID),
		),
	}

	if isAdmin {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", EmCross, "Удалить"),
			fmt.Sprintf("delete_confirm:%d", eventID),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(row)
}

// createRegButton ...
func createRegButton(eventID int64, isRegistered bool) tgbotapi.InlineKeyboardButton {
	text, icon := "Регистрация", EmReg
	if isRegistered {
		text, icon = "Зарегистрирован", EmOk
	}

	return tgbotapi.NewInlineKeyboardButtonData(
		fmt.Sprintf("%s %s", icon, text),
		fmt.Sprintf("register:%d", eventID),
	)
}
