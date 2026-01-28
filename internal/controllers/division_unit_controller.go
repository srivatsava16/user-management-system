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

type DivisionController struct {
	service      *services.DivisionService
	auditService *services.AuditService
}

func NewDivisionController(service *services.DivisionService, auditService *services.AuditService) *DivisionController {
	return &DivisionController{service: service, auditService: auditService}
}

func (c *DivisionController) CreateDivision(ctx *gin.Context) {

	var req models.Division

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.Name == "" || req.Description == "" || req.BusinessUnitID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name, Code, Description, and BusinessUnitID are required"})
		return
	}

	division, err := c.service.CreateDivision(&req)
	if err != nil {
		userMessage := shared.HandleDatabaseError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": userMessage})
		return
	}

	changeBy := shared.GetUserFromContext(ctx)
	c.auditService.LogChange("CREATE", changeBy, fmt.Sprintf("Created Division: %s", division.Name), nil, division)

	ctx.JSON(http.StatusOK, division)
}

func (c *DivisionController) GetDivisionByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := shared.StrToInt(idstr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	division, err := c.service.GetDivisionByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, division)

}

func (c *DivisionController) GetAllDivisions(ctx *gin.Context) {

	divisions, err := c.service.GetAllDivisions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, divisions)
}
func (c *DivisionController) UpdateDivision(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := shared.StrToInt(idstr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	existingDivision, err := c.service.GetDivisionByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var req models.Division
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	originalDivision := *existingDivision
	if existingDivision.IsActive != nil {
		originalValue := *existingDivision.IsActive
		originalDivision.IsActive = &originalValue
	}

	err = mergo.Merge(existingDivision, req, mergo.WithOverride, mergo.WithTransformers(&shared.BooleanTransformer{}))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updatedDivision, err := c.service.UpdateDivision(id, existingDivision)
	if err != nil {
		customErr := shared.HandleDatabaseError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": customErr})
		return
	}

	changeBy := shared.GetUserFromContext(ctx)
	c.auditService.LogChange("UPDATE", changeBy, fmt.Sprintf("Updated Division: %s", updatedDivision.Name), originalDivision, updatedDivision)

	ctx.JSON(http.StatusOK, updatedDivision)
}

func (c *DivisionController) GetDivisionsByBusinessUnit(ctx *gin.Context) {
	idstr := ctx.Param("id")
	businessUnitID, err := shared.StrToInt(idstr)
	if err != nil || businessUnitID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Business Unit ID"})
		return
	}

	divisions, err := c.service.GetDivisionsByBusinessUnit(businessUnitID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, divisions)
}
