package memorial

import (
	"context"
)

// Repository 宗祠仓储接口
type Repository interface {
	// 基础CRUD
	FindByID(ctx context.Context, id int64) (*MemorialHall, error)
	Create(ctx context.Context, hall *MemorialHall) error
	Update(ctx context.Context, hall *MemorialHall) error
	Delete(ctx context.Context, id int64) error

	// 查询
	FindAll(ctx context.Context) ([]*MemorialHall, error)
	Search(ctx context.Context, query *SearchQuery) ([]*MemorialHall, int64, error)
	FindByRegion(ctx context.Context, province, city string) ([]*MemorialHall, error)
}

// TabletRepository 牌位仓储接口
type TabletRepository interface {
	FindByID(ctx context.Context, id int64) (*MemorialTablet, error)
	FindByHallID(ctx context.Context, hallID int64) ([]*MemorialTablet, error)
	FindByPersonID(ctx context.Context, personID int64) (*MemorialTablet, error)
	Create(ctx context.Context, tablet *MemorialTablet) error
	Update(ctx context.Context, tablet *MemorialTablet) error
	Delete(ctx context.Context, id int64) error
	DeleteByHallID(ctx context.Context, hallID int64) error
	Search(ctx context.Context, query *TabletSearchQuery) ([]*MemorialTablet, int64, error)
}

// SearchQuery 宗祠搜索查询参数
type SearchQuery struct {
	Name     string `json:"name"`
	Province string `json:"province"`
	City     string `json:"city"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`
	SortDesc bool   `json:"sort_desc"`
}

// TabletSearchQuery 牌位搜索查询参数
type TabletSearchQuery struct {
	HallID     *int64
	PersonID   *int64
	TabletType *TabletType
	Floor      *int
	Page       int
	PageSize   int
}