package config

import (
	"testing"
	"time"
)

func TestConfigEntity_Validate(t *testing.T) {
	tests := []struct {
		name    string
		c       *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			c: &Config{
				Group: "system",
				Key:   "site_name",
				Type:  ConfigTypeString,
			},
			wantErr: false,
		},
		{
			name: "empty group",
			c: &Config{
				Group: "",
				Key:   "site_name",
				Type:  ConfigTypeString,
			},
			wantErr: true,
			errMsg:  "group is required",
		},
		{
			name: "empty key",
			c: &Config{
				Group: "system",
				Key:   "",
				Type:  ConfigTypeString,
			},
			wantErr: true,
			errMsg:  "key is required",
		},
		{
			name: "empty type",
			c: &Config{
				Group: "system",
				Key:   "site_name",
				Type:  "",
			},
			wantErr: true,
			errMsg:  "type is required",
		},
		{
			name: "invalid type",
			c: &Config{
				Group: "system",
				Key:   "site_name",
				Type:  ConfigType("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid type: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestConfigEntity_IsSystemConfig(t *testing.T) {
	tests := []struct {
		name     string
		c        *Config
		isSystem bool
	}{
		{
			name: "system config",
			c: &Config{
				IsSystem: true,
			},
			isSystem: true,
		},
		{
			name: "user config",
			c: &Config{
				IsSystem: false,
			},
			isSystem: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.c.IsSystemConfig() != tt.isSystem {
				t.Errorf("IsSystemConfig() = %v, want %v", tt.c.IsSystemConfig(), tt.isSystem)
			}
		})
	}
}

func TestConfigEntity_CanDelete(t *testing.T) {
	tests := []struct {
		name      string
		c         *Config
		canDelete bool
	}{
		{
			name: "system config cannot delete",
			c: &Config{
				IsSystem: true,
			},
			canDelete: false,
		},
		{
			name: "user config can delete",
			c: &Config{
				IsSystem: false,
			},
			canDelete: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.c.CanDelete() != tt.canDelete {
				t.Errorf("CanDelete() = %v, want %v", tt.c.CanDelete(), tt.canDelete)
			}
		})
	}
}

func TestConfigType_Constants(t *testing.T) {
	if ConfigTypeString != "string" {
		t.Errorf("ConfigTypeString = %v, want string", ConfigTypeString)
	}
	if ConfigTypeNumber != "number" {
		t.Errorf("ConfigTypeNumber = %v, want number", ConfigTypeNumber)
	}
	if ConfigTypeBoolean != "boolean" {
		t.Errorf("ConfigTypeBoolean = %v, want boolean", ConfigTypeBoolean)
	}
	if ConfigTypeJSON != "json" {
		t.Errorf("ConfigTypeJSON = %v, want json", ConfigTypeJSON)
	}
}

func TestIsValidType(t *testing.T) {
	tests := []struct {
		name  string
		t     ConfigType
		valid bool
	}{
		{"valid string", ConfigTypeString, true},
		{"valid number", ConfigTypeNumber, true},
		{"valid boolean", ConfigTypeBoolean, true},
		{"valid json", ConfigTypeJSON, true},
		{"invalid type", ConfigType("invalid"), false},
		{"empty type", ConfigType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if isValidType(tt.t) != tt.valid {
				t.Errorf("isValidType(%v) = %v, want %v", tt.t, isValidType(tt.t), tt.valid)
			}
		})
	}
}

func TestSearchQuery_Fields(t *testing.T) {
	q := &SearchQuery{
		Group:    "system",
		Key:      "site",
		Label:    "站点",
		Page:     1,
		PageSize: 20,
	}

	if q.Group != "system" {
		t.Errorf("SearchQuery.Group = %v, want system", q.Group)
	}
	if q.Key != "site" {
		t.Errorf("SearchQuery.Key = %v, want site", q.Key)
	}
	if q.Page != 1 {
		t.Errorf("SearchQuery.Page = %v, want 1", q.Page)
	}
	if q.PageSize != 20 {
		t.Errorf("SearchQuery.PageSize = %v, want 20", q.PageSize)
	}
}

func TestConfigEntity_Fields(t *testing.T) {
	now := time.Now()
	c := &Config{
		ID:        1,
		Group:     "system",
		Key:       "site_name",
		Value:     "族谱平台",
		Type:      ConfigTypeString,
		Label:     "站点名称",
		HelpText:  "请输入站点名称",
		SortOrder: 1,
		IsSystem:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if c.ID != 1 {
		t.Errorf("Config.ID = %v, want 1", c.ID)
	}
	if c.Value != "族谱平台" {
		t.Errorf("Config.Value = %v, want 族谱平台", c.Value)
	}
	if c.Type != ConfigTypeString {
		t.Errorf("Config.Type = %v, want string", c.Type)
	}
	if c.SortOrder != 1 {
		t.Errorf("Config.SortOrder = %v, want 1", c.SortOrder)
	}
	if !c.IsSystem {
		t.Errorf("Config.IsSystem = %v, want true", c.IsSystem)
	}
}

func TestConfigEntity_VariousTypes(t *testing.T) {
	configs := []*Config{
		{Group: "system", Key: "num", Type: ConfigTypeNumber},
		{Group: "system", Key: "bool", Type: ConfigTypeBoolean},
		{Group: "system", Key: "json", Type: ConfigTypeJSON},
	}

	for _, cfg := range configs {
		if err := cfg.Validate(); err != nil {
			t.Errorf("Validate() error = %v for type %v", err, cfg.Type)
		}
	}
}