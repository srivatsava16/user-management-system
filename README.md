# User Management System

A comprehensive, enterprise-grade user management system built with Go, featuring role-based access control, organizational hierarchy, approval workflows, and comprehensive audit logging.

## Features

- **Authentication & Authorization**
  - JWT-based stateless authentication (30-minute token expiration)
  - Bcrypt password hashing
  - Role-based access control (RBAC)
  - Token invalidation on logout and role changes
  - Rate-limited login attempts

- **User Management**
  - User registration with approval workflow
  - Email and username uniqueness validation
  - Password complexity requirements
  - Secure password reset via email
  - User status management (active/inactive)
  - Multi-role assignment per user

- **Organizational Hierarchy**
  - Business Units (top-level)
  - Divisions (sub-units within business units)
  - Hierarchical user assignment

- **Security Features**
  - Password strength validation (8+ chars, uppercase, lowercase, number, special char)
  - Rate limiting on login (5 attempts per 15 minutes)
  - Rate limiting on password reset (3 attempts per hour)
  - CORS configuration
  - Security headers (X-Frame-Options, CSP, etc.)
  - Input validation and sanitization
  - SQL injection prevention via GORM

- **Audit & Compliance**
  - Comprehensive audit logging
  - Before/after state tracking
  - All CRUD operations logged

- **API Features**
  - RESTful API with versioning (/api/v1)
  - Pagination support
  - Health check endpoint
  - JSON responses

## Tech Stack

- **Language**: Go 1.24.2
- **Framework**: Gin (HTTP framework)
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password**: Bcrypt (golang.org/x/crypto)
- **Logging**: Zerolog
- **Email**: SMTP

## Prerequisites

- Go 1.24.2 or higher
- MySQL 5.7+ or MySQL 8.0+
- SMTP server (for password reset emails)

## Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd user-management-system
```

### 2. Set up environment variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` and configure your settings:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_NAME=user_management_system

# Frontend URL (for CORS)
FRONTEND_URL=http://localhost:3000

# SMTP Configuration (for password reset)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_FROM=noreply@example.com
SMTP_USERNAME=your_smtp_username
SMTP_PASSWORD=your_smtp_password

# Server Configuration
SERVER_PORT=8083

# Security Secrets (IMPORTANT: Generate new secrets!)
# Generate with: openssl rand -base64 32
SESSION_SECRET=your_session_secret_here
JWT_SECRET=your_jwt_secret_here

# Timezone
TZ=America/New_York
```

### 3. Generate secure secrets

Generate new secrets for SESSION_SECRET and JWT_SECRET:

```bash
# Generate SESSION_SECRET
openssl rand -base64 32

# Generate JWT_SECRET
openssl rand -base64 32
```

Copy the generated values into your `.env` file.

### 4. Install dependencies

```bash
go mod download
```

### 5. Set up the database

Create a MySQL database:

```sql
CREATE DATABASE user_management_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

The application will automatically run migrations on startup.

### 6. Run the application

```bash
go run cmd/server/main.go
```

The server will start on the port specified in `SERVER_PORT` (default: 8083).

## API Documentation

### Base URL

```
http://localhost:8083/api/v1
```

### Health Check

Check if the service is healthy:

```
GET /health
```

Response:
```json
{
  "status": "healthy",
  "database": "connected"
}
```

### Authentication Endpoints

#### Login

```
POST /api/v1/auth/login
```

Request:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

Response:
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": 1,
  "roles": ["Super Admin"]
}
```

#### Forgot Password

```
POST /api/v1/auth/forgot-password
```

Request:
```json
{
  "email": "user@example.com"
}
```

#### Reset Password

```
POST /api/v1/auth/reset-password/:token
```

Request:
```json
{
  "password": "NewSecurePass123!"
}
```

### Protected Endpoints

All endpoints below require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your-token>
```

#### Logout

```
POST /api/v1/logout
```

#### User Management

```
POST   /api/v1/users              # Create user
GET    /api/v1/users              # List all users (supports pagination)
GET    /api/v1/users/:id          # Get user by ID
PUT    /api/v1/users/:id          # Update user
GET    /api/v1/users/divisions/:id/users        # Get users by division
GET    /api/v1/users/business-units/:id/users   # Get users by business unit
```

##### Pagination

Add query parameters to list endpoints:

```
GET /api/v1/users?page=1&page_size=20
```

