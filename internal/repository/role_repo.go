package repository

import (
	"user-management-system/internal/models"

	"gorm.io/gorm"
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}
func (r *RoleRepo) CreateRole(role *models.Role) (*models.Role, error) {
	result := r.db.Omit("updated_at", "updated_by").Create(role)
	if result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}
func (r *RoleRepo) GetRoleByID(id int) (*models.Role, error) {
	var role models.Role
	result := r.db.Preload("Permissions").First(&role, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &role, nil
}
func (r *RoleRepo) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	result := r.db.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}
func (r *RoleRepo) UpdateRole(id int, role *models.Role) (*models.Role, error) {
	if err := r.db.Omit("created_at", "created_by").Save(role).Error; err != nil {
		return nil, err
	}

	for _, perm := range role.Permissions {
		if err := r.db.Save(&perm).Error; err != nil {
			return nil, err
		}
	}

	return role, nil
}

func (r *RoleRepo) GetRoleNamesByIDs(roleIds []int) ([]string, error) {
	var roles []models.Role
	result := r.db.Where("id IN ?", roleIds).Find(&roles).Error
	if result != nil {
		return nil, result
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames, nil
}

func (r *RoleRepo) GetRolesByIDs(roleIDs []int) ([]models.Role, error) {
	var roles []models.Role

	err := r.db.Preload("Permissions").Where("id IN ?", roleIDs).Find(&roles).Error

	return roles, err
}
