package controllers

import (
	"fmt"
	"net/http"
	"user-management-system/internal/models"
	"user-management-system/internal/repository"
	"user-management-system/internal/services"
	"user-management-system/internal/shared"

	"dario.cat/mergo"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	service                  *services.RoleService
	auditService             *services.AuditService
	tokenInvalidationService *services.TokenInvalidationService
	userRepo                 *repository.UserRepo
}

func NewRoleController(service *services.RoleService, auditService *services.AuditService, tokenInvalidationService *services.TokenInvalidationService, userRepo *repository.UserRepo) *RoleController {
	return &RoleController{service: service, auditService: auditService, tokenInvalidationService: tokenInvalidationService, userRepo: userRepo}
}

func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req models.Role

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.Name == "" || req.Description == "" || req.DivisionID <= 0 || req.BusinessUnitID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name, Description, DivisionID, and BusinessUnitID are required"})
		return
	}

	role, err := c.service.CreateRole(&req)
	if err != nil {
		customErr := shared.HandleDatabaseError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": customErr})
		return
	}

	changeBy := shared.GetUserFromContext(ctx)
	c.auditService.LogChange("CREATE", changeBy, fmt.Sprintf("Created Role: %s", role.Name), nil, role)
	ctx.JSON(http.StatusOK, role)
}
func (c *RoleController) GetRoleByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := shared.StrToInt(idstr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	role, err := c.service.GetRoleByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, role)

}
func (c *RoleController) GetAllRoles(ctx *gin.Context) {

	roles, err := c.service.GetAllRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, roles)
}
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := shared.StrToInt(idstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var req models.Role
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	existingRole, err := c.service.GetRoleByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	originalRole := DeepCopyRole(existingRole)

	reqWithoutPermissions := req
	reqWithoutPermissions.Permissions = nil

	err = mergo.Merge(existingRole, reqWithoutPermissions, mergo.WithOverride, mergo.WithTransformers(&shared.BooleanTransformer{}))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if req.Permissions != nil {
		for _, reqPerm := range req.Permissions {
			if reqPerm.ID != 0 {
				for i, existingPerm := range existingRole.Permissions {
					if existingPerm.ID == reqPerm.ID {
						err = mergo.Merge(&existingRole.Permissions[i], reqPerm, mergo.WithOverride, mergo.WithTransformers(&shared.BooleanTransformer{}))
						if err != nil {
							ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
							return
						}
						break
					}
				}
			}
		}
	}

	updatedRole, err := c.service.UpdateRole(id, existingRole)
	if err != nil {
		customErr := shared.HandleDatabaseError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": customErr})
		return
	}

	if c.service.HasRoleChanged(originalRole, updatedRole) {
		fmt.Printf("DEBUG: Role changed detected for Role ID %d\n", id)
		userIDs, err := c.getUsersByRoleID(id)
		fmt.Printf("DEBUG: Found %d users With role ID %d: %v\n", len(userIDs), id, userIDs)

		if err == nil && len(userIDs) > 0 {

			for _, userID := range userIDs {
				fmt.Printf("DEBUG: Invalidating tokens for users with Role ID %d\n", id)
				c.tokenInvalidationService.InvalidateUserTokens(userID, "role_configuration_change")
			}
		}
	}

	changeBy := shared.GetUserFromContext(ctx)
	c.auditService.LogChange("UPDATE", changeBy, fmt.Sprintf("Updated Role: %s", updatedRole.Name), originalRole, updatedRole)

	ctx.JSON(http.StatusOK, updatedRole)
}

func DeepCopyRole(role *models.Role) *models.Role {
	copied := *role

	if role.IsActive != nil {
		val := *role.IsActive
		copied.IsActive = &val
	}

	if role.UpdatedAt != nil {
		val := *role.UpdatedAt
		copied.UpdatedAt = &val
	}

	if role.DataSetAccess != nil {
		copied.DataSetAccess = make(models.JSON)
		for k, v := range role.DataSetAccess {
			copied.DataSetAccess[k] = v
		}
	}
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

func (c *RoleController) getUsersByRoleID(roleID int) ([]int, error) {
	return c.userRepo.GetUsersByRoleID(roleID)
}

func (c *RoleController) GetUserPermissions(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := shared.StrToInt(idstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	permissions, err := c.service.GetUserPermissions(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, permissions)
}
