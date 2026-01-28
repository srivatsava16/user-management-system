package repository

import (
	"user-management-system/internal/models"

	"gorm.io/gorm"
)

type BusinessUnitRepo struct {
	db *gorm.DB
}

func NewBusinessUnitRepo(db *gorm.DB) *BusinessUnitRepo {
	return &BusinessUnitRepo{db: db}
}

func (r *BusinessUnitRepo) CreateBusinessUnit(businessUnit *models.BusinessUnit) (*models.BusinessUnit, error) {

	result := r.db.Omit("updated_at", "updated_by").Create(businessUnit)
	if result.Error != nil {
		return nil, result.Error
	}
	return businessUnit, nil
}

func (r *BusinessUnitRepo) GetBusinessUnitByID(id int) (*models.BusinessUnit, error) {
	var businessUnit models.BusinessUnit
	result := r.db.First(&businessUnit, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &businessUnit, nil
}

func (r *BusinessUnitRepo) GetAllBusinessUnits() ([]models.BusinessUnit, error) {
	var businessUnits []models.BusinessUnit
	result := r.db.Find(&businessUnits)
	if result.Error != nil {
		return nil, result.Error
	}
	return businessUnits, nil
}

func (r *BusinessUnitRepo) UpdateBusinessUnit(id int, bu *models.BusinessUnit) (*models.BusinessUnit, error) {

	if err := r.db.Omit("created_at", "created_by").Save(bu).Error; err != nil {
		return nil, err
	}
	return bu, nil
}

func (r *BusinessUnitRepo) GetPendingBusinessUnits() ([]models.BusinessUnit, error) {
	var businessUnits []models.BusinessUnit

	err := r.db.Where("approved_by IS NULL").Find(&businessUnits).Error
	return businessUnits, err
}

func (r *BusinessUnitRepo) ApproveBusinessUnit(businessUnitID int, approveEmail string) error {
	return r.db.Model(&models.BusinessUnit{}).Where("id = ?", businessUnitID).Updates(map[string]interface{}{"is_active": true, "approved_by": approveEmail}).Error
}
