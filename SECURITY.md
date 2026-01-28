# Security Policy

## Reporting Security Vulnerabilities

If you discover a security vulnerability in this project, please report it responsibly:

1. **DO NOT** create a public GitHub issue
2. Email security details to: security@example.com
3. Include steps to reproduce the vulnerability
4. Allow reasonable time for a fix before public disclosure

We take security seriously and will respond promptly to verified vulnerabilities.

## Security Features

### Authentication & Authorization

- **JWT Tokens**: Stateless authentication with 30-minute expiration
- **Password Hashing**: Bcrypt with default cost (cost factor 12)
- **Role-Based Access Control**: Three-tier permission system
- **Token Invalidation**: Automatic invalidation on logout, password change, or role modification

### Password Security

All passwords must meet these requirements:
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character

Passwords are hashed using bcrypt before storage. Plaintext passwords are never stored.

### Rate Limiting

Protection against brute force attacks:
- **Login Endpoint**: 5 attempts per 15 minutes per IP
- **Password Reset**: 3 attempts per hour per IP

### Security Headers

The following security headers are automatically added to all responses:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy: default-src 'self'`

When using HTTPS (required in production):
- `Strict-Transport-Security: max-age=31536000; includeSubDomains`

### Input Validation

- Email format validation
- Username format validation (alphanumeric, 3-50 characters)
- Password complexity validation
- SQL injection prevention via parameterized queries (GORM)
- JSON input validation

### CORS Protection

CORS is configured to only allow requests from trusted origins specified in the `FRONTEND_URL` environment variable.

## Security Best Practices

### For Deployment

1. **Environment Variables**
   - Never commit `.env` files to version control
   - Generate new secrets for each environment
   - Use strong, unique secrets (32+ characters)
   - Rotate secrets regularly (every 90 days recommended)

2. **HTTPS/TLS**
   - Always use HTTPS in production
   - Use valid TLS certificates (Let's Encrypt, etc.)
   - Enable HSTS header (uncomment in `security_headers.go`)
   - Redirect all HTTP traffic to HTTPS

3. **Database Security**
   - Use strong database passwords
   - Create a dedicated database user with minimal privileges
   - Restrict database access by IP address
   - Enable MySQL/MariaDB SSL connections
   - Regular backups with encryption

4. **Network Security**
   - Use a firewall (only expose necessary ports)
   - Run behind a reverse proxy (nginx, Apache)
   - Implement IP whitelisting for admin endpoints
   - Use VPN or bastion hosts for database access

5. **Monitoring & Logging**
   - Monitor failed login attempts
   - Set up alerts for suspicious activity
   - Log security events (not sensitive data)
   - Regular security audits
   - Keep logs for compliance requirements

6. **Dependencies**
   - Regularly update Go and dependencies
   - Use `go mod verify` to check integrity
   - Scan for vulnerabilities: `go list -json -m all | nancy sleuth`
   - Subscribe to security advisories

### Secret Management

#### Generating Secrets

Generate cryptographically secure secrets:

```bash
# For JWT_SECRET and SESSION_SECRET
openssl rand -base64 32

# Or using Go
go run -c 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'
```

#### Secret Rotation

To rotate JWT_SECRET:

1. Generate a new secret
2. Update environment variable in production
3. Restart the application
4. All existing tokens will be invalidated
5. Users must log in again

### Password Reset Security

- Reset tokens expire after 1 hour
- Tokens are single-use only
- Tokens are 64-character random hex strings
- Email delivery uses SMTP with authentication
- Failed attempts are rate-limited

### Token Security

- Tokens are signed with HS256 (HMAC SHA-256)
- Each token has a unique JTI (JWT ID) for tracking
- Tokens are automatically invalidated on:
  - User logout
  - Password change
  - Role modification
- Expired tokens are automatically cleaned up

## Known Security Considerations

### Current Limitations

1. **No Refresh Tokens**: Access tokens expire after 30 minutes. Users must re-authenticate. Consider implementing refresh tokens for better UX.

2. **IP-Based Rate Limiting**: Uses client IP for rate limiting. May need adjustment if behind a proxy (use `X-Forwarded-For` with caution).

3. **Email Security**: Password reset emails are sent via SMTP. Ensure SMTP credentials are secured and use TLS.

4. **Session Storage**: JWT tokens are stateless. Token invalidation relies on database checks on every request.

## Security Checklist for Production

- [ ] Generate new secrets (JWT_SECRET, SESSION_SECRET)
- [ ] Enable HTTPS/TLS with valid certificates
- [ ] Configure firewall (block all except 80/443)
- [ ] Set up database with strong password
- [ ] Restrict database access by IP
- [ ] Enable database SSL connections
- [ ] Configure CORS with specific frontend origins
- [ ] Set up monitoring and alerting
- [ ] Enable HSTS header
- [ ] Regular database backups (encrypted)
- [ ] Set up log rotation and retention
- [ ] Document incident response procedures
- [ ] Regular security audits
- [ ] Keep dependencies updated
- [ ] Test disaster recovery procedures

## Compliance

### Password Storage
Passwords are hashed using bcrypt (NIST recommended), making them compliant with:
- OWASP Password Storage Guidelines
- NIST SP 800-63B
- PCI DSS Requirement 8.2.1

### Audit Logging
All changes are logged with:
- User identification
- Timestamp
- Action performed
- Before/after state

Suitable for:
- SOC 2 compliance
- GDPR audit requirements
- HIPAA audit trails

### Data Protection
- Passwords never logged or transmitted in plaintext
- Sensitive data excluded from error messages
- Database queries use parameterized statements

## Security Updates

We monitor the following for security advisories:
- Go security announcements
- GORM security issues
- Gin framework advisories
- JWT library updates
- MySQL/MariaDB security bulletins

## Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [Go Security Best Practices](https://golang.org/doc/security/best-practices)
- [NIST Password Guidelines](https://pages.nist.gov/800-63-3/sp800-63b.html)

---

Last Updated: 2026-01-28
