package controller

import (
	"net/http"
	"strconv"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// MemorialController 宗祠控制器
type MemorialController struct {
	memorialService *service.MemorialService
}

// NewMemorialController 创建宗祠控制器
func NewMemorialController(memorialService *service.MemorialService) *MemorialController {
	return &MemorialController{memorialService: memorialService}
}

// RegisterRoutes 注册路由
func (c *MemorialController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// 宗祠路由
	halls := r.Group("/halls")
	{
		halls.GET("", c.List)
		halls.GET("/:id", c.Get)
		halls.GET("/:id/tablets", c.GetWithTablets)

		authRequired := halls.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionGenealogyWrite), c.Create)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.Update)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.Delete)
		}
	}

	// 牌位路由
	tablets := r.Group("/tablets")
	{
		tablets.GET("/:id", c.GetTablet)
		tablets.GET("/person/:person_id", c.GetTabletByPersonID)

		authRequired := tablets.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionGenealogyWrite), c.CreateTablet)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.UpdateTablet)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.DeleteTablet)
		}
	}
}

// List 获取宗祠列表
// @Summary 获取宗祠列表
// @Description 获取所有宗祠列表
// @Tags 宗祠
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]service.MemorialHallDTO}
// @Router /halls [get]
func (c *MemorialController) List(ctx *gin.Context) {
	halls, err := c.memorialService.ListHalls(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取宗祠列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, halls)
}

// Get 获取宗祠详情
// @Summary 获取宗祠详情
// @Description 根据ID获取宗祠详细信息
// @Tags 宗祠
// @Accept json
// @Produce json
// @Param id path int true "宗祠ID"
// @Success 200 {object} utils.Response{data=service.MemorialHallDTO}
// @Router /halls/{id} [get]
func (c *MemorialController) Get(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	hall, err := c.memorialService.GetHall(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取宗祠详情失败: "+err.Error())
		return
	}
	if hall == nil {
		utils.Error(ctx, http.StatusNotFound, "宗祠不存在")
		return
	}
	utils.Success(ctx, hall)
}

// GetWithTablets 获取宗祠详情（含牌位）
// @Summary 获取宗祠详情（含牌位）
// @Description 根据ID获取宗祠详细信息，包含所有牌位
// @Tags 宗祠
// @Accept json
// @Produce json
// @Param id path int true "宗祠ID"
// @Success 200 {object} utils.Response{data=service.MemorialHallDTO}
// @Router /halls/{id}/tablets [get]
func (c *MemorialController) GetWithTablets(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	hall, err := c.memorialService.GetHallWithTablets(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取宗祠详情失败: "+err.Error())
		return
	}
	if hall == nil {
		utils.Error(ctx, http.StatusNotFound, "宗祠不存在")
		return
	}
	utils.Success(ctx, hall)
}

// Create 创建宗祠
// @Summary 创建宗祠
// @Description 创建新宗祠
// @Tags 宗祠
// @Accept json
// @Produce json
// @Param hall body service.CreateHallRequest true "宗祠信息"
// @Success 200 {object} utils.Response{data=service.MemorialHallDTO}
// @Security BearerAuth
// @Router /halls [post]
func (c *MemorialController) Create(ctx *gin.Context) {
	var req service.CreateHallRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	hall, err := c.memorialService.CreateHall(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, hall)
}

// Update 更新宗祠
// @Summary 更新宗祠
// @Description 更新宗祠信息
// @Tags 宗祠
// @Accept json
// @Produce json
// @Param id path int true "宗祠ID"
// @Param hall body service.UpdateHallRequest true "更新的宗祠信息"
// @Success 200 {object} utils.Response{data=service.MemorialHallDTO}
// @Security BearerAuth
// @Router /halls/{id} [put]
func (c *MemorialController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateHallRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	hall, err := c.memorialService.UpdateHall(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if hall == nil {
		utils.Error(ctx, http.StatusNotFound, "宗祠不存在")
		return
	}
	utils.Success(ctx, hall)
}

// Delete 删除宗祠
// @Summary 删除宗祠
// @Description 删除指定宗祠
// @Tags 宗祠
// @Accept json
// @Produce json
// @Param id path int true "宗祠ID"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /halls/{id} [delete]
func (c *MemorialController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.memorialService.DeleteHall(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// GetTablet 获取牌位详情
// @Summary 获取牌位详情
// @Description 根据ID获取牌位详细信息
// @Tags 牌位
// @Accept json
// @Produce json
// @Param id path int true "牌位ID"
// @Success 200 {object} utils.Response{data=service.TabletDTO}
// @Router /tablets/{id} [get]
func (c *MemorialController) GetTablet(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	tablet, err := c.memorialService.GetTablet(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取牌位详情失败: "+err.Error())
		return
	}
	if tablet == nil {
		utils.Error(ctx, http.StatusNotFound, "牌位不存在")
		return
	}
	utils.Success(ctx, tablet)
}

// GetTabletByPersonID 根据人物ID获取牌位
// @Summary 根据人物ID获取牌位
// @Description 根据人物ID获取对应的牌位信息
// @Tags 牌位
// @Accept json
// @Produce json
// @Param person_id path int true "人物ID"
// @Success 200 {object} utils.Response{data=service.TabletDTO}
// @Router /tablets/person/{person_id} [get]
func (c *MemorialController) GetTabletByPersonID(ctx *gin.Context) {
	personID, err := strconv.ParseInt(ctx.Param("person_id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的人物ID")
		return
	}

	tablet, err := c.memorialService.GetTabletByPersonID(ctx, personID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取牌位详情失败: "+err.Error())
		return
	}
	if tablet == nil {
		utils.Error(ctx, http.StatusNotFound, "牌位不存在")
		return
	}
	utils.Success(ctx, tablet)
}

// CreateTablet 创建牌位
// @Summary 创建牌位
// @Description 创建新牌位
// @Tags 牌位
// @Accept json
// @Produce json
// @Param tablet body service.CreateTabletRequest true "牌位信息"
// @Success 200 {object} utils.Response{data=service.TabletDTO}
// @Security BearerAuth
// @Router /tablets [post]
func (c *MemorialController) CreateTablet(ctx *gin.Context) {
	var req service.CreateTabletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tablet, err := c.memorialService.CreateTablet(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, tablet)
}

// UpdateTablet 更新牌位
// @Summary 更新牌位
// @Description 更新牌位信息
// @Tags 牌位
// @Accept json
// @Produce json
// @Param id path int true "牌位ID"
// @Param tablet body service.UpdateTabletRequest true "更新的牌位信息"
// @Success 200 {object} utils.Response{data=service.TabletDTO}
// @Security BearerAuth
// @Router /tablets/{id} [put]
func (c *MemorialController) UpdateTablet(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateTabletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tablet, err := c.memorialService.UpdateTablet(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if tablet == nil {
		utils.Error(ctx, http.StatusNotFound, "牌位不存在")
		return
	}
	utils.Success(ctx, tablet)
}

// DeleteTablet 删除牌位
// @Summary 删除牌位
// @Description 删除指定牌位
// @Tags 牌位
// @Accept json
// @Produce json
// @Param id path int true "牌位ID"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /tablets/{id} [delete]
func (c *MemorialController) DeleteTablet(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.memorialService.DeleteTablet(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}