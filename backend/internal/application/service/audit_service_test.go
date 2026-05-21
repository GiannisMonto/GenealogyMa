package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/audit"
)

// mockAuditRepository 模拟审计仓储
type mockAuditRepository struct {
	logs   map[int64]*audit.AuditLog
	nextID int64
}

func newMockAuditRepository() *mockAuditRepository {
	return &mockAuditRepository{
		logs:   make(map[int64]*audit.AuditLog),
		nextID: 1,
	}
}

func (m *mockAuditRepository) Create(ctx context.Context, log *audit.AuditLog) error {
	log.ID = m.nextID
	m.nextID++
	log.CreatedAt = time.Now()
	m.logs[log.ID] = log
	return nil
}

func (m *mockAuditRepository) GetByID(ctx context.Context, id int64) (*audit.AuditLog, error) {
	log, ok := m.logs[id]
	if !ok {
		return nil, nil
	}
	return log, nil
}

func (m *mockAuditRepository) List(ctx context.Context, filter audit.AuditLogFilter) (*audit.AuditLogList, error) {
	var result []audit.AuditLog
	for _, log := range m.logs {
		result = append(result, *log)
	}

	total := int64(len(result))
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &audit.AuditLogList{
		Data:       result,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (m *mockAuditRepository) DeleteOldLogs(ctx context.Context, before time.Time) (int64, error) {
	var count int64
	for id, log := range m.logs {
		if log.CreatedAt.Before(before) {
			delete(m.logs, id)
			count++
		}
	}
	return count, nil
}

func newTestAuditService() (*AuditService, *mockAuditRepository) {
	repo := newMockAuditRepository()
	svc := NewAuditService(repo)
	return svc, repo
}

// ===== 测试用例 =====

func TestLogAction(t *testing.T) {
	svc, repo := newTestAuditService()
	ctx := context.Background()

	t.Run("log action successfully", func(t *testing.T) {
		req := &LogActionRequest{
			UserID:       1,
			Username:     "admin",
			Module:       "user",
			Action:       "create",
			ResourceType: "user",
			ResourceID:   "123",
			Description:  "Created new user",
		}

		err := svc.LogAction(ctx, req)
		if err != nil {
			t.Fatalf("LogAction() error = %v", err)
		}
		if len(repo.logs) != 1 {
			t.Errorf("expected 1 log in repo, got %d", len(repo.logs))
		}
	})

	t.Run("log action without optional fields", func(t *testing.T) {
		req := &LogActionRequest{
			UserID:   1,
			Username: "admin",
			Module:   "user",
			Action:   "login",
		}

		err := svc.LogAction(ctx, req)
		if err != nil {
			t.Fatalf("LogAction() error = %v", err)
		}
	})
}

func TestGetLogByID(t *testing.T) {
	svc, repo := newTestAuditService()
	ctx := context.Background()

	// 预先创建一个日志
	repo.Create(ctx, &audit.AuditLog{
		UserID:       1,
		Username:     "admin",
		Module:       "user",
		Action:       "create",
		ResourceType: "user",
		ResourceID:   "123",
	})

	t.Run("get existing log", func(t *testing.T) {
		result, err := svc.GetLogByID(ctx, 1)
		if err != nil {
			t.Fatalf("GetLogByID() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected log, got nil")
		}
		if result.Username != "admin" {
			t.Errorf("expected username 'admin', got '%s'", result.Username)
		}
		if result.Action != "create" {
			t.Errorf("expected action 'create', got '%s'", result.Action)
		}
	})

	t.Run("get non-existent log", func(t *testing.T) {
		result, err := svc.GetLogByID(ctx, 999)
		if err != nil {
			t.Fatalf("GetLogByID() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestListLogs(t *testing.T) {
	svc, repo := newTestAuditService()
	ctx := context.Background()

	// 创建多个日志
	for i := 0; i < 5; i++ {
		repo.Create(ctx, &audit.AuditLog{
			UserID:       1,
			Username:     "admin",
			Module:       "user",
			Action:       "create",
			ResourceType: "user",
		})
	}

	t.Run("list all logs with default pagination", func(t *testing.T) {
		req := &AuditFilterRequest{}

		result, err := svc.ListLogs(ctx, req)
		if err != nil {
			t.Fatalf("ListLogs() error = %v", err)
		}
		if result.Total != 5 {
			t.Errorf("expected total 5, got %d", result.Total)
		}
		if len(result.Data) != 5 {
			t.Errorf("expected 5 logs, got %d", len(result.Data))
		}
	})

	t.Run("list logs with pagination", func(t *testing.T) {
		req := &AuditFilterRequest{
			Page:     1,
			PageSize: 2,
		}

		result, err := svc.ListLogs(ctx, req)
		if err != nil {
			t.Fatalf("ListLogs() error = %v", err)
		}
		if result.TotalPages != 3 {
			t.Errorf("expected 3 total pages, got %d", result.TotalPages)
		}
	})

	t.Run("list logs filtered by module", func(t *testing.T) {
		req := &AuditFilterRequest{
			Module: "user",
			Page:   1,
			PageSize: 10,
		}

		result, err := svc.ListLogs(ctx, req)
		if err != nil {
			t.Fatalf("ListLogs() error = %v", err)
		}
		// All logs have module "user"
		if result.Total != 5 {
			t.Errorf("expected total 5, got %d", result.Total)
		}
	})
}

func TestCleanOldLogs(t *testing.T) {
	svc, repo := newTestAuditService()
	ctx := context.Background()

	now := time.Now()

	// 创建旧日志
	repo.logs[1] = &audit.AuditLog{
		ID:        1,
		UserID:    1,
		Username:  "admin",
		Module:    "user",
		Action:    "login",
		CreatedAt: now.Add(-48 * time.Hour),
	}

	// 创建新日志
	repo.logs[2] = &audit.AuditLog{
		ID:        2,
		UserID:    1,
		Username:  "admin",
		Module:    "user",
		Action:    "login",
		CreatedAt: now,
	}

	t.Run("clean logs older than 24 hours", func(t *testing.T) {
		count, err := svc.CleanOldLogs(ctx, 1)
		if err != nil {
			t.Fatalf("CleanOldLogs() error = %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 deleted log, got %d", count)
		}
		// 新日志应该保留
		if _, ok := repo.logs[2]; !ok {
			t.Error("new log should not be deleted")
		}
		// 旧日志应该被删除
		if _, ok := repo.logs[1]; ok {
			t.Error("old log should be deleted")
		}
	})

	t.Run("clean with invalid days", func(t *testing.T) {
		_, err := svc.CleanOldLogs(ctx, 0)
		if err == nil {
			t.Error("expected error for invalid days")
		}
	})
}

// 测试错误场景
func TestAuditServiceErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("repository error on log action", func(t *testing.T) {
		badRepo := &mockAuditRepositoryError{
			inner: newMockAuditRepository(),
		}
		badSvc := NewAuditService(badRepo)

		req := &LogActionRequest{
			UserID:   1,
			Username: "admin",
			Module:   "user",
			Action:   "create",
		}

		err := badSvc.LogAction(ctx, req)
		if err == nil {
			t.Error("expected error from repository")
		}
	})

	t.Run("repository error on list logs", func(t *testing.T) {
		badRepo := &mockAuditRepositoryError{
			inner: newMockAuditRepository(),
		}
		badSvc := NewAuditService(badRepo)

		_, err := badSvc.ListLogs(ctx, &AuditFilterRequest{})
		if err == nil {
			t.Error("expected error from repository")
		}
	})
}

type mockAuditRepositoryError struct {
	inner *mockAuditRepository
}

func (m *mockAuditRepositoryError) Create(ctx context.Context, log *audit.AuditLog) error {
	return errors.New("repository error")
}
func (m *mockAuditRepositoryError) GetByID(ctx context.Context, id int64) (*audit.AuditLog, error) {
	return m.inner.GetByID(ctx, id)
}
func (m *mockAuditRepositoryError) List(ctx context.Context, filter audit.AuditLogFilter) (*audit.AuditLogList, error) {
	return nil, errors.New("repository error")
}
func (m *mockAuditRepositoryError) DeleteOldLogs(ctx context.Context, before time.Time) (int64, error) {
	return m.inner.DeleteOldLogs(ctx, before)
}