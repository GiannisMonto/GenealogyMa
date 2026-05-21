package config

import (
	"context"
	"testing"
	"time"
)

// MockConfigRepository 模拟配置仓储
type MockConfigRepository struct {
	configs map[int64]*Config
	nextID  int64
}

func NewMockConfigRepository() *MockConfigRepository {
	return &MockConfigRepository{
		configs: make(map[int64]*Config),
		nextID:  1,
	}
}

func (m *MockConfigRepository) FindByID(ctx context.Context, id int64) (*Config, error) {
	if c, ok := m.configs[id]; ok {
		// Return a copy to prevent mutation issues
		copy := *c
		return &copy, nil
	}
	return nil, nil
}

func (m *MockConfigRepository) Create(ctx context.Context, config *Config) error {
	config.ID = m.nextID
	m.nextID++
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	// Store a copy to prevent mutation issues
	copy := *config
	m.configs[config.ID] = &copy
	return nil
}

func (m *MockConfigRepository) Update(ctx context.Context, config *Config) error {
	if _, ok := m.configs[config.ID]; !ok {
		return nil
	}
	config.UpdatedAt = time.Now()
	m.configs[config.ID] = config
	return nil
}

func (m *MockConfigRepository) Delete(ctx context.Context, id int64) error {
	delete(m.configs, id)
	return nil
}

func (m *MockConfigRepository) FindAll(ctx context.Context) ([]*Config, error) {
	result := make([]*Config, 0, len(m.configs))
	for _, c := range m.configs {
		result = append(result, c)
	}
	return result, nil
}

