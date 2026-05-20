package audit

import (
	"context"
	"testing"
	"time"
)

type mockRepository struct {
	logs     map[int64]*AuditLog
	nextID   int64
	onCreate func(log *AuditLog) error
}

func newMockRepo() *mockRepository {
	return &mockRepository{
		logs:   make(map[int64]*AuditLog),
		nextID: 1,
	}
}

func (m *mockRepository) Create(ctx context.Context, log *AuditLog) error {
	if m.onCreate != nil {
		return m.onCreate(log)
	}
	log.ID = m.nextID
	m.nextID++
	m.logs[log.ID] = log
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id int64) (*AuditLog, error) {
	log, ok := m.logs[id]
	if !ok {
		return nil, nil
	}
	return log, nil
}

func (m *mockRepository) List(ctx context.Context, filter AuditLogFilter) (*AuditLogList, error) {
	var result []AuditLog
	for _, log := range m.logs {
		if filter.UserID != nil && log.UserID != *filter.UserID {
			continue
		}
		if filter.Module != "" && log.Module != filter.Module {
			continue
		}
		if filter.Action != "" && log.Action != filter.Action {
			continue
		}
		result = append(result, *log)
	}

	total := int64(len(result))
	start := (filter.Page - 1) * filter.PageSize
	end := start + filter.PageSize

	if start > len(result) {
		start = len(result)
	}
	if end > len(result) {
		end = len(result)
	}

	return &AuditLogList{
		Data:       result[start:end],
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	}, nil
}

func (m *mockRepository) DeleteOldLogs(ctx context.Context, before time.Time) (int64, error) {
	var deleted int64
	for id, log := range m.logs {
		if log.CreatedAt.Before(before) {
			delete(m.logs, id)
			deleted++
		}
	}
	return deleted, nil
}

func TestAuditService_LogAction(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	cmd := &CreateAuditLogCommand{
		UserID:       1,
		Username:     "admin",
		Module:       "person",
		Action:       "create",
		ResourceType: "person",
		ResourceID:   "123",
		Description:  "创建人物",
	}

	err := svc.LogAction(context.Background(), cmd)
	if err != nil {
		t.Fatalf("LogAction failed: %v", err)
	}

	if len(repo.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(repo.logs))
	}

	for _, log := range repo.logs {
		if log.UserID != cmd.UserID {
			t.Errorf("UserID = %d, want %d", log.UserID, cmd.UserID)
		}
		if log.Username != cmd.Username {
			t.Errorf("Username = %s, want %s", log.Username, cmd.Username)
		}
		if log.Module != cmd.Module {
			t.Errorf("Module = %s, want %s", log.Module, cmd.Module)
		}
	}
}

func TestAuditService_GetLogByID(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	now := time.Now()
	log := &AuditLog{
		ID:          1,
		UserID:      1,
		Username:    "admin",
		Module:      "user",
		Action:      "login",
		CreatedAt:   now,
	}
	repo.logs[1] = log

	found, err := svc.GetLogByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetLogByID failed: %v", err)
	}
	if found == nil {
		t.Fatal("GetLogByID returned nil")
	}
	if found.ID != 1 {
		t.Errorf("ID = %d, want 1", found.ID)
	}
}

func TestAuditService_ListLogs(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	now := time.Now()
	for i := int64(1); i <= 5; i++ {
		repo.logs[i] = &AuditLog{
			ID:        i,
			UserID:    1,
			Module:    "person",
			Action:    "update",
			CreatedAt: now,
		}
	}
	repo.nextID = 6

	filter := AuditLogFilter{
		Module:   "person",
		Page:     1,
		PageSize: 10,
	}

	result, err := svc.ListLogs(context.Background(), filter)
	if err != nil {
		t.Fatalf("ListLogs failed: %v", err)
	}
	if result.Total != 5 {
		t.Errorf("Total = %d, want 5", result.Total)
	}
}

func TestAuditService_ListLogs_Pagination(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	now := time.Now()
	for i := int64(1); i <= 25; i++ {
		repo.logs[i] = &AuditLog{
			ID:        i,
			UserID:    1,
			Module:    "person",
			Action:    "update",
			CreatedAt: now,
		}
	}
	repo.nextID = 26

	filter := AuditLogFilter{
		Module:   "person",
		Page:     2,
		PageSize: 10,
	}

	result, err := svc.ListLogs(context.Background(), filter)
	if err != nil {
		t.Fatalf("ListLogs failed: %v", err)
	}
	if result.Total != 25 {
		t.Errorf("Total = %d, want 25", result.Total)
	}
	if result.Page != 2 {
		t.Errorf("Page = %d, want 2", result.Page)
	}
	if len(result.Data) != 10 {
		t.Errorf("Data length = %d, want 10", len(result.Data))
	}
}

func TestAuditService_CleanOldLogs(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	now := time.Now()
	repo.logs[1] = &AuditLog{ID: 1, CreatedAt: now.Add(-48 * time.Hour)}
	repo.logs[2] = &AuditLog{ID: 2, CreatedAt: now.Add(-72 * time.Hour)}
	repo.logs[3] = &AuditLog{ID: 3, CreatedAt: now.Add(-24 * time.Hour)}

	// days=2 means delete logs older than 48 hours
	// log 1: 48h old - not deleted (at boundary)
	// log 2: 72h old - deleted
	// log 3: 24h old - not deleted
	deleted, err := svc.CleanOldLogs(context.Background(), 2)
	if err != nil {
		t.Fatalf("CleanOldLogs failed: %v", err)
	}
	if deleted != 1 {
		t.Errorf("deleted = %d, want 1", deleted)
	}
	if len(repo.logs) != 2 {
		t.Errorf("remaining logs = %d, want 2", len(repo.logs))
	}
}

func TestAuditService_CleanOldLogs_InvalidDays(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.CleanOldLogs(context.Background(), 0)
	if err == nil {
		t.Error("expected error for days=0")
	}

	_, err = svc.CleanOldLogs(context.Background(), -1)
	if err == nil {
		t.Error("expected error for days=-1")
	}
}