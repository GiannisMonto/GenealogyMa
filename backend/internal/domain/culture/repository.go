package culture

import (
	"context"
)

// Repository 文献仓储接口
type Repository interface {
	FindByID(ctx context.Context, id int64) (*Document, error)
	Create(ctx context.Context, doc *Document) error
	Update(ctx context.Context, doc *Document) error
	Delete(ctx context.Context, id int64) error
	FindAll(ctx context.Context) ([]*Document, error)
	Search(ctx context.Context, query *SearchQuery) ([]*Document, int64, error)
	FindByCategory(ctx context.Context, category DocumentCategory) ([]*Document, error)
}

// SearchQuery 搜索查询参数
type SearchQuery struct {
	Keyword   string
	Category  DocumentCategory
	Author    string
	Page      int
	PageSize  int
	SortBy    string
	SortDesc  bool
}

// StoryRepository 故事仓储接口
type StoryRepository interface {
	FindByID(ctx context.Context, id int64) (*Story, error)
	Create(ctx context.Context, story *Story) error
	Update(ctx context.Context, story *Story) error
	Delete(ctx context.Context, id int64) error
	FindAll(ctx context.Context) ([]*Story, error)
	Search(ctx context.Context, query *StorySearchQuery) ([]*Story, int64, error)
	FindByEra(ctx context.Context, era string) ([]*Story, error)
}

// StorySearchQuery 故事搜索查询参数
type StorySearchQuery struct {
	Keyword    string
	Era        string
	Page       int
	PageSize   int
	SortBy     string
	SortDesc   bool
}

// FamilyTeachingsRepository 家训仓储接口
type FamilyTeachingsRepository interface {
	FindByID(ctx context.Context, id int64) (*FamilyTeachings, error)
	Create(ctx context.Context, ft *FamilyTeachings) error
	Update(ctx context.Context, ft *FamilyTeachings) error
	Delete(ctx context.Context, id int64) error
	FindAll(ctx context.Context) ([]*FamilyTeachings, error)
	FindByGeneration(ctx context.Context, generation int) ([]*FamilyTeachings, error)
}