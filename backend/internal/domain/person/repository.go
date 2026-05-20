package person

import (
	"context"
)

// Repository 人物仓储接口
type Repository interface {
	// 基础CRUD
	FindByID(ctx context.Context, id int64) (*Person, error)
	FindByLegacyID(ctx context.Context, legacyID string) (*Person, error)
	Create(ctx context.Context, person *Person) error
	Update(ctx context.Context, person *Person) error
	Delete(ctx context.Context, id int64) error

	// 查询
	Search(ctx context.Context, query *SearchQuery) ([]*Person, int64, error)
	FindByGeneration(ctx context.Context, generation int) ([]*Person, error)
	FindByName(ctx context.Context, name string, fuzzy bool) ([]*Person, error)

	// 世系树查询
	FindAncestors(ctx context.Context, personID int64, depth int) ([]*Person, error)
	FindDescendants(ctx context.Context, personID int64, depth int) ([]*Person, error)
	FindSiblings(ctx context.Context, personID int64) ([]*Person, error)
	GetFamilyTree(ctx context.Context, rootID int64, depth int) (*TreeResult, error)

	// 关系查询
	FindChildren(ctx context.Context, parentID int64) ([]*PersonChild, error)
	FindSpouses(ctx context.Context, personID int64) ([]*Spouse, error)
	FindParents(ctx context.Context, personID int64) ([]*Person, error)

	// 统计
	Count(ctx context.Context) (int64, error)
	CountByGeneration(ctx context.Context) (map[int]int64, error)

	// 批量操作
	BatchCreate(ctx context.Context, persons []*Person) error
	BatchUpdate(ctx context.Context, persons []*Person) error
}

// SearchQuery 搜索查询参数
type SearchQuery struct {
	Keyword    string   `json:"keyword"`
	Name       string   `json:"name"`
	StyleName  string   `json:"style_name"`
	Gender     *Gender  `json:"gender"`
	Generation *int     `json:"generation"`
	BirthYear  *int     `json:"birth_year"`
	DeathYear  *int     `json:"death_year"`
	BirthPlace string   `json:"birth_place"`

	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	SortBy   string `json:"sort_by"`
	SortDesc bool   `json:"sort_desc"`
}

// TreeResult 族谱树结果
type TreeResult struct {
	Root      *Person   `json:"root"`
	Nodes     []*Person `json:"nodes"`
	Relations []*Relation `json:"relations"`
	MaxDepth  int       `json:"max_depth"`
	TotalNodes int      `json:"total_nodes"`
}

// Relation 树关系
type Relation struct {
	ParentID int64 `json:"parent_id"`
	ChildID  int64 `json:"child_id"`
	Type     string `json:"type"`
}
