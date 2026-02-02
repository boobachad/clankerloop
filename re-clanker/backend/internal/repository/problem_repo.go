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

// ProblemRepository handles database operations for problems
type ProblemRepository struct {
	db *database.DB
}

// NewProblemRepository creates a new problem repository
func NewProblemRepository(db *database.DB) *ProblemRepository {
	return &ProblemRepository{db: db}
}

// Create creates a new problem
func (r *ProblemRepository) Create(ctx context.Context, problemText, functionSignature, problemTextReworded, generatedByUserID string) (uuid.UUID, error) {
	var id uuid.UUID
	query := `
		INSERT INTO problems (problem_text, function_signature, problem_text_reworded, generated_by_user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.Pool.QueryRow(ctx, query, problemText, functionSignature, problemTextReworded, generatedByUserID).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create problem: %w", err)
	}
	return id, nil
}

// GetByID retrieves a problem by ID with its test cases
func (r *ProblemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ProblemWithTestCases, error) {
	// Get problem
	var problem models.Problem
	var functionSignatureSchema []byte
	query := `
		SELECT id, problem_text, function_signature, function_signature_schema,
		       problem_text_reworded, solution, generated_by_model_id,
		       generated_by_user_id, easier_than, harder_than, created_at, updated_at
		FROM problems
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&problem.ID, &problem.ProblemText, &problem.FunctionSignature, &functionSignatureSchema,
		&problem.ProblemTextReworded, &problem.Solution, &problem.GeneratedByModelID,
		&problem.GeneratedByUserID, &problem.EasierThan, &problem.HarderThan,
		&problem.CreatedAt, &problem.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("problem not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	// Parse function signature schema if exists
	if functionSignatureSchema != nil {
		var schema map[string]interface{}
		if err := json.Unmarshal(functionSignatureSchema, &schema); err == nil {
			problem.FunctionSignatureSchema = schema
		}
	}

	// Get test cases
	testCases, err := r.getTestCases(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.ProblemWithTestCases{
		Problem:   problem,
		TestCases: testCases,
	}, nil
}

// Update updates a problem
func (r *ProblemRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	// Build dynamic update query
	query := "UPDATE problems SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 1

	if val, ok := updates["problemText"]; ok {
		query += fmt.Sprintf(", problem_text = $%d", argCount)
		args = append(args, val)
		argCount++
	}
	if val, ok := updates["functionSignature"]; ok {
		query += fmt.Sprintf(", function_signature = $%d", argCount)
		args = append(args, val)
		argCount++
	}
	if val, ok := updates["functionSignatureSchema"]; ok {
		jsonData, _ := json.Marshal(val)
		query += fmt.Sprintf(", function_signature_schema = $%d", argCount)
		args = append(args, jsonData)
		argCount++
	}
	if val, ok := updates["problemTextReworded"]; ok {
		query += fmt.Sprintf(", problem_text_reworded = $%d", argCount)
		args = append(args, val)
		argCount++
	}
	if val, ok := updates["solution"]; ok {
		query += fmt.Sprintf(", solution = $%d", argCount)
		args = append(args, val)
		argCount++
	}
	if val, ok := updates["generatedByModelId"]; ok {
		query += fmt.Sprintf(", generated_by_model_id = $%d", argCount)
		args = append(args, val)
		argCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, id)

	_, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update problem: %w", err)
	}
	return nil
}

// List lists all problem IDs
func (r *ProblemRepository) List(ctx context.Context) ([]uuid.UUID, error) {
	query := `SELECT id FROM problems ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list problems: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan problem id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// getTestCases retrieves test cases for a problem
func (r *ProblemRepository) getTestCases(ctx context.Context, problemID uuid.UUID) ([]models.TestCase, error) {
	query := `
		SELECT id, problem_id, description, is_edge_case, is_sample_case,
		       input_code, input, expected, created_at, updated_at
		FROM test_cases
		WHERE problem_id = $1
		ORDER BY created_at
	`
	rows, err := r.db.Pool.Query(ctx, query, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases: %w", err)
	}
	defer rows.Close()

	var testCases []models.TestCase
	for rows.Next() {
		var tc models.TestCase
		var inputJSON, expectedJSON []byte
		err := rows.Scan(
			&tc.ID, &tc.ProblemID, &tc.Description, &tc.IsEdgeCase, &tc.IsSampleCase,
			&tc.InputCode, &inputJSON, &expectedJSON, &tc.CreatedAt, &tc.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test case: %w", err)
		}

		// Parse JSON fields
		if inputJSON != nil {
			json.Unmarshal(inputJSON, &tc.Input)
		}
		if expectedJSON != nil {
			json.Unmarshal(expectedJSON, &tc.Expected)
		}

		testCases = append(testCases, tc)
	}
	return testCases, nil
}

// CreateTestCase creates a new test case
func (r *ProblemRepository) CreateTestCase(ctx context.Context, tc models.TestCase) (uuid.UUID, error) {
	var id uuid.UUID
	inputJSON, _ := json.Marshal(tc.Input)
	expectedJSON, _ := json.Marshal(tc.Expected)

	query := `
		INSERT INTO test_cases (problem_id, description, is_edge_case, is_sample_case, input_code, input, expected)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.db.Pool.QueryRow(ctx, query,
		tc.ProblemID, tc.Description, tc.IsEdgeCase, tc.IsSampleCase,
		tc.InputCode, inputJSON, expectedJSON,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create test case: %w", err)
	}
	return id, nil
}

// DeleteTestCases deletes all test cases for a problem
func (r *ProblemRepository) DeleteTestCases(ctx context.Context, problemID uuid.UUID) error {
	query := `DELETE FROM test_cases WHERE problem_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, problemID)
	if err != nil {
		return fmt.Errorf("failed to delete test cases: %w", err)
	}
	return nil
}

// GetMostRecentByUser gets the most recent problem for a user
func (r *ProblemRepository) GetMostRecentByUser(ctx context.Context, userID string) (*uuid.UUID, error) {
	var id uuid.UUID
	query := `SELECT id FROM problems WHERE generated_by_user_id = $1 ORDER BY created_at DESC LIMIT 1`
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&id)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get most recent problem: %w", err)
	}
	return &id, nil
}
