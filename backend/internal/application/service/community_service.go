package service

import (
	"context"

	"github.com/genealogy-ma/platform/internal/domain/community"
)

// CommunityService 社区应用服务
type CommunityService struct {
	domainService *community.Service
}

// NewCommunityService 创建社区应用服务
func NewCommunityService(
	profileRepo community.UserProfileRepository,
	postRepo community.PostRepository,
	commentRepo community.CommentRepository,
	messageRepo community.MessageRepository,
) *CommunityService {
	return &CommunityService{
		domainService: community.NewService(profileRepo, postRepo, commentRepo, messageRepo),
	}
}

// ===== DTO 定义 =====

// UserProfileDTO 用户资料DTO
type UserProfileDTO struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"user_id"`
	Nickname       string `json:"nickname"`
	AvatarURL      string `json:"avatar_url"`
	Bio            string `json:"bio"`
	Gender         string `json:"gender"`
	BirthDate      string `json:"birth_date,omitempty"`
	Province       string `json:"province"`
	City           string `json:"city"`
	FollowerCount  int    `json:"follower_count"`
	FollowingCount int    `json:"following_count"`
	PostCount      int    `json:"post_count"`
	PrivacyLevel   int    `json:"privacy_level"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// PostDTO 动态DTO
type PostDTO struct {
	ID           int64    `json:"id"`
	UserID       int64    `json:"user_id"`
	Content      string   `json:"content"`
	Images       []string `json:"images"`
	Tag          string   `json:"tag"`
	LikeCount    int      `json:"like_count"`
	CommentCount int      `json:"comment_count"`
	ShareCount   int      `json:"share_count"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// CommentDTO 评论DTO
type CommentDTO struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	ParentID  *int64 `json:"parent_id,omitempty"`
	Content   string `json:"content"`
	LikeCount int    `json:"like_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// MessageDTO 私信DTO
type MessageDTO struct {
	ID         int64  `json:"id"`
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	Content    string `json:"content"`
	IsRead     bool   `json:"is_read"`
	ReadAt     string `json:"read_at,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// CreateProfileRequest 创建用户资料请求
type CreateProfileRequest struct {
	UserID       int64  `json:"user_id" binding:"required"`
	Nickname     string `json:"nickname" binding:"required"`
	AvatarURL    string `json:"avatar_url"`
	Bio          string `json:"bio"`
	Gender       string `json:"gender"`
	BirthDate    string `json:"birth_date"`
	Province     string `json:"province"`
	City         string `json:"city"`
	PrivacyLevel int    `json:"privacy_level"`
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Nickname     *string `json:"nickname"`
	AvatarURL    *string `json:"avatar_url"`
	Bio          *string `json:"bio"`
	Gender       *string `json:"gender"`
	BirthDate    *string `json:"birth_date"`
	Province     *string `json:"province"`
	City         *string `json:"city"`
	PrivacyLevel *int    `json:"privacy_level"`
}

// CreatePostRequest 创建动态请求
type CreatePostRequest struct {
	UserID  int64    `json:"user_id" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Images  []string `json:"images"`
	Tag     string   `json:"tag"`
}

// UpdatePostRequest 更新动态请求
type UpdatePostRequest struct {
	Content *string  `json:"content"`
	Images  []string `json:"images"`
	Tag     *string  `json:"tag"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID   int64  `json:"post_id" binding:"required"`
	UserID   int64  `json:"user_id" binding:"required"`
	ParentID *int64 `json:"parent_id"`
	Content  string `json:"content" binding:"required"`
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content *string `json:"content"`
}

// SendMessageRequest 发送私信请求
type SendMessageRequest struct {
	SenderID   int64  `json:"sender_id" binding:"required"`
	ReceiverID int64  `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// ===== 用户资料服务方法 =====

func (s *CommunityService) CreateProfile(ctx context.Context, req *CreateProfileRequest) (*UserProfileDTO, error) {
	profile := &community.UserProfile{
		UserID:       req.UserID,
		Nickname:     req.Nickname,
		AvatarURL:    req.AvatarURL,
		Bio:          req.Bio,
		Gender:       req.Gender,
		Province:     req.Province,
		City:         req.City,
		PrivacyLevel: req.PrivacyLevel,
	}
	if req.BirthDate != "" {
		profile.BirthDate = &req.BirthDate
	}
	if err := s.domainService.CreateProfile(ctx, profile); err != nil {
		return nil, err
	}
	return profileToDTO(profile), nil
}

func (s *CommunityService) GetProfile(ctx context.Context, id int64) (*UserProfileDTO, error) {
	profile, err := s.domainService.GetProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, nil
	}
	return profileToDTO(profile), nil
}

func (s *CommunityService) GetProfileByUserID(ctx context.Context, userID int64) (*UserProfileDTO, error) {
	profile, err := s.domainService.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, nil
	}
	return profileToDTO(profile), nil
}

