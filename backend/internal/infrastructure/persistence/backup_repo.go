package persistence

import (
	"context"
	"strings"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/backup"
	"gorm.io/gorm"
)

// BackupModel GORM模型对应backups表
type BackupModel struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string    `gorm:"column:name;type:varchar(255);not null"`
	Type         string    `gorm:"column:type;type:varchar(20);not null"`
	Status       string    `gorm:"column:status;type:varchar(20);not null;index"`
	FilePath     string    `gorm:"column:file_path;type:varchar(500)"`
	FileSize     int64     `gorm:"column:file_size"`
	Database     string    `gorm:"column:database;type:varchar(100);not null"`
	Tables       string    `gorm:"column:tables;type:text"`
	ErrorMessage string    `gorm:"column:error_message;type:text"`
	StartedAt    *time.Time `gorm:"column:started_at"`
	CompletedAt  *time.Time `gorm:"column:completed_at"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	CreatedBy    int64     `gorm:"column:created_by"`
}

func (BackupModel) TableName() string {
	return "backups"
}

// BackupRepositoryImpl 备份仓储实现
type BackupRepositoryImpl struct {
	db *gorm.DB
}

// NewBackupRepository 创建备份仓储
func NewBackupRepository(db *gorm.DB) backup.Repository {
	return &BackupRepositoryImpl{db: db}
}

func backupModelToEntity(m *BackupModel) *backup.Backup {
	b := &backup.Backup{
		ID:           m.ID,
		Name:         m.Name,
		Type:         backup.BackupType(m.Type),
		Status:       backup.BackupStatus(m.Status),
		FilePath:     m.FilePath,
		FileSize:     m.FileSize,
		Database:     m.Database,
		ErrorMessage: m.ErrorMessage,
		StartedAt:    m.StartedAt,
		CompletedAt:  m.CompletedAt,
		CreatedAt:    m.CreatedAt,
		CreatedBy:    m.CreatedBy,
	}
	if m.Tables != "" {
		b.Tables = strings.Split(m.Tables, ",")
	}
	return b
}

func backupEntityToModel(b *backup.Backup) *BackupModel {
	model := &BackupModel{
		ID:           b.ID,
		Name:         b.Name,
		Type:         string(b.Type),
		Status:       string(b.Status),
		FilePath:     b.FilePath,
		FileSize:     b.FileSize,
		Database:     b.Database,
		Tables:       strings.Join(b.Tables, ","),
		ErrorMessage: b.ErrorMessage,
		StartedAt:    b.StartedAt,
		CompletedAt:  b.CompletedAt,
		CreatedAt:    b.CreatedAt,
		CreatedBy:    b.CreatedBy,
	}
	return model
}

func (r *BackupRepositoryImpl) Create(ctx context.Context, b *backup.Backup) error {
	model := backupEntityToModel(b)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *BackupRepositoryImpl) GetByID(ctx context.Context, id int64) (*backup.Backup, error) {
	var model BackupModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return backupModelToEntity(&model), nil
}

func (r *BackupRepositoryImpl) List(ctx context.Context, filter backup.BackupFilter) (*backup.BackupList, error) {
	var models []BackupModel
	var total int64

	db := r.db.WithContext(ctx).Model(&BackupModel{})

	if filter.Type != nil {
		db = db.Where("type = ?", string(*filter.Type))
	}
	if filter.Status != nil {
		db = db.Where("status = ?", string(*filter.Status))
	}
	if filter.Database != "" {
		db = db.Where("database = ?", filter.Database)
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

	backups := make([]backup.Backup, len(models))
	for i := range models {
		backups[i] = *backupModelToEntity(&models[i])
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &backup.BackupList{
		Data:       backups,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *BackupRepositoryImpl) Update(ctx context.Context, b *backup.Backup) error {
	model := backupEntityToModel(b)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *BackupRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&BackupModel{}, id).Error
}

func (r *BackupRepositoryImpl) DeleteOldBackups(ctx context.Context, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("created_at < ? AND status != ?", before, backup.BackupStatusRunning).Delete(&BackupModel{})
	return result.RowsAffected, result.Error
}