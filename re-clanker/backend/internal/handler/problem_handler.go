package handler

import (
	"encoding/json"
	"net/http"

	"github.com/boobachad/clankerloop/re-clanker/backend/internal/repository"
	"github.com/google/uuid"
)

// ProblemHandler handles problem-related HTTP requests
type ProblemHandler struct {
	problemRepo *repository.ProblemRepository
	focusRepo   *repository.FocusAreaRepository
	jobRepo     *repository.GenerationJobRepository
}

// NewProblemHandler creates a new problem handler
func NewProblemHandler(
	problemRepo *repository.ProblemRepository,
	focusRepo *repository.FocusAreaRepository,
	jobRepo *repository.GenerationJobRepository,
) *ProblemHandler {
	return &ProblemHandler{
		problemRepo: problemRepo,
		focusRepo:   focusRepo,
		jobRepo:     jobRepo,
	}
}

// CreateProblem handles POST /api/v1/problems
func (h *ProblemHandler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FocusAreaIDs []string `json:"focusAreaIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create problem with default user ID (no auth)
	problemID, err := h.problemRepo.Create(r.Context(), "", "", "", "default-user")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create problem")
		return
	}

	// Link focus areas if provided
	if len(req.FocusAreaIDs) > 0 {
		focusAreaUUIDs := make([]uuid.UUID, 0, len(req.FocusAreaIDs))
		for _, idStr := range req.FocusAreaIDs {
			id, err := uuid.Parse(idStr)
			if err == nil {
				focusAreaUUIDs = append(focusAreaUUIDs, id)
			}
		}
		if len(focusAreaUUIDs) > 0 {
			h.focusRepo.LinkToProblem(r.Context(), problemID, focusAreaUUIDs)
		}
	}

	// Create generation job
	jobID, err := h.jobRepo.Create(r.Context(), problemID, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create generation job")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success":   true,
		"problemId": problemID,
		"jobId":     jobID,
	})
}

// GetProblem handles GET /api/v1/problems/:id
func (h *ProblemHandler) GetProblem(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL (will be handled by router)
	idStr := r.PathValue("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "Missing problem ID")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid problem ID")
		return
	}

	problem, err := h.problemRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Problem not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"problem": problem,
	})
}

// ListProblems handles GET /api/v1/problems
func (h *ProblemHandler) ListProblems(w http.ResponseWriter, r *http.Request) {
	ids, err := h.problemRepo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list problems")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"problems": ids,
	})
}

// GetProblemFocusAreas handles GET /api/v1/problems/:id/focus-areas
func (h *ProblemHandler) GetProblemFocusAreas(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "Missing problem ID")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid problem ID")
		return
	}

	focusAreas, err := h.focusRepo.GetForProblem(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get focus areas")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"focusAreas": focusAreas,
	})
}

// Helper functions
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false,
		"error": map[string]string{
			"message": message,
		},
	})
}
