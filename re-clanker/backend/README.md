# ClankerLoop Backend (Go)

This is the Go-based backend for ClankerLoop, ported from the original Cloudflare Workers TypeScript backend. All authentication logic has been removed per project requirements.

## Features

- RESTful API for problem management
- PostgreSQL database with pgx driver
- Support for OpenRouter or Google Gemini AI APIs
- CORS-enabled for frontend integration
- No authentication (public API)

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- OpenRouter API key OR Google Gemini API key

## Setup

1. **Install Dependencies**

   ```bash
   go mod download
   ```

2. **Configure Environment**

   Copy `.env.example` to `.env` and fill in your values:

   ```bash
   cp .env.example .env
   ```

   Required environment variables:
   - `DATABASE_URL`: PostgreSQL connection string
   - `AI_PROVIDER`: Either "openrouter" or "gemini"
   - `OPENROUTER_API_KEY`: Your OpenRouter API key (if using OpenRouter)
   - `GEMINI_API_KEY`: Your Google Gemini API key (if using Gemini)

3. **Run Database Migrations**

   The database schema should already exist from the main clankerloop setup. If not, refer to `packages/db` for Drizzle migrations.

## Development

Run the server:

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080` by default.

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Models
- `GET /api/v1/models` - List all AI models

### Focus Areas
- `GET /api/v1/focus-areas` - List all focus areas

### Problems
- `POST /api/v1/problems` - Create a new problem
- `GET /api/v1/problems` - List all problems
- `GET /api/v1/problems/:id` - Get problem by ID
- `GET /api/v1/problems/:id/focus-areas` - Get focus areas for a problem

## Architecture

```
cmd/api/           - Application entry point
internal/
  config/          - Configuration management
  database/        - Database connection
  models/          - Data models
  repository/      - Database queries
  service/         - Business logic (AI integration)
  handler/         - HTTP request handlers
  middleware/      - HTTP middleware
```

## AI Provider Configuration

### Using OpenRouter

Set in `.env`:
```
AI_PROVIDER=openrouter
OPENROUTER_API_KEY=sk-or-v1-...
```

Default model: `anthropic/claude-3.5-sonnet`

### Using Google Gemini

Set in `.env`:
```
AI_PROVIDER=gemini
GEMINI_API_KEY=AIza...
```

Default model: `gemini-1.5-pro-latest`

## Building

Build for production:

```bash
go build -o bin/api cmd/api/main.go
```

Run the binary:

```bash
./bin/api
```

## Differences from Original Backend

1. **No Authentication**: All auth middleware removed, endpoints are public
2. **AI Providers**: Uses OpenRouter or Gemini instead of OpenAI
3. **Language**: Go instead of TypeScript/Cloudflare Workers
4. **Database**: Direct PostgreSQL with pgx instead of Drizzle ORM
5. **Framework**: Standard library http instead of Hono.js
