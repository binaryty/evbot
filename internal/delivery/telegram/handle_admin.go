package telegram

// isAdmin ...
func (h *Handler) isAdmin(userID int64) bool {
	for _, id := range h.cfg.AdminIDs {
		if id == userID {
			return true
		}
	}

	return false
}
