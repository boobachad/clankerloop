package models

import (
	"time"

	"github.com/google/uuid"
)

// Model represents an AI model in the system
type Model struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

// Problem represents a coding problem
type Problem struct {
	ID                      uuid.UUID              `json:"id" db:"id"`
	ProblemText             string                 `json:"problemText" db:"problem_text"`
	FunctionSignature       string                 `json:"functionSignature" db:"function_signature"`
	FunctionSignatureSchema map[string]interface{} `json:"functionSignatureSchema,omitempty" db:"function_signature_schema"`
	ProblemTextReworded     string                 `json:"problemTextReworded" db:"problem_text_reworded"`
	Solution                *string                `json:"solution,omitempty" db:"solution"`
	GeneratedByModelID      *uuid.UUID             `json:"generatedByModelId,omitempty" db:"generated_by_model_id"`
	GeneratedByUserID       string                 `json:"generatedByUserId" db:"generated_by_user_id"`
	EasierThan              *uuid.UUID             `json:"easierThan,omitempty" db:"easier_than"`
	HarderThan              *uuid.UUID             `json:"harderThan,omitempty" db:"harder_than"`
	CreatedAt               time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt               time.Time              `json:"updatedAt" db:"updated_at"`
}

// TestCase represents a test case for a problem
type TestCase struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	ProblemID    uuid.UUID              `json:"problemId" db:"problem_id"`
	Description  string                 `json:"description" db:"description"`
	IsEdgeCase   bool                   `json:"isEdgeCase" db:"is_edge_case"`
	IsSampleCase bool                   `json:"isSampleCase" db:"is_sample_case"`
	InputCode    *string                `json:"inputCode,omitempty" db:"input_code"`
	Input        map[string]interface{} `json:"input,omitempty" db:"input"`
	Expected     map[string]interface{} `json:"expected,omitempty" db:"expected"`
	CreatedAt    time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time              `json:"updatedAt" db:"updated_at"`
}

// GenerationJob represents a problem generation job
type GenerationJob struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	ProblemID      uuid.UUID  `json:"problemId" db:"problem_id"`
	ModelID        *uuid.UUID `json:"modelId,omitempty" db:"model_id"`
	Status         string     `json:"status" db:"status"` // pending, in_progress, completed, failed
	CurrentStep    *string    `json:"currentStep,omitempty" db:"current_step"`
	CompletedSteps []string   `json:"completedSteps" db:"completed_steps"`
	Error          *string    `json:"error,omitempty" db:"error"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

// FocusArea represents a problem focus area
type FocusArea struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Slug           string    `json:"slug" db:"slug"`
	Description    *string   `json:"description,omitempty" db:"description"`
	PromptGuidance string    `json:"promptGuidance" db:"prompt_guidance"`
	DisplayOrder   int       `json:"displayOrder" db:"display_order"`
	IsActive       bool      `json:"isActive" db:"is_active"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

// ProblemFocusArea represents a many-to-many relationship
type ProblemFocusArea struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ProblemID   uuid.UUID `json:"problemId" db:"problem_id"`
	FocusAreaID uuid.UUID `json:"focusAreaId" db:"focus_area_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

// UserProblemAttempt represents a user's attempt at solving a problem
type UserProblemAttempt struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             string    `json:"userId" db:"user_id"`
	ProblemID          uuid.UUID `json:"problemId" db:"problem_id"`
	SubmissionCode     string    `json:"submissionCode" db:"submission_code"`
	SubmissionLanguage string    `json:"submissionLanguage" db:"submission_language"`
	Status             string    `json:"status" db:"status"` // attempt, run, pass
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time `json:"updatedAt" db:"updated_at"`
}

// ProblemWithTestCases is a problem with its test cases
type ProblemWithTestCases struct {
	Problem
	TestCases []TestCase `json:"testCases"`
}
