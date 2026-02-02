# ClankerLoop Frontend (Next.js)

This is the Next.js frontend for ClankerLoop, ported from the original `apps/web` with all authentication logic removed per project requirements.

## Features

- Next.js 16 with App Router
- TypeScript
- Tailwind CSS with design tokens (no hardcoded colors)
- shadcn/ui components
- Monaco Editor for code editing
- React Query for data fetching
- No authentication (public access)

## Prerequisites

- Node.js 18+ or Bun
- The Go backend running on port 8080

## Setup

1. **Install Dependencies**

   ```bash
   npm install
   ```

2. **Configure Environment**

   Copy `.env.example` to `.env.local`:

   ```bash
   cp .env.example .env.local
   ```

   Update the API URL if your backend is running on a different port:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8080
   ```

## Development

Run the development server:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Building

Build for production:

```bash
npm run build
```

Run the production build:

```bash
npm start
```

## Project Structure

```
src/
  app/
    layout.tsx           - Root layout (no AuthKitProvider)
    page.tsx             - Home page (shows new problem interface)
    problem/             - Problem pages
    api/                 - API route handlers (if any)
    globals.css          - Global styles with design tokens
  components/
    ui/                  - shadcn/ui components
    ai-elements/         - AI-specific components
    ...                  - Other application components
  lib/
    api-client.ts        - API client (no auth headers)
    utils.ts             - Utility functions
  hooks/                 - Custom React hooks
```

## Design Tokens

This project uses CSS custom properties (design tokens) for colors, spacing, and other design values. **Do not hardcode color values** - always use the design tokens defined in `globals.css`.

Example:
```css
/* ✅ Good - uses design token */
background-color: var(--background);
color: var(--foreground);

/* ❌ Bad - hardcoded color */
background-color: #ffffff;
color: #000000;
```

## Key Changes from Original

1. **No Authentication**:
   - Removed `@workos-inc/authkit-nextjs`
   - Removed `AuthKitProvider` from layout
   - Removed all `withAuth()` calls
   - Removed login/callback/signout routes
   - Removed API key handling
   - No user context or session

2. **API Configuration**:
   - Points to Go backend instead of Cloudflare Workers
   - Uses `NEXT_PUBLIC_API_URL` instead of `NEXT_PUBLIC_BACKEND_URL`
   - No X-API-Key headers in requests

3. **Simplified Pages**:
   - Home page always shows new problem interface
   - No user-specific redirects or data
   - No profile/settings pages

## API Endpoints Used

The frontend communicates with the Go backend at these endpoints:

- `GET /api/v1/models` - List AI models
- `GET /api/v1/focus-areas` - List focus areas
- `POST /api/v1/problems` - Create new problem
- `GET /api/v1/problems/:id` - Get problem details
- `GET /api/v1/problems/:id/test-cases` - Get test cases

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm start` - Start production server
- `npm run lint` - Run ESLint
