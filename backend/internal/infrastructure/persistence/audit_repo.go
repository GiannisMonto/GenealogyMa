package persistence

import (
	"context"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/audit"
	"gorm.io/gorm"
)

// AuditLogModel GORM模型对应audit_logs表
type AuditLogModel struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64  `gorm:"column:user_id;not null;index"`
	Username     string `gorm:"column:username;type:varchar(100);not null"`
	Module       string `gorm:"column:module;type:varchar(50);not null;index"`
	Action       string `gorm:"column:action;type:varchar(50);not null;index"`
	ResourceType string `gorm:"column:resource_type;type:varchar(50);index"`
	ResourceID   string `gorm:"column:resource_id;type:varchar(100);index"`
	OldValue     string `gorm:"column:old_value;type:text"`
	NewValue     string `gorm:"column:new_value;type:text"`
	IPAddress    string `gorm:"column:ip_address;type:varchar(50)"`
	UserAgent    string `gorm:"column:user_agent;type:varchar(500)"`
	Description  string `gorm:"column:description;type:text"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;index"`
}

func (AuditLogModel) TableName() string {
	return "audit_logs"
}

// ===== 审计日志仓储实现 =====

type AuditRepositoryImpl struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) audit.Repository {
	return &AuditRepositoryImpl{db: db}
}

func auditModelToEntity(m *AuditLogModel) *audit.AuditLog {
	return &audit.AuditLog{
		ID:           m.ID,
		UserID:       m.UserID,
		Username:     m.Username,
		Module:       m.Module,
		Action:       m.Action,
		ResourceType: m.ResourceType,
		ResourceID:   m.ResourceID,
		OldValue:     m.OldValue,
		NewValue:     m.NewValue,
		IPAddress:    m.IPAddress,
		UserAgent:    m.UserAgent,
		Description:  m.Description,
		CreatedAt:    m.CreatedAt,
	}
}

func auditEntityToModel(log *audit.AuditLog) *AuditLogModel {
	return &AuditLogModel{
		ID:           log.ID,
		UserID:       log.UserID,
		Username:     log.Username,
		Module:       log.Module,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		OldValue:     log.OldValue,
		NewValue:     log.NewValue,
		IPAddress:    log.IPAddress,
		UserAgent:    log.UserAgent,
		Description:  log.Description,
	}
}

func (r *AuditRepositoryImpl) Create(ctx context.Context, log *audit.AuditLog) error {
	model := auditEntityToModel(log)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *AuditRepositoryImpl) GetByID(ctx context.Context, id int64) (*audit.AuditLog, error) {
	var model AuditLogModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return auditModelToEntity(&model), nil
}

func (r *AuditRepositoryImpl) List(ctx context.Context, filter audit.AuditLogFilter) (*audit.AuditLogList, error) {
	var models []AuditLogModel
	var total int64

	db := r.db.WithContext(ctx).Model(&AuditLogModel{})

	if filter.UserID != nil {
		db = db.Where("user_id = ?", *filter.UserID)
	}
	if filter.Module != "" {
		db = db.Where("module = ?", filter.Module)
	}
	if filter.Action != "" {
		db = db.Where("action = ?", filter.Action)
	}
	if filter.ResourceType != "" {
		db = db.Where("resource_type = ?", filter.ResourceType)
	}
	if filter.StartDate != nil {
		db = db.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		db = db.Where("created_at <= ?", *filter.EndDate)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&models).Error; err != nil {
		return nil, err
	}

	logs := make([]audit.AuditLog, len(models))
	for i := range models {
		logs[i] = *auditModelToEntity(&models[i])
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &audit.AuditLogList{
		Data:       logs,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *AuditRepositoryImpl) DeleteOldLogs(ctx context.Context, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("created_at < ?", before).Delete(&AuditLogModel{})
	return result.RowsAffected, result.Error
}
