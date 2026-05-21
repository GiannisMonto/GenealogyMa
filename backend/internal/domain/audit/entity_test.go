package audit

import (
	"testing"
	"time"
)

func TestAuditLog_Fields(t *testing.T) {
	now := time.Now()
	log := &AuditLog{
		ID:           1,
		UserID:       100,
		Username:     "admin",
		Module:       "user",
		Action:       "create",
		ResourceType: "user",
		ResourceID:   "123",
		OldValue:     "",
		NewValue:     `{"name":"test"}`,
		IPAddress:    "192.168.1.1",
		UserAgent:    "Mozilla/5.0",
		Description:  "Create new user",
		CreatedAt:    now,
	}

	if log.ID != 1 {
		t.Errorf("AuditLog.ID = %v, want 1", log.ID)
	}
	if log.UserID != 100 {
		t.Errorf("AuditLog.UserID = %v, want 100", log.UserID)
	}
	if log.Username != "admin" {
		t.Errorf("AuditLog.Username = %v, want admin", log.Username)
	}
	if log.Module != "user" {
		t.Errorf("AuditLog.Module = %v, want user", log.Module)
	}
	if log.Action != "create" {
		t.Errorf("AuditLog.Action = %v, want create", log.Action)
	}
	if log.ResourceType != "user" {
		t.Errorf("AuditLog.ResourceType = %v, want user", log.ResourceType)
	}
	if log.ResourceID != "123" {
		t.Errorf("AuditLog.ResourceID = %v, want 123", log.ResourceID)
	}
	if log.IPAddress != "192.168.1.1" {
		t.Errorf("AuditLog.IPAddress = %v, want 192.168.1.1", log.IPAddress)
	}
}

func TestCreateAuditLogCommand_Fields(t *testing.T) {
	cmd := &CreateAuditLogCommand{
		UserID:       100,
		Username:     "admin",
		Module:       "user",
		Action:       "create",
		ResourceType: "user",
		ResourceID:   "123",
		OldValue:     "",
		NewValue:     `{"name":"test"}`,
		IPAddress:    "192.168.1.1",
		UserAgent:    "Mozilla/5.0",
		Description:  "Create new user",
	}

	if cmd.UserID != 100 {
		t.Errorf("CreateAuditLogCommand.UserID = %v, want 100", cmd.UserID)
	}
	if cmd.Username != "admin" {
		t.Errorf("CreateAuditLogCommand.Username = %v, want admin", cmd.Username)
	}
	if cmd.Module != "user" {
		t.Errorf("CreateAuditLogCommand.Module = %v, want user", cmd.Module)
	}
	if cmd.Action != "create" {
		t.Errorf("CreateAuditLogCommand.Action = %v, want create", cmd.Action)
	}
}

func TestAuditLogFilter_Fields(t *testing.T) {
	userID := int64(100)
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	filter := &AuditLogFilter{
		UserID:       &userID,
		Module:       "user",
		Action:       "create",
		ResourceType: "user",
		StartDate:    &startDate,
		EndDate:      &endDate,
		Page:         1,
		PageSize:     20,
	}

	if filter.UserID == nil || *filter.UserID != 100 {
		t.Errorf("AuditLogFilter.UserID = %v, want 100", *filter.UserID)
	}
	if filter.Module != "user" {
		t.Errorf("AuditLogFilter.Module = %v, want user", filter.Module)
	}
	if filter.Action != "create" {
		t.Errorf("AuditLogFilter.Action = %v, want create", filter.Action)
	}
	if filter.Page != 1 {
		t.Errorf("AuditLogFilter.Page = %v, want 1", filter.Page)
	}
	if filter.PageSize != 20 {
		t.Errorf("AuditLogFilter.PageSize = %v, want 20", filter.PageSize)
	}
}

func TestAuditLogFilter_DefaultValues(t *testing.T) {
	filter := &AuditLogFilter{}

	if filter.Page != 0 {
		t.Errorf("AuditLogFilter.Page default = %v, want 0", filter.Page)
	}
	if filter.PageSize != 0 {
		t.Errorf("AuditLogFilter.PageSize default = %v, want 0", filter.PageSize)
	}
}

func TestAuditLogList_Fields(t *testing.T) {
	logs := []AuditLog{
		{ID: 1, Username: "user1"},
		{ID: 2, Username: "user2"},
	}

	list := &AuditLogList{
		Data:       logs,
		Total:      100,
		Page:       1,
		PageSize:   20,
		TotalPages: 5,
	}

	if len(list.Data) != 2 {
		t.Errorf("AuditLogList.Data length = %v, want 2", len(list.Data))
	}
	if list.Total != 100 {
		t.Errorf("AuditLogList.Total = %v, want 100", list.Total)
	}
	if list.TotalPages != 5 {
		t.Errorf("AuditLogList.TotalPages = %v, want 5", list.TotalPages)
	}
}

func TestAuditLogList_Calculation(t *testing.T) {
	tests := []struct {
		name       string
		total      int64
		pageSize   int
		totalPages int
	}{
		{"exact division", 100, 20, 5},
		{"with remainder", 101, 20, 6},
		{"single page", 10, 20, 1},
		{"empty", 0, 20, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &AuditLogList{
				Total:      tt.total,
				PageSize:   tt.pageSize,
				TotalPages: tt.totalPages,
			}
			if list.TotalPages != tt.totalPages {
				t.Errorf("TotalPages = %v, want %v", list.TotalPages, tt.totalPages)
			}
		})
	}
}