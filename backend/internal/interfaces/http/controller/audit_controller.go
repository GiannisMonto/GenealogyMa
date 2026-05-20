package controller

import (
	"net/http"
	"strconv"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuditController 审计控制器
type AuditController struct {
	auditService *service.AuditService
}

// NewAuditController 创建审计控制器
func NewAuditController(auditService *service.AuditService) *AuditController {
	return &AuditController{auditService: auditService}
}

// RegisterRoutes 注册路由
func (c *AuditController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	logs := r.Group("/audit-logs")
	{
		logs.GET("", c.ListLogs)
		logs.GET("/:id", c.GetLogByID)
		logs.GET("/export", c.ExportLogs)

		authRequired := logs.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionAdminConfig), c.LogAction)
			authRequired.DELETE("/clean", auth.RequirePermission(auth.PermissionAdminConfig), c.CleanOldLogs)
		}
	}
}

// ListLogs 获取审计日志列表
func (c *AuditController) ListLogs(ctx *gin.Context) {
	var req service.AuditFilterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	result, err := c.auditService.ListLogs(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取审计日志失败: "+err.Error())
		return
	}
	utils.SuccessWithPage(ctx, result.Data, result.Total)
}

// GetLogByID 获取审计日志详情
func (c *AuditController) GetLogByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	log, err := c.auditService.GetLogByID(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取审计日志详情失败: "+err.Error())
		return
	}
	if log == nil {
		utils.Error(ctx, http.StatusNotFound, "审计日志不存在")
		return
	}
	utils.Success(ctx, log)
}

// LogAction 记录审计日志
func (c *AuditController) LogAction(ctx *gin.Context) {
	var req service.LogActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := c.auditService.LogAction(ctx, &req); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "记录审计日志失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "记录成功", nil)
}

// CleanOldLogs 清理旧日志
func (c *AuditController) CleanOldLogs(ctx *gin.Context) {
	var req service.CleanOldLogsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	count, err := c.auditService.CleanOldLogs(ctx, req.Days)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "清理旧日志失败: "+err.Error())
		return
	}
	utils.Success(ctx, gin.H{"deleted_count": count})
}

// ExportLogs 导出审计日志
func (c *AuditController) ExportLogs(ctx *gin.Context) {
	var req service.AuditFilterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := c.auditService.ListLogs(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "导出审计日志失败: "+err.Error())
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Disposition", "attachment; filename=audit_logs.json")
	ctx.JSON(http.StatusOK, result.Data)
}
