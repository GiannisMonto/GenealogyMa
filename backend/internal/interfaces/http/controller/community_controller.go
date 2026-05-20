package controller

import (
	"net/http"
	"strconv"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/genealogy-ma/platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// CommunityController 社区控制器
type CommunityController struct {
	communityService *service.CommunityService
}

// NewCommunityController 创建社区控制器
func NewCommunityController(communityService *service.CommunityService) *CommunityController {
	return &CommunityController{communityService: communityService}
}

// RegisterRoutes 注册路由
func (c *CommunityController) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// 用户资料路由
	profiles := r.Group("/profiles")
	{
		profiles.GET("/:id", c.GetProfile)
		profiles.GET("/user/:user_id", c.GetProfileByUserID)

		authRequired := profiles.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionCommunityWrite), c.CreateProfile)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionCommunityWrite), c.UpdateProfile)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionCommunityWrite), c.DeleteProfile)
		}
	}

	// 动态路由
	posts := r.Group("/posts")
	{
		posts.GET("/:id", c.GetPost)
		posts.GET("/user/:user_id", c.GetUserPosts)
		posts.GET("/tag/:tag", c.GetPostsByTag)

		authRequired := posts.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionCommunityWrite), c.CreatePost)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionCommunityWrite), c.UpdatePost)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionCommunityWrite), c.DeletePost)
			authRequired.POST("/:id/like", auth.RequirePermission(auth.PermissionCommunityWrite), c.LikePost)
		}
	}

	// 评论路由
	comments := r.Group("/comments")
	{
		comments.GET("/:id", c.GetComment)
		comments.GET("/post/:post_id", c.GetPostComments)

		authRequired := comments.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionCommunityWrite), c.CreateComment)
			authRequired.PUT("/:id", auth.RequirePermission(auth.PermissionCommunityWrite), c.UpdateComment)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionCommunityWrite), c.DeleteComment)
		}
	}

	// 私信路由
	messages := r.Group("/messages")
	{
		authRequired := messages.Group("")
		authRequired.Use(authMiddleware)
		{
			authRequired.POST("", auth.RequirePermission(auth.PermissionCommunityWrite), c.SendMessage)
			authRequired.GET("/conversation/:user_id", c.GetConversation)
			authRequired.GET("/inbox", c.GetUserMessages)
			authRequired.PUT("/:id/read", auth.RequirePermission(auth.PermissionCommunityWrite), c.MarkMessageAsRead)
			authRequired.DELETE("/:id", auth.RequirePermission(auth.PermissionCommunityWrite), c.DeleteMessage)
		}
	}
}

// ===== 用户资料接口 =====

func (c *CommunityController) GetProfile(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	profile, err := c.communityService.GetProfile(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取用户资料失败: "+err.Error())
		return
	}
	if profile == nil {
		utils.Error(ctx, http.StatusNotFound, "用户资料不存在")
		return
	}
	utils.Success(ctx, profile)
}

func (c *CommunityController) GetProfileByUserID(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的用户ID")
		return
	}

	profile, err := c.communityService.GetProfileByUserID(ctx, userID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取用户资料失败: "+err.Error())
		return
	}
	if profile == nil {
		utils.Error(ctx, http.StatusNotFound, "用户资料不存在")
		return
	}
	utils.Success(ctx, profile)
}

func (c *CommunityController) CreateProfile(ctx *gin.Context) {
	var req service.CreateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	profile, err := c.communityService.CreateProfile(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, profile)
}

func (c *CommunityController) UpdateProfile(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	profile, err := c.communityService.UpdateProfile(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if profile == nil {
		utils.Error(ctx, http.StatusNotFound, "用户资料不存在")
		return
	}
	utils.Success(ctx, profile)
}

func (c *CommunityController) DeleteProfile(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.communityService.DeleteProfile(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// ===== 动态接口 =====

func (c *CommunityController) GetPost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	post, err := c.communityService.GetPost(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取动态失败: "+err.Error())
		return
	}
	if post == nil {
		utils.Error(ctx, http.StatusNotFound, "动态不存在")
		return
	}
	utils.Success(ctx, post)
}

func (c *CommunityController) GetUserPosts(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的用户ID")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	posts, total, err := c.communityService.GetUserPosts(ctx, userID, page, pageSize)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取动态列表失败: "+err.Error())
		return
	}
	utils.SuccessWithPage(ctx, posts, total)
}

func (c *CommunityController) GetPostsByTag(ctx *gin.Context) {
	tag := ctx.Param("tag")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	posts, total, err := c.communityService.GetPostsByTag(ctx, tag, page, pageSize)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取动态列表失败: "+err.Error())
		return
	}
	utils.SuccessWithPage(ctx, posts, total)
}

func (c *CommunityController) CreatePost(ctx *gin.Context) {
	var req service.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	post, err := c.communityService.CreatePost(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, post)
}

func (c *CommunityController) UpdatePost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	post, err := c.communityService.UpdatePost(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if post == nil {
		utils.Error(ctx, http.StatusNotFound, "动态不存在")
		return
	}
	utils.Success(ctx, post)
}

func (c *CommunityController) DeletePost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.communityService.DeletePost(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

func (c *CommunityController) LikePost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.communityService.LikePost(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "点赞失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "点赞成功", nil)
}

// ===== 评论接口 =====

func (c *CommunityController) GetComment(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	comment, err := c.communityService.GetComment(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取评论失败: "+err.Error())
		return
	}
	if comment == nil {
		utils.Error(ctx, http.StatusNotFound, "评论不存在")
		return
	}
	utils.Success(ctx, comment)
}

func (c *CommunityController) GetPostComments(ctx *gin.Context) {
	postID, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的动态ID")
		return
	}

	comments, err := c.communityService.GetPostComments(ctx, postID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取评论列表失败: "+err.Error())
		return
	}
	utils.Success(ctx, comments)
}

func (c *CommunityController) CreateComment(ctx *gin.Context) {
	var req service.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	comment, err := c.communityService.CreateComment(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	utils.Success(ctx, comment)
}

func (c *CommunityController) UpdateComment(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req service.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	comment, err := c.communityService.UpdateComment(ctx, id, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	if comment == nil {
		utils.Error(ctx, http.StatusNotFound, "评论不存在")
		return
	}
	utils.Success(ctx, comment)
}

func (c *CommunityController) DeleteComment(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.communityService.DeleteComment(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// ===== 私信接口 =====

func (c *CommunityController) SendMessage(ctx *gin.Context) {
	var req service.SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	msg, err := c.communityService.SendMessage(ctx, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "发送失败: "+err.Error())
		return
	}
	utils.Success(ctx, msg)
}

func (c *CommunityController) GetConversation(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的用户ID")
		return
	}

	currentUserID := int64(1)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	messages, total, err := c.communityService.GetConversation(ctx, currentUserID, userID, page, pageSize)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取会话失败: "+err.Error())
		return
	}
	utils.SuccessWithPage(ctx, messages, total)
}

func (c *CommunityController) GetUserMessages(ctx *gin.Context) {
	userID := int64(1)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	messages, total, err := c.communityService.GetUserMessages(ctx, userID, page, pageSize)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取私信列表失败: "+err.Error())
		return
	}
	utils.SuccessWithPage(ctx, messages, total)
}

func (c *CommunityController) MarkMessageAsRead(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.communityService.MarkMessageAsRead(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "标记已读失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "标记已读成功", nil)
}

func (c *CommunityController) DeleteMessage(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.communityService.DeleteMessage(ctx, id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}
	utils.SuccessWithMessage(ctx, "删除成功", nil)
}
