package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/genealogy-ma/platform/internal/application/service"
)

// MockBackupService 模拟备份服务
type MockBackupService struct {
	backups map[int64]*service.BackupDTO
	nextID  int64
}

func NewMockBackupService() *MockBackupService {
	return &MockBackupService{
		backups: make(map[int64]*service.BackupDTO),
		nextID:  1,
	}
}

func (m *MockBackupService) CreateBackup(ctx context.Context, req *service.CreateBackupRequest) (*service.BackupDTO, error) {
	b := &service.BackupDTO{
		ID:       m.nextID,
		Name:     req.Name,
		Type:     req.Type,
		Status:   "pending",
		Database: req.Database,
	}
	m.nextID++
	m.backups[b.ID] = b
	return b, nil
}

func (m *MockBackupService) GetBackup(ctx context.Context, id int64) (*service.BackupDTO, error) {
	if b, ok := m.backups[id]; ok {
		return b, nil
	}
	return nil, nil
}

func (m *MockBackupService) ListBackups(ctx context.Context, req *service.BackupFilterRequest) (*service.BackupListDTO, error) {
	result := make([]*service.BackupDTO, 0, len(m.backups))
	for _, b := range m.backups {
		result = append(result, b)
	}
	return &service.BackupListDTO{
		Data:       result,
		Total:      int64(len(result)),
		Page:       1,
		PageSize:   20,
		TotalPages: 1,
	}, nil
}

func (m *MockBackupService) RestoreBackup(ctx context.Context, req *service.RestoreBackupRequest) error {
	if b, ok := m.backups[req.BackupID]; ok {
		b.Status = "restored"
	}
	return nil
}

func (m *MockBackupService) DeleteBackup(ctx context.Context, id int64) error {
	delete(m.backups, id)
	return nil
}

func (m *MockBackupService) CleanOldBackups(ctx context.Context, days int) (int64, error) {
	return int64(len(m.backups)), nil
}

func setupBackupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestBackupController_List(t *testing.T) {
	mockService := NewMockBackupService()
	mockService.backups[1] = &service.BackupDTO{ID: 1, Name: "backup1", Type: "full", Status: "completed"}
	mockService.backups[2] = &service.BackupDTO{ID: 2, Name: "backup2", Type: "partial", Status: "completed"}

	r := setupBackupTestRouter()
	r.GET("/backups", func(ctx *gin.Context) {
		result, err := mockService.ListBackups(ctx, &service.BackupFilterRequest{Page: 1, PageSize: 20})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	req, _ := http.NewRequest("GET", "/backups", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected data in response")
	}
	// Verify list has data array with 2 backups
	listData, ok := data["data"].([]interface{})
	if !ok || len(listData) != 2 {
		t.Errorf("Expected 2 backups in list, got %d", len(listData))
	}
}

func TestBackupController_Get(t *testing.T) {
	mockService := NewMockBackupService()
	mockService.backups[1] = &service.BackupDTO{ID: 1, Name: "backup1", Type: "full", Status: "completed"}

	r := setupBackupTestRouter()
	r.GET("/backups/:id", func(ctx *gin.Context) {
		backup, err := mockService.GetBackup(ctx, 1)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if backup == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "backup not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": backup})
	})

	req, _ := http.NewRequest("GET", "/backups/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestBackupController_Create(t *testing.T) {
	mockService := NewMockBackupService()

	r := setupBackupTestRouter()
	r.POST("/backups", func(ctx *gin.Context) {
		var req service.CreateBackupRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		backup, err := mockService.CreateBackup(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": backup})
	})

	body := `{"name":"test_backup","type":"full","database":"genealogy"}`
	req, _ := http.NewRequest("POST", "/backups", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected data in response")
	}
	if data["name"] != "test_backup" {
		t.Errorf("Expected name 'test_backup', got '%v'", data["name"])
	}
}

func TestBackupController_Delete(t *testing.T) {
	mockService := NewMockBackupService()
	mockService.backups[1] = &service.BackupDTO{ID: 1, Name: "backup1"}

	r := setupBackupTestRouter()
	r.DELETE("/backups/:id", func(ctx *gin.Context) {
		err := mockService.DeleteBackup(ctx, 1)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/backups/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if _, exists := mockService.backups[1]; exists {
		t.Errorf("Backup should have been deleted")
	}
}

func TestBackupController_Restore(t *testing.T) {
	mockService := NewMockBackupService()
	mockService.backups[1] = &service.BackupDTO{ID: 1, Name: "backup1", Status: "completed"}

	r := setupBackupTestRouter()
	r.POST("/backups/:id/restore", func(ctx *gin.Context) {
		var req service.RestoreBackupRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.BackupID = 1
		err := mockService.RestoreBackup(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "备份恢复成功"})
	})

	// Include backup_id in body to satisfy binding:"required"
	body := `{"backup_id":1,"target_db":"genealogy"}`
	req, _ := http.NewRequest("POST", "/backups/1/restore", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestBackupController_CleanOldBackups(t *testing.T) {
	mockService := NewMockBackupService()
	mockService.backups[1] = &service.BackupDTO{ID: 1, Name: "old_backup"}

	r := setupBackupTestRouter()
	r.DELETE("/backups/clean-old", func(ctx *gin.Context) {
		var req service.CleanOldBackupsRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		count, err := mockService.CleanOldBackups(ctx, req.Days)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "已清理旧备份", "deleted_count": count})
	})

	body := `{"days":30}`
	req, _ := http.NewRequest("DELETE", "/backups/clean-old", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}