package service

import (
	"context"
	"fmt"

	"github.com/boobachad/clankerloop/re-clanker/backend/internal/repository"
	"github.com/google/uuid"
)

// ProblemService handles problem generation logic
type ProblemService struct {
	problemRepo *repository.ProblemRepository
	jobRepo     *repository.GenerationJobRepository
	aiService   *AIService
}

// NewProblemService creates a new problem service
func NewProblemService(
	problemRepo *repository.ProblemRepository,
	jobRepo *repository.GenerationJobRepository,
	aiService *AIService,
) *ProblemService {
	return &ProblemService{
		problemRepo: problemRepo,
		jobRepo:     jobRepo,
		aiService:   aiService,
	}
}

// GenerateProblemText generates problem text using AI
func (s *ProblemService) GenerateProblemText(ctx context.Context, problemID uuid.UUID, focusAreas []string, model string) error {
	// Build prompt based on focus areas
	prompt := "Generate a coding interview problem"
	if len(focusAreas) > 0 {
		prompt += fmt.Sprintf(" focusing on: %v", focusAreas)
	}
	prompt += ". Include a clear problem statement, input/output format, and examples."

	// Generate text using AI
	problemText, err := s.aiService.GenerateText(ctx, prompt, model)
	if err != nil {
		return fmt.Errorf("failed to generate problem text: %w", err)
	}

	// Update problem in database
	updates := map[string]interface{}{
		"problemText": problemText,
	}
	if err := s.problemRepo.Update(ctx, problemID, updates); err != nil {
		return fmt.Errorf("failed to update problem: %w", err)
	}

	return nil
}

// GenerateSolution generates a solution for the problem using AI
func (s *ProblemService) GenerateSolution(ctx context.Context, problemID uuid.UUID, model string) error {
	// Get problem
	problem, err := s.problemRepo.GetByID(ctx, problemID)
	if err != nil {
		return fmt.Errorf("failed to get problem: %w", err)
	}

	// Build prompt
	prompt := fmt.Sprintf("Generate a solution in Python for this problem:\n\n%s\n\nProvide only the code.", problem.ProblemText)

	// Generate solution using AI
	solution, err := s.aiService.GenerateText(ctx, prompt, model)
	if err != nil {
		return fmt.Errorf("failed to generate solution: %w", err)
	}

	// Update problem with solution
	updates := map[string]interface{}{
		"solution": solution,
	}
	if err := s.problemRepo.Update(ctx, problemID, updates); err != nil {
		return fmt.Errorf("failed to update problem with solution: %w", err)
	}

	return nil
}
