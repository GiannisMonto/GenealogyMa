package backup

import (
	"time"
)

// BackupType 备份类型
type BackupType string

const (
	BackupTypeFull    BackupType = "full"    // 全量备份
	BackupTypePartial BackupType = "partial"  // 部分备份
)

// BackupStatus 备份状态
type BackupStatus string

const (
	BackupStatusPending   BackupStatus = "pending"   // 待处理
	BackupStatusRunning   BackupStatus = "running"   // 执行中
	BackupStatusCompleted BackupStatus = "completed" // 已完成
	BackupStatusFailed    BackupStatus = "failed"    // 失败
)

// Backup 备份记录实体
type Backup struct {
	ID           int64        `json:"id"`
	Name         string       `json:"name"`
	Type         BackupType   `json:"type"`
	Status       BackupStatus `json:"status"`
	FilePath     string       `json:"file_path"`
	FileSize     int64        `json:"file_size"`
	Database     string       `json:"database"`
	Tables       []string     `json:"tables,omitempty"`
	ErrorMessage string       `json:"error_message,omitempty"`
	StartedAt    *time.Time   `json:"started_at,omitempty"`
	CompletedAt  *time.Time   `json:"completed_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	CreatedBy    int64        `json:"created_by"`
}

// BackupFilter 备份查询过滤器
type BackupFilter struct {
	Type     *BackupType
	Status   *BackupStatus
	Database string
	StartDate *time.Time
	EndDate   *time.Time
	Page     int
	PageSize int
}

// BackupList 备份列表
type BackupList struct {
	Data       []Backup `json:"data"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}

// CreateBackupCommand 创建备份命令
type CreateBackupCommand struct {
	Name      string     `json:"name"`
	Type      BackupType `json:"type"`
	Database  string     `json:"database"`
	Tables    []string   `json:"tables,omitempty"`
	CreatedBy int64      `json:"created_by"`
}

// RestoreBackupCommand 恢复备份命令
type RestoreBackupCommand struct {
	BackupID  int64  `json:"backup_id"`
	TargetDB  string `json:"target_db,omitempty"`
	CreatedBy int64  `json:"created_by"`
}