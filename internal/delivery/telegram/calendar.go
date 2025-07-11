package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	dateFormat = "02.01.2006"
)

func generateCalendar(currentDate time.Time, selectedDate time.Time) tgbotapi.InlineKeyboardMarkup {
	now := time.Now()

	var keyboard [][]tgbotapi.InlineKeyboardButton

	// Заголовок с названием месяца и навигацией
	header := generateHeader(currentDate)
	keyboard = append(keyboard, header)

	// Заголовок дней недели
	weekRow := generateWeakRow()
	keyboard = append(keyboard, weekRow)

	// Сетка с днями
	grid := generateGrid(&keyboard, now, currentDate, selectedDate)
	if len(grid) > 0 {
		keyboard = append(keyboard, grid)
	}

	// кнопка подтверждения
	keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Готово", "calendar:confirm"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

// generateHeader ...
func generateHeader(currentDate time.Time) []tgbotapi.InlineKeyboardButton {
	return []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			EmPrev,
			fmt.Sprintf("calendar:prev:%s", currentDate.Format(dateFormat)),
		),
		tgbotapi.NewInlineKeyboardButtonData(
			currentDate.Format("January 2006"),
			"ignore",
		),
		tgbotapi.NewInlineKeyboardButtonData(
			EmNext,
			fmt.Sprintf("calendar:next:%s", currentDate.Format(dateFormat)),
		),
	}
}

// generateWeakRow ...
func generateWeakRow() []tgbotapi.InlineKeyboardButton {
	weekDays := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	var weekRow []tgbotapi.InlineKeyboardButton
	for _, day := range weekDays {
		weekRow = append(weekRow, tgbotapi.NewInlineKeyboardButtonData(day, "ignore"))
	}

	return weekRow
}

// generateGrid ...
func generateGrid(keyboard *[][]tgbotapi.InlineKeyboardButton, now time.Time, currentDate time.Time, selectedDate time.Time) []tgbotapi.InlineKeyboardButton {
	var row []tgbotapi.InlineKeyboardButton

	year, month, _ := currentDate.Date()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < int(firstDay.Weekday()-1); i++ {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
	}

	for d := firstDay; d.Month() == month; d = d.AddDate(0, 0, 1) {
		dayText := fmt.Sprintf("%d", d.Day())

		if !selectedDate.IsZero() &&
			d.Year() == selectedDate.Year() &&
			d.Month() == selectedDate.Month() &&
			d.Day() == selectedDate.Day() {
			dayText = "✅ " + dayText
		} else if d.Year() == now.Year() &&
			d.Month() == now.Month() &&
			d.Day() == now.Day() {
			dayText = "🟢 " + dayText
		}

		callbackData := fmt.Sprintf("calendar:select:%s", d.Format(dateFormat))

		btn := tgbotapi.NewInlineKeyboardButtonData(dayText, callbackData)
		row = append(row, btn)

		if len(row) == 7 {
			*keyboard = append(*keyboard, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	return row
}
