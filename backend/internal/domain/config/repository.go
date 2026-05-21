package config

import (
	"context"
)

// Repository 配置仓储接口
type Repository interface {
	// 基础CRUD
	FindByID(ctx context.Context, id int64) (*Config, error)
	Create(ctx context.Context, config *Config) error
	Update(ctx context.Context, config *Config) error
	Delete(ctx context.Context, id int64) error

	// 查询
	FindAll(ctx context.Context) ([]*Config, error)
	FindByGroup(ctx context.Context, group string) ([]*Config, error)
	FindByKey(ctx context.Context, group, key string) (*Config, error)
	Search(ctx context.Context, query *SearchQuery) ([]*Config, int64, error)
}
