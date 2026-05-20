package community

import (
	"fmt"
	"time"
)

// UserProfile 用户资料实体
type UserProfile struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Nickname     string    `json:"nickname"`
	AvatarURL    string    `json:"avatar_url"`
	Bio          string    `json:"bio"`
	Gender       string    `json:"gender"`
	BirthDate    *string   `json:"birth_date"`
	Province     string    `json:"province"`
	City         string    `json:"city"`
	FollowerCount int      `json:"follower_count"`
	FollowingCount int     `json:"following_count"`
	PostCount     int      `json:"post_count"`
	PrivacyLevel int       `json:"privacy_level"` // 0: 公开, 1: 仅好友, 2: 私密
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Validate 验证用户资料
func (p *UserProfile) Validate() error {
	if p.UserID <= 0 {
		return fmt.Errorf("user_id is required")
	}
	if p.Nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}
	if len(p.Nickname) > 50 {
		return fmt.Errorf("nickname too long (max 50)")
	}
	return nil
}

// Post 动态实体
type Post struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Content    string    `json:"content"`
	Images     []string  `json:"images"`
	Tag        string    `json:"tag"`
	LikeCount  int       `json:"like_count"`
	CommentCount int     `json:"comment_count"`
	ShareCount  int       `json:"share_count"`
	IsDeleted  bool      `json:"is_deleted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate 验证动态
func (p *Post) Validate() error {
	if p.UserID <= 0 {
		return fmt.Errorf("user_id is required")
	}
	if p.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	if len(p.Content) > 5000 {
		return fmt.Errorf("content too long (max 5000)")
	}
	return nil
}

// Comment 评论实体
type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	ParentID  *int64    `json:"parent_id"` // 回复的评论ID
	Content   string    `json:"content"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 验证评论
func (c *Comment) Validate() error {
	if c.PostID <= 0 {
		return fmt.Errorf("post_id is required")
	}
	if c.UserID <= 0 {
		return fmt.Errorf("user_id is required")
	}
	if c.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	if len(c.Content) > 1000 {
		return fmt.Errorf("content too long (max 1000)")
	}
	return nil
}

// Message 私信实体
type Message struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	ReadAt     *time.Time `json:"read_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Validate 验证私信
func (m *Message) Validate() error {
	if m.SenderID <= 0 {
		return fmt.Errorf("sender_id is required")
	}
	if m.ReceiverID <= 0 {
		return fmt.Errorf("receiver_id is required")
	}
	if m.SenderID == m.ReceiverID {
		return fmt.Errorf("cannot send message to yourself")
	}
	if m.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	if len(m.Content) > 2000 {
		return fmt.Errorf("content too long (max 2000)")
	}
	return nil
}

// IsRead 判断是否已读
func (m *Message) IsReadMessage() bool {
	return m.IsRead
}

// MarkAsRead 标记为已读
func (m *Message) MarkAsRead() {
	m.IsRead = true
	now := time.Now()
	m.ReadAt = &now
}