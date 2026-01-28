package services

import (
	"fmt"
	"user-management-system/internal/models"
	"user-management-system/internal/repository"
)

type DivisionService struct {
	repo *repository.DivisionRepo
}

func NewDivisionService(repo *repository.DivisionRepo) *DivisionService {
	return &DivisionService{repo: repo}
}

func (s *DivisionService) CreateDivision(req *models.Division) (*models.Division, error) {
	if req == nil {
		return nil, fmt.Errorf("Division request cannot be nil")
	}

	division, err := s.repo.CreateDivision(req)
	if err != nil {
		return nil, err
	}
	return division, nil

}

func (s *DivisionService) GetDivisionByID(id int) (*models.Division, error) {
	if id == 0 {
		return nil, fmt.Errorf("Invalid Division ID")
	}
	division, err := s.repo.GetDivisionByID(id)
	if err != nil {
		return nil, err
	}
	return division, nil
}

func (s *DivisionService) GetAllDivisions() ([]models.Division, error) {
	divisions, err := s.repo.GetAllDivisions()
	if err != nil {
		return nil, err
	}
	return divisions, nil
}

func (s *DivisionService) UpdateDivision(id int, req *models.Division) (*models.Division, error) {

	if id == 0 || req == nil {
		return nil, fmt.Errorf("Invalid input for updating Division")
	}
	division, err := s.repo.UpdateDivision(id, req)
	if err != nil {
		return nil, err
	}
	return division, nil

}

func (s *DivisionService) GetDivisionsByBusinessUnit(businessUnitID int) ([]models.Division, error) {
	if businessUnitID == 0 || businessUnitID < 0 {
		return nil, fmt.Errorf("Invalid Business Unit ID")
	}

	return s.repo.GetDivisionsByBusinessUnit(businessUnitID)
}
