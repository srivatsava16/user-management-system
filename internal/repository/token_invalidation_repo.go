package repository

import (
	"time"
	"user-management-system/internal/models"

	"gorm.io/gorm"
)

type TokenInvalidationRepo struct {
	db *gorm.DB
}

func NewTokenInvalidationRepo(db *gorm.DB) *TokenInvalidationRepo {
	return &TokenInvalidationRepo{db: db}
}

func (r *TokenInvalidationRepo) CreateInvalidation(userID int, reason string) error {
	invalidation := &models.TokenInvalidation{
		UserID: userID,
		Reason: reason,
	}
	return r.db.Create(invalidation).Error
}

func (r *TokenInvalidationRepo) HasInvalidation(userID int) (bool, error) {
	var count int64
	err := r.db.Model(&models.TokenInvalidation{}).Where("user_id = ?", userID).Count(&count).Error

	return count > 0, err
}

func (r *TokenInvalidationRepo) ClearInvalidations(userID int) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.TokenInvalidation{}).Error
}

func (r *TokenInvalidationRepo) DeleteExpiredInvalidations(cutoffTime time.Time) error {
	return r.db.Where("created_at < ?", cutoffTime).Delete(&models.TokenInvalidation{}).Error
}

func (r *TokenInvalidationRepo) CreateTokenInvalidation(jti string, userID int, reason string) error {
	invalidation := &models.TokenInvalidation{

		JTI:    jti,
		UserID: userID,
		Reason: reason,
	}
	return r.db.Create(invalidation).Error
}

func (r *TokenInvalidationRepo) HasInvalidationByJTI(jti string) bool {
	var count int64
	r.db.Model(&models.TokenInvalidation{}).Where("jti = ?", jti).Count(&count)
	return count > 0
}
