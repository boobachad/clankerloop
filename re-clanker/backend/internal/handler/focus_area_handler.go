package handler

import (
	"net/http"

	"github.com/boobachad/clankerloop/re-clanker/backend/internal/repository"
)

// FocusAreaHandler handles focus area-related HTTP requests
type FocusAreaHandler struct {
	focusRepo *repository.FocusAreaRepository
}

// NewFocusAreaHandler creates a new focus area handler
func NewFocusAreaHandler(focusRepo *repository.FocusAreaRepository) *FocusAreaHandler {
	return &FocusAreaHandler{focusRepo: focusRepo}
}

// ListFocusAreas handles GET /api/v1/focus-areas
func (h *FocusAreaHandler) ListFocusAreas(w http.ResponseWriter, r *http.Request) {
	focusAreas, err := h.focusRepo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list focus areas")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"focusAreas": focusAreas,
	})
}
