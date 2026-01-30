# AGENTS.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

Chirpy is a Twitter-like social media API built with Go. It features JWT-based authentication, user management, short message posting (chirps), and premium user upgrades. The application uses PostgreSQL for data persistence and follows a handler-based HTTP routing pattern.

## Essential Commands

### Build and Run
```bash
# Run the application
go run .

# Build the binary
go build -o Chirpy .
```

### Database Operations
```bash
# Apply database migrations (requires goose)
goose -dir sql/schema postgres $DB_URL up

# Rollback last migration
goose -dir sql/schema postgres $DB_URL down

# Regenerate database code after modifying SQL queries
sqlc generate
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/auth

# Run tests with verbose output
go test -v ./...

# Run a single test
go test -run TestCheckPasswordHash ./internal/auth
```

### Dependency Management
```bash
# Install/update dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Architecture

### Request Flow
1. HTTP requests arrive at `main.go`'s ServeMux router
2. Requests are routed to handler functions (handler_*.go files in root)
3. Handlers extract and validate JWT tokens via `internal/auth` package
4. Handlers interact with database through `internal/database.Queries` (sqlc-generated)
5. Responses are sent via `respondWithJSON()` or `respondWithError()` helper functions

### Authentication System
- **Password Security**: Bcrypt hashing via `internal/auth/hash_password.go`
- **Access Tokens**: JWT tokens with 1-hour expiry, validated in handlers
- **Refresh Tokens**: Stored in database, allow obtaining new access tokens without re-login
- **Token Flow**: Extract from Authorization header → Validate → Get user ID → Proceed with request

### Database Layer (sqlc Pattern)
- **Schema**: Located in `sql/schema/`, numbered migration files (001_users.sql, 002_chirps.sql, etc.)
- **Queries**: SQL definitions in `sql/queries/`, organized by entity (users.sql, chirps.sql, refresh_tokens.sql)
- **Generated Code**: `sqlc generate` produces type-safe Go functions in `internal/database/`
- **Usage**: Handlers call methods on `*database.Queries` (e.g., `cfg.db.CreateChirp()`)

### Code Organization
- **Root directory**: HTTP handlers (handler_*.go), main.go, helper functions (json.go, metrics.go)
- **internal/auth**: Authentication utilities (JWT, password hashing, token extraction)
- **internal/database**: sqlc-generated database code (DO NOT edit manually)
- **sql/schema**: Database migration files managed by goose
- **sql/queries**: SQL query definitions for sqlc to generate Go code

## Important Development Patterns

### Modifying Database Schema
1. Create new migration: `sql/schema/00X_description.sql`
2. Apply migration: `goose -dir sql/schema postgres $DB_URL up`
3. Add queries to appropriate file in `sql/queries/`
4. Run `sqlc generate` to regenerate Go code
5. Never manually edit files in `internal/database/` - they're auto-generated

### Adding New Endpoints
1. Create handler function in new or existing `handler_*.go` file
2. Follow pattern: extract params → validate auth (if needed) → call database → respond
3. Register route in `main.go`'s ServeMux
4. Use `respondWithJSON()` and `respondWithError()` for consistent responses

### Authentication Requirements
- Protected endpoints: Extract token with `auth.GetBearerToken(r.Header)`
- Validate with `auth.ValidateJWT(token, cfg.secret)` to get user ID
- User ID is then used for authorization checks (e.g., verifying chirp ownership)

## Environment Configuration

Required environment variables in `.env`:
- `DB_URL`: PostgreSQL connection string (format: `postgres://username:password@localhost:5432/chirpy?sslmode=disable`)
- `PLATFORM`: Environment identifier (typically "dev" or "prod")
- `SECRET`: JWT signing secret key
- `POLKA_KEY`: API key for Polka webhook integration

## Key Constraints

- Chirps are limited to 140 characters
- Profanity filter replaces: kerfuffle, sharbert, fornax with `****`
- JWT access tokens expire after 1 hour
- Users can only delete their own chirps
- Admin reset endpoint only works when `PLATFORM=dev`
