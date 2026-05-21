package controller

import (
	"net/http"
	"strconv"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// PersonController 人物控制器
type PersonController struct {
	personService *service.PersonService
}

// NewPersonController 创建人物控制器
func NewPersonController(personService *service.PersonService) *PersonController {
	return &PersonController{personService: personService}
}

// RegisterRoutes 注册路由
func (c *PersonController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	persons := r.Group("/persons")
	{
		// 公开路由（游客可访问）
		persons.GET("", c.List)
		persons.GET("/:id", c.Get)
		persons.GET("/:id/tree", c.GetFamilyTree)
		persons.GET("/statistics", c.GetStatistics)

		// 需要认证的路由
		authRequired := persons.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionPersonWrite), c.Create)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionPersonWrite), c.Update)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionPersonDelete), c.Delete)
			// 批量操作
			authRequired.POST("/batch", auth.RequirePermission(auth.PermissionPersonWrite), c.BatchCreate)
			authRequired.PUT("/batch", auth.RequirePermission(auth.PermissionPersonWrite), c.BatchUpdate)
			authRequired.DELETE("/batch", auth.RequirePermission(auth.PermissionPersonDelete), c.BatchDelete)
		}
	}
}

// ===== 处理函数 =====

// Get 获取人物详情
// @Summary 获取人物详情
// @Description 根据ID获取人物详细信息
// @Tags 人物
// @Accept json
// @Produce json
// @Param id path int true "人物ID"
// @Param with_relations query bool false "是否加载关联数据"
// @Success 200 {object} utils.Response{data=service.PersonDTO}
// @Router /persons/{id} [get]
func (c *PersonController) Get(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	withRelations := ctx.Query("with_relations") == "true"

	person, err := c.personService.GetPerson(ctx, id, withRelations)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "人物不存在")
		return
	}

	utils.Success(ctx, person)
}

// List 搜索人物列表
// @Summary 搜索人物列表
// @Description 根据条件搜索人物列表
// @Tags 人物
// @Accept json
// @Produce json
// @Param keyword query string false "关键字"
// @Param name query string false "姓名"
// @Param gender query string false "性别"
// @Param generation query int false "世代"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段"
// @Param sort_desc query bool false "是否降序"
// @Success 200 {object} utils.Response{data=[]service.PersonDTO,total=int}
// @Router /persons [get]
func (c *PersonController) List(ctx *gin.Context) {
	var req service.SearchPersonRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	persons, total, err := c.personService.SearchPersons(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	utils.SuccessWithPage(ctx, persons, total)
}

// Create 创建人物
// @Summary 创建人物
// @Description 创建新人物
// @Tags 人物
// @Accept json
// @Produce json
// @Param person body service.CreatePersonRequest true "人物信息"
// @Success 200 {object} utils.Response{data=service.PersonDTO}
// @Security BearerAuth
// @Router /persons [post]
func (c *PersonController) Create(ctx *gin.Context) {
	var req service.CreatePersonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	person, err := c.personService.CreatePerson(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	utils.Success(ctx, person)
}

// Update 更新人物
// @Summary 更新人物
// @Description 更新人物信息
// @Tags 人物
// @Accept json
// @Produce json
// @Param id path int true "人物ID"
// @Param person body service.UpdatePersonRequest true "更新的人物信息"
// @Success 200 {object} utils.Response{data=service.PersonDTO}
// @Security BearerAuth
// @Router /persons/{id} [put]
func (c *PersonController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdatePersonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	person, err := c.personService.UpdatePerson(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	utils.Success(ctx, person)
}

// Delete 删除人物
// @Summary 删除人物
// @Description 删除指定人物
// @Tags 人物
// @Accept json
// @Produce json
// @Param id path int true "人物ID"
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /persons/{id} [delete]
func (c *PersonController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.personService.DeletePerson(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// GetFamilyTree 获取族谱树
// @Summary 获取族谱树
// @Description 获取指定人物的上下五代族谱树
// @Tags 人物
// @Accept json
// @Produce json
// @Param id path int true "人物ID"
// @Param up_depth query int false "向上几代" default(5)
// @Param down_depth query int false "向下几代" default(5)
// @Success 200 {object} utils.Response{data=service.FamilyTreeResponse}
// @Router /persons/{id}/tree [get]
func (c *PersonController) GetFamilyTree(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	upDepth, _ := strconv.Atoi(ctx.DefaultQuery("up_depth", "5"))
	downDepth, _ := strconv.Atoi(ctx.DefaultQuery("down_depth", "5"))

	if upDepth < 0 {
		upDepth = 0
	}
	if downDepth < 0 {
		downDepth = 0
	}
	if upDepth > 10 {
		upDepth = 10
	}
	if downDepth > 10 {
		downDepth = 10
	}

	tree, err := c.personService.GetFamilyTree(ctx, id, upDepth, downDepth)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取族谱树失败: "+err.Error())
		return
	}

	utils.Success(ctx, tree)
}

// GetStatistics 获取统计数据
// @Summary 获取统计数据
// @Description 获取人物统计数据
// @Tags 人物
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=service.StatisticsResponse}
// @Router /persons/statistics [get]
func (c *PersonController) GetStatistics(ctx *gin.Context) {
	stats, err := c.personService.GetStatistics(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取统计数据失败: "+err.Error())
		return
	}

	utils.Success(ctx, stats)
}

// BatchCreate 批量创建人物
// @Summary 批量创建人物
// @Description 批量创建多个新人物
// @Tags 人物
// @Accept json
// @Produce json
// @Param persons body service.BatchCreatePersonRequest true "人物列表"
// @Success 200 {object} utils.Response{data=service.BatchResult}
// @Security BearerAuth
// @Router /persons/batch [post]
func (c *PersonController) BatchCreate(ctx *gin.Context) {
	var req service.BatchCreatePersonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if len(req.Persons) == 0 {
		utils.Error(ctx, http.StatusBadRequest, "人物列表不能为空")
		return
	}

	result, err := c.personService.BatchCreatePersons(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "批量创建失败: "+err.Error())
		return
	}

	utils.Success(ctx, result)
}

// BatchUpdate 批量更新人物
// @Summary 批量更新人物
// @Description 批量更新多个人的信息
// @Tags 人物
// @Accept json
// @Produce json
// @Param persons body service.BatchUpdatePersonRequest true "更新列表"
// @Success 200 {object} utils.Response{data=service.BatchResult}
// @Security BearerAuth
// @Router /persons/batch [put]
func (c *PersonController) BatchUpdate(ctx *gin.Context) {
	var req service.BatchUpdatePersonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if len(req.Persons) == 0 {
		utils.Error(ctx, http.StatusBadRequest, "更新列表不能为空")
		return
	}

	result, err := c.personService.BatchUpdatePersons(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "批量更新失败: "+err.Error())
		return
	}

	utils.Success(ctx, result)
}

// BatchDelete 批量删除人物
// @Summary 批量删除人物
// @Description 批量删除指定人物
// @Tags 人物
// @Accept json
// @Produce json
// @Param ids body service.BatchDeleteRequest true "ID列表"
// @Success 200 {object} utils.Response{data=service.BatchResult}
// @Security BearerAuth
// @Router /persons/batch [delete]
func (c *PersonController) BatchDelete(ctx *gin.Context) {
	var req service.BatchDeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if len(req.IDs) == 0 {
		utils.Error(ctx, http.StatusBadRequest, "ID列表不能为空")
		return
	}

	result, err := c.personService.BatchDeletePersons(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "批量删除失败: "+err.Error())
		return
	}

	utils.Success(ctx, result)
}
