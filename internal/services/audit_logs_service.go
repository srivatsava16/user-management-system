package services

import (
	"encoding/json"
	"time"
	"user-management-system/internal/models"
	"user-management-system/internal/repository"
)

type AuditService struct {
	repo *repository.AuditRepo
}

func NewAuditService(repo *repository.AuditRepo) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) LogChange(changeType, changeBy, changeDetails string, beforeData, afterData interface{}) (*models.AuditLog, error) {

	var beforeJson, afterJson []byte
	var err error

	if beforeData != nil {
		beforeJson, err = json.Marshal(beforeData)
		if err != nil {
			return nil, err
		}
	}

	if afterData != nil {
		afterJson, err = json.Marshal(afterData)
		if err != nil {
			return nil, err
		}
	}

	audit := &models.AuditLog{
		ChangeType:    changeType,
		ChangeDetails: changeDetails,
		ChangeBy:      changeBy,
		BeforeJson:    beforeJson,
		AfterJson:     afterJson,
		CreatedAt:     time.Now(),
	}

	return s.repo.CreateAuditLog(audit)

}

func (s *AuditService) GetAllAuditLogs() ([]models.AuditLog, error) {
	return s.repo.GetAllAuditLogs()
}
