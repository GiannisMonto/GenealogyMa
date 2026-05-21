package genealogy

import (
	"context"
)

// Repository 族谱仓储接口
type Repository interface {
	// 基础CRUD
	FindByID(ctx context.Context, id int64) (*Genealogy, error)
	Create(ctx context.Context, genealogy *Genealogy) error
	Update(ctx context.Context, genealogy *Genealogy) error
	Delete(ctx context.Context, id int64) error

	// 查询
	FindAll(ctx context.Context) ([]*Genealogy, error)
	Search(ctx context.Context, query *SearchQuery) ([]*Genealogy, int64, error)
	FindBySurname(ctx context.Context, surname string) ([]*Genealogy, error)
}

// BranchRepository 分支仓储接口
type BranchRepository interface {
	FindByID(ctx context.Context, id int64) (*Branch, error)
	FindByGenealogyID(ctx context.Context, genealogyID int64) ([]*Branch, error)
	Create(ctx context.Context, branch *Branch) error
	Update(ctx context.Context, branch *Branch) error
	Delete(ctx context.Context, id int64) error
	DeleteByGenealogyID(ctx context.Context, genealogyID int64) error
	Search(ctx context.Context, query *BranchSearchQuery) ([]*Branch, int64, error)
}

// GenerationRepository 世代仓储接口
type GenerationRepository interface {
	FindByID(ctx context.Context, id int64) (*Generation, error)
	FindByGenealogyID(ctx context.Context, genealogyID int64) ([]*Generation, error)
	FindByGeneration(ctx context.Context, genealogyID int64, generation int) (*Generation, error)
	Create(ctx context.Context, generation *Generation) error
	Update(ctx context.Context, generation *Generation) error
	Delete(ctx context.Context, id int64) error
	DeleteByGenealogyID(ctx context.Context, genealogyID int64) error
}

// SearchQuery 族谱搜索查询参数
type SearchQuery struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`
	SortDesc bool   `json:"sort_desc"`
}

// BranchSearchQuery 分支搜索查询参数
type BranchSearchQuery struct {
	GenealogyID *int64
	Name        string
	Code        string
	Page        int
	PageSize    int
}