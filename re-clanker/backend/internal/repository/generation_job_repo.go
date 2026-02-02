package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/boobachad/clankerloop/re-clanker/backend/internal/database"
	"github.com/boobachad/clankerloop/re-clanker/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GenerationJobRepository handles database operations for generation jobs
type GenerationJobRepository struct {
	db *database.DB
}

// NewGenerationJobRepository creates a new generation job repository
func NewGenerationJobRepository(db *database.DB) *GenerationJobRepository {
	return &GenerationJobRepository{db: db}
}

// Create creates a new generation job
func (r *GenerationJobRepository) Create(ctx context.Context, problemID uuid.UUID, modelID *uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	query := `
		INSERT INTO generation_jobs (problem_id, model_id, status, completed_steps)
		VALUES ($1, $2, 'pending', '[]'::jsonb)
		RETURNING id
	`
	err := r.db.Pool.QueryRow(ctx, query, problemID, modelID).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create generation job: %w", err)
	}
	return id, nil
}

// GetByID retrieves a generation job by ID
func (r *GenerationJobRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.GenerationJob, error) {
	var job models.GenerationJob
	var completedStepsJSON []byte
	query := `
		SELECT id, problem_id, model_id, status, current_step, completed_steps, error, created_at, updated_at
		FROM generation_jobs
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&job.ID, &job.ProblemID, &job.ModelID, &job.Status, &job.CurrentStep,
		&completedStepsJSON, &job.Error, &job.CreatedAt, &job.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get generation job: %w", err)
	}

	// Parse completed steps
	if completedStepsJSON != nil {
		json.Unmarshal(completedStepsJSON, &job.CompletedSteps)
	}
	if job.CompletedSteps == nil {
		job.CompletedSteps = []string{}
	}

	return &job, nil
}

// GetLatestForProblem retrieves the latest generation job for a problem
func (r *GenerationJobRepository) GetLatestForProblem(ctx context.Context, problemID uuid.UUID) (*models.GenerationJob, error) {
	var job models.GenerationJob
	var completedStepsJSON []byte
	query := `
		SELECT id, problem_id, model_id, status, current_step, completed_steps, error, created_at, updated_at
		FROM generation_jobs
		WHERE problem_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := r.db.Pool.QueryRow(ctx, query, problemID).Scan(
		&job.ID, &job.ProblemID, &job.ModelID, &job.Status, &job.CurrentStep,
		&completedStepsJSON, &job.Error, &job.CreatedAt, &job.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest generation job: %w", err)
	}

	// Parse completed steps
	if completedStepsJSON != nil {
		json.Unmarshal(completedStepsJSON, &job.CompletedSteps)
	}
	if job.CompletedSteps == nil {
		job.CompletedSteps = []string{}
	}

	return &job, nil
}

// UpdateStatus updates a generation job's status
func (r *GenerationJobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status, currentStep string, errorMsg *string) error {
	query := `
		UPDATE generation_jobs
		SET status = $1, current_step = $2, error = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.Pool.Exec(ctx, query, status, currentStep, errorMsg, id)
	if err != nil {
		return fmt.Errorf("failed to update generation job status: %w", err)
	}
	return nil
}

// MarkStepComplete marks a step as completed
func (r *GenerationJobRepository) MarkStepComplete(ctx context.Context, id uuid.UUID, step string) error {
	// Get current job
	job, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if job == nil {
		return fmt.Errorf("generation job not found: %s", id)
	}

	// Add step to completed steps
	completedSteps := append(job.CompletedSteps, step)
	completedStepsJSON, _ := json.Marshal(completedSteps)

	query := `
		UPDATE generation_jobs
		SET completed_steps = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = r.db.Pool.Exec(ctx, query, completedStepsJSON, id)
	if err != nil {
		return fmt.Errorf("failed to mark step complete: %w", err)
	}
	return nil
}
