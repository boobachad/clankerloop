package repository

import (
	"context"
	"fmt"

	"github.com/boobachad/clankerloop/re-clanker/backend/internal/database"
	"github.com/boobachad/clankerloop/re-clanker/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// FocusAreaRepository handles database operations for focus areas
type FocusAreaRepository struct {
	db *database.DB
}

// NewFocusAreaRepository creates a new focus area repository
func NewFocusAreaRepository(db *database.DB) *FocusAreaRepository {
	return &FocusAreaRepository{db: db}
}

// List lists all active focus areas
func (r *FocusAreaRepository) List(ctx context.Context) ([]models.FocusArea, error) {
	query := `
		SELECT id, name, slug, description, prompt_guidance, display_order, is_active, created_at, updated_at
		FROM focus_areas
		WHERE is_active = true
		ORDER BY display_order
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list focus areas: %w", err)
	}
	defer rows.Close()

	var focusAreas []models.FocusArea
	for rows.Next() {
		var fa models.FocusArea
		if err := rows.Scan(&fa.ID, &fa.Name, &fa.Slug, &fa.Description, &fa.PromptGuidance, &fa.DisplayOrder, &fa.IsActive, &fa.CreatedAt, &fa.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan focus area: %w", err)
		}
		focusAreas = append(focusAreas, fa)
	}
	return focusAreas, nil
}

// GetByIDs retrieves focus areas by IDs
func (r *FocusAreaRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.FocusArea, error) {
	if len(ids) == 0 {
		return []models.FocusArea{}, nil
	}

	query := `
		SELECT id, name, slug, description, prompt_guidance, display_order, is_active, created_at, updated_at
		FROM focus_areas
		WHERE id = ANY($1)
	`
	rows, err := r.db.Pool.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get focus areas by IDs: %w", err)
	}
	defer rows.Close()

	var focusAreas []models.FocusArea
	for rows.Next() {
		var fa models.FocusArea
		if err := rows.Scan(&fa.ID, &fa.Name, &fa.Slug, &fa.Description, &fa.PromptGuidance, &fa.DisplayOrder, &fa.IsActive, &fa.CreatedAt, &fa.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan focus area: %w", err)
		}
		focusAreas = append(focusAreas, fa)
	}
	return focusAreas, nil
}

// GetForProblem retrieves focus areas for a problem
func (r *FocusAreaRepository) GetForProblem(ctx context.Context, problemID uuid.UUID) ([]models.FocusArea, error) {
	query := `
		SELECT fa.id, fa.name, fa.slug, fa.description, fa.prompt_guidance, fa.display_order, fa.is_active, fa.created_at, fa.updated_at
		FROM focus_areas fa
		INNER JOIN problem_focus_areas pfa ON pfa.focus_area_id = fa.id
		WHERE pfa.problem_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, query, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get focus areas for problem: %w", err)
	}
	defer rows.Close()

	var focusAreas []models.FocusArea
	for rows.Next() {
		var fa models.FocusArea
		if err := rows.Scan(&fa.ID, &fa.Name, &fa.Slug, &fa.Description, &fa.PromptGuidance, &fa.DisplayOrder, &fa.IsActive, &fa.CreatedAt, &fa.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan focus area: %w", err)
		}
		focusAreas = append(focusAreas, fa)
	}
	return focusAreas, nil
}

// LinkToProblem links focus areas to a problem
func (r *FocusAreaRepository) LinkToProblem(ctx context.Context, problemID uuid.UUID, focusAreaIDs []uuid.UUID) error {
	if len(focusAreaIDs) == 0 {
		return nil
	}

	query := `INSERT INTO problem_focus_areas (problem_id, focus_area_id) VALUES ($1, $2)`
	batch := &pgx.Batch{}
	for _, faID := range focusAreaIDs {
		batch.Queue(query, problemID, faID)
	}

	br := r.db.Pool.SendBatch(ctx, batch)
	defer br.Close()

	for range focusAreaIDs {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("failed to link focus area to problem: %w", err)
		}
	}
	return nil
}
