package service

import (
	"context"
	"testing"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/config"
)

// MockConfigRepository 模拟配置仓储
type MockConfigRepository struct {
	configs map[int64]*config.Config
	nextID  int64
}

func NewMockConfigRepository() *MockConfigRepository {
	return &MockConfigRepository{
		configs: make(map[int64]*config.Config),
		nextID:  1,
	}
}

func (m *MockConfigRepository) FindByID(ctx context.Context, id int64) (*config.Config, error) {
	if c, ok := m.configs[id]; ok {
		copy := *c
		return &copy, nil
	}
	return nil, nil
}

func (m *MockConfigRepository) Create(ctx context.Context, c *config.Config) error {
	c.ID = m.nextID
	m.nextID++
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	copy := *c
	m.configs[c.ID] = &copy
	return nil
}

func (m *MockConfigRepository) Update(ctx context.Context, c *config.Config) error {
	if _, ok := m.configs[c.ID]; !ok {
		return nil
	}
	c.UpdatedAt = time.Now()
	m.configs[c.ID] = c
	return nil
}

func (m *MockConfigRepository) Delete(ctx context.Context, id int64) error {
	delete(m.configs, id)
	return nil
}

func (m *MockConfigRepository) FindAll(ctx context.Context) ([]*config.Config, error) {
	result := make([]*config.Config, 0, len(m.configs))
	for _, c := range m.configs {
		result = append(result, c)
	}
	return result, nil
}

