package persistence

import (
	"context"
	"encoding/json"

	"github.com/genealogy-ma/platform/internal/domain/community"
	"gorm.io/gorm"
)

// UserProfileModel GORM模型对应user_profiles表
type UserProfileModel struct {
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement"`
	UserID         int64  `gorm:"column:user_id;not null;uniqueIndex"`
	Nickname       string `gorm:"column:nickname;type:varchar(50);not null"`
	AvatarURL      string `gorm:"column:avatar_url;type:varchar(500)"`
	Bio            string `gorm:"column:bio;type:text"`
	Gender         string `gorm:"column:gender;type:varchar(10)"`
	BirthDate      string `gorm:"column:birth_date;type:varchar(20)"`
	Province       string `gorm:"column:province;type:varchar(50)"`
	City           string `gorm:"column:city;type:varchar(50)"`
	FollowerCount  int    `gorm:"column:follower_count;default:0"`
	FollowingCount int    `gorm:"column:following_count;default:0"`
	PostCount      int    `gorm:"column:post_count;default:0"`
	PrivacyLevel   int    `gorm:"column:privacy_level;default:0"`
	CreatedAt      string `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      string `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserProfileModel) TableName() string {
	return "user_profiles"
}

// PostModel GORM模型对应posts表
type PostModel struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64  `gorm:"column:user_id;not null;index"`
	Content      string `gorm:"column:content;type:text;not null"`
	Images       string `gorm:"column:images;type:text"`
	Tag          string `gorm:"column:tag;type:varchar(50);index"`
	LikeCount    int    `gorm:"column:like_count;default:0"`
	CommentCount int    `gorm:"column:comment_count;default:0"`
	ShareCount   int    `gorm:"column:share_count;default:0"`
	IsDeleted    bool   `gorm:"column:is_deleted;default:false;index"`
	CreatedAt    string `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    string `gorm:"column:updated_at;autoUpdateTime"`
}

func (PostModel) TableName() string {
	return "posts"
}

// CommentModel GORM模型对应comments表
type CommentModel struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	PostID    int64  `gorm:"column:post_id;not null;index"`
	UserID    int64  `gorm:"column:user_id;not null;index"`
	ParentID  *int64 `gorm:"column:parent_id;index"`
	Content   string `gorm:"column:content;type:text;not null"`
	LikeCount int    `gorm:"column:like_count;default:0"`
	CreatedAt string `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt string `gorm:"column:updated_at;autoUpdateTime"`
}

func (CommentModel) TableName() string {
	return "comments"
}

// MessageModel GORM模型对应messages表
type MessageModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	SenderID   int64  `gorm:"column:sender_id;not null;index"`
	ReceiverID int64  `gorm:"column:receiver_id;not null;index"`
	Content    string `gorm:"column:content;type:text;not null"`
	IsRead     bool   `gorm:"column:is_read;default:false"`
	ReadAt     string `gorm:"column:read_at"`
	CreatedAt  string `gorm:"column:created_at;autoCreateTime"`
}

func (MessageModel) TableName() string {
	return "messages"
}

// ===== 用户资料仓储实现 =====

type UserProfileRepositoryImpl struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) community.UserProfileRepository {
	return &UserProfileRepositoryImpl{db: db}
}

func profileModelToEntity(m *UserProfileModel) *community.UserProfile {
	p := &community.UserProfile{
		ID:             m.ID,
		UserID:         m.UserID,
		Nickname:       m.Nickname,
		AvatarURL:      m.AvatarURL,
		Bio:            m.Bio,
		Gender:         m.Gender,
		Province:       m.Province,
		City:           m.City,
		FollowerCount:  m.FollowerCount,
		FollowingCount: m.FollowingCount,
		PostCount:      m.PostCount,
		PrivacyLevel:   m.PrivacyLevel,
	}
	if m.BirthDate != "" {
		p.BirthDate = &m.BirthDate
	}
	return p
}

func profileEntityToModel(p *community.UserProfile) *UserProfileModel {
	m := &UserProfileModel{
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
	}
	if p.BirthDate != nil {
		m.BirthDate = *p.BirthDate
	}
	return m
}