func (s *CommunityService) UpdateProfile(ctx context.Context, id int64, req *UpdateProfileRequest) (*UserProfileDTO, error) {
	profile, err := s.domainService.GetProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, nil
	}
	if req.Nickname != nil {
		profile.Nickname = *req.Nickname
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = *req.AvatarURL
	}
	if req.Bio != nil {
		profile.Bio = *req.Bio
	}
	if req.Gender != nil {
		profile.Gender = *req.Gender
	}
	if req.BirthDate != nil {
		profile.BirthDate = req.BirthDate
	}
	if req.Province != nil {
		profile.Province = *req.Province
	}
	if req.City != nil {
		profile.City = *req.City
	}
	if req.PrivacyLevel != nil {
		profile.PrivacyLevel = *req.PrivacyLevel
	}
	if err := s.domainService.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}
	return profileToDTO(profile), nil
}

func (s *CommunityService) DeleteProfile(ctx context.Context, id int64) error {
	return s.domainService.DeleteProfile(ctx, id)
}

// ===== 动态服务方法 =====

func (s *CommunityService) CreatePost(ctx context.Context, req *CreatePostRequest) (*PostDTO, error) {
	post := &community.Post{
		UserID:  req.UserID,
		Content: req.Content,
		Images:  req.Images,
		Tag:     req.Tag,
	}
	if err := s.domainService.CreatePost(ctx, post); err != nil {
		return nil, err
	}
	return postToDTO(post), nil
}

func (s *CommunityService) GetPost(ctx context.Context, id int64) (*PostDTO, error) {
	post, err := s.domainService.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, nil
	}
	return postToDTO(post), nil
}

func (s *CommunityService) GetUserPosts(ctx context.Context, userID int64, page, pageSize int) ([]*PostDTO, int64, error) {
	posts, total, err := s.domainService.GetUserPosts(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]*PostDTO, len(posts))
	for i, p := range posts {
		dtos[i] = postToDTO(p)
	}
	return dtos, total, nil
}

func (s *CommunityService) GetPostsByTag(ctx context.Context, tag string, page, pageSize int) ([]*PostDTO, int64, error) {
	posts, total, err := s.domainService.GetPostsByTag(ctx, tag, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]*PostDTO, len(posts))
	for i, p := range posts {
		dtos[i] = postToDTO(p)
	}
	return dtos, total, nil
}

func (s *CommunityService) UpdatePost(ctx context.Context, id int64, req *UpdatePostRequest) (*PostDTO, error) {
	post, err := s.domainService.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, nil
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Images != nil {
		post.Images = req.Images
	}
	if req.Tag != nil {
		post.Tag = *req.Tag
	}
	if err := s.domainService.UpdatePost(ctx, post); err != nil {
		return nil, err
	}
	return postToDTO(post), nil
}

func (s *CommunityService) DeletePost(ctx context.Context, id int64) error {
	return s.domainService.DeletePost(ctx, id)
}

func (s *CommunityService) LikePost(ctx context.Context, id int64) error {
	return s.domainService.IncrementPostLikeCount(ctx, id)
}

// ===== 评论服务方法 =====

func (s *CommunityService) CreateComment(ctx context.Context, req *CreateCommentRequest) (*CommentDTO, error) {
	comment := &community.Comment{
		PostID:   req.PostID,
		UserID:   req.UserID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}
	if err := s.domainService.CreateComment(ctx, comment); err != nil {
		return nil, err
	}
	return commentToDTO(comment), nil
}

