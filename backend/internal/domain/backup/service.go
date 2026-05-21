package backup

import (
	"context"
	"fmt"
	"time"
)

// Service 备份领域服务
type Service struct {
	repo Repository
}

// NewService 创建备份领域服务
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateBackup 创建新备份
func (s *Service) CreateBackup(ctx context.Context, cmd *CreateBackupCommand) (*Backup, error) {
	if cmd.Name == "" {
		return nil, fmt.Errorf("backup name is required")
	}
	if cmd.Type == "" {
		cmd.Type = BackupTypeFull
	}
	if cmd.Database == "" {
		return nil, fmt.Errorf("database name is required")
	}

	backup := &Backup{
		Name:      cmd.Name,
		Type:      cmd.Type,
		Status:    BackupStatusPending,
		Database:  cmd.Database,
		Tables:    cmd.Tables,
		CreatedAt: time.Now(),
		CreatedBy: cmd.CreatedBy,
	}

	if err := s.repo.Create(ctx, backup); err != nil {
		return nil, err
	}

	return backup, nil
}

// GetBackupByID 获取备份详情
func (s *Service) GetBackupByID(ctx context.Context, id int64) (*Backup, error) {
	return s.repo.GetByID(ctx, id)
}

// ListBackups 获取备份列表
func (s *Service) ListBackups(ctx context.Context, filter BackupFilter) (*BackupList, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	return s.repo.List(ctx, filter)
}

// UpdateBackupStatus 更新备份状态
func (s *Service) UpdateBackupStatus(ctx context.Context, id int64, status BackupStatus, filePath string, fileSize int64, errorMsg string) error {
	backup, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if backup == nil {
		return fmt.Errorf("backup not found")
	}

	backup.Status = status
	if filePath != "" {
		backup.FilePath = filePath
	}
	if fileSize > 0 {
		backup.FileSize = fileSize
	}
	if errorMsg != "" {
		backup.ErrorMessage = errorMsg
	}

	now := time.Now()
	if status == BackupStatusRunning {
		backup.StartedAt = &now
	}
	if status == BackupStatusCompleted || status == BackupStatusFailed {
		backup.CompletedAt = &now
	}

	return s.repo.Update(ctx, backup)
}

// DeleteBackup 删除备份记录
func (s *Service) DeleteBackup(ctx context.Context, id int64) error {
	backup, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if backup == nil {
		return fmt.Errorf("backup not found")
	}
	if backup.Status == BackupStatusRunning {
		return fmt.Errorf("cannot delete running backup")
	}
	return s.repo.Delete(ctx, id)
}

// CleanOldBackups 清理旧备份
func (s *Service) CleanOldBackups(ctx context.Context, days int) (int64, error) {
	if days < 1 {
		return 0, fmt.Errorf("days must be positive")
	}
	before := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	return s.repo.DeleteOldBackups(ctx, before)
}