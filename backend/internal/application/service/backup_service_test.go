package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/backup"
)

// mockBackupRepository 模拟备份仓储
type mockBackupRepository struct {
	backups   map[int64]*backup.Backup
	nextID    int64
}

func newMockBackupRepository() *mockBackupRepository {
	return &mockBackupRepository{
		backups: make(map[int64]*backup.Backup),
		nextID:  1,
	}
}

func (m *mockBackupRepository) Create(ctx context.Context, b *backup.Backup) error {
	b.ID = m.nextID
	m.nextID++
	b.CreatedAt = time.Now()
	m.backups[b.ID] = b
	return nil
}

func (m *mockBackupRepository) GetByID(ctx context.Context, id int64) (*backup.Backup, error) {
	b, ok := m.backups[id]
	if !ok {
		return nil, nil
	}
	return b, nil
}

func (m *mockBackupRepository) List(ctx context.Context, filter backup.BackupFilter) (*backup.BackupList, error) {
	var result []backup.Backup
	for _, b := range m.backups {
		result = append(result, *b)
	}

	total := int64(len(result))
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &backup.BackupList{
		Data:       result,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (m *mockBackupRepository) Update(ctx context.Context, b *backup.Backup) error {
	m.backups[b.ID] = b
	return nil
}

func (m *mockBackupRepository) Delete(ctx context.Context, id int64) error {
	delete(m.backups, id)
	return nil
}

func (m *mockBackupRepository) DeleteOldBackups(ctx context.Context, before time.Time) (int64, error) {
	var count int64
	for id, b := range m.backups {
		if b.CreatedAt.Before(before) && b.Status != backup.BackupStatusRunning {
			delete(m.backups, id)
			count++
		}
	}
	return count, nil
}

func newTestBackupService() (*BackupService, *mockBackupRepository) {
	repo := newMockBackupRepository()
	svc := &BackupService{
		domainService: backup.NewService(repo),
		backupDir:    "/tmp/backups",
	}
	return svc, repo
}

// ===== 测试用例 =====

func TestCreateBackup(t *testing.T) {
	svc, repo := newTestBackupService()
	ctx := context.Background()

	t.Run("create full backup successfully", func(t *testing.T) {
		req := &CreateBackupRequest{
			Name:     "test_backup",
			Type:     "full",
			Database: "genealogy",
		}

		result, err := svc.CreateBackup(ctx, req)
		if err != nil {
			t.Fatalf("CreateBackup() error = %v", err)
		}
		if result.Name != "test_backup" {
			t.Errorf("expected name 'test_backup', got '%s'", result.Name)
		}
		if result.Type != "full" {
			t.Errorf("expected type 'full', got '%s'", result.Type)
		}
		// Status is "running" because CreateBackup starts a goroutine that updates status
		if result.Status != "running" {
			t.Errorf("expected status 'running', got '%s'", result.Status)
		}
		if result.Database != "genealogy" {
			t.Errorf("expected database 'genealogy', got '%s'", result.Database)
		}
		// 备份记录应该已保存
		if len(repo.backups) != 1 {
			t.Errorf("expected 1 backup in repo, got %d", len(repo.backups))
		}
	})

	t.Run("create partial backup", func(t *testing.T) {
		req := &CreateBackupRequest{
			Name:     "partial_backup",
			Type:     "partial",
			Database: "genealogy",
			Tables:   []string{"members", "parent_child_relations"},
		}

		result, err := svc.CreateBackup(ctx, req)
		if err != nil {
			t.Fatalf("CreateBackup() error = %v", err)
		}
		if result.Type != "partial" {
			t.Errorf("expected type 'partial', got '%s'", result.Type)
		}
	})

	t.Run("create backup with empty name", func(t *testing.T) {
		req := &CreateBackupRequest{
			Name:     "",
			Database: "genealogy",
		}

		_, err := svc.CreateBackup(ctx, req)
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("create backup with empty database", func(t *testing.T) {
		req := &CreateBackupRequest{
			Name:     "test",
			Database: "",
		}

		_, err := svc.CreateBackup(ctx, req)
		if err == nil {
			t.Error("expected error for empty database")
		}
	})
}

func TestGetBackup(t *testing.T) {
	svc, repo := newTestBackupService()
	ctx := context.Background()

	// 预先创建一个备份
	repo.Create(ctx, &backup.Backup{
		Name:     "existing_backup",
		Type:     backup.BackupTypeFull,
		Status:   backup.BackupStatusCompleted,
		Database: "genealogy",
	})

	t.Run("get existing backup", func(t *testing.T) {
		result, err := svc.GetBackup(ctx, 1)
		if err != nil {
			t.Fatalf("GetBackup() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected backup, got nil")
		}
		if result.Name != "existing_backup" {
			t.Errorf("expected name 'existing_backup', got '%s'", result.Name)
		}
	})

	t.Run("get non-existent backup", func(t *testing.T) {
		result, err := svc.GetBackup(ctx, 999)
		if err != nil {
			t.Fatalf("GetBackup() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestListBackups(t *testing.T) {
	svc, repo := newTestBackupService()
	ctx := context.Background()

	// 创建多个备份
	for i := 0; i < 5; i++ {
		repo.Create(ctx, &backup.Backup{
			Name:     "backup_" + string(rune('a'+i)),
			Type:     backup.BackupTypeFull,
			Status:   backup.BackupStatusCompleted,
			Database: "genealogy",
		})
	}

	t.Run("list all backups", func(t *testing.T) {
		req := &BackupFilterRequest{
			Page:     1,
			PageSize: 10,
		}

		result, err := svc.ListBackups(ctx, req)
		if err != nil {
			t.Fatalf("ListBackups() error = %v", err)
		}
		if result.Total != 5 {
			t.Errorf("expected total 5, got %d", result.Total)
		}
		if len(result.Data) != 5 {
			t.Errorf("expected 5 backups, got %d", len(result.Data))
		}
	})

	t.Run("list with pagination", func(t *testing.T) {
		req := &BackupFilterRequest{
			Page:     1,
			PageSize: 2,
		}

		result, err := svc.ListBackups(ctx, req)
		if err != nil {
			t.Fatalf("ListBackups() error = %v", err)
		}
		if result.TotalPages != 3 {
			t.Errorf("expected 3 total pages, got %d", result.TotalPages)
		}
	})
}

func TestDeleteBackup(t *testing.T) {
	svc, repo := newTestBackupService()
	ctx := context.Background()

	now := time.Now()
	repo.backups[1] = &backup.Backup{
		ID:        1,
		Name:      "to_delete",
		Type:      backup.BackupTypeFull,
		Status:    backup.BackupStatusCompleted,
		Database:  "genealogy",
		CreatedAt: now,
	}

	t.Run("delete existing backup", func(t *testing.T) {
		err := svc.DeleteBackup(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteBackup() error = %v", err)
		}
		if len(repo.backups) != 0 {
			t.Errorf("expected 0 backups after delete, got %d", len(repo.backups))
		}
	})

	t.Run("delete non-existent backup", func(t *testing.T) {
		err := svc.DeleteBackup(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent backup")
		}
	})

	t.Run("delete running backup should fail", func(t *testing.T) {
		repo.backups[2] = &backup.Backup{
			ID:        2,
			Name:      "running",
			Type:      backup.BackupTypeFull,
			Status:    backup.BackupStatusRunning,
			Database:  "genealogy",
			CreatedAt: now,
		}

		err := svc.DeleteBackup(ctx, 2)
		if err == nil {
			t.Error("expected error when deleting running backup")
		}
	})
}

func TestCleanOldBackups(t *testing.T) {
	svc, repo := newTestBackupService()
	ctx := context.Background()

	now := time.Now()

	// 创建旧备份
	repo.backups[1] = &backup.Backup{
		ID:        1,
		Name:      "old_backup",
		Type:      backup.BackupTypeFull,
		Status:    backup.BackupStatusCompleted,
		Database:  "genealogy",
		CreatedAt: now.Add(-48 * time.Hour),
	}

	// 创建新备份
	repo.backups[2] = &backup.Backup{
		ID:        2,
		Name:      "new_backup",
		Type:      backup.BackupTypeFull,
		Status:    backup.BackupStatusCompleted,
		Database:  "genealogy",
		CreatedAt: now,
	}

	// 创建正在运行的备份
	repo.backups[3] = &backup.Backup{
		ID:        3,
		Name:      "running_backup",
		Type:      backup.BackupTypeFull,
		Status:    backup.BackupStatusRunning,
		Database:  "genealogy",
		CreatedAt: now.Add(-48 * time.Hour),
	}

	t.Run("clean backups older than 24 hours", func(t *testing.T) {
		count, err := svc.CleanOldBackups(ctx, 1)
		if err != nil {
			t.Fatalf("CleanOldBackups() error = %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 deleted backup, got %d", count)
		}
		// running backup 应该保留
		if _, ok := repo.backups[3]; !ok {
			t.Error("running backup should not be deleted")
		}
		// 新备份应该保留
		if _, ok := repo.backups[2]; !ok {
			t.Error("new backup should not be deleted")
		}
	})

	t.Run("clean with invalid days", func(t *testing.T) {
		_, err := svc.CleanOldBackups(ctx, 0)
		if err == nil {
			t.Error("expected error for invalid days")
		}
	})
}

func TestBackupServiceDomainService(t *testing.T) {
	ctx := context.Background()
	repo := newMockBackupRepository()
	domainSvc := backup.NewService(repo)

	t.Run("create backup through domain service", func(t *testing.T) {
		cmd := &backup.CreateBackupCommand{
			Name:     "domain_test",
			Type:     backup.BackupTypeFull,
			Database: "testdb",
		}

		b, err := domainSvc.CreateBackup(ctx, cmd)
		if err != nil {
			t.Fatalf("CreateBackup() error = %v", err)
		}
		if b.Name != "domain_test" {
			t.Errorf("expected name 'domain_test', got '%s'", b.Name)
		}
	})

	t.Run("update backup status", func(t *testing.T) {
		repo.Create(ctx, &backup.Backup{
			Name:     "status_test",
			Type:     backup.BackupTypeFull,
			Status:   backup.BackupStatusPending,
			Database: "testdb",
		})

		err := domainSvc.UpdateBackupStatus(ctx, 1, backup.BackupStatusRunning, "", 0, "")
		if err != nil {
			t.Fatalf("UpdateBackupStatus() error = %v", err)
		}

		b, _ := domainSvc.GetBackupByID(ctx, 1)
		if b.Status != backup.BackupStatusRunning {
			t.Errorf("expected status 'running', got '%s'", b.Status)
		}
	})

	t.Run("get non-existent backup", func(t *testing.T) {
		b, err := domainSvc.GetBackupByID(ctx, 999)
		if err != nil {
			t.Fatalf("GetBackupByID() error = %v", err)
		}
		if b != nil {
			t.Error("expected nil for non-existent backup")
		}
	})

	t.Run("list backups with pagination", func(t *testing.T) {
		filter := backup.BackupFilter{
			Page:     1,
			PageSize: 5,
		}

		result, err := domainSvc.ListBackups(ctx, filter)
		if err != nil {
			t.Fatalf("ListBackups() error = %v", err)
		}
		if result.Page != 1 {
			t.Errorf("expected page 1, got %d", result.Page)
		}
	})

	t.Run("list backups with invalid pagination", func(t *testing.T) {
		filter := backup.BackupFilter{
			Page:     0,
			PageSize: 0,
		}

		result, err := domainSvc.ListBackups(ctx, filter)
		if err != nil {
			t.Fatalf("ListBackups() error = %v", err)
		}
		// 默认值应该被应用
		if result.Page != 1 {
			t.Errorf("expected default page 1, got %d", result.Page)
		}
	})
}

// 测试错误场景
func TestBackupServiceErrors(t *testing.T) {
	_, repo := newTestBackupService()
	ctx := context.Background()

	t.Run("repository error on create", func(t *testing.T) {
		// 创建一个会返回错误的 mock
		badRepo := &mockBackupRepositoryWithError{
			inner: repo,
		}
		badDomainSvc := backup.NewService(badRepo)
		badSvc := &BackupService{
			domainService: badDomainSvc,
			backupDir:     "/tmp/backups",
		}

		req := &CreateBackupRequest{
			Name:     "test",
			Database: "genealogy",
		}

		_, err := badSvc.CreateBackup(ctx, req)
		if err == nil {
			t.Error("expected error from repository")
		}
	})
}

type mockBackupRepositoryWithError struct {
	inner *mockBackupRepository
}

func (m *mockBackupRepositoryWithError) Create(ctx context.Context, b *backup.Backup) error {
	return errors.New("repository error")
}
func (m *mockBackupRepositoryWithError) GetByID(ctx context.Context, id int64) (*backup.Backup, error) {
	return m.inner.GetByID(ctx, id)
}
func (m *mockBackupRepositoryWithError) List(ctx context.Context, filter backup.BackupFilter) (*backup.BackupList, error) {
	return m.inner.List(ctx, filter)
}
func (m *mockBackupRepositoryWithError) Update(ctx context.Context, b *backup.Backup) error {
	return errors.New("update error")
}
func (m *mockBackupRepositoryWithError) Delete(ctx context.Context, id int64) error {
	return m.inner.Delete(ctx, id)
}
func (m *mockBackupRepositoryWithError) DeleteOldBackups(ctx context.Context, before time.Time) (int64, error) {
	return m.inner.DeleteOldBackups(ctx, before)
}