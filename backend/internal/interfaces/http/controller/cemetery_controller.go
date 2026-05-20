package controller

import (
	"net/http"
	"strconv"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// CemeteryController 墓园控制器
type CemeteryController struct {
	cemeteryService *service.CemeteryService
}

// NewCemeteryController 创建墓园控制器
func NewCemeteryController(cemeteryService *service.CemeteryService) *CemeteryController {
	return &CemeteryController{cemeteryService: cemeteryService}
}

// RegisterRoutes 注册路由
func (c *CemeteryController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// 墓园路由
	cemeteries := r.Group("/cemeteries")
	{
		// 公开路由
		cemeteries.GET("", c.List)
		cemeteries.GET("/:id", c.Get)

		// 需要认证的路由
		authRequired := cemeteries.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionGenealogyWrite), c.Create)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.Update)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.Delete)
			authRequired.GET("/:id/graves", c.ListGrave)
		}
	}

	// 墓位路由
	graves := r.Group("/graves")
	graves.Use(authMiddleware)
	{
		graves.GET("/:id", c.GetGrave)
		graves.POST("", auth.RequirePermission(auth.PermissionGenealogyWrite), c.CreateGrave)
		graves.PUT("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.UpdateGrave)
		graves.DELETE("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.DeleteGrave)
	}
}

// List 获取墓园列表
// @Summary 获取墓园列表
// @Description 获取所有墓园列表
// @Tags 墓园
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]service.CemeteryDTO}
// @Router /cemeteries [get]
func (c *CemeteryController) List(ctx *gin.Context) {
	cemeteries, err := c.cemeteryService.ListCemeteries(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取墓园列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, cemeteries)
}

// Get 获取墓园详情
// @Summary 获取墓园详情
// @Description 根据ID获取墓园详细信息
// @Tags 墓园
// @Accept json
// @Produce json
// @Param id path int true "墓园ID"
// @Success 200 {object} utils.Response{data=service.CemeteryDTO}
// @Router /cemeteries/{id} [get]
func (c *CemeteryController) Get(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	cemetery, err := c.cemeteryService.GetCemetery(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取墓园详情失败: "+err.Error())
		return
	}
	if cemetery == nil {
		utils.Error(ctx, http.StatusNotFound, "墓园不存在")
		return
	}
	utils.Success(ctx, cemetery)
}

// Create 创建墓园
// @Summary 创建墓园
// @Description 创建新墓园
// @Tags 墓园
// @Accept json
// @Produce json
// @Param cemetery body service.CreateCemeteryRequest true "墓园信息"
// @Success 200 {object} utils.Response{data=service.CemeteryDTO}
// @Security BearerAuth
// @Router /cemeteries [post]
func (c *CemeteryController) Create(ctx *gin.Context) {
	var req service.CreateCemeteryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	cemetery, err := c.cemeteryService.CreateCemetery(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, cemetery)
}

// Update 更新墓园
// @Summary 更新墓园
// @Description 更新墓园信息
// @Tags 墓园
// @Accept json
// @Produce json
// @Param id path int true "墓园ID"
// @Param cemetery body service.UpdateCemeteryRequest true "更新的墓园信息"
// @Success 200 {object} utils.Response{data=service.CemeteryDTO}
// @Security BearerAuth
// @Router /cemeteries/{id} [put]
func (c *CemeteryController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateCemeteryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	cemetery, err := c.cemeteryService.UpdateCemetery(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if cemetery == nil {
		utils.Error(ctx, http.StatusNotFound, "墓园不存在")
		return
	}
	utils.Success(ctx, cemetery)
}

// Delete 删除墓园
// @Summary 删除墓园
// @Description 删除指定墓园
// @Tags 墓园
// @Accept json
// @Produce json
// @Param id path int true "墓园ID"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /cemeteries/{id} [delete]
func (c *CemeteryController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.cemeteryService.DeleteCemetery(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// ListGrave 获取墓园下的墓位列表
// @Summary 获取墓园下的墓位列表
// @Description 获取指定墓园的所有墓位
// @Tags 墓园
// @Accept json
// @Produce json
// @Param id path int true "墓园ID"
// @Success 200 {object} utils.Response{data=[]service.GraveDTO}
// @Router /cemeteries/{id}/graves [get]
func (c *CemeteryController) ListGrave(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	graves, err := c.cemeteryService.ListGravesByCemetery(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取墓位列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, graves)
}

// GetGrave 获取墓位详情
// @Summary 获取墓位详情
// @Description 根据ID获取墓位详细信息
// @Tags 墓位
// @Accept json
// @Produce json
// @Param id path int true "墓位ID"
// @Success 200 {object} utils.Response{data=service.GraveDTO}
// @Router /graves/{id} [get]
func (c *CemeteryController) GetGrave(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	grave, err := c.cemeteryService.GetGrave(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取墓位详情失败: "+err.Error())
		return
	}
	if grave == nil {
		utils.Error(ctx, http.StatusNotFound, "墓位不存在")
		return
	}
	utils.Success(ctx, grave)
}

// CreateGrave 创建墓位
// @Summary 创建墓位
// @Description 创建新墓位
// @Tags 墓位
// @Accept json
// @Produce json
// @Param grave body service.CreateGraveRequest true "墓位信息"
// @Success 200 {object} utils.Response{data=service.GraveDTO}
// @Security BearerAuth
// @Router /graves [post]
func (c *CemeteryController) CreateGrave(ctx *gin.Context) {
	var req service.CreateGraveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	grave, err := c.cemeteryService.CreateGrave(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, grave)
}

// UpdateGrave 更新墓位
// @Summary 更新墓位
// @Description 更新墓位信息
// @Tags 墓位
// @Accept json
// @Produce json
// @Param id path int true "墓位ID"
// @Param grave body service.UpdateGraveRequest true "更新的墓位信息"
// @Success 200 {object} utils.Response{data=service.GraveDTO}
// @Security BearerAuth
// @Router /graves/{id} [put]
func (c *CemeteryController) UpdateGrave(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateGraveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	grave, err := c.cemeteryService.UpdateGrave(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if grave == nil {
		utils.Error(ctx, http.StatusNotFound, "墓位不存在")
		return
	}
	utils.Success(ctx, grave)
}

// DeleteGrave 删除墓位
// @Summary 删除墓位
// @Description 删除指定墓位
// @Tags 墓位
// @Accept json
// @Produce json
// @Param id path int true "墓位ID"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /graves/{id} [delete]
func (c *CemeteryController) DeleteGrave(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.cemeteryService.DeleteGrave(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}
