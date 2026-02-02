package handler

import (
	"net/http"

	"github.com/boobachad/clankerloop/re-clanker/backend/internal/repository"
)

// ModelHandler handles model-related HTTP requests
type ModelHandler struct {
	modelRepo *repository.ModelRepository
}

// NewModelHandler creates a new model handler
func NewModelHandler(modelRepo *repository.ModelRepository) *ModelHandler {
	return &ModelHandler{modelRepo: modelRepo}
}

// ListModels handles GET /api/v1/models
func (h *ModelHandler) ListModels(w http.ResponseWriter, r *http.Request) {
	models, err := h.modelRepo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list models")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"models":  models,
	})
}
