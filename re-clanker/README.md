# Re-Clanker Project

This directory contains the re-scaffolded version of ClankerLoop with the following modifications:

## Key Changes

1. **No Authentication**: All authentication logic has been removed from both frontend and backend
2. **Go Backend**: Ported from TypeScript/Cloudflare Workers to Go
3. **AI Providers**: Uses OpenRouter or Google Gemini instead of OpenAI
4. **Design Tokens**: Frontend uses CSS custom properties (no hardcoded colors)

## Project Structure

```
re-clanker/
├── backend/           # Go-based REST API
│   ├── cmd/api/       # Application entry point
│   ├── internal/      # Internal packages
│   │   ├── config/    # Configuration
│   │   ├── database/  # Database connection
│   │   ├── models/    # Data models
│   │   ├── repository/# Database queries
│   │   ├── service/   # Business logic & AI
│   │   ├── handler/   # HTTP handlers
│   │   └── middleware/# HTTP middleware
│   ├── .env.example   # Environment template
│   └── README.md      # Backend documentation
│
└── frontend/          # Next.js 16 application
    ├── src/           # Source directory
    │   ├── app/       # Next.js app router
    │   ├── components/# React components (shadcn/ui)
    │   ├── lib/       # Utilities & API client
    │   └── hooks/     # Custom React hooks
    ├── public/        # Static assets
    ├── .env.example   # Environment template
    └── README.md      # Frontend documentation
```

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+ or Bun
- PostgreSQL database
- OpenRouter API key OR Google Gemini API key

### Backend Setup

```bash
cd backend

# Copy and configure environment
cp .env.example .env
# Edit .env with your database URL and AI provider credentials

# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

Backend runs on `http://localhost:8080`

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Copy and configure environment
cp .env.example .env.local
# Edit .env.local if needed (default points to localhost:8080)

# Run development server
npm run dev
```

Frontend runs on `http://localhost:3000`

## Architecture

### Backend (Go)

- **Framework**: Standard library `net/http`
- **Database**: PostgreSQL with `pgx` driver
- **AI Providers**: 
  - OpenRouter (default model: `anthropic/claude-3.5-sonnet`)
  - Google Gemini (default model: `gemini-1.5-pro-latest`)
- **Architecture**: Repository pattern with service layer
- **Authentication**: None (all endpoints public)

### Frontend (Next.js)

- **Framework**: Next.js 16 with App Router
- **Styling**: Tailwind CSS with design tokens
- **UI Components**: shadcn/ui (Radix UI primitives)
- **State Management**: Jotai
- **Data Fetching**: TanStack Query (React Query)
- **Code Editor**: Monaco Editor
- **Authentication**: None (public access)

## API Endpoints

### Backend REST API

- `GET /health` - Health check
- `GET /api/v1/models` - List AI models
- `GET /api/v1/focus-areas` - List focus areas
- `POST /api/v1/problems` - Create new problem
- `GET /api/v1/problems` - List all problems
- `GET /api/v1/problems/:id` - Get problem details
- `GET /api/v1/problems/:id/focus-areas` - Get problem focus areas

## Environment Variables

### Backend

```bash
DATABASE_URL=postgresql://user:pass@localhost:5432/clankerloop
AI_PROVIDER=openrouter  # or "gemini"
OPENROUTER_API_KEY=sk-or-v1-...
GEMINI_API_KEY=AIza...
PORT=8080
CORS_ORIGINS=http://localhost:3000
LOG_LEVEL=info
```

### Frontend

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Design Principles

### Backend

- **No Global State**: All state managed in database or request context
- **Repository Pattern**: Clean separation of data access
- **Service Layer**: Business logic isolated from HTTP layer
- **Error Handling**: Go idiomatic error handling
- **Performance**: Connection pooling, efficient queries

### Frontend

- **Design Tokens**: No hardcoded colors/values
- **Type Safety**: Full TypeScript coverage
- **Component Reusability**: shadcn/ui pattern
- **Performance**: React Query caching, code splitting
- **Accessibility**: Radix UI accessible primitives

## Differences from Original

### Removed Features

1. **Authentication**:
   - No WorkOS AuthKit
   - No API key encryption/decryption
   - No user sessions
   - No protected routes

2. **Analytics**:
   - No Vercel Analytics
   - No PostHog tracking (can be re-added if needed)

3. **Database ORM**:
   - No Drizzle ORM (replaced with pgx)

### Modified Features

1. **AI Integration**:
   - OpenAI → OpenRouter or Gemini
   - Configurable AI provider

2. **Backend Platform**:
   - Cloudflare Workers → Standard Go HTTP server
   - TypeScript → Go

3. **Data Access**:
   - Drizzle queries → Go repository pattern with SQL

## Database Schema

The database schema remains the same as the original ClankerLoop:

- `models` - AI models
- `problems` - Coding problems
- `test_cases` - Problem test cases
- `generation_jobs` - Problem generation jobs
- `focus_areas` - Problem focus areas
- `problem_focus_areas` - Many-to-many relationship
- `user_problem_attempts` - User submissions (no auth validation)

## Development Workflow

1. Start PostgreSQL database
2. Run backend: `cd backend && go run cmd/api/main.go`
3. Run frontend: `cd frontend && npm run dev`
4. Access application at `http://localhost:3000`

## Production Deployment

### Backend

```bash
cd backend
go build -o bin/api cmd/api/main.go
./bin/api
```

### Frontend

```bash
cd frontend
npm run build
npm start
```

## Contributing

Follow the existing code style and architecture patterns:
- Go: Standard library patterns, repository pattern
- Frontend: shadcn/ui patterns, design tokens
- No authentication logic
- Use AI providers (OpenRouter/Gemini)

## License

Same as original ClankerLoop project.
