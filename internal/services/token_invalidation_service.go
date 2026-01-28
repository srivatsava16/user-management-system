package services

import (
	"time"
	"user-management-system/internal/repository"
)

type TokenInvalidationService struct {
	repo *repository.TokenInvalidationRepo
}

func NewTokenInvalidationService(repo *repository.TokenInvalidationRepo) *TokenInvalidationService {
	return &TokenInvalidationService{repo: repo}
}

func (s *TokenInvalidationService) InvalidateUserTokens(userID int, reason string) error {

	return s.repo.CreateInvalidation(userID, reason)
}

func (s *TokenInvalidationService) IsUserTokenInvalidated(userID int) bool {

	hasInvalidation, err := s.repo.HasInvalidation(userID)
	if err != nil {
		return false
	}
	return hasInvalidation
}

func (s *TokenInvalidationService) ClearUserInvalidations(userID int) error {
	return s.repo.ClearInvalidations(userID)
}

func (s *TokenInvalidationService) CleanupExpiredInvalidations() error {

	cutoffTime := time.Now().Add(-2 * time.Hour)

	return s.repo.DeleteExpiredInvalidations(cutoffTime)

}

func (s *TokenInvalidationService) InvalidateToken(jti string, userID int, reason string) error {
	return s.repo.CreateTokenInvalidation(jti, userID, reason)
}

func (s *TokenInvalidationService) IsTokenInvalidated(jti string) bool {
	return s.repo.HasInvalidationByJTI(jti)
}
