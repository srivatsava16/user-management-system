package services

import (
	"errors"
	"fmt"
	"os"
	"time"
	"user-management-system/internal/constants"
	"user-management-system/internal/shared"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey []byte
}

type JWTClaims struct {
	UserID         int      `json:"user_id"`
	Email          string   `json:"email"`
	Roles          []string `json:"roles"`
	BusinessUnitID int      `json:"business_unit_id"`
	DivisionID     int      `json:"division_id"`
	JTI            string   `json:"jti"`
	jwt.RegisteredClaims
}

func NewJWTService() (*JWTService, error) {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable not set")
	}

	return &JWTService{
		secretKey: []byte(secret),
	}, nil
}

func (j *JWTService) GenerateToken(userID int, email string, roles []string, businessUnitID, divisionID int) (string, error) {
	claims := JWTClaims{
		UserID:         userID,
		Email:          email,
		Roles:          roles,
		BusinessUnitID: businessUnitID,
		DivisionID:     divisionID,
		JTI:            uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(constants.AccessTokenExpiration) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-management-system"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secretKey)

}

func (j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.secretKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, shared.ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, shared.ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, shared.ErrTokenInvalidSignature
		default:
			return nil, shared.ErrTokenInvalid
		}
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, shared.ErrTokenInvalid
	}

	return claims, nil
}
