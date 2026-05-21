package controller

import (
	"net/http"
	"strconv"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// BackupController 备份控制器
type BackupController struct {
	backupService *service.BackupService
}

// NewBackupController 创建备份控制器
func NewBackupController(backupService *service.BackupService) *BackupController {
	return &BackupController{backupService: backupService}
}

// RegisterRoutes 注册路由
func (c *BackupController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	backups := r.Group("/backups")
	{
		backups.Use(authMiddleware)
		{
			backups.GET("", c.List)
			backups.GET("/:id", c.Get)
			backups.POST("", auth.RequirePermission(auth.PermissionAdminConfig), c.Create)
			backups.POST("/:id/restore", auth.RequirePermission(auth.PermissionAdminConfig), c.Restore)
			backups.DELETE("/:id", auth.RequirePermission(auth.PermissionAdminConfig), c.Delete)
			backups.DELETE("/clean-old", auth.RequirePermission(auth.PermissionAdminConfig), c.CleanOldBackups)
		}
	}
}

// List 获取备份列表
// @Summary 获取备份列表
// @Description 获取所有备份记录
// @Tags 备份
// @Accept json
// @Produce json
// @Param type query string false "备份类型"
// @Param status query string false "备份状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} utils.Response{data=service.BackupListDTO}
// @Router /backups [get]
func (c *BackupController) List(ctx *gin.Context) {
	var req service.BackupFilterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	result, err := c.backupService.ListBackups(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取备份列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, result)
}

// Get 获取备份详情
// @Summary 获取备份详情
// @Description 根据ID获取备份详细信息
// @Tags 备份
// @Accept json
// @Produce json
// @Param id path int true "备份ID"
// @Success 200 {object} utils.Response{data=service.BackupDTO}
// @Router /backups/{id} [get]
func (c *BackupController) Get(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	backup, err := c.backupService.GetBackup(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取备份详情失败: "+err.Error())
		return
	}
	if backup == nil {
		utils.Error(ctx, http.StatusNotFound, "备份不存在")
		return
	}
	utils.Success(ctx, backup)
}

// Create 创建备份
// @Summary 创建备份
// @Description 创建新备份
// @Tags 备份
// @Accept json
// @Produce json
// @Param backup body service.CreateBackupRequest true "备份信息"
// @Success 200 {object} utils.Response{data=service.BackupDTO}
// @Security BearerAuth
// @Router /backups [post]
func (c *BackupController) Create(ctx *gin.Context) {
	var req service.CreateBackupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	backup, err := c.backupService.CreateBackup(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建备份失败: "+err.Error())
		return
	}
	utils.Success(ctx, backup)
}

// Restore 恢复备份
// @Summary 恢复备份
// @Description 恢复指定备份
// @Tags 备份
// @Accept json
// @Produce json
// @Param id path int true "备份ID"
// @Param request body service.RestoreBackupRequest true "恢复信息"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /backups/{id}/restore [post]
func (c *BackupController) Restore(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.RestoreBackupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	req.BackupID = id
	if err := c.backupService.RestoreBackup(ctx, &req); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "恢复备份失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "备份恢复成功", nil)
}

// Delete 删除备份
// @Summary 删除备份
// @Description 删除指定备份
// @Tags 备份
// @Accept json
// @Produce json
// @Param id path int true "备份ID"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /backups/{id} [delete]
func (c *BackupController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.backupService.DeleteBackup(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除备份失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// CleanOldBackups 清理旧备份
// @Summary 清理旧备份
// @Description 清理指定天数之前的备份
// @Tags 备份
// @Accept json
// @Produce json
// @Param request body service.CleanOldBackupsRequest true "清理参数"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /backups/clean-old [delete]
func (c *BackupController) CleanOldBackups(ctx *gin.Context) {
	var req service.CleanOldBackupsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	count, err := c.backupService.CleanOldBackups(ctx, req.Days)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "清理旧备份失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "已清理旧备份", map[string]int64{"deleted_count": count})
}