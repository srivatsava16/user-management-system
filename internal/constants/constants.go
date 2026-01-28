package constants

// Role names constants
const (
	RoleSuperAdmin = "Super Admin"
	RoleBUAdmin    = "BU Admin"
	RoleDVAdmin    = "DV Admin"
)

// Password requirements
const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
)

// Pagination defaults
const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Token expiration
const (
	AccessTokenExpiration  = 30 // minutes
	RefreshTokenExpiration = 24 // hours
	ResetTokenExpiration   = 1  // hour
)