Response:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total_items": 100,
    "total_pages": 5
  }
}
```

#### Business Unit Management

```
POST   /api/v1/business-units     # Create (Super Admin only)
GET    /api/v1/business-units     # List all
GET    /api/v1/business-units/:id # Get by ID
PUT    /api/v1/business-units/:id # Update (Super Admin only)
```

#### Division Management

```
POST   /api/v1/divisions          # Create (Super Admin, BU Admin)
GET    /api/v1/divisions          # List all
GET    /api/v1/divisions/:id      # Get by ID
PUT    /api/v1/divisions/:id      # Update (Super Admin, BU Admin)
GET    /api/v1/business-units/:id/divisions  # Get divisions by business unit
```

#### Role Management

```
POST   /api/v1/roles              # Create role
GET    /api/v1/roles              # List all roles
GET    /api/v1/roles/:id          # Get role by ID
PUT    /api/v1/roles/:id          # Update role
GET    /api/v1/roles/users/:id/permissions  # Get user permissions
```

#### Approval Workflow

```
PUT    /api/v1/approvals/approve                # Approve entity
GET    /api/v1/approvals/pending/:entity_type  # Get pending approvals
```

#### Audit Logs

```
GET    /api/v1/audit-logs         # Retrieve audit logs
```

## Security Best Practices

### 🔒 Critical Security Notes

1. **Never commit the .env file to version control**
   - The `.env` file is already in `.gitignore`
   - Use `.env.example` as a template

2. **Generate strong secrets**
   - Use `openssl rand -base64 32` to generate secrets
   - Never reuse secrets across environments
   - Rotate secrets regularly

3. **Use HTTPS in production**
   - Never run production services over HTTP
   - Configure TLS certificates
   - Enable HSTS header (uncomment in security_headers.go)

4. **Protect your database**
   - Use strong database passwords
   - Limit database access by IP
   - Use a dedicated database user with minimal permissions
   - Regular backups

5. **Rate limiting**
   - Login: 5 attempts per 15 minutes
   - Password reset: 3 attempts per hour
   - Adjust rates in `internal/middleware/rate_limit.go` if needed

6. **Password requirements**
   - Minimum 8 characters
   - At least one uppercase letter
   - At least one lowercase letter
   - At least one number
   - At least one special character

## Password Requirements

Users must create passwords that meet these requirements:
- Minimum 8 characters
- At least one uppercase letter (A-Z)
- At least one lowercase letter (a-z)
- At least one number (0-9)
- At least one special character (!@#$%^&*()_+-=[]{}|;:,.<>?)

## Role Types

- **Super Admin**: Full system access
- **BU Admin**: Business unit level access
- **DV Admin**: Division level access

## Project Structure

```
user-management-system/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── constants/               # Application constants
│   ├── controllers/             # HTTP request handlers
│   ├── middleware/              # JWT auth, rate limiting, security
│   ├── models/                  # Data models
│   ├── repository/              # Database access layer
│   ├── routes/                  # Route definitions
│   ├── services/                # Business logic
│   └── shared/                  # Shared utilities
├── pkg/
│   └── database/                # Database configuration
├── .env.example                 # Environment template
├── .gitignore                   # Git ignore file
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Development

### Running in Development Mode

```bash
go run cmd/server/main.go
```

### Building for Production

```bash
go build -o user-management-system cmd/server/main.go
./user-management-system
```

### Testing Database Connection

The health check endpoint can verify database connectivity:

```bash
curl http://localhost:8083/health
```

## Database Schema

The application uses the following main tables:

- `APT_CUSTOM_ZXDS_USERS` - User accounts
- `APT_CUSTOM_ZXDS_ROLES` - User roles
- `APT_CUSTOM_ZXDS_MODULE_PERMISSIONS` - Module-level permissions
- `APT_CUSTOM_ZXDS_BUSINESS_UNITS` - Business units
- `APT_CUSTOM_ZXDS_DIVISIONS` - Divisions
- `APT_CUSTOM_ZXDS_AUDIT_LOGS` - Audit trail
- `APT_CUSTOM_ZXDS_PASSWORD_RESET_TOKENS` - Password reset tokens
- `APT_CUSTOM_ZXDS_TOKEN_INVALIDATIONS` - Invalidated JWT tokens

All tables are automatically created via GORM migrations on startup.

## Configuration

### Database Connection Pool

Default settings (configurable in `pkg/database/database.go`):
- Max Idle Connections: 10
- Max Open Connections: 100
- Connection Max Lifetime: 1 hour
- Connection Max Idle Time: 10 minutes

### Token Expiration

Default settings (configurable in `internal/constants/constants.go`):
- Access Token: 30 minutes
- Password Reset Token: 1 hour

### Pagination

Default settings:
- Default Page Size: 20 items
- Maximum Page Size: 100 items

## Troubleshooting

### Database Connection Fails

1. Verify MySQL is running
2. Check database credentials in `.env`
3. Ensure database exists
4. Check network connectivity and firewall rules

### Token Expired Errors

JWT tokens expire after 30 minutes. Users need to log in again to get a new token.

### Rate Limit Errors

If you see "Too many requests" errors:
- Login: Wait 15 minutes or adjust rate limits
- Password Reset: Wait 1 hour or adjust rate limits

### CORS Errors

Update `FRONTEND_URL` in `.env` to include your frontend domain. Multiple origins can be comma-separated.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write tests (if applicable)
5. Submit a pull request

## Security Vulnerabilities

If you discover a security vulnerability, please email security@example.com instead of using the issue tracker.

## License

[Your License Here]

## Support

For support, please contact support@example.com or create an issue in the repository.

## Credits

Built with:
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Zerolog](https://github.com/rs/zerolog)

---

**Note**: This is a production application. Ensure all security best practices are followed before deploying to production environments.
