# Changelog

All notable changes to this project are documented in this file.

## [2.0.0] - 2026-01-28

### 🔒 CRITICAL Security Fixes

- **Added `.env` to `.gitignore`** - Prevents accidental exposure of secrets in version control
- **Created `.env.example`** - Template for environment configuration without sensitive data
- **IMPORTANT**: All production secrets must be regenerated and rotated immediately

### ✨ New Features

#### Security Enhancements
- **Password Validation** - Enforced password complexity requirements:
  - Minimum 8 characters
  - At least one uppercase letter
  - At least one lowercase letter
  - At least one number
  - At least one special character
- **Rate Limiting** - Protection against brute force attacks:
  - Login endpoint: 5 attempts per 15 minutes per IP
  - Password reset: 3 attempts per hour per IP
- **Security Headers Middleware** - Added security headers to all responses:
  - X-Content-Type-Options
  - X-Frame-Options
  - X-XSS-Protection
  - Referrer-Policy
  - Content-Security-Policy
- **CORS Configuration** - Proper CORS setup with configurable origins
- **Input Validation** - Email, username, and password validation
- **Token Expiration Improvement** - Reduced from 1 hour to 30 minutes

#### API Improvements
- **API Versioning** - All routes now under `/api/v1`
- **Health Check Endpoint** - `/health` endpoint for monitoring
- **Pagination Support** - List endpoints now support pagination:
  - Query parameters: `?page=1&page_size=20`
  - Returns total count and page metadata
  - Backward compatible (works without pagination params)

#### Infrastructure
- **Database Connection Pool** - Configured connection pooling:
  - Max Idle Connections: 10
  - Max Open Connections: 100
  - Connection Max Lifetime: 1 hour
  - Connection Max Idle Time: 10 minutes
- **Graceful Shutdown** - Server now handles SIGINT/SIGTERM gracefully
  - 10-second timeout for in-flight requests
  - Clean database connection closure
  - Proper logging of shutdown process

#### Code Quality
- **Constants Package** - Centralized constants for:
  - Role names (no more magic strings)
  - Password requirements
  - Pagination defaults
  - Token expiration times
- **Validation Utilities** - New helper functions for:
  - Password strength validation
  - Email format validation
  - Username format validation

### 🐛 Bug Fixes

- **Fixed Password Update Logic** - Password is now only hashed when actually being changed
- **Fixed Password Reset Token Validation** - Added explicit token expiration and usage checks
- **Fixed User Model Validation** - Added proper validation tags to User model fields

### 📚 Documentation

- **README.md** - Comprehensive documentation including:
  - Features overview
  - Installation instructions
  - API documentation
  - Security best practices
  - Configuration guide
  - Troubleshooting
- **SECURITY.md** - Security policy and guidelines:
  - Vulnerability reporting process
  - Security features documentation
  - Deployment best practices
  - Secret management guidelines
  - Compliance information
- **CHANGELOG.md** - This file

### 🔄 Breaking Changes

- **API Versioning** - All API routes moved from `/api/*` to `/api/v1/*`
  - Old: `POST /api/auth/login`
  - New: `POST /api/v1/auth/login`
- **Token Expiration** - JWT tokens now expire after 30 minutes (was 1 hour)
- **Password Requirements** - All passwords must now meet complexity requirements

### 📦 Dependencies

- **Added** `github.com/gin-contrib/cors v1.7.2` - CORS middleware

### 🗂️ File Structure Changes

#### New Files
- `.env.example` - Environment configuration template
- `internal/constants/constants.go` - Application constants
- `internal/shared/password_validator.go` - Password validation utilities
- `internal/shared/pagination.go` - Pagination utilities
- `internal/middleware/rate_limit.go` - Rate limiting middleware
- `internal/middleware/security_headers.go` - Security headers middleware
- `README.md` - Project documentation
- `SECURITY.md` - Security documentation
- `CHANGELOG.md` - This changelog

#### Modified Files
- `.gitignore` - Added `.env` and `.env.local`
- `go.mod` - Added CORS dependency
- `cmd/server/main.go` - Added graceful shutdown
- `pkg/database/database.go` - Added connection pool configuration
- `internal/models/user.go` - Added validation tags
- `internal/services/user_service.go` - Added password validation
- `internal/services/jwt_service.go` - Updated token expiration
- `internal/repository/user_repo.go` - Added pagination support
- `internal/controllers/user_controller.go` - Added pagination support
- `internal/routes/routes.go` - Added CORS, security headers, rate limiting, API versioning

### ⚠️ Migration Guide

#### For Existing Deployments

1. **Update API Endpoints**
   ```
   # Update all API calls from:
   /api/auth/login → /api/v1/auth/login
   /api/users → /api/v1/users
   # etc.
   ```

2. **Rotate All Secrets**
   ```bash
   # Generate new JWT_SECRET
   openssl rand -base64 32

   # Generate new SESSION_SECRET
   openssl rand -base64 32

   # Update .env file with new values
   ```

3. **Update Dependencies**
   ```bash
   go mod download
   ```

4. **Update Frontend CORS Configuration**
   ```
   # Update FRONTEND_URL in .env to match your frontend domain
   FRONTEND_URL=https://your-frontend-domain.com
   ```

5. **Test Password Requirements**
   - All new passwords must meet complexity requirements
   - Existing users will need to update passwords on next reset
   - Update documentation/UI to reflect new requirements

6. **Monitor Rate Limiting**
   - Check for false positives
   - Adjust limits if needed in `internal/middleware/rate_limit.go`

### 📈 Performance Improvements

- **Database Connection Pooling** - Reduces connection overhead
- **Pagination** - Reduces memory usage and response size for large datasets
- **Rate Limiting** - In-memory implementation with automatic cleanup

### 🔐 Security Improvements Summary

| Issue | Severity | Status |
|-------|----------|--------|
| Secrets in repository | Critical | ✅ Fixed |
| Password validation missing | High | ✅ Fixed |
| No rate limiting | High | ✅ Fixed |
| Token expiration too long | Medium | ✅ Fixed |
| No CORS configuration | Medium | ✅ Fixed |
| No security headers | Medium | ✅ Fixed |
| No password update validation | High | ✅ Fixed |
| No token reset validation | Medium | ✅ Fixed |
| Magic strings for roles | Low | ✅ Fixed |
| No pagination | Medium | ✅ Fixed |
| No health check | Low | ✅ Fixed |
| No API versioning | Medium | ✅ Fixed |
| No connection pooling | Low | ✅ Fixed |
| No graceful shutdown | Low | ✅ Fixed |

### 🎯 Next Steps / Recommendations

1. **Testing** - Add comprehensive unit and integration tests
2. **Monitoring** - Implement Prometheus metrics
3. **HTTPS** - Configure TLS certificates for production
4. **Refresh Tokens** - Implement refresh token mechanism
5. **2FA** - Consider adding two-factor authentication
6. **API Documentation** - Generate Swagger/OpenAPI docs
7. **Docker** - Create Dockerfile and docker-compose.yml
8. **CI/CD** - Set up automated testing and deployment

### 📝 Notes

- This release focuses on security hardening and production readiness
- All critical security issues have been addressed
- Backward compatibility is maintained except for API versioning
- Comprehensive documentation has been added
- No database schema changes required

---

## [1.0.0] - Initial Release

- Initial implementation of user management system
- JWT authentication
- Role-based access control
- Organizational hierarchy (Business Units, Divisions)
- Audit logging
- Password reset functionality
- Approval workflows
