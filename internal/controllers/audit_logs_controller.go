package controllers

import (
	"net/http"
	"user-management-system/internal/services"

	"github.com/gin-gonic/gin"
)

type AuditController struct {
	service *services.AuditService
}

func NewAuditController(service *services.AuditService) *AuditController {
	return &AuditController{service: service}
}

func (c *AuditController) GetAllAuditLogs(ctx *gin.Context) {
	auditLogs, err := c.service.GetAllAuditLogs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, auditLogs)
}
