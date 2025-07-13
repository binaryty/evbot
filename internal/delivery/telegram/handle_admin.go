package telegram

import "slices"

// isAdmin ...
func (h *Handler) isAdmin(userID int64) bool {
	return slices.Contains(h.cfg.AdminIDs, userID)
}
