package services

import (
	"fmt"
	"user-management-system/internal/models"
	"user-management-system/internal/repository"
	"user-management-system/internal/shared"
)

type RoleService struct {
	repo     *repository.RoleRepo
	UserRepo *repository.UserRepo
}

func NewRoleService(repo *repository.RoleRepo, userRepo *repository.UserRepo) *RoleService {
	return &RoleService{repo: repo, UserRepo: userRepo}
}

func (s *RoleService) CreateRole(req *models.Role) (*models.Role, error) {
	if req == nil {
		return nil, fmt.Errorf("Role request cannot be nil")
	}

	role, err := s.repo.CreateRole(req)
	if err != nil {
		return nil, err
	}
	return role, nil
}
func (s *RoleService) GetRoleByID(id int) (*models.Role, error) {
	if id == 0 {
		return nil, fmt.Errorf("Invalid Role ID")
	}
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	roles, err := s.repo.GetAllRoles()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
func (s *RoleService) UpdateRole(id int, req *models.Role) (*models.Role, error) {
	if id == 0 || req == nil {
		return nil, fmt.Errorf("Invalid input for updating Role")
	}
	role, err := s.repo.UpdateRole(id, req)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// In internal/shared/deep_copy.go or similar
func DeepCopyRole(role *models.Role) *models.Role {
	copied := *role

	// Deep copy pointer fields
	if role.IsActive != nil {
		val := *role.IsActive
		copied.IsActive = &val
	}
	if role.CreatedBy != nil {
		val := *role.CreatedBy
		copied.CreatedBy = &val
	}
	if role.UpdatedBy != nil {
		val := *role.UpdatedBy
		copied.UpdatedBy = &val
	}
	if role.UpdatedAt != nil {
		val := *role.UpdatedAt
		copied.UpdatedAt = &val
	}

	// Deep copy permissions
	copied.Permissions = make([]models.ModulePermissions, len(role.Permissions))
	for i, perm := range role.Permissions {
		copied.Permissions[i] = deepCopyPermission(perm)
	}

	return &copied
}

func deepCopyPermission(perm models.ModulePermissions) models.ModulePermissions {
	copied := perm

	if perm.WithinDivision != nil {
		copied.WithinDivision = make(models.JSON)
		for k, v := range perm.WithinDivision {
			copied.WithinDivision[k] = v
		}
	}

	if perm.AcrossDivisions != nil {
		copied.AcrossDivisions = make(models.JSON)
		for k, v := range perm.AcrossDivisions {
			copied.AcrossDivisions[k] = v
		}
	}

	if perm.WithinBusinessUnit != nil {
		copied.WithinBusinessUnit = make(models.JSON)
		for k, v := range perm.WithinBusinessUnit {
			copied.WithinBusinessUnit[k] = v
		}
	}

	if perm.AcrossBusinessUnits != nil {
		copied.AcrossBusinessUnits = make(models.JSON)
		for k, v := range perm.AcrossBusinessUnits {
			copied.AcrossBusinessUnits[k] = v
		}
	}

	if perm.SpecialPermissions != nil {
		copied.SpecialPermissions = make(models.JSON)
		for k, v := range perm.SpecialPermissions {
			copied.SpecialPermissions[k] = v
		}
	}

	return copied
}

func (s *RoleService) GetRoleNamesByIDs(roleIds []int) ([]string, error) {
	if roleIds == nil {
		return nil, fmt.Errorf("Role IDs cannot be nil or empty")
	}
	return s.repo.GetRoleNamesByIDs(roleIds)
}

func (s *RoleService) HasRoleChanged(original, updated *models.Role) bool {

	if original.Name != updated.Name || original.Description != updated.Description {
		return true
	}

	if original.BusinessUnitID != updated.BusinessUnitID || original.DivisionID != updated.DivisionID {
		return true
	}

	if !shared.CompareJSON(original.DataSetAccess, updated.DataSetAccess) {
		return true
	}

	if (original.IsActive == nil) != (updated.IsActive == nil) {
		return true
	}
	if original.IsActive != nil && updated.IsActive != nil && *original.IsActive != *updated.IsActive {
		return true
	}

	return s.hasPermissionsChanged(original.Permissions, updated.Permissions)
}

func (s *RoleService) hasPermissionsChanged(original, updated []models.ModulePermissions) bool {

	if len(original) != len(updated) {
		return true
	}

	for i, origPerm := range original {
		if i >= len(updated) {
			return true
		}

		updPerm := updated[i]
		if !shared.CompareJSON(origPerm.WithinDivision, updPerm.WithinDivision) ||
			!shared.CompareJSON(origPerm.AcrossDivisions, updPerm.AcrossDivisions) ||
			!shared.CompareJSON(origPerm.WithinBusinessUnit, updPerm.WithinBusinessUnit) ||
			!shared.CompareJSON(origPerm.AcrossBusinessUnits, updPerm.AcrossBusinessUnits) ||
			!shared.CompareJSON(origPerm.SpecialPermissions, updPerm.SpecialPermissions) {
			return true
		}
	}
	return false
}

// func (s *RoleService) compareJSON(a, b models.JSON) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}

// 	for k, v := range a {
// 		if b[k] != v {
// 			return false
// 		}
// 	}
// 	return true
// }

func (s *RoleService) GetUsersByRoleID(roleID int) ([]int, error) {
	return s.UserRepo.GetUsersByRoleID(roleID)
}

func (s *RoleService) GetRolesByIDs(roleIDs []int) ([]models.Role, error) {
	if len(roleIDs) == 0 {
		return []models.Role{}, nil
	}

	roles, err := s.repo.GetRolesByIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *RoleService) GetUserPermissions(userID int) (*models.UserPermisssionsResponse, error) {

	if userID == 0 {
		return nil, fmt.Errorf("Invalid User ID")
	}

	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	roles, err := s.GetRolesByIDs(user.RoleIds)
	if err != nil {
		return nil, err
	}

	return &models.UserPermisssionsResponse{
		UserID:    user.ID,
		UserEmail: user.Email,
		Roles:     roles,
	}, nil

}
