package config

import (
	"context"
	"fmt"
)

// Service 配置领域服务
type Service struct {
	repo Repository
}

// NewService 创建配置服务
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateConfig 创建配置
func (s *Service) CreateConfig(ctx context.Context, config *Config) error {
	if err := config.Validate(); err != nil {
		return err
	}
	// 检查是否已存在相同的 group+key
	existing, err := s.repo.FindByKey(ctx, config.Group, config.Key)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("config with key %s already exists in group %s", config.Key, config.Group)
	}
	return s.repo.Create(ctx, config)
}

// UpdateConfig 更新配置
func (s *Service) UpdateConfig(ctx context.Context, config *Config) error {
	if err := config.Validate(); err != nil {
		return err
	}
	existing, err := s.repo.FindByID(ctx, config.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("config not found")
	}
	// 系统配置不允许修改 group 和 key
	if existing.IsSystem && (existing.Group != config.Group || existing.Key != config.Key) {
		return fmt.Errorf("cannot change group or key of system config")
	}
	return s.repo.Update(ctx, config)
}

// DeleteConfig 删除配置
func (s *Service) DeleteConfig(ctx context.Context, id int64) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("config not found")
	}
	if !existing.CanDelete() {
		return fmt.Errorf("cannot delete system config")
	}
	return s.repo.Delete(ctx, id)
}

// GetConfig 获取配置详情
func (s *Service) GetConfig(ctx context.Context, id int64) (*Config, error) {
	return s.repo.FindByID(ctx, id)
}

// ListConfigs 获取所有配置
func (s *Service) ListConfigs(ctx context.Context) ([]*Config, error) {
	return s.repo.FindAll(ctx)
}

// ListConfigsByGroup 获取指定分组的配置
func (s *Service) ListConfigsByGroup(ctx context.Context, group string) ([]*Config, error) {
	return s.repo.FindByGroup(ctx, group)
}

// GetConfigByKey 根据 group 和 key 获取配置
func (s *Service) GetConfigByKey(ctx context.Context, group, key string) (*Config, error) {
	return s.repo.FindByKey(ctx, group, key)
}

// SearchConfigs 搜索配置
func (s *Service) SearchConfigs(ctx context.Context, query *SearchQuery) ([]*Config, int64, error) {
	return s.repo.Search(ctx, query)
}
