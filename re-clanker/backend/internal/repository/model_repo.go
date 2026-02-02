package repository

import (
	"context"
	"fmt"

	"github.com/boobachad/clankerloop/re-clanker/backend/internal/database"
	"github.com/boobachad/clankerloop/re-clanker/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ModelRepository handles database operations for models
type ModelRepository struct {
	db *database.DB
}

// NewModelRepository creates a new model repository
func NewModelRepository(db *database.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

// Create creates a new model
func (r *ModelRepository) Create(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	query := `INSERT INTO models (name) VALUES ($1) RETURNING id`
	err := r.db.Pool.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create model: %w", err)
	}
	return id, nil
}

// GetByID retrieves a model by ID
func (r *ModelRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Model, error) {
	var model models.Model
	query := `SELECT id, name FROM models WHERE id = $1`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&model.ID, &model.Name)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}
	return &model, nil
}

// GetByName retrieves a model by name
func (r *ModelRepository) GetByName(ctx context.Context, name string) (*models.Model, error) {
	var model models.Model
	query := `SELECT id, name FROM models WHERE name = $1`
	err := r.db.Pool.QueryRow(ctx, query, name).Scan(&model.ID, &model.Name)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get model by name: %w", err)
	}
	return &model, nil
}

// List lists all models
func (r *ModelRepository) List(ctx context.Context) ([]models.Model, error) {
	query := `SELECT id, name FROM models ORDER BY name`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer rows.Close()

	var modelsList []models.Model
	for rows.Next() {
		var model models.Model
		if err := rows.Scan(&model.ID, &model.Name); err != nil {
			return nil, fmt.Errorf("failed to scan model: %w", err)
		}
		modelsList = append(modelsList, model)
	}
	return modelsList, nil
}
