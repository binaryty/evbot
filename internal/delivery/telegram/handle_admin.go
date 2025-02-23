package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (h *Handler) isAdmin(userID int64) bool {
	for _, id := range h.cfg.AdminIDs {
		if id == userID {
			return true
		}
	}

	return false
}

func (h *Handler) AdminOnly(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, update *tgbotapi.Update) error {
		userID := GetUserIDFromUpdate(update)

		isAdmin := false
		for _, id := range h.cfg.AdminIDs {
			if id == userID {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			h.sendAccessDenied(update)
			log.Printf("Unauthorized access attempt by user %d", userID)
			return fmt.Errorf("access denied for user %d", userID)
		}

		return next(ctx, update)
	}
}

func (h *Handler) sendAccessDenied(update *tgbotapi.Update) {
	icon := "🚫"
	text := "У вас нет прав для выполнения этой команды"
	if update.CallbackQuery != nil {
		h.sendCallback(update.CallbackQuery.ID, icon, text)
	} else if update.Message != nil {
		h.sendMsg(update.Message.Chat.ID, icon, text)
	}
}
