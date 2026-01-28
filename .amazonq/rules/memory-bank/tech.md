# Technology Stack

## Programming Language
- **Go 1.24.2**: Modern, statically-typed language with excellent concurrency support and performance

## Core Framework & Libraries

### Web Framework
- **Gin v1.11.0**: High-performance HTTP web framework with middleware support
- **gin-contrib/sessions v1.0.4**: Session management middleware for Gin

### Database & ORM
- **GORM v1.31.1**: Feature-rich ORM for Go with associations, hooks, and migrations
- **gorm.io/driver/mysql v1.6.0**: MySQL driver for GORM
- **gorm.io/datatypes v1.2.7**: Additional data types for GORM (JSON support)

### Security & Authentication
- **golang.org/x/crypto v0.46.0**: Cryptographic libraries including bcrypt for password hashing
- **gorilla/sessions v1.4.0**: Session management (used via gin-contrib/sessions)
- **gorilla/securecookie v1.1.2**: Secure cookie encoding/decoding

### Validation
- **go-playground/validator/v10 v10.27.0**: Struct and field validation with tag support

### Logging
- **rs/zerolog v1.34.0**: Zero-allocation JSON logger for structured logging

### Configuration
- **joho/godotenv v1.5.1**: Environment variable loading from .env files

### Utilities
- **dario.cat/mergo v1.0.2**: Library for merging structs and maps
- **google/uuid v1.6.0**: UUID generation and parsing

## Database
- **MySQL**: Primary relational database for data persistence
- Connection pooling and transaction support via GORM

## Development Tools

### Dependency Management
- **Go Modules**: Native Go dependency management
- `go.mod`: Module definition and direct dependencies
- `go.sum`: Dependency checksums for verification

### Environment Configuration
- **.env file**: Local environment variables
- Required variables:
  - `SERVER_PORT`: HTTP server port
  - Database connection parameters
  - Email service configuration

## Key Development Commands

### Running the Application
```bash
# Run the server
go run cmd/server/main.go

# Build the application
go build -o bin/server cmd/server/main.go

# Run the built binary
./bin/server
```

### Dependency Management
```bash
# Download dependencies
go mod download

# Add a new dependency
go get github.com/package/name

# Update dependencies
go get -u ./...

# Tidy up dependencies
go mod tidy

# Verify dependencies
go mod verify
```

### Database Operations
- Database migrations run automatically on application startup via `database.DB.AutoMigrate()`
- Models migrated: BusinessUnit, Division, Role, ModulePermissions, User, AuditLog, PasswordResetToken

### Testing
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## Project Configuration

### Timezone
- Application timezone set to **America/New_York** (EST/EDT)
- Configured in main.go init function

### Session Management
- Cookie-based sessions via gin-contrib/sessions
- Secure cookie storage with gorilla/securecookie

### Logging
- Structured JSON logging via zerolog
- Configured in shared/logger.go
- Log levels: Debug, Info, Warn, Error, Fatal

## Architecture Characteristics

### Performance
- Gin framework provides high-performance HTTP routing
- GORM connection pooling for efficient database access
- Zero-allocation logging with zerolog

### Security
- Bcrypt password hashing with configurable cost
- Session-based authentication with secure cookies
- Role-based access control middleware
- SQL injection protection via GORM parameterized queries

### Scalability
- Stateless service layer enables horizontal scaling
- Database connection pooling
- Middleware-based request processing pipeline

### Maintainability
- Clean layered architecture
- Dependency injection patterns
- Structured logging for debugging
- Validation tags for data integrity
