package cemetery

import (
	"context"
)

// Repository 墓园仓储接口
type Repository interface {
	// 基础CRUD
	FindByID(ctx context.Context, id int64) (*Cemetery, error)
	Create(ctx context.Context, cemetery *Cemetery) error
	Update(ctx context.Context, cemetery *Cemetery) error
	Delete(ctx context.Context, id int64) error

	// 查询
	FindAll(ctx context.Context) ([]*Cemetery, error)
	Search(ctx context.Context, query *SearchQuery) ([]*Cemetery, int64, error)
	FindByRegion(ctx context.Context, province, city string) ([]*Cemetery, error)
}

// SearchQuery 搜索查询参数
type SearchQuery struct {
	Name       string `json:"name"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	SortBy     string `json:"sort_by"`
	SortDesc   bool   `json:"sort_desc"`
}

// GraveRepository 墓位仓储接口
type GraveRepository interface {
	FindByID(ctx context.Context, id int64) (*Grave, error)
	FindByCemeteryID(ctx context.Context, cemeteryID int64) ([]*Grave, error)
	FindByPersonID(ctx context.Context, personID int64) (*Grave, error)
	Create(ctx context.Context, grave *Grave) error
	Update(ctx context.Context, grave *Grave) error
	Delete(ctx context.Context, id int64) error
	DeleteByCemeteryID(ctx context.Context, cemeteryID int64) error
	Search(ctx context.Context, query *GraveSearchQuery) ([]*Grave, int64, error)
}

// GraveSearchQuery 墓位搜索查询参数
type GraveSearchQuery struct {
	CemeteryID *int64
	PersonID   *int64
	Section    string
	Row        *int
	Number     *int
	Page       int
	PageSize   int
}