package backup

import (
	"testing"
	"time"
)

func TestBackupTypeConstants(t *testing.T) {
	if BackupTypeFull != "full" {
		t.Errorf("BackupTypeFull = %v, want full", BackupTypeFull)
	}
	if BackupTypePartial != "partial" {
		t.Errorf("BackupTypePartial = %v, want partial", BackupTypePartial)
	}
}

func TestBackupStatusConstants(t *testing.T) {
	if BackupStatusPending != "pending" {
		t.Errorf("BackupStatusPending = %v, want pending", BackupStatusPending)
	}
	if BackupStatusRunning != "running" {
		t.Errorf("BackupStatusRunning = %v, want running", BackupStatusRunning)
	}
	if BackupStatusCompleted != "completed" {
		t.Errorf("BackupStatusCompleted = %v, want completed", BackupStatusCompleted)
	}
	if BackupStatusFailed != "failed" {
		t.Errorf("BackupStatusFailed = %v, want failed", BackupStatusFailed)
	}
}

func TestBackupFilterDefaultValues(t *testing.T) {
	f := &BackupFilter{
		Database: "genealogy",
		Page:     1,
		PageSize: 20,
	}

	if f.Page < 1 {
		t.Error("default filter should have valid page")
	}
	if f.PageSize < 1 {
		t.Error("default filter should have valid page size")
	}
}

func TestBackupListPagination(t *testing.T) {
	bl := &BackupList{
		Data:       []Backup{},
		Total:      100,
		Page:       1,
		PageSize:   20,
		TotalPages: 5,
	}

	if bl.TotalPages != 5 {
		t.Errorf("BackupList.TotalPages = %d, want 5", bl.TotalPages)
	}
}

func TestBackupWithAllFields(t *testing.T) {
	now := time.Now()
	startedAt := now.Add(-1 * time.Hour)
	completedAt := now

	b := &Backup{
		ID:           1,
		Name:         "full_backup_2024",
		Type:         BackupTypeFull,
		Status:       BackupStatusCompleted,
		FilePath:     "/backups/full_backup_2024.sql",
		FileSize:     1024 * 1024 * 100,
		Database:     "genealogy",
		Tables:       []string{"members", "spouses", "parent_child_relations"},
		ErrorMessage: "",
		StartedAt:    &startedAt,
		CompletedAt:  &completedAt,
		CreatedAt:    now,
		CreatedBy:    1,
	}

	if b.Type != BackupTypeFull {
		t.Errorf("Backup.Type = %v, want full", b.Type)
	}
	if b.Status != BackupStatusCompleted {
		t.Errorf("Backup.Status = %v, want completed", b.Status)
	}
	if b.Database != "genealogy" {
		t.Errorf("Backup.Database = %v, want genealogy", b.Database)
	}
	if len(b.Tables) != 3 {
		t.Errorf("Backup.Tables length = %d, want 3", len(b.Tables))
	}
}

func TestCreateBackupCommand(t *testing.T) {
	cmd := &CreateBackupCommand{
		Name:      "full_backup",
		Type:      BackupTypeFull,
		Database:  "genealogy",
		Tables:    []string{"members"},
		CreatedBy: 1,
	}

	if cmd.Name != "full_backup" {
		t.Errorf("CreateBackupCommand.Name = %v, want full_backup", cmd.Name)
	}
	if cmd.Type != BackupTypeFull {
		t.Errorf("CreateBackupCommand.Type = %v, want full", cmd.Type)
	}
	if cmd.Database != "genealogy" {
		t.Errorf("CreateBackupCommand.Database = %v, want genealogy", cmd.Database)
	}
	if cmd.CreatedBy != 1 {
		t.Errorf("CreateBackupCommand.CreatedBy = %v, want 1", cmd.CreatedBy)
	}
}

func TestRestoreBackupCommand(t *testing.T) {
	cmd := &RestoreBackupCommand{
		BackupID:  1,
		TargetDB:  "genealogy_restore",
		CreatedBy: 1,
	}

	if cmd.BackupID != 1 {
		t.Errorf("RestoreBackupCommand.BackupID = %v, want 1", cmd.BackupID)
	}
	if cmd.TargetDB != "genealogy_restore" {
		t.Errorf("RestoreBackupCommand.TargetDB = %v, want genealogy_restore", cmd.TargetDB)
	}
	if cmd.CreatedBy != 1 {
		t.Errorf("RestoreBackupCommand.CreatedBy = %v, want 1", cmd.CreatedBy)
	}
}

func TestBackupEmptyTables(t *testing.T) {
	b := &Backup{
		Name:     "backup",
		Type:     BackupTypeFull,
		Database: "genealogy",
		Tables:   []string{},
	}

	if len(b.Tables) != 0 {
		t.Errorf("Backup.Tables should be empty slice, got length %d", len(b.Tables))
	}
}

func TestBackupStatusTransitions(t *testing.T) {
	b := &Backup{
		Status: BackupStatusPending,
	}

	if b.Status != BackupStatusPending {
		t.Error("Backup should start with Pending status")
	}

	b.Status = BackupStatusRunning
	if b.Status != BackupStatusRunning {
		t.Error("Backup status should transition to Running")
	}

	b.Status = BackupStatusCompleted
	if b.Status != BackupStatusCompleted {
		t.Error("Backup status should transition to Completed")
	}
}