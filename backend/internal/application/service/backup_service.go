package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/backup"
)

// BackupService 备份应用服务
type BackupService struct {
	domainService *backup.Service
	backupDir     string
}

// NewBackupService 创建备份应用服务
func NewBackupService(repo backup.Repository, backupDir string) *BackupService {
	return &BackupService{
		domainService: backup.NewService(repo),
		backupDir:     backupDir,
	}
}

// ===== DTO 定义 =====

// BackupDTO 备份DTO
type BackupDTO struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	FilePath     string `json:"file_path"`
	FileSize     int64  `json:"file_size"`
	FileSizeStr  string `json:"file_size_str"`
	Database     string `json:"database"`
	ErrorMessage string `json:"error_message,omitempty"`
	StartedAt    string `json:"started_at,omitempty"`
	CompletedAt  string `json:"completed_at,omitempty"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    int64  `json:"created_by"`
}

// BackupListDTO 备份列表DTO
type BackupListDTO struct {
	Data       []*BackupDTO `json:"data"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}

// CreateBackupRequest 创建备份请求
type CreateBackupRequest struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Database string   `json:"database"`
	Tables   []string `json:"tables"`
}

// RestoreBackupRequest 恢复备份请求
type RestoreBackupRequest struct {
	BackupID int64  `json:"backup_id" binding:"required"`
	TargetDB string `json:"target_db"`
}

// BackupFilterRequest 备份筛选请求
type BackupFilterRequest struct {
	Type      string `form:"type"`
	Status    string `form:"status"`
	Database  string `form:"database"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

// CleanOldBackupsRequest 清理旧备份请求
type CleanOldBackupsRequest struct {
	Days int `json:"days" binding:"required,min=1"`
}

// ===== 服务方法 =====

// CreateBackup 创建新备份
func (s *BackupService) CreateBackup(ctx context.Context, req *CreateBackupRequest) (*BackupDTO, error) {
	backupType := backup.BackupTypeFull
	if req.Type == "partial" {
		backupType = backup.BackupTypePartial
	}

	cmd := &backup.CreateBackupCommand{
		Name:      req.Name,
		Type:      backupType,
		Database:  req.Database,
		Tables:    req.Tables,
		CreatedBy: 1, // TODO: 从上下文获取用户ID
	}

	b, err := s.domainService.CreateBackup(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// 异步执行备份
	go s.executeBackup(b.ID, req.Database, req.Tables)

	return backupToDTO(b), nil
}

// GetBackup 获取备份详情
func (s *BackupService) GetBackup(ctx context.Context, id int64) (*BackupDTO, error) {
	b, err := s.domainService.GetBackupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}
	return backupToDTO(b), nil
}

// ListBackups 获取备份列表
func (s *BackupService) ListBackups(ctx context.Context, req *BackupFilterRequest) (*BackupListDTO, error) {
	filter := backup.BackupFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	if req.Type != "" {
		t := backup.BackupType(req.Type)
		filter.Type = &t
	}
	if req.Status != "" {
		st := backup.BackupStatus(req.Status)
		filter.Status = &st
	}
	filter.Database = req.Database

	result, err := s.domainService.ListBackups(ctx, filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]*BackupDTO, len(result.Data))
	for i, b := range result.Data {
		dtos[i] = backupToDTO(&b)
	}

	return &BackupListDTO{
		Data:       dtos,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// RestoreBackup 恢复备份
func (s *BackupService) RestoreBackup(ctx context.Context, req *RestoreBackupRequest) error {
	b, err := s.domainService.GetBackupByID(ctx, req.BackupID)
	if err != nil {
		return err
	}
	if b == nil {
		return fmt.Errorf("backup not found")
	}
	if b.Status != backup.BackupStatusCompleted {
		return fmt.Errorf("can only restore completed backup")
	}
	if b.FilePath == "" {
		return fmt.Errorf("backup file not found")
	}

	targetDB := req.TargetDB
	if targetDB == "" {
		targetDB = b.Database
	}

	return s.executeRestore(b.FilePath, targetDB)
}

// DeleteBackup 删除备份
func (s *BackupService) DeleteBackup(ctx context.Context, id int64) error {
	b, err := s.domainService.GetBackupByID(ctx, id)
	if err != nil {
		return err
	}
	if b == nil {
		return fmt.Errorf("backup not found")
	}

	// 删除物理文件
	if b.FilePath != "" {
		os.Remove(b.FilePath)
	}

	return s.domainService.DeleteBackup(ctx, id)
}

// CleanOldBackups 清理旧备份
func (s *BackupService) CleanOldBackups(ctx context.Context, days int) (int64, error) {
	return s.domainService.CleanOldBackups(ctx, days)
}

// ===== 内部方法 =====

// executeBackup 执行备份操作
func (s *BackupService) executeBackup(backupID int64, database string, tables []string) {
	ctx := context.Background()

	// 更新状态为运行中
	s.domainService.UpdateBackupStatus(ctx, backupID, backup.BackupStatusRunning, "", 0, "")

	// 生成备份文件名
	filename := fmt.Sprintf("backup_%d_%s.sql", backupID, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(s.backupDir, filename)

	// 确保备份目录存在
	os.MkdirAll(s.backupDir, 0755)

	var err error
	var fileSize int64

	if len(tables) > 0 {
		// 部分表备份
		err = s.backupTables(database, tables, filePath)
	} else {
		// 全量备份
		err = s.backupDatabase(database, filePath)
	}

	if err != nil {
		s.domainService.UpdateBackupStatus(ctx, backupID, backup.BackupStatusFailed, filePath, 0, err.Error())
		return
	}

	// 获取文件大小
	if info, statErr := os.Stat(filePath); statErr == nil {
		fileSize = info.Size()
	}

	s.domainService.UpdateBackupStatus(ctx, backupID, backup.BackupStatusCompleted, filePath, fileSize, "")
}

// backupDatabase 执行全量数据库备份
func (s *BackupService) backupDatabase(database, outputPath string) error {
	cmd := exec.Command("pg_dump", "-Fc", "-f", outputPath, database)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %v, output: %s", err, string(output))
	}
	return nil
}

// backupTables 执行部分表备份
func (s *BackupService) backupTables(database string, tables []string, outputPath string) error {
	args := append([]string{"-Fc", "-f", outputPath, "-t"}, tables...)
	args = append(args, database)
	cmd := exec.Command("pg_dump", args...)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %v, output: %s", err, string(output))
	}
	return nil
}

// executeRestore 执行数据库恢复
func (s *BackupService) executeRestore(backupFile, targetDB string) error {
	// 使用 pg_restore 恢复
	cmd := exec.Command("pg_restore", "-d", targetDB, backupFile)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_restore failed: %v, output: %s", err, string(output))
	}
	return nil
}

// ===== 辅助函数 =====

func backupToDTO(b *backup.Backup) *BackupDTO {
	dto := &BackupDTO{
		ID:           b.ID,
		Name:         b.Name,
		Type:         string(b.Type),
		Status:       string(b.Status),
		FilePath:     b.FilePath,
		FileSize:     b.FileSize,
		FileSizeStr:  formatFileSize(b.FileSize),
		Database:     b.Database,
		ErrorMessage: b.ErrorMessage,
		CreatedAt:    b.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:    b.CreatedBy,
	}

	if b.StartedAt != nil {
		dto.StartedAt = b.StartedAt.Format("2006-01-02 15:04:05")
	}
	if b.CompletedAt != nil {
		dto.CompletedAt = b.CompletedAt.Format("2006-01-02 15:04:05")
	}

	return dto
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

