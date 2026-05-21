package controller

import (
	"net/http"
	"strconv"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// GenealogyController 族谱控制器
type GenealogyController struct {
	genealogyService *service.GenealogyService
}

// NewGenealogyController 创建族谱控制器
func NewGenealogyController(genealogyService *service.GenealogyService) *GenealogyController {
	return &GenealogyController{genealogyService: genealogyService}
}

// RegisterRoutes 注册路由
func (c *GenealogyController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// 族谱路由
	genealogies := r.Group("/genealogies")
	{
		genealogies.GET("", c.List)
		genealogies.GET("/:id", c.Get)
		genealogies.GET("/:id/details", c.GetWithDetails)

		authRequired := genealogies.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionGenealogyWrite), c.Create)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.Update)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.Delete)
		}
	}

	// 分支路由
	branches := r.Group("/branches")
	{
		branches.GET("/:id", c.GetBranch)

		authRequired := branches.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionGenealogyWrite), c.CreateBranch)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.UpdateBranch)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.DeleteBranch)
		}
	}

	// 世代路由
	generations := r.Group("/generations")
	{
		generations.GET("/:id", c.GetGeneration)
		generations.GET("/by-number", c.GetGenerationByNumber)

		authRequired := generations.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionGenealogyWrite), c.CreateGeneration)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.UpdateGeneration)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionGenealogyWrite), c.DeleteGeneration)
		}
	}
}

// List 获取族谱列表
// @Summary 获取族谱列表
// @Description 获取所有族谱列表
// @Tags 族谱
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]service.GenealogyDTO}
// @Router /genealogies [get]
func (c *GenealogyController) List(ctx *gin.Context) {
	genealogies, err := c.genealogyService.ListGenealogies(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取族谱列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, genealogies)
}

// Get 获取族谱详情
// @Summary 获取族谱详情
// @Description 根据ID获取族谱详细信息
// @Tags 族谱
// @Accept json
// @Produce json
// @Param id path int true "族谱ID"
// @Success 200 {object} utils.Response{data=service.GenealogyDTO}
// @Router /genealogies/{id} [get]
func (c *GenealogyController) Get(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	genealogy, err := c.genealogyService.GetGenealogy(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取族谱详情失败: "+err.Error())
		return
	}
	if genealogy == nil {
		utils.Error(ctx, http.StatusNotFound, "族谱不存在")
		return
	}
	utils.Success(ctx, genealogy)
}

// GetWithDetails 获取族谱详情（含分支和世代）
// @Summary 获取族谱详情（含分支和世代）
// @Description 根据ID获取族谱详细信息，包含所有分支和世代
// @Tags 族谱
// @Accept json
// @Produce json
// @Param id path int true "族谱ID"
// @Success 200 {object} utils.Response{data=service.GenealogyDTO}
// @Router /genealogies/{id}/details [get]
func (c *GenealogyController) GetWithDetails(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	genealogy, err := c.genealogyService.GetGenealogyWithDetails(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取族谱详情失败: "+err.Error())
		return
	}
	if genealogy == nil {
		utils.Error(ctx, http.StatusNotFound, "族谱不存在")
		return
	}
	utils.Success(ctx, genealogy)
}

// Create 创建族谱
// @Summary 创建族谱
// @Description 创建新族谱
// @Tags 族谱
// @Accept json
// @Produce json
// @Param genealogy body service.CreateGenealogyRequest true "族谱信息"
// @Success 200 {object} utils.Response{data=service.GenealogyDTO}
// @Security BearerAuth
// @Router /genealogies [post]
func (c *GenealogyController) Create(ctx *gin.Context) {
	var req service.CreateGenealogyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	genealogy, err := c.genealogyService.CreateGenealogy(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, genealogy)
}

