package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"

	"github.com/binaryty/evbot/internal/domain/entities"
)

const (
	dateFormat = "02.01.2006"
)

func generateCalendar(calendar *domain.Calendar) tgbotapi.InlineKeyboardMarkup {
	now := time.Now()
	current := calendar.CurrentDate
	year, month, _ := current.Date()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	var keyboard [][]tgbotapi.InlineKeyboardButton

	// Заголовок с названием месяца и навигацией
	header := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			"◀️",
			fmt.Sprintf("calendar:prev:%s", current.Format(dateFormat)),
		),
		tgbotapi.NewInlineKeyboardButtonData(
			current.Format("January 2006"),
			"ignore",
		),
		tgbotapi.NewInlineKeyboardButtonData(
			"▶️",
			fmt.Sprintf("calendar:next:%s", current.Format(dateFormat)),
		),
	}
	keyboard = append(keyboard, header)

	// Заголовок дней недели
	weekDays := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	var weekRow []tgbotapi.InlineKeyboardButton
	for _, day := range weekDays {
		weekRow = append(weekRow, tgbotapi.NewInlineKeyboardButtonData(day, "ignore"))
	}
	keyboard = append(keyboard, weekRow)

	// Сетка с днями
	var row []tgbotapi.InlineKeyboardButton
	for i := 0; i < int(firstDay.Weekday()-1); i++ {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
	}

	for d := firstDay; d.Month() == month; d = d.AddDate(0, 0, 1) {
		dayText := fmt.Sprintf("%d", d.Day())
		callbackData := fmt.Sprintf("calendar:select:%s", d.Format(dateFormat))

		// подсветка текущего дня
		if d.Year() == now.Year() && d.Month() == now.Month() && d.Day() == now.Day() {
			dayText = "✅ " + dayText
		}

		btn := tgbotapi.NewInlineKeyboardButtonData(dayText, callbackData)
		row = append(row, btn)

		if len(row) == 7 {
			keyboard = append(keyboard, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}

	// кнопка подтверждения
	keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Готово", "calendar:confirm"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
