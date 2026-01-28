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

type BusinessUnitController struct {
	service      *services.BusinessUnitService
	auditService *services.AuditService
}

func NewBusinessUnitController(service *services.BusinessUnitService, auditService *services.AuditService) *BusinessUnitController {
	return &BusinessUnitController{service: service, auditService: auditService}
}

func (c *BusinessUnitController) CreateBusinessUnit(ctx *gin.Context) {

	var req models.BusinessUnit

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.Name == "" || req.Code == "" || req.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name, Code, and Description are required"})
		return
	}

	businessUnit, err := c.service.CreateBusinessUnit(&req)
	if err != nil {
		customErr := shared.HandleDatabaseError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": customErr})
		return
	}
	changeBy := shared.GetUserFromContext(ctx)

	c.auditService.LogChange("CREATE", changeBy, fmt.Sprintf("Created Business Unit: %s", businessUnit.Name), nil, businessUnit)
	ctx.JSON(http.StatusOK, businessUnit)
}

func (c *BusinessUnitController) GetBusinessUnitByID(ctx *gin.Context) {
	idstr := ctx.Param("id")

	id, err := shared.StrToInt(idstr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	businessUnit, err := c.service.GetBusinessUnitByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, businessUnit)

}

func (c *BusinessUnitController) GetAllBusinessUnits(ctx *gin.Context) {

	businessUnits, err := c.service.GetAllBusinessUnits()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, businessUnits)
}

func (c *BusinessUnitController) UpdateBusinessUnit(ctx *gin.Context) {
	idstr := ctx.Param("id")

	id, err := shared.StrToInt(idstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.BusinessUnit

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	businessUnit, err := c.service.GetBusinessUnitByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	originalBusinessUnit := *businessUnit

	if businessUnit.IsActive != nil {
		originalValue := *businessUnit.IsActive
		originalBusinessUnit.IsActive = &originalValue
	}

	mergo.Merge(businessUnit, req, mergo.WithOverride, mergo.WithTransformers(&shared.BooleanTransformer{}))

	updatedBusinessUnit, err := c.service.UpdateBusinessUnit(id, businessUnit)
	if err != nil {
		customErr := shared.HandleDatabaseError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": customErr})
		return
	}
	changeBy := shared.GetUserFromContext(ctx)
	c.auditService.LogChange("UPDATE", changeBy, fmt.Sprintf("Updated Business Unit: %s", businessUnit.Name), originalBusinessUnit, updatedBusinessUnit)
	ctx.JSON(http.StatusOK, updatedBusinessUnit)
}
