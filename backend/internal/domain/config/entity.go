package config

import (
	"fmt"
	"time"
)

// Config 配置项聚合根
type Config struct {
	ID        int64     `json:"id"`
	Group     string    `json:"group"`      // 配置分组：system, cemetery, culture, etc.
	Key       string    `json:"key"`        // 配置键
	Value     string    `json:"value"`      // 配置值
	Type      ConfigType `json:"type"`      // 配置类型：string, number, boolean, json
	Label     string    `json:"label"`      // 配置标签/名称
	HelpText  string    `json:"help_text"`  // 帮助文本
	SortOrder int       `json:"sort_order"` // 排序
	IsSystem  bool      `json:"is_system"` // 是否系统配置（不可删除）
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ConfigType 配置值类型
type ConfigType string

const (
	ConfigTypeString  ConfigType = "string"  // 字符串
	ConfigTypeNumber  ConfigType = "number"  // 数字
	ConfigTypeBoolean ConfigType = "boolean" // 布尔
	ConfigTypeJSON    ConfigType = "json"    // JSON
)

// Validate 验证配置信息
func (c *Config) Validate() error {
	if c.Group == "" {
		return fmt.Errorf("group is required")
	}
	if c.Key == "" {
		return fmt.Errorf("key is required")
	}
	if c.Type == "" {
		return fmt.Errorf("type is required")
	}
	if !isValidType(c.Type) {
		return fmt.Errorf("invalid type: %s", c.Type)
	}
	return nil
}

// isValidType 验证配置类型
func isValidType(t ConfigType) bool {
	switch t {
	case ConfigTypeString, ConfigTypeNumber, ConfigTypeBoolean, ConfigTypeJSON:
		return true
	default:
		return false
	}
}

// IsSystemConfig 判断是否为系统配置
func (c *Config) IsSystemConfig() bool {
	return c.IsSystem
}

// CanDelete 判断配置是否可以删除
func (c *Config) CanDelete() bool {
	return !c.IsSystem
}

// SearchQuery 搜索查询参数
type SearchQuery struct {
	Group   string `json:"group"`
	Key     string `json:"key"`
	Label   string `json:"label"`
	Page    int    `json:"page"`
	PageSize int   `json:"page_size"`
}