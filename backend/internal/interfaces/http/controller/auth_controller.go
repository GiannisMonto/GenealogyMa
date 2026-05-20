package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/internal/infrastructure/cache"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	userService     *service.UserService
	tokenBlacklist *cache.TokenBlacklist
}

// NewAuthController 创建认证控制器
func NewAuthController(userService *service.UserService, tokenBlacklist *cache.TokenBlacklist) *AuthController {
	return &AuthController{
		userService:     userService,
		tokenBlacklist: tokenBlacklist,
	}
}

// RegisterRoutes 注册路由
func (c *AuthController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	authGroup := r.Group("/auth")
	{
		// 公开路由
		authGroup.POST("/register", c.Register)
		authGroup.POST("/login", c.Login)
		authGroup.POST("/refresh", c.RefreshToken)

		// 需要认证的路由
		authRequired := authGroup.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("/logout", c.Logout)
			authRequired.GET("/profile", c.GetProfile)
			authRequired.PUT("/profile", c.UpdateProfile)
			authRequired.POST("/change-password", c.ChangePassword)
		}
	}

	// 用户管理路由
	users := r.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("", auth.RequirePermission(auth.PermissionAdminUser), c.ListUsers)
		users.GET("/:id", auth.RequirePermission(auth.PermissionAdminUser), c.GetUser)
		users.PUT("/:id", auth.RequirePermission(auth.PermissionAdminUser), c.UpdateUser)
		users.DELETE("/:id", auth.RequirePermission(auth.PermissionAdminUser), c.DeleteUser)
		users.POST("/:id/roles", auth.RequirePermission(auth.PermissionAdminUser), c.AssignRole)
		users.DELETE("/:id/roles/:role_code", auth.RequirePermission(auth.PermissionAdminUser), c.RemoveRole)
	}
}

// ===== 认证相关处理函数 =====

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 201 {object} utils.Response{data=service.UserDTO}
// @Failure 400 {object} utils.Response
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req service.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	user, err := c.userService.Register(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithStatus(ctx, http.StatusCreated, user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} utils.Response{data=service.LoginResponse}
// @Failure 401 {object} utils.Response
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	response, err := c.userService.Login(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	utils.Success(ctx, response)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 登出当前用户，将令牌加入黑名单
// @Tags 认证
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	token := auth.ExtractToken(ctx)
	if token != "" {
		// 将令牌加入黑名单，设置过期时间为24小时
		expiresAt := time.Now().Add(24 * time.Hour)
		c.tokenBlacklist.Add(context.Background(), token, expiresAt)
	}

	utils.Success(ctx, gin.H{"message": "登出成功"})
}

// RefreshToken 刷新令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param refresh_token body string true "刷新令牌"
// @Success 200 {object} utils.Response{data=service.LoginResponse}
// @Failure 401 {object} utils.Response
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 检查令牌是否在黑名单中
	if c.tokenBlacklist.IsBlacklisted(context.Background(), req.RefreshToken) {
		utils.Error(ctx, http.StatusUnauthorized, "令牌已失效")
		return
	}

	// 解析刷新令牌获取用户ID
	claims, err := auth.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, "无效的刷新令牌")
		return
	}

	// 获取用户信息
	user, err := c.userService.GetByID(ctx, claims.UserID)
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, "用户不存在")
		return
	}

	// 生成新令牌
	roleCodes := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roleCodes[i] = r.Code
	}

	accessToken, err := auth.GenerateToken(user.ID, user.Username, roleCodes)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "生成刷新令牌失败")
		return
	}

	utils.Success(ctx, service.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    24 * 60 * 60,
		User:         *user,
	})
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=service.UserDTO}
// @Router /auth/profile [get]
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	if userID == 0 {
		utils.Error(ctx, http.StatusUnauthorized, "未登录")
		return
	}

	user, err := c.userService.GetByID(ctx, userID)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "用户不存在")
		return
	}

	utils.Success(ctx, user)
}

// UpdateProfile 更新当前用户信息
// @Summary 更新当前用户信息
// @Description 更新当前登录用户的个人信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.UpdateUserRequest true "用户信息"
// @Success 200 {object} utils.Response{data=service.UserDTO}
// @Router /auth/profile [put]
func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	if userID == 0 {
		utils.Error(ctx, http.StatusUnauthorized, "未登录")
		return
	}

	var req service.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	user, err := c.userService.Update(ctx, userID, &req)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, user)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.ChangePasswordRequest true "密码信息"
// @Success 200 {object} utils.Response
// @Router /auth/change-password [post]
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	if userID == 0 {
		utils.Error(ctx, http.StatusUnauthorized, "未登录")
		return
	}

	var req service.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := c.userService.ChangePassword(ctx, userID, &req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "密码修改成功"})
}

// ===== 用户管理处理函数 =====

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 搜索和分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param keyword query string false "关键字"
// @Param status query string false "状态"
// @Param role_id query int false "角色ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Success 200 {object} utils.Response{data=[]service.UserDTO}
// @Router /users [get]
func (c *AuthController) ListUsers(ctx *gin.Context) {
	var req service.SearchUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	users, total, err := c.userService.Search(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	utils.PageSuccess(ctx, users, total, req.Page, req.PageSize)
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response{data=service.UserDTO}
// @Router /users/{id} [get]
func (c *AuthController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	user, err := c.userService.GetByID(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "用户不存在")
		return
	}

	utils.Success(ctx, user)
}

// UpdateUser 更新用户
// @Summary 更新用户信息
// @Description 管理员更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Param request body service.UpdateUserRequest true "用户信息"
// @Success 200 {object} utils.Response{data=service.UserDTO}
// @Router /users/{id} [put]
func (c *AuthController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	user, err := c.userService.Update(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, user)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response
// @Router /users/{id} [delete]
func (c *AuthController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.userService.Delete(ctx, id); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "删除成功"})
}

// AssignRole 分配角色
// @Summary 分配角色给用户
// @Description 给指定用户分配角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Param request body service.AssignRoleRequest true "角色信息"
// @Success 200 {object} utils.Response
// @Router /users/{id}/roles [post]
func (c *AuthController) AssignRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.AssignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	assignedBy := auth.GetUserID(ctx)
	if err := c.userService.AssignRole(ctx, id, assignedBy, req.RoleCode); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "角色分配成功"})
}

// RemoveRole 移除用户角色
// @Summary 移除用户角色
// @Description 移除指定用户的角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Param role_code path string true "角色代码"
// @Success 200 {object} utils.Response
// @Router /users/{id}/roles/{role_code} [delete]
func (c *AuthController) RemoveRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	roleCode := ctx.Param("role_code")
	if roleCode == "" {
		utils.Error(ctx, http.StatusBadRequest, "角色代码不能为空")
		return
	}

	if err := c.userService.RemoveRole(ctx, id, roleCode); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "角色移除成功"})
}


