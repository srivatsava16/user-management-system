# Migration Guide - v1.0 to v2.0

## 🚨 CRITICAL: Immediate Actions Required

### 1. Rotate All Secrets (URGENT)

Your current `.env` file contains exposed secrets. You MUST rotate all secrets immediately:

```bash
# Generate new JWT_SECRET
openssl rand -base64 32

# Generate new SESSION_SECRET
openssl rand -base64 32
```

Update your `.env` file with these new values:
```env
SESSION_SECRET=<new-secret-here>
JWT_SECRET=<new-secret-here>
```

### 2. Update Database Credentials

Change your production database password:
```sql
ALTER USER 'techuser'@'%' IDENTIFIED BY 'new-secure-password';
```

Then update `.env`:
```env
DB_PASSWORD=new-secure-password
```

### 3. Rotate SMTP Credentials

Generate new SMTP credentials in Mailtrap and update `.env`.

## 📋 Step-by-Step Migration

### Step 1: Backup

```bash
# Backup your database
mysqldump -u techuser -p ZX_REQUESTS > backup_$(date +%Y%m%d).sql

# Backup your current application
cp -r user-management-system user-management-system-backup
```

### Step 2: Update Code

The code changes have already been applied. Verify by checking:

```bash
# Check .gitignore includes .env
cat .gitignore | grep ".env"

# Verify new files exist
ls -la internal/constants/
ls -la internal/middleware/rate_limit.go
ls -la internal/middleware/security_headers.go
```

### Step 3: Update Environment Variables

```bash
# Copy the new template
cp .env.example .env.new

# Migrate your settings from old .env to .env.new
# Add your new secrets (generated in Critical Actions above)

# Replace old file
mv .env.new .env
```

### Step 4: Update Dependencies

```bash
go mod download
go mod tidy
```

### Step 5: Update Frontend API Calls

All API endpoints have moved from `/api/*` to `/api/v1/*`:

**Before:**
```javascript
// Old endpoints
POST /api/auth/login
GET /api/users
POST /api/business-units
```

**After:**
```javascript
// New endpoints
POST /api/v1/auth/login
GET /api/v1/users
POST /api/v1/business-units
```

### Step 6: Test Locally

```bash
# Start the server
go run cmd/server/main.go

# In another terminal, test health check
curl http://localhost:8083/health

# Test login (should be rate limited after 5 attempts)
curl -X POST http://localhost:8083/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!@#"}'
```

### Step 7: Update Documentation

Update your API documentation and frontend integration guides with:
- New API versioning (`/api/v1`)
- Password requirements
- Pagination support
- Rate limiting information

### Step 8: Deploy to Production

```bash
# Build the application
go build -o user-management-system cmd/server/main.go

# Stop old server gracefully
kill -SIGTERM <old-pid>

# Start new server
./user-management-system
```

## 🔄 API Changes

### Endpoint Changes

| Old Endpoint | New Endpoint | Status |
|-------------|--------------|--------|
| `/api/auth/login` | `/api/v1/auth/login` | Rate limited (5/15min) |
| `/api/auth/forgot-password` | `/api/v1/auth/forgot-password` | Rate limited (3/hour) |
| `/api/users` | `/api/v1/users` | Supports pagination |
| All other `/api/*` | `/api/v1/*` | No functional change |

### New Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check endpoint |

### Pagination

List endpoints now support optional pagination:

```bash
# Get all users (paginated)
GET /api/v1/users?page=1&page_size=20

# Response format
{
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total_items": 100,
    "total_pages": 5
  }
}

# Without pagination (backward compatible)
GET /api/v1/users
# Returns: [...]
```

## 🔐 Security Changes

### Password Requirements (NEW)