func (r *UserProfileRepositoryImpl) FindByID(ctx context.Context, id int64) (*community.UserProfile, error) {
	var model UserProfileModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return profileModelToEntity(&model), nil
}

func (r *UserProfileRepositoryImpl) FindByUserID(ctx context.Context, userID int64) (*community.UserProfile, error) {
	var model UserProfileModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return profileModelToEntity(&model), nil
}

func (r *UserProfileRepositoryImpl) Create(ctx context.Context, p *community.UserProfile) error {
	model := profileEntityToModel(p)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *UserProfileRepositoryImpl) Update(ctx context.Context, p *community.UserProfile) error {
	model := profileEntityToModel(p)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *UserProfileRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&UserProfileModel{}, id).Error
}

// ===== 动态仓储实现 =====

type PostRepositoryImpl struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) community.PostRepository {
	return &PostRepositoryImpl{db: db}
}

func postModelToEntity(m *PostModel) *community.Post {
	var images []string
	if m.Images != "" {
		_ = json.Unmarshal([]byte(m.Images), &images)
	}
	return &community.Post{
		ID:           m.ID,
		UserID:       m.UserID,
		Content:      m.Content,
		Images:       images,
		Tag:          m.Tag,
		LikeCount:    m.LikeCount,
		CommentCount: m.CommentCount,
		ShareCount:   m.ShareCount,
		IsDeleted:    m.IsDeleted,
	}
}

func postEntityToModel(p *community.Post) *PostModel {
	var images string
	if len(p.Images) > 0 {
		b, _ := json.Marshal(p.Images)
		images = string(b)
	}
	return &PostModel{
		ID:           p.ID,
		UserID:       p.UserID,
		Content:      p.Content,
		Images:       images,
		Tag:          p.Tag,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		ShareCount:   p.ShareCount,
		IsDeleted:    p.IsDeleted,
	}
}

func (r *PostRepositoryImpl) FindByID(ctx context.Context, id int64) (*community.Post, error) {
	var model PostModel
	if err := r.db.WithContext(ctx).Where("is_deleted = false").First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return postModelToEntity(&model), nil
}

func (r *PostRepositoryImpl) Create(ctx context.Context, p *community.Post) error {
	model := postEntityToModel(p)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *PostRepositoryImpl) Update(ctx context.Context, p *community.Post) error {
	model := postEntityToModel(p)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *PostRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&PostModel{}, id).Error
}

func (r *PostRepositoryImpl) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*community.Post, int64, error) {
	var models []PostModel
	var total int64

	db := r.db.WithContext(ctx).Model(&PostModel{}).Where("user_id = ? AND is_deleted = false", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*community.Post, len(models))
	for i := range models {
		result[i] = postModelToEntity(&models[i])
	}
	return result, total, nil
}

func (r *PostRepositoryImpl) FindByTag(ctx context.Context, tag string, page, pageSize int) ([]*community.Post, int64, error) {
	var models []PostModel
	var total int64

	db := r.db.WithContext(ctx).Model(&PostModel{}).Where("tag = ? AND is_deleted = false", tag)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*community.Post, len(models))
	for i := range models {
		result[i] = postModelToEntity(&models[i])
	}
	return result, total, nil
}

func (r *PostRepositoryImpl) Search(ctx context.Context, query *community.PostSearchQuery) ([]*community.Post, int64, error) {
	var models []PostModel
	var total int64

	db := r.db.WithContext(ctx).Model(&PostModel{}).Where("is_deleted = false")

	if query.Keyword != "" {
		db = db.Where("content LIKE ?", "%"+query.Keyword+"%")
	}
	if query.UserID != nil {
		db = db.Where("user_id = ?", *query.UserID)
	}
	if query.Tag != "" {
		db = db.Where("tag = ?", query.Tag)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.SortBy != "" {
		order := query.SortBy
		if query.SortDesc {
			order += " DESC"
		} else {
			order += " ASC"
		}
		db = db.Order(order)
	} else {
		db = db.Order("created_at DESC")
	}

	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*community.Post, len(models))
	for i := range models {
		result[i] = postModelToEntity(&models[i])
	}
	return result, total, nil
}

