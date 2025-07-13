package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

// GenerateTimePicker ...
func GenerateTimePicker(tp *domain.TimePicker, userID int64) tgbotapi.InlineKeyboardMarkup {
	var timePicker [][]tgbotapi.InlineKeyboardButton

	switch tp.Step {
	case domain.StepHours:
		timePicker = generateHours(userID)
	case domain.StepMinutes:
		timePicker = generateMinutes(userID)
	}

	timePicker = append(timePicker, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", EmOk, "Готово"),
			fmt.Sprintf("time_picker:time_confirm:%d", userID),
		),
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", EmCross, "Отмена"),
			fmt.Sprintf("time_picker:time_cancel:%d", userID),
		),
	})

	return tgbotapi.NewInlineKeyboardMarkup(timePicker...)
}

// generateHours ...
func generateHours(userID int64) [][]tgbotapi.InlineKeyboardButton {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for h := 0; h < 24; h++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%02d", h),
			fmt.Sprintf("time_picker:time_h_select:%02d:%d", h, userID),
		)
		row = append(row, btn)

		if len(row) == 4 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	return rows
}

// generateMinutes ...
func generateMinutes(userID int64) [][]tgbotapi.InlineKeyboardButton {
	minutes := []int{0, 15, 30, 45}
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for _, m := range minutes {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%02d", m),
			fmt.Sprintf("time_picker:time_m_select:%02d:%d", m, userID),
		)
		row = append(row, btn)
	}
	rows = append(rows, row)

	return rows
}