All passwords must now meet these requirements:
- Minimum 8 characters
- At least one uppercase letter (A-Z)
- At least one lowercase letter (a-z)
- At least one number (0-9)
- At least one special character (!@#$%^&*()_+-=[]{}|;:,.<>?)

**Impact:**
- New user registration will enforce these rules
- Existing users are not affected until password change
- Password reset will enforce these rules

### Token Expiration

| Token Type | Old | New | Impact |
|-----------|-----|-----|--------|
| Access Token | 60 minutes | 30 minutes | Users need to re-login more frequently |
| Reset Token | 60 minutes | 60 minutes | No change |

### Rate Limiting (NEW)

| Endpoint | Limit | Window | Action if Exceeded |
|----------|-------|--------|-------------------|
| Login | 5 attempts | 15 minutes | HTTP 429 (Too Many Requests) |
| Password Reset | 3 attempts | 60 minutes | HTTP 429 (Too Many Requests) |

**Testing Rate Limits:**
```bash
# Test login rate limiting
for i in {1..6}; do
  curl -X POST http://localhost:8083/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"wrong"}'
  echo "\nAttempt $i"
done
# 6th attempt should return 429
```

## 🧪 Testing Checklist

- [ ] Health check endpoint works
- [ ] Login with correct credentials
- [ ] Login fails with weak password
- [ ] Rate limiting triggers after 5 failed logins
- [ ] Password reset with weak password rejected
- [ ] Password reset rate limiting works
- [ ] Pagination works with `?page=1&page_size=20`
- [ ] All existing API calls work with `/api/v1` prefix
- [ ] CORS allows frontend origin
- [ ] Token expires after 30 minutes
- [ ] Graceful shutdown works (Ctrl+C)

## 🐛 Troubleshooting

### Issue: "Too many requests" on first login

**Cause:** Rate limiter is tracking by IP. If behind a proxy, all requests may appear from same IP.

**Solution:** If behind a reverse proxy, configure it to pass `X-Forwarded-For` header.

### Issue: CORS errors from frontend

**Cause:** `FRONTEND_URL` in `.env` doesn't match your frontend domain.

**Solution:** Update `.env`:
```env
FRONTEND_URL=https://your-frontend-domain.com
```

For multiple origins:
```env
FRONTEND_URL=http://localhost:3000,https://app.example.com
```

### Issue: Password reset emails not sending

**Cause:** SMTP credentials may be invalid after rotation.

**Solution:** Verify SMTP settings in `.env` and test:
```bash
# Check SMTP connection
telnet sandbox.smtp.mailtrap.io 2525
```

### Issue: Existing passwords don't work

**Cause:** You may have accidentally re-hashed already hashed passwords.

**Solution:** This shouldn't happen with the fix, but if it does:
1. Use password reset for affected users
2. Check logs for password update operations

### Issue: Database connection fails

**Cause:** Connection pool settings may be too aggressive.

**Solution:** Adjust in `pkg/database/database.go`:
```go
sqlDB.SetMaxOpenConns(50)  // Reduce from 100
sqlDB.SetMaxIdleConns(5)   // Reduce from 10
```

## 📊 Monitoring

### Key Metrics to Monitor

1. **Failed Login Attempts**
   - Normal: < 5% of login attempts
   - Alert if: > 20% failed attempts

2. **Rate Limit Triggers**
   - Check logs for frequent 429 responses
   - May indicate attack or misconfiguration

3. **Token Expirations**
   - Monitor re-login frequency
   - Adjust token expiration if too many complaints

4. **Database Connection Pool**
   - Monitor active connections
   - Alert if consistently near max

5. **Response Times**
   - Health check should be < 100ms
   - Login should be < 500ms
   - List endpoints should be < 1s

### Log Queries

```bash
# Count failed login attempts
grep "Failed login attempt" /var/log/app.log | wc -l

# Count rate limit triggers
grep "Too many requests" /var/log/app.log | wc -l

# Check graceful shutdowns
grep "Server exited gracefully" /var/log/app.log
```

## 🔙 Rollback Procedure

If you need to rollback:

```bash
# Stop new server
kill -SIGTERM <new-pid>

# Restore backup
cp -r user-management-system-backup/* user-management-system/

# Restore database (if needed)
mysql -u techuser -p ZX_REQUESTS < backup_YYYYMMDD.sql

# Start old server
cd user-management-system
go run cmd/server/main.go

# Update frontend to use old API paths (/api/*)
```

## 📞 Support

If you encounter issues:

1. Check logs: `tail -f /var/log/app.log`
2. Test health: `curl http://localhost:8083/health`
3. Review `SECURITY.md` for configuration issues
4. Check `CHANGELOG.md` for breaking changes

## ✅ Post-Migration Checklist

- [ ] All secrets rotated
- [ ] Database password changed
- [ ] SMTP credentials rotated
- [ ] `.env` file not in git (check: `git status`)
- [ ] Frontend updated to use `/api/v1`
- [ ] Password requirements documented for users
- [ ] Monitoring set up for rate limits
- [ ] Load testing completed
- [ ] Backup and rollback plan tested
- [ ] Team trained on new password requirements
- [ ] Documentation updated

---

**Migration completed?** Update your production deployment documentation with these changes.
