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

// MockAuditService 模拟审计服务
type MockAuditService struct {
	logs   map[int64]*service.AuditLogDTO
	nextID int64
}

func NewMockAuditService() *MockAuditService {
	return &MockAuditService{
		logs:   make(map[int64]*service.AuditLogDTO),
		nextID: 1,
	}
}

func (m *MockAuditService) LogAction(ctx context.Context, req *service.LogActionRequest) error {
	log := &service.AuditLogDTO{
		ID:           m.nextID,
		UserID:       req.UserID,
		Username:     req.Username,
		Module:       req.Module,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Description:  req.Description,
	}
	m.nextID++
	m.logs[log.ID] = log
	return nil
}

func (m *MockAuditService) GetLogByID(ctx context.Context, id int64) (*service.AuditLogDTO, error) {
	if log, ok := m.logs[id]; ok {
		return log, nil
	}
	return nil, nil
}

func (m *MockAuditService) ListLogs(ctx context.Context, req *service.AuditFilterRequest) (*service.AuditLogListDTO, error) {
	result := make([]*service.AuditLogDTO, 0, len(m.logs))
	for _, log := range m.logs {
		if req.Module != "" && log.Module != req.Module {
			continue
		}
		if req.Action != "" && log.Action != req.Action {
			continue
		}
		result = append(result, log)
	}
	return &service.AuditLogListDTO{
		Data:       result,
		Total:      int64(len(result)),
		Page:       1,
		PageSize:   20,
		TotalPages: 1,
	}, nil
}

func (m *MockAuditService) CleanOldLogs(ctx context.Context, days int) (int64, error) {
	return int64(len(m.logs)), nil
}

// ===== 审计日志测试 =====

func TestAuditController_ListLogs(t *testing.T) {
	mockService := NewMockAuditService()
	mockService.logs[1] = &service.AuditLogDTO{ID: 1, UserID: 1, Username: "admin", Module: "person", Action: "create"}
	mockService.logs[2] = &service.AuditLogDTO{ID: 2, UserID: 1, Username: "admin", Module: "person", Action: "update"}

	r := setupTestRouter()
	r.GET("/audit-logs", func(ctx *gin.Context) {
		var req service.AuditFilterRequest
		req.Page = 1
		req.PageSize = 20
		result, err := mockService.ListLogs(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": result.Data, "total": result.Total})
	})

	req, _ := http.NewRequest("GET", "/audit-logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 audit logs, got %d", len(data))
	}
}

func TestAuditController_GetLogByID(t *testing.T) {
	mockService := NewMockAuditService()
	mockService.logs[1] = &service.AuditLogDTO{ID: 1, UserID: 1, Username: "admin", Module: "person", Action: "create"}

	r := setupTestRouter()
	r.GET("/audit-logs/:id", func(ctx *gin.Context) {
		id := int64(1)
		log, err := mockService.GetLogByID(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if log == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": log})
	})

	req, _ := http.NewRequest("GET", "/audit-logs/1", nil)
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
	if data["module"] != "person" {
		t.Errorf("Expected module 'person', got '%v'", data["module"])
	}
}

func TestAuditController_LogAction(t *testing.T) {
	mockService := NewMockAuditService()

	r := setupTestRouter()
	r.POST("/audit-logs", func(ctx *gin.Context) {
		var req service.LogActionRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := mockService.LogAction(ctx, &req); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "记录成功"})
	})

	body := `{"user_id":1,"username":"admin","module":"person","action":"create","resource_type":"member","description":"创建了新成员"}`
	req, _ := http.NewRequest("POST", "/audit-logs", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuditController_CleanOldLogs(t *testing.T) {
	mockService := NewMockAuditService()
	mockService.logs[1] = &service.AuditLogDTO{ID: 1, UserID: 1, Username: "admin", Module: "person", Action: "create"}

	r := setupTestRouter()
	r.DELETE("/audit-logs/clean", func(ctx *gin.Context) {
		days := 90
		count, err := mockService.CleanOldLogs(ctx, days)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": gin.H{"deleted_count": count}})
	})

	req, _ := http.NewRequest("DELETE", "/audit-logs/clean", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuditController_ListLogsByModule(t *testing.T) {
	mockService := NewMockAuditService()
	mockService.logs[1] = &service.AuditLogDTO{ID: 1, UserID: 1, Username: "admin", Module: "person", Action: "create"}
	mockService.logs[2] = &service.AuditLogDTO{ID: 2, UserID: 1, Username: "admin", Module: "auth", Action: "login"}

	r := setupTestRouter()
	r.GET("/audit-logs", func(ctx *gin.Context) {
		var req service.AuditFilterRequest
		req.Module = "person"
		req.Page = 1
		req.PageSize = 20
		result, err := mockService.ListLogs(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": result.Data, "total": result.Total})
	})

	req, _ := http.NewRequest("GET", "/audit-logs?module=person", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 1 {
		t.Errorf("Expected 1 audit log filtered by module, got %d", len(data))
	}
}
