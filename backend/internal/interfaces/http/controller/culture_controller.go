package controller

import (
	"net/http"
	"strconv"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// CultureController 文化控制器
type CultureController struct {
	cultureService *service.CultureService
}

// NewCultureController 创建文化控制器
func NewCultureController(cultureService *service.CultureService) *CultureController {
	return &CultureController{cultureService: cultureService}
}

// RegisterRoutes 注册路由
func (c *CultureController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// 文献路由
	documents := r.Group("/documents")
	{
		documents.GET("", c.ListDocuments)
		documents.GET("/:id", c.GetDocument)

		authRequired := documents.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionCultureWrite), c.CreateDocument)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionCultureWrite), c.UpdateDocument)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionCultureWrite), c.DeleteDocument)
		}
	}

	// 故事路由
	stories := r.Group("/stories")
	{
		stories.GET("", c.ListStories)
		stories.GET("/:id", c.GetStory)

		authRequired := stories.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionCultureWrite), c.CreateStory)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionCultureWrite), c.UpdateStory)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionCultureWrite), c.DeleteStory)
		}
	}

	// 家训路由
	teachings := r.Group("/family-teachings")
	{
		teachings.GET("", c.ListFamilyTeachings)
		teachings.GET("/:id", c.GetFamilyTeachings)

		authRequired := teachings.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionCultureWrite), c.CreateFamilyTeachings)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionCultureWrite), c.UpdateFamilyTeachings)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionCultureWrite), c.DeleteFamilyTeachings)
		}
	}
}

// ===== 文献接口 =====

// ListDocuments 获取文献列表
func (c *CultureController) ListDocuments(ctx *gin.Context) {
	documents, err := c.cultureService.ListDocuments(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取文献列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, documents)
}

// GetDocument 获取文献详情
func (c *CultureController) GetDocument(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	doc, err := c.cultureService.GetDocument(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取文献详情失败: "+err.Error())
		return
	}
	if doc == nil {
		utils.Error(ctx, http.StatusNotFound, "文献不存在")
		return
	}
	utils.Success(ctx, doc)
}

// CreateDocument 创建文献
func (c *CultureController) CreateDocument(ctx *gin.Context) {
	var req service.CreateDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	doc, err := c.cultureService.CreateDocument(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, doc)
}

// UpdateDocument 更新文献
func (c *CultureController) UpdateDocument(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	doc, err := c.cultureService.UpdateDocument(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if doc == nil {
		utils.Error(ctx, http.StatusNotFound, "文献不存在")
		return
	}
	utils.Success(ctx, doc)
}

// DeleteDocument 删除文献
func (c *CultureController) DeleteDocument(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.cultureService.DeleteDocument(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// ===== 故事接口 =====

// ListStories 获取故事列表
func (c *CultureController) ListStories(ctx *gin.Context) {
	stories, err := c.cultureService.ListStories(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取故事列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, stories)
}

// GetStory 获取故事详情
func (c *CultureController) GetStory(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	story, err := c.cultureService.GetStory(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取故事详情失败: "+err.Error())
		return
	}
	if story == nil {
		utils.Error(ctx, http.StatusNotFound, "故事不存在")
		return
	}
	utils.Success(ctx, story)
}

// CreateStory 创建故事
func (c *CultureController) CreateStory(ctx *gin.Context) {
	var req service.CreateStoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	story, err := c.cultureService.CreateStory(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, story)
}

// UpdateStory 更新故事
func (c *CultureController) UpdateStory(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateStoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	story, err := c.cultureService.UpdateStory(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if story == nil {
		utils.Error(ctx, http.StatusNotFound, "故事不存在")
		return
	}
	utils.Success(ctx, story)
}

// DeleteStory 删除故事
func (c *CultureController) DeleteStory(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.cultureService.DeleteStory(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// ===== 家训接口 =====

// ListFamilyTeachings 获取家训列表
func (c *CultureController) ListFamilyTeachings(ctx *gin.Context) {
	items, err := c.cultureService.ListFamilyTeachings(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取家训列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, items)
}

// GetFamilyTeachings 获取家训详情
func (c *CultureController) GetFamilyTeachings(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	ft, err := c.cultureService.GetFamilyTeachings(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取家训详情失败: "+err.Error())
		return
	}
	if ft == nil {
		utils.Error(ctx, http.StatusNotFound, "家训不存在")
		return
	}
	utils.Success(ctx, ft)
}

// CreateFamilyTeachings 创建家训
func (c *CultureController) CreateFamilyTeachings(ctx *gin.Context) {
	var req service.CreateFamilyTeachingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ft, err := c.cultureService.CreateFamilyTeachings(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, ft)
}

// UpdateFamilyTeachings 更新家训
func (c *CultureController) UpdateFamilyTeachings(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateFamilyTeachingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ft, err := c.cultureService.UpdateFamilyTeachings(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if ft == nil {
		utils.Error(ctx, http.StatusNotFound, "家训不存在")
		return
	}
	utils.Success(ctx, ft)
}

// DeleteFamilyTeachings 删除家训
func (c *CultureController) DeleteFamilyTeachings(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.cultureService.DeleteFamilyTeachings(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}
