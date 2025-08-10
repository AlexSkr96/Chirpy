# Chirpy 🐦

A Twitter-like social media API built with Go, featuring user authentication, JWT tokens, and short message posting functionality.

## Features

- **User Management**
  - User registration and authentication
  - Password hashing with bcrypt
  - JWT-based authentication
  - Refresh token functionality
  - Premium user upgrades (Chirpy Red)

- **Chirps (Messages)**
  - Create short messages (max 140 characters)
  - Retrieve chirps with sorting options (ascending/descending by date)
  - Filter chirps by author
  - Delete chirps (users can only delete their own)
  - Profanity filtering

- **Admin Features**
  - Metrics tracking for API usage
  - Admin panel with hit counter
  - Database reset functionality (dev environment)

## Tech Stack

- **Backend**: Go 1.24.3
- **Database**: PostgreSQL
- **Authentication**: JWT tokens with refresh token mechanism
- **Password Security**: bcrypt hashing
- **Database Migrations**: SQL migrations with goose
- **Code Generation**: sqlc for type-safe database queries

## Project Structure

```
Chirpy/
├── main.go                    # Application entry point
├── handler_*.go               # HTTP request handlers
├── internal/
│   ├── auth/                  # Authentication utilities
│   │   ├── jwt.go            # JWT token management
│   │   ├── hash_password.go  # Password hashing
│   │   ├── refresh_token.go  # Refresh token logic
│   │   └── token_extractor.go # Token extraction from headers
│   └── database/             # Database layer
│       ├── db.go             # Database connection
│       ├── models.go         # Generated database models
│       └── *.sql.go          # Generated query functions
├── sql/
│   ├── schema/               # Database schema migrations
│   └── queries/              # SQL query definitions
├── admin/                    # Admin interface
└── assets/                   # Static assets
```

## API Endpoints

### Authentication
- `POST /api/users` - Create new user account
- `PUT /api/users` - Update user information
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh JWT token
- `POST /api/revoke` - Revoke refresh token

### Chirps
- `POST /api/chirps` - Create a new chirp (authenticated)
- `GET /api/chirps` - Get all chirps with optional sorting and filtering
- `GET /api/chirps/{chirpID}` - Get specific chirp by ID
- `DELETE /api/chirps/{chirpID}` - Delete chirp (author only)

### Premium Features
- `POST /api/polka/webhooks` - Upgrade user to Chirpy Red (webhook)

### Admin
- `GET /admin/metrics` - View API usage metrics
- `POST /admin/reset` - Reset database (development only)

### System
- `GET /api/healthz` - Health check endpoint

## Getting Started

### Prerequisites

- Go 1.24.3 or higher
- PostgreSQL database
- Environment variables configured

### Environment Variables

Create a `.env` file in the root directory:

```env
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
SECRET=your-jwt-secret-key
POLKA_KEY=your-polka-api-key
```

### Installation

1. Clone the repository:
```bash
git clone https://github.com/AlexSkr96/Chirpy.git
cd Chirpy
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up the database:
```bash
# Run migrations to create the database schema
# (Assuming you have goose installed)
goose -dir sql/schema postgres $DB_URL up
```

4. Generate database code:
```bash
# If you have sqlc installed
sqlc generate
```

5. Run the application:
```bash
go run .
```

The server will start on port 8080. Visit `http://localhost:8080/app/` to see the welcome page.

## Usage Examples

### Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'
```

### Create a Chirp
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"body": "Hello, Chirpy world!"}'
```

### Get All Chirps
```bash
# Get chirps sorted by date (ascending - default)
curl http://localhost:8080/api/chirps

# Get chirps sorted by date (descending)
curl "http://localhost:8080/api/chirps?sort=desc"

# Get chirps by specific author
curl "http://localhost:8080/api/chirps?author_id=USER_UUID"
```

## Features in Detail

### Profanity Filter
Chirpy automatically filters out inappropriate words and replaces them with `****`. The filtered words include:
- kerfuffle
- sharbert  
- fornax

### Authentication Flow
1. User registers with email and password
2. Password is hashed using bcrypt
3. User logs in and receives JWT access token (1 hour expiry) and refresh token
4. Access token is used for authenticated requests
5. Refresh token can be used to get new access tokens
6. Tokens can be revoked for security

### Premium Features
Users can be upgraded to "Chirpy Red" premium status through webhook integration with external payment systems.

## Development

### Database Migrations
Add new migrations to `sql/schema/` directory and run:
```bash
goose -dir sql/schema postgres $DB_URL up
```

### Regenerate Database Code
After modifying SQL queries in `sql/queries/`, regenerate the Go code:
```bash
sqlc generate
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is part of a learning exercise and is available for educational purposes.

---

**Chirpy** - Tweet-sized messages, Go-sized performance! 🚀
