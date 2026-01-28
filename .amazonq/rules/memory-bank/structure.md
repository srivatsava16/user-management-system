# Project Structure

## Directory Organization

```
user-management-system/
├── cmd/                    # Application entry points
│   └── server/            # Main server application
│       └── main.go        # Server initialization and startup
├── internal/              # Private application code
│   ├── controllers/       # HTTP request handlers
│   ├── middleware/        # HTTP middleware (auth, logging)
│   ├── models/           # Data models and database schemas
│   ├── repository/       # Data access layer
│   ├── routes/           # Route definitions and setup
│   ├── services/         # Business logic layer
│   └── shared/           # Shared utilities and helpers
├── pkg/                  # Public reusable packages
│   └── database/         # Database connection management
├── .env                  # Environment configuration
├── go.mod               # Go module dependencies
└── go.sum               # Dependency checksums
```

## Core Components

### Entry Point (cmd/server)
- **main.go**: Application bootstrap, database initialization, route setup, and server startup
- Configures timezone to America/New_York
- Loads environment variables from .env file
- Performs database migrations for all models
- Starts Gin HTTP server on configured port

### Controllers Layer (internal/controllers)
HTTP request handlers that process incoming requests and return responses:
- **auth_controller.go**: Login, logout, password reset request/confirmation
- **user_controller.go**: User CRUD operations, profile management
- **role_controller.go**: Role and permission management
- **business_unit_controller.go**: Business unit operations
- **division_unit_controller.go**: Division management
- **approval_controller.go**: User approval workflow
- **audit_logs_controller.go**: Audit log retrieval and filtering

### Services Layer (internal/services)
Business logic implementation and orchestration:
- **user_service.go**: User management logic, validation, password handling
- **role_service.go**: Role and permission business rules
- **business_unit_service.go**: Business unit operations
- **division_service.go**: Division management logic
- **approval_service.go**: Approval workflow processing
- **audit_logs_service.go**: Audit log creation and retrieval
- **email_service.go**: Email notification delivery

### Repository Layer (internal/repository)
Data access and database operations:
- **user_repo.go**: User database queries
- **role_repo.go**: Role and permission data access
- **business_unit_repo.go**: Business unit database operations
- **division_repo.go**: Division data access
- **audit_repo.go**: Audit log persistence and queries

### Models (internal/models)
Database schema definitions using GORM:
- **user.go**: User entity with relationships to roles, business units, divisions
- **roles.go**: Role and ModulePermissions entities
- **business_unit.go**: Business unit entity
- **division.go**: Division entity with business unit relationship
- **audit_logs.go**: Audit log entity for tracking system activities
- **password_reset.go**: Password reset token entity

### Middleware (internal/middleware)
HTTP middleware for cross-cutting concerns:
- **auth.go**: Authentication verification (RequireAuth) and role-based authorization (RequireRole)

### Routes (internal/routes)
Route configuration and grouping:
- **routes.go**: Centralized route definitions with middleware application

### Shared Utilities (internal/shared)
Common helper functions and utilities:
- **error_handler.go**: Centralized error handling and response formatting
- **logger.go**: Structured logging configuration using zerolog
- **password.go**: Password hashing and verification using bcrypt
- **context.go**: Context value extraction helpers
- **convert.go**: Data type conversion utilities
- **mergo_transformer.go**: Custom transformers for mergo library

### Database Package (pkg/database)
Database connection management:
- **database.go**: Database initialization, connection pooling, and cleanup

## Architectural Patterns

### Layered Architecture
The application follows a clean layered architecture:
1. **Controllers**: Handle HTTP concerns (request/response)
2. **Services**: Implement business logic
3. **Repository**: Manage data persistence
4. **Models**: Define data structures

### Dependency Flow
- Controllers depend on Services
- Services depend on Repositories
- Repositories depend on Models
- All layers can use Shared utilities

### Separation of Concerns
- HTTP handling separated from business logic
- Business logic separated from data access
- Database schema definitions isolated in models
- Reusable utilities centralized in shared package

### Session-Based Authentication
- Uses gin-contrib/sessions for session management
- Session data stored with user_id, user_email, user_roles, business_unit_id, division_id
- Middleware validates sessions and injects user context

### Repository Pattern
- Abstracts database operations behind repository interfaces
- Enables easier testing and potential database switching
- Centralizes query logic

### Service Layer Pattern
- Encapsulates business rules and validation
- Orchestrates multiple repository calls
- Provides transaction boundaries
