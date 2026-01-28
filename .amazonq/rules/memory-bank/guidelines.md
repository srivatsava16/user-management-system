# Development Guidelines

## Code Quality Standards

### Package Organization
- Use lowercase package names matching directory names (e.g., `package controllers`, `package services`)
- Keep packages focused on single responsibilities
- Internal packages should be in `internal/` directory for encapsulation
- Public reusable packages should be in `pkg/` directory

### Import Grouping
Organize imports in three groups with blank lines between:
1. Standard library imports
2. External dependencies
3. Internal project imports

Example:
```go
import (
    "os"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    
    "user-management-system/internal/models"
    "user-management-system/internal/services"
)
```

### Naming Conventions
- **Structs**: PascalCase (e.g., `UserService`, `RoleController`)
- **Functions/Methods**: PascalCase for exported, camelCase for unexported (e.g., `CreateUser`, `checkUserRoles`)
- **Variables**: camelCase (e.g., `userID`, `hashedPassword`)
- **Constants**: PascalCase or UPPER_SNAKE_CASE for exported constants
- **Interfaces**: PascalCase, often ending with "er" suffix (e.g., `Repository`, `Handler`)

### Error Handling
- Always check errors immediately after function calls
- Return errors up the call stack rather than handling at low levels
- Use custom error handlers for consistent error responses:
  - `shared.HandleValidationError(err)` for validation errors
  - `shared.HandleDatabaseError(err)` for database errors
- Provide descriptive error messages using `fmt.Errorf()` or `errors.New()`
- Log errors at appropriate levels using zerolog

Example:
```go
user, err := s.repo.GetUserByID(id)
if err != nil {
    return nil, err
}
```

### HTTP Response Patterns
- Use consistent JSON response format with `gin.H{}`
- Return appropriate HTTP status codes:
  - 200 OK for successful operations
  - 400 Bad Request for validation errors
  - 401 Unauthorized for authentication failures
  - 403 Forbidden for authorization failures
  - 500 Internal Server Error for server errors
- Always call `ctx.Abort()` after error responses in middleware

Example:
```go
if err != nil {
    customErr := shared.HandleDatabaseError(err)
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": customErr})
    return
}
ctx.JSON(http.StatusOK, user)
```

## Architectural Patterns

### Constructor Pattern
Use constructor functions for all structs that have dependencies:
```go
func NewUserService(repo *repository.UserRepo) *UserService {
    return &UserService{repo: repo}
}
```

### Dependency Injection
- Inject dependencies through constructors
- Pass dependencies from outer layers (routes) to inner layers (controllers → services → repositories)
- Initialize all dependencies in `routes.SetupRoutes()` function

Example from routes.go:
```go
userRepo := repository.NewUserRepo(db)
userService := services.NewUserService(userRepo)
userController := controllers.NewUserController(userService, auditService)
```

### Repository Pattern
- All database operations go through repository layer
- Repositories accept GORM DB instance in constructor
- Services call repository methods, never direct database access
- Repository methods return domain models, not database-specific types

### Service Layer Pattern
- Business logic resides in service layer
- Services orchestrate multiple repository calls
- Services handle validation and business rules
- Services return domain models to controllers

### Controller Pattern
- Controllers handle HTTP concerns only
- Extract request data using `ctx.ShouldBindJSON()` or `ctx.Param()`
- Call service methods for business logic
- Format and return HTTP responses
- Controllers should be thin, delegating to services

## Common Implementation Patterns

### Parameter Extraction from URL
```go
idstr := ctx.Param("id")
id, err := shared.StrToInt(idstr)
if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
    return
}
```

### Request Binding and Validation
```go
var req models.User
if err := ctx.ShouldBindJSON(&req); err != nil {
    customErr := shared.HandleValidationError(err)
    ctx.JSON(http.StatusBadRequest, gin.H{"error": customErr})
    return
}
```

### Password Hashing
Always hash passwords before storing:
```go
hashedPassword, err := shared.HashPassword(req.Password)
if err != nil {
    return nil, err
}
req.Password = hashedPassword
```

### Password Verification
```go
if !shared.CheckPassword(password, user.Password) {
    return nil, errors.New("Invalid credentials")
}
```

### Audit Logging Pattern
Log all CREATE, UPDATE, DELETE operations:
```go
changeBy := shared.GetUserFromContext(ctx)
c.auditService.LogChange("CREATE", changeBy, fmt.Sprintf("Created Role: %s", role.Name), nil, role)
```

### Deep Copy for Audit Trail
Create deep copies before updates to preserve original state:
```go
originalRole := DeepCopyRole(existingRole)
// ... perform updates ...
c.auditService.LogChange("UPDATE", changeBy, description, originalRole, updatedRole)
```

