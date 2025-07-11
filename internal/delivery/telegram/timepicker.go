package timepicker

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/binaryty/evbot/internal/delivery/telegram"
	domain "github.com/binaryty/evbot/internal/domain/entities"
)

func GenerateTimePicker(tp *domain.TimePicker) tgbotapi.InlineKeyboardMarkup {
	var timePicker [][]tgbotapi.InlineKeyboardButton

	if tp.Step == "hours" {
		timePicker = generateHours()
	} else {
		timePicker = generateMinutes()
	}

	timePicker = append(timePicker, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", telegram.EmOk, "Готово"),
			"time_confirm",
		),
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", telegram.EmCross, "Отмена"),
			"time_cancel",
		),
	})

	t := tgbotapi.NewInlineKeyboardMarkup(timePicker...)

	return t
}

func generateHours() [][]tgbotapi.InlineKeyboardButton {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for h := 0; h < 24; h++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%02d", h),
			fmt.Sprintf("time_h:%02d", h),
		)
		row = append(row, btn)

		if len(row) == 4 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	return rows
}

func generateMinutes() [][]tgbotapi.InlineKeyboardButton {
	minutes := []int{0, 15, 30, 45}
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for _, m := range minutes {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%02d", m),
			fmt.Sprintf("time:%02d", m),
		)
		row = append(row, btn)
	}
	rows = append(rows, row)

	return rows
}
