package repository

import (
	"user-management-system/internal/models"

	"gorm.io/gorm"
)

type DivisionRepo struct {
	db *gorm.DB
}

func NewDivisionRepo(db *gorm.DB) *DivisionRepo {
	return &DivisionRepo{db: db}
}

func (r *DivisionRepo) CreateDivision(division *models.Division) (*models.Division, error) {

	result := r.db.Omit("updated_at", "updated_by").Create(division)
	if result.Error != nil {
		return nil, result.Error
	}
	return division, nil

}
func (r *DivisionRepo) GetDivisionByID(id int) (*models.Division, error) {
	var division models.Division
	result := r.db.First(&division, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &division, nil
}
func (r *DivisionRepo) GetAllDivisions() ([]models.Division, error) {
	var divisions []models.Division
	result := r.db.Find(&divisions)
	if result.Error != nil {
		return nil, result.Error
	}
	return divisions, nil
}
func (r *DivisionRepo) UpdateDivision(id int, div *models.Division) (*models.Division, error) {
	if err := r.db.Omit("created_at", "created_by").Save(div).Error; err != nil {
		return nil, err
	}
	return div, nil
}

func (r *DivisionRepo) GetDivisionsByBusinessUnit(businessUnitID int) ([]models.Division, error) {
	var divisions []models.Division
	result := r.db.Where("business_unit_id = ?", businessUnitID).Find(&divisions)
	if result.Error != nil {
		return nil, result.Error
	}
	return divisions, nil
}

func (r *DivisionRepo) GetPendingDivisions() ([]models.Division, error) {

	var divisions []models.Division

	err := r.db.Where("approved_by IS NULL").Find(&divisions).Error

	return divisions, err
}

func (r *DivisionRepo) ApproveDivision(divisionID int, approveEmail string) error {
	return r.db.Model(&models.Division{}).Where("id = ?", divisionID).Updates(map[string]interface{}{"is_active": true, "approved_by": approveEmail}).Error
}