// ===== 评论仓储实现 =====

type CommentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) community.CommentRepository {
	return &CommentRepositoryImpl{db: db}
}

func commentModelToEntity(m *CommentModel) *community.Comment {
	return &community.Comment{
		ID:        m.ID,
		PostID:    m.PostID,
		UserID:    m.UserID,
		ParentID:  m.ParentID,
		Content:   m.Content,
		LikeCount: m.LikeCount,
	}
}

func commentEntityToModel(c *community.Comment) *CommentModel {
	return &CommentModel{
		ID:        c.ID,
		PostID:    c.PostID,
		UserID:    c.UserID,
		ParentID:  c.ParentID,
		Content:   c.Content,
		LikeCount: c.LikeCount,
	}
}

func (r *CommentRepositoryImpl) FindByID(ctx context.Context, id int64) (*community.Comment, error) {
	var model CommentModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return commentModelToEntity(&model), nil
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, c *community.Comment) error {
	model := commentEntityToModel(c)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *CommentRepositoryImpl) Update(ctx context.Context, c *community.Comment) error {
	model := commentEntityToModel(c)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *CommentRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&CommentModel{}, id).Error
}

func (r *CommentRepositoryImpl) FindByPostID(ctx context.Context, postID int64) ([]*community.Comment, error) {
	var models []CommentModel
	if err := r.db.WithContext(ctx).Where("post_id = ?", postID).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*community.Comment, len(models))
	for i := range models {
		result[i] = commentModelToEntity(&models[i])
	}
	return result, nil
}

func (r *CommentRepositoryImpl) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*community.Comment, int64, error) {
	var models []CommentModel
	var total int64

	db := r.db.WithContext(ctx).Model(&CommentModel{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*community.Comment, len(models))
	for i := range models {
		result[i] = commentModelToEntity(&models[i])
	}
	return result, total, nil
}

// ===== 私信仓储实现 =====

type MessageRepositoryImpl struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) community.MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

func messageModelToEntity(m *MessageModel) *community.Message {
	msg := &community.Message{
		ID:         m.ID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Content:    m.Content,
		IsRead:     m.IsRead,
	}
	if m.ReadAt != "" {
		msg.ReadAt = nil
	}
	return msg
}

func messageEntityToModel(m *community.Message) *MessageModel {
	return &MessageModel{
		ID:         m.ID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Content:    m.Content,
		IsRead:     m.IsRead,
	}
}

func (r *MessageRepositoryImpl) FindByID(ctx context.Context, id int64) (*community.Message, error) {
	var model MessageModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return messageModelToEntity(&model), nil
}

func (r *MessageRepositoryImpl) Create(ctx context.Context, m *community.Message) error {
	model := messageEntityToModel(m)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *MessageRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&MessageModel{}, id).Error
}

func (r *MessageRepositoryImpl) FindConversation(ctx context.Context, userID1, userID2 int64, page, pageSize int) ([]*community.Message, int64, error) {
	var models []MessageModel
	var total int64

	db := r.db.WithContext(ctx).Model(&MessageModel{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID1, userID2, userID2, userID1)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*community.Message, len(models))
	for i := range models {
		result[i] = messageModelToEntity(&models[i])
	}
	return result, total, nil
}

func (r *MessageRepositoryImpl) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*community.Message, int64, error) {
	var models []MessageModel
	var total int64

	db := r.db.WithContext(ctx).Model(&MessageModel{}).
		Where("sender_id = ? OR receiver_id = ?", userID, userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*community.Message, len(models))
	for i := range models {
		result[i] = messageModelToEntity(&models[i])
	}
	return result, total, nil
}

func (r *MessageRepositoryImpl) MarkAsRead(ctx context.Context, messageID int64) error {
	return r.db.WithContext(ctx).Model(&MessageModel{}).Where("id = ?", messageID).
		Updates(map[string]interface{}{"is_read": true, "read_at": gorm.Expr("NOW()")}).Error
}