### Merging Updates with Mergo
Use mergo library for partial updates with custom transformers:
```go
err = mergo.Merge(existingRole, reqWithoutPermissions, mergo.WithOverride, mergo.WithTransformers(&shared.BooleanTransformer{}))
```

### Nil Pointer Handling
Check for nil pointers before dereferencing:
```go
if user.IsActive == nil || !*user.IsActive {
    return nil, errors.New("Account is not active")
}
```

### Input Validation in Services
Validate inputs at service layer entry points:
```go
if id == 0 || req == nil {
    return nil, fmt.Errorf("Invalid input for updating User")
}
```

## Middleware Patterns

### Authentication Middleware
```go
func RequireAuth() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        session := sessions.Default(ctx)
        userID := session.Get("user_id")
        if userID == nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            ctx.Abort()
            return
        }
        ctx.Set("user_id", userID)
        ctx.Next()
    }
}
```

### Authorization Middleware
```go
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        userRoles, exists := ctx.Get("user_roles")
        if !exists {
            ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
            ctx.Abort()
            return
        }
        hasPermission := checkUserRoles(userRoles, allowedRoles)
        if !hasPermission {
            ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            ctx.Abort()
            return
        }
        ctx.Next()
    }
}
```

### Middleware Application in Routes
```go
api := router.Group("/api")
api.Use(middleware.RequireAuth())
{
    userRoutes := api.Group("/users")
    userRoutes.Use(middleware.RequireRole("Super Admin", "BU Admin", "DV Admin"))
    {
        userRoutes.POST("", userController.CreateUser)
        userRoutes.GET("/:id", userController.GetUserByID)
    }
}
```

## Session Management

### Session Storage
Store essential user information in session:
```go
session.Set("user_id", user.ID)
session.Set("user_email", user.Email)
session.Set("user_roles", roleNames)
session.Set("business_unit_id", user.BusinessUnitID)
session.Set("division_id", user.DivisionID)
session.Save()
```

### Context Propagation
Middleware extracts session data and sets in context:
```go
ctx.Set("user_id", session.Get("user_id"))
ctx.Set("user_email", session.Get("user_email"))
ctx.Set("user_roles", session.Get("user_roles"))
```

## Logging Standards

### Structured Logging with Zerolog
Use zerolog for all logging with appropriate levels:
```go
log.Info().Msg("Server starting")
log.Warn().Msgf("Could not load .env file: %v", err)
log.Error().Err(err).Msg("Error loading timezone")
log.Fatal().Err(err).Msg("Failed to initialize database")
```

### Log Levels
- **Fatal**: Application cannot continue (e.g., database connection failure)
- **Error**: Errors that need attention but don't stop the application
- **Warn**: Warning conditions (e.g., missing .env file)
- **Info**: General informational messages
- **Debug**: Detailed debugging information

## Database Patterns

### Auto-Migration
Run migrations on application startup:
```go
err = database.DB.AutoMigrate(&models.BusinessUnit{}, &models.Division{}, &models.Role{})
```

### Resource Cleanup
Always defer cleanup operations:
```go
defer database.CloseDB()
```

## Configuration Management

### Environment Variables
- Load environment variables using godotenv
- Fail gracefully if .env file is missing (use system environment variables)
- Validate critical environment variables at startup:
```go
secretKey := os.Getenv("SESSION_SECRET")
if secretKey == "" {
    log.Fatal("SESSION_SECRET environment variable is not set")
}
```

## Initialization Patterns

### Init Function for Global Setup
Use init() for package-level initialization:
```go
func init() {
    est, err := time.LoadLocation("America/New_York")
    if err != nil {
        log.Error().Err(err).Msg("Error loading timezone")
    }
    time.Local = est
}
```

### Main Function Structure
1. Initialize logger
2. Load environment variables
3. Initialize database
4. Run migrations
5. Setup routes
6. Start server

## Security Best Practices

### Password Security
- Always hash passwords using bcrypt before storage
- Never log or return passwords in responses
- Use secure password reset tokens with expiration

### Session Security
- Use secure session storage with secret keys
- Store minimal data in sessions
- Validate session data on every request

### Authorization Checks
- Apply authentication middleware to all protected routes
- Use role-based middleware for fine-grained access control
- Check permissions at both route and controller levels

### Input Validation
- Use struct tags for validation (e.g., `validate:"required"`)
- Validate all user inputs at service layer
- Sanitize error messages to avoid information leakage

## Testing Considerations

### Testable Code Structure
- Use dependency injection for easy mocking
- Keep functions small and focused
- Separate business logic from HTTP handling
- Use interfaces for repositories to enable mocking

## Code Comments

### When to Comment
- Complex business logic that isn't immediately obvious
- Public APIs and exported functions
- Workarounds or non-obvious solutions
- TODO items for future improvements

### When NOT to Comment
- Self-explanatory code
- Obvious operations
- Redundant descriptions of what code does

Keep code self-documenting through clear naming and structure rather than relying on comments.
