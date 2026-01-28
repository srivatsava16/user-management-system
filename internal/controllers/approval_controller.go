package controllers

import (
	"fmt"
	"net/http"
	"user-management-system/internal/services"
	"user-management-system/internal/shared"

	"github.com/gin-gonic/gin"
)

type ApprovalController struct {
	approvalService *services.ApprovalService
	auditService    *services.AuditService
}

type ApprovalRequest struct {
	EntityType string `json:"entity_type" binding:"required"`
	EntityID   int    `json:"entity_id" binding:"required"`
}

var allowedEntityTypes = map[string]bool{
	"user":          true,
	"division":      true,
	"business_unit": true,
}

func NewApprovalController(approvalService *services.ApprovalService, auditService *services.AuditService) *ApprovalController {
	return &ApprovalController{approvalService: approvalService, auditService: auditService}
}

func (c *ApprovalController) ApproveEntity(ctx *gin.Context) {
	var req ApprovalRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if !allowedEntityTypes[req.EntityType] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity_type"})
		return
	}

	if req.EntityID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity_id"})
		return
	}

	approverEmail := shared.GetUserFromContext(ctx)
	if approverEmail == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := c.approvalService.ApproveEntity(req.EntityType, req.EntityID, approverEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	changeBy := shared.GetUserFromContext(ctx)

	c.auditService.LogChange("APPROVE", changeBy, fmt.Sprintf("Approved %s with ID %d", req.EntityType, req.EntityID), nil, map[string]interface{}{"entity_type": req.EntityType, "entity_id": req.EntityID, "approver": approverEmail})

	ctx.JSON(http.StatusOK, gin.H{"message": "Entity approved successfully"})
}

func (c *ApprovalController) GetPendingEntities(ctx *gin.Context) {

	entityType := ctx.Param("entity_type")

	if entityType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "entity_type is required"})
		return
	}

	if !allowedEntityTypes[entityType] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity_type"})
		return
	}

	entities, err := c.approvalService.GetPendingEntities(entityType)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"pending_entities": entities})
}