// Update 更新族谱
// @Summary 更新族谱
// @Description 更新族谱信息
// @Tags 族谱
// @Accept json
// @Produce json
// @Param id path int true "族谱ID"
// @Param genealogy body service.UpdateGenealogyRequest true "更新的族谱信息"
// @Success 200 {object} utils.Response{data=service.GenealogyDTO}
// @Security BearerAuth
// @Router /genealogies/{id} [put]
func (c *GenealogyController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateGenealogyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	genealogy, err := c.genealogyService.UpdateGenealogy(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if genealogy == nil {
		utils.Error(ctx, http.StatusNotFound, "族谱不存在")
		return
	}
	utils.Success(ctx, genealogy)
}

// Delete 删除族谱
// @Summary 删除族谱
// @Description 删除指定族谱
// @Tags 族谱
// @Accept json
// @Produce json
// @Param id path int true "族谱ID"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /genealogies/{id} [delete]
func (c *GenealogyController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.genealogyService.DeleteGenealogy(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// GetBranch 获取分支详情
// @Summary 获取分支详情
// @Description 根据ID获取分支详细信息
// @Tags 分支
// @Accept json
// @Produce json
// @Param id path int true "分支ID"
// @Success 200 {object} utils.Response{data=service.BranchDTO}
// @Router /branches/{id} [get]
func (c *GenealogyController) GetBranch(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	branch, err := c.genealogyService.GetBranch(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取分支详情失败: "+err.Error())
		return
	}
	if branch == nil {
		utils.Error(ctx, http.StatusNotFound, "分支不存在")
		return
	}
	utils.Success(ctx, branch)
}

// CreateBranch 创建分支
// @Summary 创建分支
// @Description 创建新分支
// @Tags 分支
// @Accept json
// @Produce json
// @Param branch body service.CreateBranchRequest true "分支信息"
// @Success 200 {object} utils.Response{data=service.BranchDTO}
// @Security BearerAuth
// @Router /branches [post]
func (c *GenealogyController) CreateBranch(ctx *gin.Context) {
	var req service.CreateBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	branch, err := c.genealogyService.CreateBranch(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, branch)
}

// UpdateBranch 更新分支
// @Summary 更新分支
// @Description 更新分支信息
// @Tags 分支
// @Accept json
// @Produce json
// @Param id path int true "分支ID"
// @Param branch body service.UpdateBranchRequest true "更新的分支信息"
// @Success 200 {object} utils.Response{data=service.BranchDTO}
// @Security BearerAuth
// @Router /branches/{id} [put]
func (c *GenealogyController) UpdateBranch(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	branch, err := c.genealogyService.UpdateBranch(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if branch == nil {
		utils.Error(ctx, http.StatusNotFound, "分支不存在")
		return
	}
	utils.Success(ctx, branch)
}

// DeleteBranch 删除分支
// @Summary 删除分支
// @Description 删除指定分支
// @Tags 分支
// @Accept json
// @Produce json
// @Param id path int true "分支ID"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /branches/{id} [delete]
func (c *GenealogyController) DeleteBranch(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.genealogyService.DeleteBranch(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// GetGeneration 获取世代详情
// @Summary 获取世代详情
// @Description 根据ID获取世代详细信息
// @Tags 世代
// @Accept json
// @Produce json
// @Param id path int true "世代ID"
// @Success 200 {object} utils.Response{data=service.GenerationDTO}
// @Router /generations/{id} [get]
func (c *GenealogyController) GetGeneration(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	generation, err := c.genealogyService.GetGeneration(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取世代详情失败: "+err.Error())
		return
	}
	if generation == nil {
		utils.Error(ctx, http.StatusNotFound, "世代不存在")
		return
	}
	utils.Success(ctx, generation)
}

// GetGenerationByNumber 根据世代序号获取世代
// @Summary 根据世代序号获取世代
// @Description 根据族谱ID和世代序号获取世代信息
// @Tags 世代
// @Accept json
// @Produce json
// @Param genealogy_id query int true "族谱ID"
// @Param generation query int true "世代序号"
// @Success 200 {object} utils.Response{data=service.GenerationDTO}
// @Router /generations/by-number [get]
func (c *GenealogyController) GetGenerationByNumber(ctx *gin.Context) {
	genealogyID, err := strconv.ParseInt(ctx.Query("genealogy_id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的族谱ID")
		return
	}

	generationNum, err := strconv.Atoi(ctx.Query("generation"))
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的世代序号")
		return
	}

	generation, err := c.genealogyService.GetGenerationByNumber(ctx, genealogyID, generationNum)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取世代详情失败: "+err.Error())
		return
	}
	if generation == nil {
		utils.Error(ctx, http.StatusNotFound, "世代不存在")
		return
	}
	utils.Success(ctx, generation)
}

// CreateGeneration 创建世代
// @Summary 创建世代
// @Description 创建新世代
// @Tags 世代
// @Accept json
// @Produce json
// @Param generation body service.CreateGenerationRequest true "世代信息"
// @Success 200 {object} utils.Response{data=service.GenerationDTO}
// @Security BearerAuth
// @Router /generations [post]
func (c *GenealogyController) CreateGeneration(ctx *gin.Context) {
	var req service.CreateGenerationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	generation, err := c.genealogyService.CreateGeneration(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, generation)
}

// UpdateGeneration 更新世代
// @Summary 更新世代
// @Description 更新世代信息
// @Tags 世代
// @Accept json
// @Produce json
// @Param id path int true "世代ID"
// @Param generation body service.UpdateGenerationRequest true "更新的世代信息"
// @Success 200 {object} utils.Response{data=service.GenerationDTO}
// @Security BearerAuth
// @Router /generations/{id} [put]
func (c *GenealogyController) UpdateGeneration(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateGenerationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	generation, err := c.genealogyService.UpdateGeneration(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if generation == nil {
		utils.Error(ctx, http.StatusNotFound, "世代不存在")
		return
	}
	utils.Success(ctx, generation)
}

// DeleteGeneration 删除世代
// @Summary 删除世代
// @Description 删除指定世代
// @Tags 世代
// @Accept json
// @Produce json
// @Param id path int true "世代ID"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /generations/{id} [delete]
func (c *GenealogyController) DeleteGeneration(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.genealogyService.DeleteGeneration(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}