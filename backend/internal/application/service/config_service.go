package service

import (
	"context"

	"github.com/genealogy-ma/platform/internal/domain/config"
)

// ConfigService 配置应用服务
type ConfigService struct {
	domainService *config.Service
}

// NewConfigService 创建配置应用服务
func NewConfigService(repo config.Repository) *ConfigService {
	return &ConfigService{
		domainService: config.NewService(repo),
	}
}

// ===== DTO 定义 =====

// ConfigDTO 配置DTO
type ConfigDTO struct {
	ID        int64  `json:"id"`
	Group     string `json:"group"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Type      string `json:"type"`
	Label     string `json:"label"`
	HelpText  string `json:"help_text"`
	SortOrder int    `json:"sort_order"`
	IsSystem  bool   `json:"is_system"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateConfigRequest 创建配置请求
type CreateConfigRequest struct {
	Group     string `json:"group" binding:"required"`
	Key       string `json:"key" binding:"required"`
	Value     string `json:"value"`
	Type      string `json:"type" binding:"required"`
	Label     string `json:"label"`
	HelpText  string `json:"help_text"`
	SortOrder int    `json:"sort_order"`
	IsSystem  bool   `json:"is_system"`
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	Value     *string `json:"value"`
	Type      *string `json:"type"`
	Label     *string `json:"label"`
	HelpText  *string `json:"help_text"`
	SortOrder *int    `json:"sort_order"`
}

// ConfigSearchRequest 配置搜索请求
type ConfigSearchRequest struct {
	Group   string `json:"group"`
	Key     string `json:"key"`
	Label   string `json:"label"`
	Page    int    `json:"page"`
	PageSize int   `json:"page_size"`
}

// ===== 服务方法 =====

// CreateConfig 创建配置
func (s *ConfigService) CreateConfig(ctx context.Context, req *CreateConfigRequest) (*ConfigDTO, error) {
	c := &config.Config{
		Group:     req.Group,
		Key:       req.Key,
		Value:     req.Value,
		Type:      config.ConfigType(req.Type),
		Label:     req.Label,
		HelpText:  req.HelpText,
		SortOrder: req.SortOrder,
		IsSystem:  req.IsSystem,
	}

	if err := s.domainService.CreateConfig(ctx, c); err != nil {
		return nil, err
	}

	return configToDTO(c), nil
}

// UpdateConfig 更新配置
func (s *ConfigService) UpdateConfig(ctx context.Context, id int64, req *UpdateConfigRequest) (*ConfigDTO, error) {
	existing, err := s.domainService.GetConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	// 应用更新
	if req.Value != nil {
		existing.Value = *req.Value
	}
	if req.Type != nil {
		existing.Type = config.ConfigType(*req.Type)
	}
	if req.Label != nil {
		existing.Label = *req.Label
	}
	if req.HelpText != nil {
		existing.HelpText = *req.HelpText
	}
	if req.SortOrder != nil {
		existing.SortOrder = *req.SortOrder
	}

	if err := s.domainService.UpdateConfig(ctx, existing); err != nil {
		return nil, err
	}

	return configToDTO(existing), nil
}

// DeleteConfig 删除配置
func (s *ConfigService) DeleteConfig(ctx context.Context, id int64) error {
	return s.domainService.DeleteConfig(ctx, id)
}

// GetConfig 获取配置详情
func (s *ConfigService) GetConfig(ctx context.Context, id int64) (*ConfigDTO, error) {
	c, err := s.domainService.GetConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	return configToDTO(c), nil
}

// ListConfigs 获取所有配置
func (s *ConfigService) ListConfigs(ctx context.Context) ([]*ConfigDTO, error) {
	configs, err := s.domainService.ListConfigs(ctx)
	if err != nil {
		return nil, err
	}
	dtos := make([]*ConfigDTO, len(configs))
	for i, c := range configs {
		dtos[i] = configToDTO(c)
	}
	return dtos, nil
}

// ListConfigsByGroup 获取指定分组的配置
func (s *ConfigService) ListConfigsByGroup(ctx context.Context, group string) ([]*ConfigDTO, error) {
	configs, err := s.domainService.ListConfigsByGroup(ctx, group)
	if err != nil {
		return nil, err
	}
	dtos := make([]*ConfigDTO, len(configs))
	for i, c := range configs {
		dtos[i] = configToDTO(c)
	}
	return dtos, nil
}

// GetConfigByKey 根据 group 和 key 获取配置
func (s *ConfigService) GetConfigByKey(ctx context.Context, group, key string) (*ConfigDTO, error) {
	c, err := s.domainService.GetConfigByKey(ctx, group, key)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	return configToDTO(c), nil
}

// SearchConfigs 搜索配置
func (s *ConfigService) SearchConfigs(ctx context.Context, req *ConfigSearchRequest) ([]*ConfigDTO, int64, error) {
	query := &config.SearchQuery{
		Group:   req.Group,
		Key:     req.Key,
		Label:   req.Label,
		Page:    req.Page,
		PageSize: req.PageSize,
	}
	configs, total, err := s.domainService.SearchConfigs(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]*ConfigDTO, len(configs))
	for i, c := range configs {
		dtos[i] = configToDTO(c)
	}
	return dtos, total, nil
}

// ===== 辅助函数 =====

func configToDTO(c *config.Config) *ConfigDTO {
	return &ConfigDTO{
		ID:        c.ID,
		Group:     c.Group,
		Key:       c.Key,
		Value:     c.Value,
		Type:      string(c.Type),
		Label:     c.Label,
		HelpText:  c.HelpText,
		SortOrder: c.SortOrder,
		IsSystem:  c.IsSystem,
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
