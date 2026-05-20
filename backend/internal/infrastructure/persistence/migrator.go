package persistence

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// MigrationRecord 迁移记录表
type MigrationRecord struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string    `gorm:"unique;not null"`
	AppliedAt time.Time `gorm:"not null"`
}

// TableName 设置表名
func (MigrationRecord) TableName() string {
	return "schema_migrations"
}

// Migrator 数据库迁移执行器
type Migrator struct {
	db *gorm.DB
}

// NewMigrator 创建迁移执行器
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// Run 执行所有未执行的迁移
func (m *Migrator) Run() error {
	// 确保迁移记录表存在
	if err := m.ensureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	// 获取所有迁移文件
	files, err := fs.Glob(migrationFS, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to list migration files: %w", err)
	}

	// 按文件名排序（确保按顺序执行）
	sort.Strings(files)

	// 获取已执行的迁移
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	// 执行未执行的迁移
	for _, file := range files {
		name := m.extractMigrationName(file)
		if !applied[name] {
			if err := m.applyMigration(file, name); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", name, err)
			}
		}
	}

	return nil
}

// ensureMigrationTable 确保迁移记录表存在
func (m *Migrator) ensureMigrationTable() error {
	return m.db.AutoMigrate(&MigrationRecord{})
}

// getAppliedMigrations 获取已执行的迁移
func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	var records []MigrationRecord
	if err := m.db.Find(&records).Error; err != nil {
		return nil, err
	}

	applied := make(map[string]bool)
	for _, r := range records {
		applied[r.Name] = true
	}
	return applied, nil
}

// extractMigrationName 从文件路径提取迁移名称
func (m *Migrator) extractMigrationName(file string) string {
	parts := strings.Split(file, "/")
	return parts[len(parts)-1]
}

// applyMigration 执行单个迁移
func (m *Migrator) applyMigration(file, name string) error {
	content, err := migrationFS.ReadFile(file)
	if err != nil {
		return err
	}

	// 在事务中执行迁移
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 执行 SQL
		if err := tx.Exec(string(content)).Error; err != nil {
			return err
		}

		// 记录迁移
		return tx.Create(&MigrationRecord{
			Name:      name,
			AppliedAt: time.Now(),
		}).Error
	})
}

// Reset 重置数据库（仅用于测试）
func (m *Migrator) Reset() error {
	tables := []string{
		"role_permissions",
		"user_roles",
		"permissions",
		"roles",
		"users",
		"schema_migrations",
	}

	for _, table := range tables {
		if err := m.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			return err
		}
	}
	return nil
}
