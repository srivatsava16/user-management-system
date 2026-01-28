package services

import (
	"fmt"
	"user-management-system/internal/models"
	"user-management-system/internal/repository"
)

type BusinessUnitService struct {
	repo *repository.BusinessUnitRepo
}

func NewBusinessUnitService(repo *repository.BusinessUnitRepo) *BusinessUnitService {
	return &BusinessUnitService{repo: repo}
}

func (s *BusinessUnitService) CreateBusinessUnit(req *models.BusinessUnit) (*models.BusinessUnit, error) {
	if req == nil {
		return nil, fmt.Errorf("BusinessUnit request cannot be nil")
	}

	businessUnit, err := s.repo.CreateBusinessUnit(req)
	if err != nil {
		return nil, err
	}

	return businessUnit, nil

}

func (s *BusinessUnitService) GetBusinessUnitByID(id int) (*models.BusinessUnit, error) {
	if id == 0 {
		return nil, fmt.Errorf("Invalid BusinessUnit ID")
	}
	businessUnit, err := s.repo.GetBusinessUnitByID(id)
	if err != nil {
		return nil, err
	}

	return businessUnit, nil
}

func (s *BusinessUnitService) GetAllBusinessUnits() ([]models.BusinessUnit, error) {
	businessUnits, err := s.repo.GetAllBusinessUnits()
	if err != nil {
		return nil, err
	}
	return businessUnits, nil
}

func (s *BusinessUnitService) UpdateBusinessUnit(id int, req *models.BusinessUnit) (*models.BusinessUnit, error) {

	if id == 0 || req == nil {
		return nil, fmt.Errorf("Invalid input for updating BusinessUnit")
	}

	BusinessUnit, err := s.repo.UpdateBusinessUnit(id, req)
	if err != nil {
		return nil, err
	}
	return BusinessUnit, nil

}
