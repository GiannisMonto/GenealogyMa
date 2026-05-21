package backup

import (
	"context"
	"time"
)

// Repository 备份仓储接口
type Repository interface {
	// Create 创建备份记录
	Create(ctx context.Context, backup *Backup) error

	// GetByID 根据ID获取备份记录
	GetByID(ctx context.Context, id int64) (*Backup, error)

	// List 获取备份列表
	List(ctx context.Context, filter BackupFilter) (*BackupList, error)

	// Update 更新备份记录
	Update(ctx context.Context, backup *Backup) error

	// Delete 删除备份记录
	Delete(ctx context.Context, id int64) error

	// DeleteOldBackups 删除旧备份
	DeleteOldBackups(ctx context.Context, before time.Time) (int64, error)
}