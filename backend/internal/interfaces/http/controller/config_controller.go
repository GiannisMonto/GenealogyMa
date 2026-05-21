package controller

import (
	"net/http"
	"strconv"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ConfigController 配置控制器
type ConfigController struct {
	configService *service.ConfigService
}

// NewConfigController 创建配置控制器
func NewConfigController(configService *service.ConfigService) *ConfigController {
	return &ConfigController{configService: configService}
}

// RegisterRoutes 注册路由
func (c *ConfigController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	configs := r.Group("/configs")
	{
		// 公开路由 - 获取配置
		configs.GET("", c.List)
		configs.GET("/:id", c.Get)
		configs.GET("/group/:group", c.ListByGroup)
		configs.GET("/key/:group/:key", c.GetByKey)

		// 需要认证的路由
		authRequired := configs.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionAdminConfig), c.Create)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionAdminConfig), c.Update)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionAdminConfig), c.Delete)
		}
	}
}

// List 获取所有配置
// @Summary 获取所有配置
// @Description 获取所有配置列表
// @Tags 配置
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]service.ConfigDTO}
// @Router /configs [get]
func (c *ConfigController) List(ctx *gin.Context) {
	configs, err := c.configService.ListConfigs(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取配置列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, configs)
}

// Get 获取配置详情
// @Summary 获取配置详情
// @Description 根据ID获取配置详细信息
// @Tags 配置
// @Accept json
// @Produce json
// @Param id path int true "配置ID"
// @Success 200 {object} utils.Response{data=service.ConfigDTO}
// @Router /configs/{id} [get]
func (c *ConfigController) Get(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	config, err := c.configService.GetConfig(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取配置详情失败: "+err.Error())
		return
	}
	if config == nil {
		utils.Error(ctx, http.StatusNotFound, "配置不存在")
		return
	}
	utils.Success(ctx, config)
}

// ListByGroup 获取指定分组的配置
// @Summary 获取指定分组的配置
// @Description 获取指定分组的配置列表
// @Tags 配置
// @Accept json
// @Produce json
// @Param group path string true "配置分组"
// @Success 200 {object} utils.Response{data=[]service.ConfigDTO}
// @Router /configs/group/{group} [get]
func (c *ConfigController) ListByGroup(ctx *gin.Context) {
	group := ctx.Param("group")
	if group == "" {
		utils.Error(ctx, http.StatusBadRequest, "分组不能为空")
		return
	}

	configs, err := c.configService.ListConfigsByGroup(ctx, group)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取配置列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, configs)
}

// GetByKey 根据group和key获取配置
// @Summary 根据group和key获取配置
// @Description 根据配置分组和键名获取配置
// @Tags 配置
// @Accept json
// @Produce json
// @Param group path string true "配置分组"
// @Param key path string true "配置键"
// @Success 200 {object} utils.Response{data=service.ConfigDTO}
// @Router /configs/key/{group}/{key} [get]
func (c *ConfigController) GetByKey(ctx *gin.Context) {
	group := ctx.Param("group")
	key := ctx.Param("key")
	if group == "" || key == "" {
		utils.Error(ctx, http.StatusBadRequest, "分组和键不能为空")
		return
	}

	config, err := c.configService.GetConfigByKey(ctx, group, key)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取配置失败: "+err.Error())
		return
	}
	if config == nil {
		utils.Error(ctx, http.StatusNotFound, "配置不存在")
		return
	}
	utils.Success(ctx, config)
}

// Create 创建配置
// @Summary 创建配置
// @Description 创建新配置
// @Tags 配置
// @Accept json
// @Produce json
// @Param config body service.CreateConfigRequest true "配置信息"
// @Success 200 {object} utils.Response{data=service.ConfigDTO}
// @Security BearerAuth
// @Router /configs [post]
func (c *ConfigController) Create(ctx *gin.Context) {
	var req service.CreateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	config, err := c.configService.CreateConfig(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, config)
}

// Update 更新配置
// @Summary 更新配置
// @Description 更新配置信息
// @Tags 配置
// @Accept json
// @Produce json
// @Param id path int true "配置ID"
// @Param config body service.UpdateConfigRequest true "更新的配置信息"
// @Success 200 {object} utils.Response{data=service.ConfigDTO}
// @Security BearerAuth
// @Router /configs/{id} [put]
func (c *ConfigController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	config, err := c.configService.UpdateConfig(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if config == nil {
		utils.Error(ctx, http.StatusNotFound, "配置不存在")
		return
	}
	utils.Success(ctx, config)
}

// Delete 删除配置
// @Summary 删除配置
// @Description 删除指定配置
// @Tags 配置
// @Accept json
// @Produce json
// @Param id path int true "配置ID"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /configs/{id} [delete]
func (c *ConfigController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.configService.DeleteConfig(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}