func (m *MockConfigRepository) FindByGroup(ctx context.Context, group string) ([]*config.Config, error) {
	result := make([]*config.Config, 0)
	for _, c := range m.configs {
		if c.Group == group {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *MockConfigRepository) FindByKey(ctx context.Context, group, key string) (*config.Config, error) {
	for _, c := range m.configs {
		if c.Group == group && c.Key == key {
			copy := *c
			return &copy, nil
		}
	}
	return nil, nil
}

func (m *MockConfigRepository) Search(ctx context.Context, query *config.SearchQuery) ([]*config.Config, int64, error) {
	result := make([]*config.Config, 0)
	for _, c := range m.configs {
		if query.Group != "" && c.Group != query.Group {
			continue
		}
		if query.Key != "" && c.Key != query.Key {
			continue
		}
		result = append(result, c)
	}

	// Apply pagination
	total := int64(len(result))
	if query.Page > 0 && query.PageSize > 0 {
		start := (query.Page - 1) * query.PageSize
		if start >= len(result) {
			return []*config.Config{}, total, nil
		}
		end := start + query.PageSize
		if end > len(result) {
			end = len(result)
		}
		return result[start:end], total, nil
	}
	return result, total, nil
}

// ===== 单元测试 =====

func TestConfigService_CreateConfig(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewConfigService(repo)
	ctx := context.Background()

	t.Run("create valid config", func(t *testing.T) {
		req := &CreateConfigRequest{
			Group:     "system",
			Key:       "site_name",
			Value:     "My Site",
			Type:      "string",
			Label:     "网站名称",
			HelpText:  "网站显示名称",
			SortOrder: 1,
			IsSystem:  false,
		}
		result, err := svc.CreateConfig(ctx, req)
		if err != nil {
			t.Errorf("CreateConfig() error = %v", err)
			return
		}
		if result == nil {
			t.Error("CreateConfig() should return config")
			return
		}
		if result.ID == 0 {
			t.Error("CreateConfig() should set ID")
		}
		if result.Key != req.Key {
			t.Errorf("CreateConfig() key = %v, want %v", result.Key, req.Key)
		}
		if result.Group != req.Group {
			t.Errorf("CreateConfig() group = %v, want %v", result.Group, req.Group)
		}
		if result.Type != req.Type {
			t.Errorf("CreateConfig() type = %v, want %v", result.Type, req.Type)
		}
		if result.Label != req.Label {
			t.Errorf("CreateConfig() label = %v, want %v", result.Label, req.Label)
		}
	})

	t.Run("create system config", func(t *testing.T) {
		req := &CreateConfigRequest{
			Group:    "system",
			Key:      "system_mode",
			Value:    "true",
			Type:     "boolean",
			IsSystem: true,
		}
		result, err := svc.CreateConfig(ctx, req)
		if err != nil {
			t.Errorf("CreateConfig() system config error = %v", err)
			return
		}
		if result == nil {
			t.Error("CreateConfig() should return config")
			return
		}
		if !result.IsSystem {
			t.Error("CreateConfig() IsSystem should be true")
		}
	})

	t.Run("create config without optional fields", func(t *testing.T) {
		req := &CreateConfigRequest{
			Group: "custom",
			Key:   "custom_key",
			Type:  "string",
		}
		result, err := svc.CreateConfig(ctx, req)
		if err != nil {
			t.Errorf("CreateConfig() without optional fields error = %v", err)
			return
		}
		if result == nil {
			t.Error("CreateConfig() should return config")
			return
		}
		if result.Value != "" {
			t.Errorf("CreateConfig() value = %v, want empty", result.Value)
		}
	})
}

func TestConfigService_UpdateConfig(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewConfigService(repo)
	ctx := context.Background()

	// 创建初始配置
	createReq := &CreateConfigRequest{
		Group:     "system",
		Key:       "site_name",
		Value:     "Old Site",
		Type:      "string",
		Label:     "旧网站名称",
		SortOrder: 1,
	}
	created, _ := svc.CreateConfig(ctx, createReq)
	configID := created.ID

	t.Run("update existing config value", func(t *testing.T) {
		newValue := "New Site"
		req := &UpdateConfigRequest{
			Value: &newValue,
		}
		result, err := svc.UpdateConfig(ctx, configID, req)
		if err != nil {
			t.Errorf("UpdateConfig() error = %v", err)
			return
		}
		if result == nil {
			t.Error("UpdateConfig() should return config")
			return
		}
		if result.Value != newValue {
			t.Errorf("UpdateConfig() value = %v, want %v", result.Value, newValue)
		}
	})

	t.Run("update config label", func(t *testing.T) {
		newLabel := "新网站名称"
		req := &UpdateConfigRequest{
			Label: &newLabel,
		}
		result, err := svc.UpdateConfig(ctx, configID, req)
		if err != nil {
			t.Errorf("UpdateConfig() label error = %v", err)
			return
		}
		if result == nil {
			t.Error("UpdateConfig() should return config")
			return
		}
		if result.Label != newLabel {
			t.Errorf("UpdateConfig() label = %v, want %v", result.Label, newLabel)
		}
	})

	t.Run("update config type", func(t *testing.T) {
		newType := "number"
		req := &UpdateConfigRequest{
			Type: &newType,
		}
		result, err := svc.UpdateConfig(ctx, configID, req)
		if err != nil {
			t.Errorf("UpdateConfig() type error = %v", err)
			return
		}
		if result == nil {
			t.Error("UpdateConfig() should return config")
			return
		}
		if result.Type != newType {
			t.Errorf("UpdateConfig() type = %v, want %v", result.Type, newType)
		}
	})

	t.Run("update non-existent config", func(t *testing.T) {
		newValue := "test"
		req := &UpdateConfigRequest{
			Value: &newValue,
		}
		result, err := svc.UpdateConfig(ctx, 9999, req)
		if err != nil {
			t.Errorf("UpdateConfig() non-existent error = %v", err)
			return
		}
		if result != nil {
			t.Error("UpdateConfig() should return nil for non-existent config")
		}
	})
}

func TestConfigService_DeleteConfig(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewConfigService(repo)
	ctx := context.Background()

	t.Run("delete existing config", func(t *testing.T) {
		req := &CreateConfigRequest{
			Group: "custom",
			Key:   "to_delete",
			Type:  "string",
		}
		created, _ := svc.CreateConfig(ctx, req)

		err := svc.DeleteConfig(ctx, created.ID)
		if err != nil {
			t.Errorf("DeleteConfig() error = %v", err)
			return
		}

		// Verify deleted
		deleted, _ := svc.GetConfig(ctx, created.ID)
		if deleted != nil {
			t.Error("DeleteConfig() should remove config")
		}
	})

	t.Run("delete non-existent config", func(t *testing.T) {
		err := svc.DeleteConfig(ctx, 9999)
		if err == nil {
			t.Error("DeleteConfig() should return error for non-existent config")
		}
	})
}

func TestConfigService_GetConfig(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewConfigService(repo)
	ctx := context.Background()

	t.Run("get existing config", func(t *testing.T) {
		req := &CreateConfigRequest{
			Group:  "system",
			Key:    "site_name",
			Value:  "My Site",
			Type:   "string",
			Label:  "网站名称",
		}
		created, _ := svc.CreateConfig(ctx, req)

		result, err := svc.GetConfig(ctx, created.ID)
		if err != nil {
			t.Errorf("GetConfig() error = %v", err)
			return
		}
		if result == nil {
			t.Error("GetConfig() should return config")
			return
		}
		if result.Key != req.Key {
			t.Errorf("GetConfig() key = %v, want %v", result.Key, req.Key)
		}
		if result.Value != req.Value {
			t.Errorf("GetConfig() value = %v, want %v", result.Value, req.Value)
		}
	})

	t.Run("get non-existent config", func(t *testing.T) {
		result, err := svc.GetConfig(ctx, 9999)
		if err != nil {
			t.Errorf("GetConfig() non-existent error = %v", err)
			return
		}
		if result != nil {
			t.Error("GetConfig() should return nil for non-existent config")
		}
	})
}

func TestConfigService_ListConfigs(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewConfigService(repo)
	ctx := context.Background()

	t.Run("list all configs", func(t *testing.T) {
		_, _ = svc.CreateConfig(ctx, &CreateConfigRequest{Group: "system", Key: "k1", Type: "string"})
		_, _ = svc.CreateConfig(ctx, &CreateConfigRequest{Group: "system", Key: "k2", Type: "string"})
		_, _ = svc.CreateConfig(ctx, &CreateConfigRequest{Group: "user", Key: "k3", Type: "string"})

		result, err := svc.ListConfigs(ctx)
		if err != nil {
			t.Errorf("ListConfigs() error = %v", err)
			return
		}
		if len(result) != 3 {
			t.Errorf("ListConfigs() count = %v, want 3", len(result))
		}
	})

	t.Run("list configs when empty", func(t *testing.T) {
		// Create a fresh service with empty repo
		emptyRepo := NewMockConfigRepository()
		emptySvc := NewConfigService(emptyRepo)

		result, err := emptySvc.ListConfigs(ctx)
		if err != nil {
			t.Errorf("ListConfigs() empty error = %v", err)
			return
		}
		if len(result) != 0 {
			t.Errorf("ListConfigs() empty count = %v, want 0", len(result))
		}
	})
}

func TestConfigService_ListConfigsByGroup(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewConfigService(repo)
	ctx := context.Background()

	_, _ = svc.CreateConfig(ctx, &CreateConfigRequest{Group: "system", Key: "k1", Type: "string"})
	_, _ = svc.CreateConfig(ctx, &CreateConfigRequest{Group: "system", Key: "k2", Type: "string"})
	_, _ = svc.CreateConfig(ctx, &CreateConfigRequest{Group: "cemetery", Key: "k3", Type: "string"})

	t.Run("list configs by group", func(t *testing.T) {
		result, err := svc.ListConfigsByGroup(ctx, "system")
		if err != nil {
			t.Errorf("ListConfigsByGroup() error = %v", err)
			return
		}
		if len(result) != 2 {
			t.Errorf("ListConfigsByGroup() count = %v, want 2", len(result))
		}
	})

	t.Run("list configs by group with no results", func(t *testing.T) {
		result, err := svc.ListConfigsByGroup(ctx, "non_existent")
		if err != nil {
			t.Errorf("ListConfigsByGroup() no results error = %v", err)
			return
		}
		if len(result) != 0 {
			t.Errorf("ListConfigsByGroup() no results count = %v, want 0", len(result))
		}
	})
}

func TestConfigService_GetConfigByKey(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewConfigService(repo)
	ctx := context.Background()

	_, _ = svc.CreateConfig(ctx, &CreateConfigRequest{
		Group: "system",
		Key:   "site_name",
		Value: "My Site",
		Type:  "string",
	})

	t.Run("get config by group and key", func(t *testing.T) {
		result, err := svc.GetConfigByKey(ctx, "system", "site_name")
		if err != nil {
			t.Errorf("GetConfigByKey() error = %v", err)
			return
		}
		if result == nil {
			t.Error("GetConfigByKey() should return config")
			return
		}
		if result.Value != "My Site" {
			t.Errorf("GetConfigByKey() value = %v, want %v", result.Value, "My Site")
		}
	})

	t.Run("get non-existent config by key", func(t *testing.T) {
		result, err := svc.GetConfigByKey(ctx, "system", "non_existent")
		if err != nil {
			t.Errorf("GetConfigByKey() non-existent error = %v", err)
			return
		}
		if result != nil {
			t.Error("GetConfigByKey() should return nil for non-existent key")
		}
	})

	t.Run("get config with wrong group", func(t *testing.T) {
		result, err := svc.GetConfigByKey(ctx, "wrong_group", "site_name")
		if err != nil {
			t.Errorf("GetConfigByKey() wrong group error = %v", err)
			return
		}
		if result != nil {
			t.Error("GetConfigByKey() should return nil for wrong group")
		}
	})
}

func TestConfigService_SearchConfigs(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewConfigService(repo)
	ctx := context.Background()

	_, _ = svc.CreateConfig(ctx, &CreateConfigRequest{Group: "system", Key: "site_name", Type: "string", Label: "网站名称"})
	_, _ = svc.CreateConfig(ctx, &CreateConfigRequest{Group: "system", Key: "site_url", Type: "string", Label: "网站地址"})
	_, _ = svc.CreateConfig(ctx, &CreateConfigRequest{Group: "cemetery", Key: "cemetery_name", Type: "string", Label: "墓园名称"})

	t.Run("search by group", func(t *testing.T) {
		req := &ConfigSearchRequest{Group: "system"}
		results, total, err := svc.SearchConfigs(ctx, req)
		if err != nil {
			t.Errorf("SearchConfigs() error = %v", err)
			return
		}
		if total != 2 {
			t.Errorf("SearchConfigs() total = %v, want 2", total)
		}
		if len(results) != 2 {
			t.Errorf("SearchConfigs() count = %v, want 2", len(results))
		}
	})

	t.Run("search by key", func(t *testing.T) {
		req := &ConfigSearchRequest{Key: "site_name"}
		results, total, err := svc.SearchConfigs(ctx, req)
		if err != nil {
			t.Errorf("SearchConfigs() by key error = %v", err)
			return
		}
		if total != 1 {
			t.Errorf("SearchConfigs() by key total = %v, want 1", total)
		}
		if len(results) != 1 {
			t.Errorf("SearchConfigs() by key count = %v, want 1", len(results))
		}
	})

	t.Run("search with pagination", func(t *testing.T) {
		req := &ConfigSearchRequest{Group: "system", Page: 1, PageSize: 1}
		results, total, err := svc.SearchConfigs(ctx, req)
		if err != nil {
			t.Errorf("SearchConfigs() pagination error = %v", err)
			return
		}
		if total != 2 {
			t.Errorf("SearchConfigs() pagination total = %v, want 2", total)
		}
		if len(results) != 1 {
			t.Errorf("SearchConfigs() pagination count = %v, want 1", len(results))
		}
	})

	t.Run("search with no results", func(t *testing.T) {
		req := &ConfigSearchRequest{Group: "non_existent"}
		results, total, err := svc.SearchConfigs(ctx, req)
		if err != nil {
			t.Errorf("SearchConfigs() no results error = %v", err)
			return
		}
		if total != 0 {
			t.Errorf("SearchConfigs() no results total = %v, want 0", total)
		}
		if len(results) != 0 {
			t.Errorf("SearchConfigs() no results count = %v, want 0", len(results))
		}
	})
}

func TestConfigService_DTOConversion(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewConfigService(repo)
	ctx := context.Background()

	t.Run("DTO contains all fields", func(t *testing.T) {
		req := &CreateConfigRequest{
			Group:     "system",
			Key:       "test_key",
			Value:     "test_value",
			Type:      "string",
			Label:     "Test Label",
			HelpText:  "Help text",
			SortOrder: 10,
			IsSystem:  true,
		}
		created, err := svc.CreateConfig(ctx, req)
		if err != nil {
			t.Errorf("CreateConfig() error = %v", err)
			return
		}

		if created.Group != req.Group {
			t.Errorf("DTO group = %v, want %v", created.Group, req.Group)
		}
		if created.Key != req.Key {
			t.Errorf("DTO key = %v, want %v", created.Key, req.Key)
		}
		if created.Value != req.Value {
			t.Errorf("DTO value = %v, want %v", created.Value, req.Value)
		}
		if created.Type != req.Type {
			t.Errorf("DTO type = %v, want %v", created.Type, req.Type)
		}
		if created.Label != req.Label {
			t.Errorf("DTO label = %v, want %v", created.Label, req.Label)
		}
		if created.HelpText != req.HelpText {
			t.Errorf("DTO help_text = %v, want %v", created.HelpText, req.HelpText)
		}
		if created.SortOrder != req.SortOrder {
			t.Errorf("DTO sort_order = %v, want %v", created.SortOrder, req.SortOrder)
		}
		if created.IsSystem != req.IsSystem {
			t.Errorf("DTO is_system = %v, want %v", created.IsSystem, req.IsSystem)
		}
		if created.CreatedAt == "" {
			t.Error("DTO created_at should not be empty")
		}
		if created.UpdatedAt == "" {
			t.Error("DTO updated_at should not be empty")
		}
	})
}
