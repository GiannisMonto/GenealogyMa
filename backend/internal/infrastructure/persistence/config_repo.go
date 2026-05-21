package persistence

import (
	"context"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/config"
	"gorm.io/gorm"
)

// ConfigModel GORM模型对应configs表
type ConfigModel struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Group     string `gorm:"column:group;type:varchar(100);not null;index:idx_group_key,unique"`
	Key       string `gorm:"column:key;type:varchar(100);not null;index:idx_group_key,unique"`
	Value     string `gorm:"column:value;type:text"`
	Type      string `gorm:"column:type;type:varchar(20);not null;default:'string'"`
	Label     string `gorm:"column:label;type:varchar(200)"`
	HelpText  string `gorm:"column:help_text;type:text"`
	SortOrder int    `gorm:"column:sort_order;default:0"`
	IsSystem  bool   `gorm:"column:is_system;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ConfigModel) TableName() string {
	return "configs"
}

// ConfigRepositoryImpl 配置仓储实现
type ConfigRepositoryImpl struct {
	db *gorm.DB
}

// NewConfigRepository 创建配置仓储
func NewConfigRepository(db *gorm.DB) config.Repository {
	return &ConfigRepositoryImpl{db: db}
}

// modelToConfig converts a database model to domain entity
func modelToConfig(m *ConfigModel) *config.Config {
	return &config.Config{
		ID:        m.ID,
		Group:     m.Group,
		Key:       m.Key,
		Value:     m.Value,
		Type:      config.ConfigType(m.Type),
		Label:     m.Label,
		HelpText:  m.HelpText,
		SortOrder: m.SortOrder,
		IsSystem:  m.IsSystem,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// configToModel converts a domain entity to database model
func configToModel(c *config.Config) *ConfigModel {
	return &ConfigModel{
		ID:        c.ID,
		Group:     c.Group,
		Key:       c.Key,
		Value:     c.Value,
		Type:      string(c.Type),
		Label:     c.Label,
		HelpText:  c.HelpText,
		SortOrder: c.SortOrder,
		IsSystem:  c.IsSystem,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// FindByID implements config.Repository
func (r *ConfigRepositoryImpl) FindByID(ctx context.Context, id int64) (*config.Config, error) {
	var model ConfigModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToConfig(&model), nil
}

// Create implements config.Repository
func (r *ConfigRepositoryImpl) Create(ctx context.Context, c *config.Config) error {
	model := configToModel(c)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update implements config.Repository
func (r *ConfigRepositoryImpl) Update(ctx context.Context, c *config.Config) error {
	model := configToModel(c)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete implements config.Repository
func (r *ConfigRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&ConfigModel{}, id).Error
}

// FindAll implements config.Repository
func (r *ConfigRepositoryImpl) FindAll(ctx context.Context) ([]*config.Config, error) {
	var models []ConfigModel
	if err := r.db.WithContext(ctx).Order("`group`, sort_order").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*config.Config, len(models))
	for i := range models {
		result[i] = modelToConfig(&models[i])
	}
	return result, nil
}

// FindByGroup implements config.Repository
func (r *ConfigRepositoryImpl) FindByGroup(ctx context.Context, group string) ([]*config.Config, error) {
	var models []ConfigModel
	if err := r.db.WithContext(ctx).Where("`group` = ?", group).Order("sort_order").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*config.Config, len(models))
	for i := range models {
		result[i] = modelToConfig(&models[i])
	}
	return result, nil
}

// FindByKey implements config.Repository
func (r *ConfigRepositoryImpl) FindByKey(ctx context.Context, group, key string) (*config.Config, error) {
	var model ConfigModel
	if err := r.db.WithContext(ctx).Where("`group` = ? AND `key` = ?", group, key).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToConfig(&model), nil
}

// Search implements config.Repository
func (r *ConfigRepositoryImpl) Search(ctx context.Context, query *config.SearchQuery) ([]*config.Config, int64, error) {
	var models []ConfigModel
	var total int64

	db := r.db.WithContext(ctx).Model(&ConfigModel{})

	// Apply filters
	if query.Group != "" {
		db = db.Where("`group` = ?", query.Group)
	}
	if query.Key != "" {
		db = db.Where("`key` = ?", query.Key)
	}
	if query.Label != "" {
		db = db.Where("label LIKE ?", "%"+query.Label+"%")
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	db = db.Order("`group`, sort_order")

	// Apply pagination
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*config.Config, len(models))
	for i := range models {
		result[i] = modelToConfig(&models[i])
	}
	return result, total, nil
}