func (s *CommunityService) GetComment(ctx context.Context, id int64) (*CommentDTO, error) {
	comment, err := s.domainService.GetComment(ctx, id)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, nil
	}
	return commentToDTO(comment), nil
}

func (s *CommunityService) GetPostComments(ctx context.Context, postID int64) ([]*CommentDTO, error) {
	comments, err := s.domainService.GetPostComments(ctx, postID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*CommentDTO, len(comments))
	for i, c := range comments {
		dtos[i] = commentToDTO(c)
	}
	return dtos, nil
}

func (s *CommunityService) UpdateComment(ctx context.Context, id int64, req *UpdateCommentRequest) (*CommentDTO, error) {
	comment, err := s.domainService.GetComment(ctx, id)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, nil
	}
	if req.Content != nil {
		comment.Content = *req.Content
	}
	if err := s.domainService.UpdateComment(ctx, comment); err != nil {
		return nil, err
	}
	return commentToDTO(comment), nil
}

func (s *CommunityService) DeleteComment(ctx context.Context, id int64) error {
	return s.domainService.DeleteComment(ctx, id)
}

// ===== 私信服务方法 =====

func (s *CommunityService) SendMessage(ctx context.Context, req *SendMessageRequest) (*MessageDTO, error) {
	msg := &community.Message{
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}
	if err := s.domainService.SendMessage(ctx, msg); err != nil {
		return nil, err
	}
	return messageToDTO(msg), nil
}

func (s *CommunityService) GetConversation(ctx context.Context, userID1, userID2 int64, page, pageSize int) ([]*MessageDTO, int64, error) {
	messages, total, err := s.domainService.GetConversation(ctx, userID1, userID2, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]*MessageDTO, len(messages))
	for i, m := range messages {
		dtos[i] = messageToDTO(m)
	}
	return dtos, total, nil
}

func (s *CommunityService) GetUserMessages(ctx context.Context, userID int64, page, pageSize int) ([]*MessageDTO, int64, error) {
	messages, total, err := s.domainService.GetUserMessages(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]*MessageDTO, len(messages))
	for i, m := range messages {
		dtos[i] = messageToDTO(m)
	}
	return dtos, total, nil
}

func (s *CommunityService) MarkMessageAsRead(ctx context.Context, id int64) error {
	return s.domainService.MarkMessageAsRead(ctx, id)
}

func (s *CommunityService) DeleteMessage(ctx context.Context, id int64) error {
	return s.domainService.DeleteMessage(ctx, id)
}

// ===== 辅助函数 =====

func profileToDTO(p *community.UserProfile) *UserProfileDTO {
	dto := &UserProfileDTO{
		ID:             p.ID,
		UserID:         p.UserID,
		Nickname:       p.Nickname,
		AvatarURL:      p.AvatarURL,
		Bio:            p.Bio,
		Gender:         p.Gender,
		Province:       p.Province,
		City:           p.City,
		FollowerCount:  p.FollowerCount,
		FollowingCount: p.FollowingCount,
		PostCount:      p.PostCount,
		PrivacyLevel:   p.PrivacyLevel,
		CreatedAt:      p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if p.BirthDate != nil {
		dto.BirthDate = *p.BirthDate
	}
	return dto
}

func postToDTO(p *community.Post) *PostDTO {
	return &PostDTO{
		ID:           p.ID,
		UserID:       p.UserID,
		Content:      p.Content,
		Images:       p.Images,
		Tag:          p.Tag,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		ShareCount:   p.ShareCount,
		CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func commentToDTO(c *community.Comment) *CommentDTO {
	return &CommentDTO{
		ID:        c.ID,
		PostID:    c.PostID,
		UserID:    c.UserID,
		ParentID:  c.ParentID,
		Content:   c.Content,
		LikeCount: c.LikeCount,
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func messageToDTO(m *community.Message) *MessageDTO {
	dto := &MessageDTO{
		ID:         m.ID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Content:    m.Content,
		IsRead:     m.IsRead,
		CreatedAt:  m.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if m.ReadAt != nil {
		dto.ReadAt = m.ReadAt.Format("2006-01-02 15:04:05")
	}
	return dto
}
