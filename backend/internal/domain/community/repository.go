package community

import (
	"context"
)

// UserProfileRepository 用户资料仓储接口
type UserProfileRepository interface {
	FindByID(ctx context.Context, id int64) (*UserProfile, error)
	FindByUserID(ctx context.Context, userID int64) (*UserProfile, error)
	Create(ctx context.Context, profile *UserProfile) error
	Update(ctx context.Context, profile *UserProfile) error
	Delete(ctx context.Context, id int64) error
}

// PostRepository 动态仓储接口
type PostRepository interface {
	FindByID(ctx context.Context, id int64) (*Post, error)
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int64) error
	FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*Post, int64, error)
	FindByTag(ctx context.Context, tag string, page, pageSize int) ([]*Post, int64, error)
	Search(ctx context.Context, query *PostSearchQuery) ([]*Post, int64, error)
}

// CommentRepository 评论仓储接口
type CommentRepository interface {
	FindByID(ctx context.Context, id int64) (*Comment, error)
	Create(ctx context.Context, comment *Comment) error
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id int64) error
	FindByPostID(ctx context.Context, postID int64) ([]*Comment, error)
	FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*Comment, int64, error)
}

// MessageRepository 私信仓储接口
type MessageRepository interface {
	FindByID(ctx context.Context, id int64) (*Message, error)
	Create(ctx context.Context, msg *Message) error
	Delete(ctx context.Context, id int64) error
	FindConversation(ctx context.Context, userID1, userID2 int64, page, pageSize int) ([]*Message, int64, error)
	FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*Message, int64, error)
	MarkAsRead(ctx context.Context, messageID int64) error
}

// PostSearchQuery 动态搜索查询
type PostSearchQuery struct {
	Keyword   string
	UserID    *int64
	Tag       string
	Page      int
	PageSize  int
	SortBy    string
	SortDesc  bool
}