package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boobachad/clankerloop/re-clanker/backend/internal/config"
	"github.com/boobachad/clankerloop/re-clanker/backend/internal/database"
	"github.com/boobachad/clankerloop/re-clanker/backend/internal/handler"
	"github.com/boobachad/clankerloop/re-clanker/backend/internal/middleware"
	"github.com/boobachad/clankerloop/re-clanker/backend/internal/repository"
	"github.com/boobachad/clankerloop/re-clanker/backend/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	ctx := context.Background()
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Initialize repositories
	problemRepo := repository.NewProblemRepository(db)
	modelRepo := repository.NewModelRepository(db)
	focusRepo := repository.NewFocusAreaRepository(db)
	jobRepo := repository.NewGenerationJobRepository(db)

	// Initialize AI service
	aiService, err := service.NewAIService(cfg.AIProvider, cfg.OpenRouterAPIKey, cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}
	log.Printf("AI service initialized with provider: %s", cfg.AIProvider)

	// Initialize services
	problemService := service.NewProblemService(problemRepo, jobRepo, aiService)

	// Initialize handlers
	problemHandler := handler.NewProblemHandler(problemRepo, focusRepo, jobRepo)
	modelHandler := handler.NewModelHandler(modelRepo)
	focusHandler := handler.NewFocusAreaHandler(focusRepo)
	
	_ = problemService // Will be used in future handlers

	// Setup router
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})

	// API routes
	mux.HandleFunc("GET /api/v1/models", modelHandler.ListModels)
	mux.HandleFunc("GET /api/v1/focus-areas", focusHandler.ListFocusAreas)
	mux.HandleFunc("POST /api/v1/problems", problemHandler.CreateProblem)
	mux.HandleFunc("GET /api/v1/problems", problemHandler.ListProblems)
	mux.HandleFunc("GET /api/v1/problems/{id}", problemHandler.GetProblem)
	mux.HandleFunc("GET /api/v1/problems/{id}/focus-areas", problemHandler.GetProblemFocusAreas)

	// Apply middleware
	handler := middleware.Logging(mux)
	handler = middleware.CORS(cfg.CORSOrigins)(handler)

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
