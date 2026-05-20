package controller

import (
	"net/http"
	"strconv"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// RBACController 角色权限管理控制器
type RBACController struct {
	userService *service.UserService
}

// NewRBACController 创建角色权限管理控制器
func NewRBACController(userService *service.UserService) *RBACController {
	return &RBACController{userService: userService}
}

// RegisterRoutes 注册路由
func (c *RBACController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	roles := r.Group("/roles")
	roles.Use(authMiddleware)
	{
		roles.GET("", auth.RequirePermission(auth.PermissionRoleRead), c.ListRoles)
		roles.GET("/:id", auth.RequirePermission(auth.PermissionRoleRead), c.GetRole)
		roles.POST("", auth.RequirePermission(auth.PermissionRoleWrite), c.CreateRole)
		roles.PUT("/:id", auth.RequirePermission(auth.PermissionRoleWrite), c.UpdateRole)
		roles.DELETE("/:id", auth.RequirePermission(auth.PermissionRoleDelete), c.DeleteRole)
		roles.GET("/:id/permissions", auth.RequirePermission(auth.PermissionRoleRead), c.GetRolePermissions)
		roles.POST("/:id/permissions", auth.RequirePermission(auth.PermissionRoleWrite), c.AssignPermissions)
		roles.DELETE("/:id/permissions/:permission_id", auth.RequirePermission(auth.PermissionRoleWrite), c.RemovePermission)
	}

	permissions := r.Group("/permissions")
	permissions.Use(authMiddleware)
	{
		permissions.GET("", auth.RequirePermission(auth.PermissionPermRead), c.ListPermissions)
		permissions.POST("", auth.RequirePermission(auth.PermissionPermWrite), c.CreatePermission)
		permissions.PUT("/:id", auth.RequirePermission(auth.PermissionPermWrite), c.UpdatePermission)
		permissions.DELETE("/:id", auth.RequirePermission(auth.PermissionPermDelete), c.DeletePermission)
	}
}

// ListRoles 获取角色列表
func (c *RBACController) ListRoles(ctx *gin.Context) {
	roles, err := c.userService.GetAllRoles(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取角色列表失败")
		return
	}

	utils.Success(ctx, roles)
}

// GetRole 获取角色详情
func (c *RBACController) GetRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}

	role, err := c.userService.GetRoleByID(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(ctx, role)
}

// CreateRole 创建角色
func (c *RBACController) CreateRole(ctx *gin.Context) {
	var req service.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	role, err := c.userService.CreateRole(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithStatus(ctx, http.StatusCreated, role)
}

// UpdateRole 更新角色
func (c *RBACController) UpdateRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req service.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	role, err := c.userService.UpdateRole(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, role)
}

// DeleteRole 删除角色
func (c *RBACController) DeleteRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}

	if err := c.userService.DeleteRole(ctx, id); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "删除成功"})
}

// GetRolePermissions 获取角色权限
func (c *RBACController) GetRolePermissions(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}

	perms, err := c.userService.GetRolePermissions(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(ctx, perms)
}

// AssignPermissions 分配权限给角色
func (c *RBACController) AssignPermissions(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req service.AssignPermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := c.userService.AssignPermissionsToRole(ctx, id, req.PermissionIDs); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "权限分配成功"})
}

// RemovePermission 移除角色权限
func (c *RBACController) RemovePermission(ctx *gin.Context) {
	roleID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}

	permID, err := strconv.ParseInt(ctx.Param("permission_id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的权限ID")
		return
	}

	if err := c.userService.RemovePermissionFromRole(ctx, roleID, permID); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "权限移除成功"})
}

// ListPermissions 获取权限列表
func (c *RBACController) ListPermissions(ctx *gin.Context) {
	perms, err := c.userService.GetAllPermissions(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取权限列表失败")
		return
	}

	utils.Success(ctx, perms)
}

// CreatePermission 创建权限
func (c *RBACController) CreatePermission(ctx *gin.Context) {
	var req service.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	perm, err := c.userService.CreatePermission(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithStatus(ctx, http.StatusCreated, perm)
}

// UpdatePermission 更新权限
func (c *RBACController) UpdatePermission(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的权限ID")
		return
	}

	var req service.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	perm, err := c.userService.UpdatePermission(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, perm)
}

// DeletePermission 删除权限
func (c *RBACController) DeletePermission(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的权限ID")
		return
	}

	if err := c.userService.DeletePermission(ctx, id); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "删除成功"})
}
