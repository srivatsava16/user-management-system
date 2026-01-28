package repository

import (
	"user-management-system/internal/models"

	"gorm.io/gorm"
)

type AuditRepo struct {
	db *gorm.DB
}

func NewAuditRepo(db *gorm.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) CreateAuditLog(log *models.AuditLog) (*models.AuditLog, error) {
	err := r.db.Create(log).Error
	if err != nil {
		return nil, err
	}
	return log, nil
}

func (r *AuditRepo) GetAllAuditLogs() ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}
