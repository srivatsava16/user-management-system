package controllers

import (
	"fmt"
	"net/http"
	"user-management-system/internal/models"
	"user-management-system/internal/services"
	"user-management-system/internal/shared"

	"dario.cat/mergo"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service                  *services.UserService
	auditService             *services.AuditService
	tokenInvalidationService *services.TokenInvalidationService
}

func NewUserController(service *services.UserService, auditService *services.AuditService, tokenInvalidationService *services.TokenInvalidationService) *UserController {
	return &UserController{service: service, auditService: auditService, tokenInvalidationService: tokenInvalidationService}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req models.User
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.UserName == "" || req.Email == "" || req.Password == "" || req.DivisionID == 0 || req.BusinessUnitID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "UserName, Email, Password, DivisionID, and BusinessUnitID are required"})
		return
	}
	user, err := c.service.CreateUser(&req)
	if err != nil {
		customErr := shared.HandleDatabaseError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": customErr})
		return
	}
	changeBy := shared.GetUserFromContext(ctx)

	c.auditService.LogChange("CREATE", changeBy, fmt.Sprintf("Created User: %s", user.UserName), nil, user)
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := shared.StrToInt(idstr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	user, err := c.service.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := shared.StrToInt(idstr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var req models.User
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	existingUser, err := c.service.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	originalUser := *existingUser

	err = mergo.Merge(existingUser, req, mergo.WithOverride, mergo.WithTransformers(&shared.BooleanTransformer{}))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updatedUser, err := c.service.UpdateUser(id, existingUser)
	if err != nil {
		customErr := shared.HandleDatabaseError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": customErr})
		return
	}

	if !shared.EqualIntArrays(originalUser.RoleIds, updatedUser.RoleIds) {
		err = c.tokenInvalidationService.InvalidateUserTokens(id, "role_change")

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate user tokens"})
			return
		}
	}

	changeBy := shared.GetUserFromContext(ctx)

	c.auditService.LogChange("UPDATE", changeBy, fmt.Sprintf("Updated User: %s", updatedUser.UserName), originalUser, updatedUser)
	ctx.JSON(http.StatusOK, updatedUser)
}

func (c *UserController) GetUsersByDivision(ctx *gin.Context) {
	idstr := ctx.Param("id")
	divisionID, err := shared.StrToInt(idstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Division ID"})
		return
	}

	users, err := c.service.GetUsersByDivision(divisionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetUsersByBusinessUnit(ctx *gin.Context) {
	idstr := ctx.Param("id")
	businessUnitID, err := shared.StrToInt(idstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Business Unit ID"})
		return
	}

	users, err := c.service.GetUsersByBusinessUnit(businessUnitID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
