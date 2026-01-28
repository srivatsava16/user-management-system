package services

import (
	"errors"
	"user-management-system/internal/repository"
)

type ApprovalService struct {
	userRepo         *repository.UserRepo
	divisionRepo     *repository.DivisionRepo
	businessUnitRepo *repository.BusinessUnitRepo
}

func NewApprovalService(userRepo *repository.UserRepo, divisionRepo *repository.DivisionRepo, businessUnitRepo *repository.BusinessUnitRepo) *ApprovalService {
	return &ApprovalService{
		userRepo:         userRepo,
		divisionRepo:     divisionRepo,
		businessUnitRepo: businessUnitRepo,
	}
}

func (s *ApprovalService) ApproveEntity(entityType string, entityID int, approverEmail string) error {

	switch entityType {
	case "user":
		return s.userRepo.ApproveUser(entityID, approverEmail)
	case "division":
		return s.divisionRepo.ApproveDivision(entityID, approverEmail)
	case "business_unit":
		return s.businessUnitRepo.ApproveBusinessUnit(entityID, approverEmail)
	default:
		return errors.New("invalid entity type for approval")

	}
}

func (s *ApprovalService) GetPendingEntities(entityType string) (interface{}, error) {

	switch entityType {
	case "user":
		return s.userRepo.GetPendingUsers()
	case "division":
		return s.divisionRepo.GetPendingDivisions()
	case "business_unit":
		return s.businessUnitRepo.GetPendingBusinessUnits()
	default:
		return nil, errors.New("invalid entity type for fetching pending approvals")
	}
}