func (m *MockConfigRepository) FindByGroup(ctx context.Context, group string) ([]*Config, error) {
	result := make([]*Config, 0)
	for _, c := range m.configs {
		if c.Group == group {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *MockConfigRepository) FindByKey(ctx context.Context, group, key string) (*Config, error) {
	for _, c := range m.configs {
		if c.Group == group && c.Key == key {
			// Return a copy to prevent mutation issues
			copy := *c
			return &copy, nil
		}
	}
	return nil, nil
}

func (m *MockConfigRepository) Search(ctx context.Context, query *SearchQuery) ([]*Config, int64, error) {
	result := make([]*Config, 0)
	for _, c := range m.configs {
		if query.Group != "" && c.Group != query.Group {
			continue
		}
		if query.Key != "" && c.Key != query.Key {
			continue
		}
		result = append(result, c)
	}
	return result, int64(len(result)), nil
}

// ===== 单元测试 =====

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid string config",
			config: &Config{
				Group: "system",
				Key:   "site_name",
				Type:  ConfigTypeString,
			},
			wantErr: false,
		},
		{
			name: "valid number config",
			config: &Config{
				Group: "system",
				Key:   "max_upload_size",
				Type:  ConfigTypeNumber,
			},
			wantErr: false,
		},
		{
			name: "valid boolean config",
			config: &Config{
				Group: "system",
				Key:   "enable_debug",
				Type:  ConfigTypeBoolean,
			},
			wantErr: false,
		},
		{
			name: "valid json config",
			config: &Config{
				Group: "cemetery",
				Key:   "display_settings",
				Type:  ConfigTypeJSON,
			},
			wantErr: false,
		},
		{
			name: "missing group",
			config: &Config{
				Key:  "site_name",
				Type: ConfigTypeString,
			},
			wantErr: true,
		},
		{
			name: "missing key",
			config: &Config{
				Group: "system",
				Type:  ConfigTypeString,
			},
			wantErr: true,
		},
		{
			name: "missing type",
			config: &Config{
				Group: "system",
				Key:   "site_name",
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			config: &Config{
				Group: "system",
				Key:   "site_name",
				Type:  ConfigType("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_IsSystemConfig(t *testing.T) {
	systemConfig := &Config{IsSystem: true}
	userConfig := &Config{IsSystem: false}

	if !systemConfig.IsSystemConfig() {
		t.Error("systemConfig should return true for IsSystemConfig()")
	}
	if userConfig.IsSystemConfig() {
		t.Error("userConfig should return false for IsSystemConfig()")
	}
}

func TestConfig_CanDelete(t *testing.T) {
	systemConfig := &Config{IsSystem: true}
	userConfig := &Config{IsSystem: false}

	if systemConfig.CanDelete() {
		t.Error("systemConfig should not be deletable")
	}
	if !userConfig.CanDelete() {
		t.Error("userConfig should be deletable")
	}
}

func TestService_CreateConfig(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewService(repo)
	ctx := context.Background()

	t.Run("create valid config", func(t *testing.T) {
		config := &Config{
			Group: "system",
			Key:   "site_name",
			Type:  ConfigTypeString,
			Label: "网站名称",
		}
		err := svc.CreateConfig(ctx, config)
		if err != nil {
			t.Errorf("CreateConfig() error = %v", err)
		}
		if config.ID == 0 {
			t.Error("CreateConfig() should set ID")
		}
	})

	t.Run("create duplicate config", func(t *testing.T) {
		config := &Config{
			Group: "system",
			Key:   "site_name",
			Type:  ConfigTypeString,
		}
		err := svc.CreateConfig(ctx, config)
		if err == nil {
			t.Error("CreateConfig() should return error for duplicate key")
		}
	})

	t.Run("create config with invalid data", func(t *testing.T) {
		config := &Config{
			Group: "",
			Key:   "site_name",
			Type:  ConfigTypeString,
		}
		err := svc.CreateConfig(ctx, config)
		if err == nil {
			t.Error("CreateConfig() should return error for invalid data")
		}
	})
}

func TestService_UpdateConfig(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewService(repo)
	ctx := context.Background()

	// 创建初始配置
	config := &Config{
		Group: "system",
		Key:   "site_name",
		Type:  ConfigTypeString,
		Label: "网站名称",
	}
	_ = svc.CreateConfig(ctx, config)

	t.Run("update existing config", func(t *testing.T) {
		newLabel := "新网站名称"
		config.Label = newLabel
		err := svc.UpdateConfig(ctx, config)
		if err != nil {
			t.Errorf("UpdateConfig() error = %v", err)
		}
		updated, _ := svc.GetConfig(ctx, config.ID)
		if updated.Label != newLabel {
			t.Errorf("UpdateConfig() label = %v, want %v", updated.Label, newLabel)
		}
	})

	t.Run("update non-existent config", func(t *testing.T) {
		nonExistent := &Config{ID: 9999, Group: "system", Key: "test", Type: ConfigTypeString}
		err := svc.UpdateConfig(ctx, nonExistent)
		if err == nil {
			t.Error("UpdateConfig() should return error for non-existent config")
		}
	})

	t.Run("update system config group/key", func(t *testing.T) {
		systemConfig := &Config{
			Group:    "system",
			Key:      "system_mode",
			Type:     ConfigTypeBoolean,
			IsSystem: true,
		}
		_ = svc.CreateConfig(ctx, systemConfig)

		systemConfig.Group = "user"
		err := svc.UpdateConfig(ctx, systemConfig)
		if err == nil {
			t.Error("UpdateConfig() should prevent changing group of system config")
		}
	})
}

func TestService_DeleteConfig(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewService(repo)
	ctx := context.Background()

	t.Run("delete user config", func(t *testing.T) {
		config := &Config{
			Group: "custom",
			Key:   "custom_key",
			Type:  ConfigTypeString,
		}
		_ = svc.CreateConfig(ctx, config)

		err := svc.DeleteConfig(ctx, config.ID)
		if err != nil {
			t.Errorf("DeleteConfig() error = %v", err)
		}

		deleted, _ := svc.GetConfig(ctx, config.ID)
		if deleted != nil {
			t.Error("DeleteConfig() should remove config")
		}
	})

	t.Run("delete system config", func(t *testing.T) {
		systemConfig := &Config{
			Group:    "system",
			Key:      "system_key",
			Type:     ConfigTypeString,
			IsSystem: true,
		}
		_ = svc.CreateConfig(ctx, systemConfig)

		err := svc.DeleteConfig(ctx, systemConfig.ID)
		if err == nil {
			t.Error("DeleteConfig() should prevent deleting system config")
		}
	})

	t.Run("delete non-existent config", func(t *testing.T) {
		err := svc.DeleteConfig(ctx, 9999)
		if err == nil {
			t.Error("DeleteConfig() should return error for non-existent config")
		}
	})
}

func TestService_GetConfig(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewService(repo)
	ctx := context.Background()

	t.Run("get existing config", func(t *testing.T) {
		config := &Config{
			Group: "system",
			Key:   "site_name",
			Type:  ConfigTypeString,
		}
		_ = svc.CreateConfig(ctx, config)

		found, err := svc.GetConfig(ctx, config.ID)
		if err != nil {
			t.Errorf("GetConfig() error = %v", err)
		}
		if found == nil {
			t.Error("GetConfig() should return config")
		}
		if found.Key != config.Key {
			t.Errorf("GetConfig() key = %v, want %v", found.Key, config.Key)
		}
	})

	t.Run("get non-existent config", func(t *testing.T) {
		found, _ := svc.GetConfig(ctx, 9999)
		if found != nil {
			t.Error("GetConfig() should return nil for non-existent config")
		}
	})
}

func TestService_ListConfigs(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewService(repo)
	ctx := context.Background()

	t.Run("list all configs", func(t *testing.T) {
		_ = svc.CreateConfig(ctx, &Config{Group: "system", Key: "k1", Type: ConfigTypeString})
		_ = svc.CreateConfig(ctx, &Config{Group: "system", Key: "k2", Type: ConfigTypeString})
		_ = svc.CreateConfig(ctx, &Config{Group: "user", Key: "k3", Type: ConfigTypeString})

		all, err := svc.ListConfigs(ctx)
		if err != nil {
			t.Errorf("ListConfigs() error = %v", err)
		}
		if len(all) != 3 {
			t.Errorf("ListConfigs() count = %v, want 3", len(all))
		}
	})
}

func TestService_ListConfigsByGroup(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewService(repo)
	ctx := context.Background()

	_ = svc.CreateConfig(ctx, &Config{Group: "system", Key: "k1", Type: ConfigTypeString})
	_ = svc.CreateConfig(ctx, &Config{Group: "system", Key: "k2", Type: ConfigTypeString})
	_ = svc.CreateConfig(ctx, &Config{Group: "user", Key: "k3", Type: ConfigTypeString})

	t.Run("list configs by group", func(t *testing.T) {
		systemConfigs, err := svc.ListConfigsByGroup(ctx, "system")
		if err != nil {
			t.Errorf("ListConfigsByGroup() error = %v", err)
		}
		if len(systemConfigs) != 2 {
			t.Errorf("ListConfigsByGroup() count = %v, want 2", len(systemConfigs))
		}
	})
}

func TestService_GetConfigByKey(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewService(repo)
	ctx := context.Background()

	_ = svc.CreateConfig(ctx, &Config{Group: "system", Key: "site_name", Type: ConfigTypeString, Value: "My Site"})

	t.Run("get config by group and key", func(t *testing.T) {
		found, err := svc.GetConfigByKey(ctx, "system", "site_name")
		if err != nil {
			t.Errorf("GetConfigByKey() error = %v", err)
		}
		if found == nil {
			t.Error("GetConfigByKey() should return config")
		}
		if found.Value != "My Site" {
			t.Errorf("GetConfigByKey() value = %v, want %v", found.Value, "My Site")
		}
	})

	t.Run("get non-existent config by key", func(t *testing.T) {
		found, _ := svc.GetConfigByKey(ctx, "system", "non_existent")
		if found != nil {
			t.Error("GetConfigByKey() should return nil for non-existent key")
		}
	})
}

func TestService_SearchConfigs(t *testing.T) {
	repo := NewMockConfigRepository()
	svc := NewService(repo)
	ctx := context.Background()

	_ = svc.CreateConfig(ctx, &Config{Group: "system", Key: "site_name", Type: ConfigTypeString, Label: "网站名称"})
	_ = svc.CreateConfig(ctx, &Config{Group: "system", Key: "site_url", Type: ConfigTypeString, Label: "网站地址"})
	_ = svc.CreateConfig(ctx, &Config{Group: "cemetery", Key: "cemetery_name", Type: ConfigTypeString, Label: "墓园名称"})

	t.Run("search by group", func(t *testing.T) {
		query := &SearchQuery{Group: "system"}
		results, total, err := svc.SearchConfigs(ctx, query)
		if err != nil {
			t.Errorf("SearchConfigs() error = %v", err)
		}
		if total != 2 {
			t.Errorf("SearchConfigs() total = %v, want 2", total)
		}
		if len(results) != 2 {
			t.Errorf("SearchConfigs() count = %v, want 2", len(results))
		}
	})

	t.Run("search by key", func(t *testing.T) {
		query := &SearchQuery{Key: "site_name"}
		results, _, err := svc.SearchConfigs(ctx, query)
		if err != nil {
			t.Errorf("SearchConfigs() error = %v", err)
		}
		if len(results) != 1 {
			t.Errorf("SearchConfigs() count = %v, want 1", len(results))
		}
	})
}
