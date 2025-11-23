Fullstack example

Backend (GO)

Frontend (TypeScript)

# may — Fullstack Example

A small fullstack example showing a Go backend and a React + TypeScript frontend.

## Features

- Backend: Go, Gin (HTTP router), GORM (ORM)
- Frontend: React + TypeScript (Vite)

## Technologies

- Go
- Gin
- GORM
- React
- TypeScript
- Vite
- Bun / npm (frontend tooling)

## Prerequisites

- Go (1.20+ recommended)
- Node.js / Bun (for frontend) — either `npm`/`pnpm`/`yarn` works; this repo uses `bun` in examples
- PostgreSQL (or compatible DB) if you want to run the backend against a real database

## Quickstart

Run backend and frontend in separate terminals.

### Backend

1. Open a terminal and change to the backend directory:

	`cd backend`

2. Set your database DSN as an environment variable (example):

	Windows (PowerShell):

	`$env:POSTGRES_DSN = "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"`

3. Run the backend:

	`go run .`

### Frontend

1. Open a second terminal and change to the frontend directory:

	`cd frontend`

2. Install dependencies (if needed). With Bun:

	`bun install`

	Or with npm:

	`npm install`

3. Start the dev server:

	`bun dev`

The frontend dev server (Vite) typically runs on `http://localhost:5173` and the backend on the port configured in the Go server.

## Environment variables

- `POSTGRES_DSN` — PostgreSQL DSN used by the backend (recommended). Example:

- `POSTGRES_DSN` or `DATABASE_URL` — PostgreSQL DSN used by the backend (recommended). Example connection string values:

	`host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable`
	or
	`postgresql://user:password@localhost:5432/mydb?sslmode=disable`

Tips:

- Do not commit credentials. Use environment variables or a secrets manager.
- For local development you can use a `.env` file and a loader (or set env vars in your shell).

## Project structure

- `backend/` — Go backend source (Gin + GORM)
- `frontend/` — React + TypeScript app (Vite)

## Development notes

- If you want a package-level `*gorm.DB` instance, initialize it inside `main()` or an `init()` function — avoid assignment statements at package scope.
- To gracefully close DB connections: `sqlDB, _ := db.DB(); defer sqlDB.Close()`.

## Contributing

- Feel free to open issues or PRs. Include steps to reproduce and any logs.

## License

This repo does not include a license file. Add one if you plan to open-source the project.
4. `bun dev`
