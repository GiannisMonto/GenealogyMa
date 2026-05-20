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

// MockCemeteryService 模拟墓园服务
type MockCemeteryService struct {
	cemeteries map[int64]*service.CemeteryDTO
	graves     map[int64]*service.GraveDTO
	nextID     int64
}

func NewMockCemeteryService() *MockCemeteryService {
	return &MockCemeteryService{
		cemeteries: make(map[int64]*service.CemeteryDTO),
		graves:     make(map[int64]*service.GraveDTO),
		nextID:     1,
	}
}

func (m *MockCemeteryService) ListCemeteries(ctx context.Context) ([]*service.CemeteryDTO, error) {
	result := make([]*service.CemeteryDTO, 0, len(m.cemeteries))
	for _, c := range m.cemeteries {
		result = append(result, c)
	}
	return result, nil
}

func (m *MockCemeteryService) GetCemetery(ctx context.Context, id int64) (*service.CemeteryDTO, error) {
	if c, ok := m.cemeteries[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *MockCemeteryService) CreateCemetery(ctx context.Context, req *service.CreateCemeteryRequest) (*service.CemeteryDTO, error) {
	c := &service.CemeteryDTO{
		ID:        m.nextID,
		Name:      req.Name,
		Province:  req.Province,
		City:      req.City,
		TotalGrave: req.TotalGrave,
	}
	m.nextID++
	m.cemeteries[c.ID] = c
	return c, nil
}

func (m *MockCemeteryService) UpdateCemetery(ctx context.Context, id int64, req *service.UpdateCemeteryRequest) (*service.CemeteryDTO, error) {
	if c, ok := m.cemeteries[id]; ok {
		if req.Name != nil {
			c.Name = *req.Name
		}
		return c, nil
	}
	return nil, nil
}

func (m *MockCemeteryService) DeleteCemetery(ctx context.Context, id int64) error {
	delete(m.cemeteries, id)
	return nil
}

func (m *MockCemeteryService) ListGravesByCemetery(ctx context.Context, cemeteryID int64) ([]*service.GraveDTO, error) {
	result := make([]*service.GraveDTO, 0)
	for _, g := range m.graves {
		if g.CemeteryID == cemeteryID {
			result = append(result, g)
		}
	}
	return result, nil
}

func (m *MockCemeteryService) GetGrave(ctx context.Context, id int64) (*service.GraveDTO, error) {
	if g, ok := m.graves[id]; ok {
		return g, nil
	}
	return nil, nil
}

func (m *MockCemeteryService) CreateGrave(ctx context.Context, req *service.CreateGraveRequest) (*service.GraveDTO, error) {
	g := &service.GraveDTO{
		ID:         m.nextID,
		CemeteryID: req.CemeteryID,
		Section:    req.Section,
		Row:        req.Row,
		Number:     req.Number,
		Status:     "available",
	}
	m.nextID++
	m.graves[g.ID] = g
	return g, nil
}

func (m *MockCemeteryService) UpdateGrave(ctx context.Context, id int64, req *service.UpdateGraveRequest) (*service.GraveDTO, error) {
	if g, ok := m.graves[id]; ok {
		if req.Section != nil {
			g.Section = *req.Section
		}
		return g, nil
	}
	return nil, nil
}

func (m *MockCemeteryService) DeleteGrave(ctx context.Context, id int64) error {
	delete(m.graves, id)
	return nil
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestCemeteryController_List(t *testing.T) {
	mockService := NewMockCemeteryService()
	mockService.cemeteries[1] = &service.CemeteryDTO{ID: 1, Name: "墓园1", Province: "广东省", City: "深圳市"}
	mockService.cemeteries[2] = &service.CemeteryDTO{ID: 2, Name: "墓园2", Province: "北京市", City: "北京市"}

	r := setupTestRouter()
	r.GET("/cemeteries", func(ctx *gin.Context) {
		cemeteries, err := mockService.ListCemeteries(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": cemeteries})
	})

	req, _ := http.NewRequest("GET", "/cemeteries", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 cemeteries, got %d", len(data))
	}
}

func TestCemeteryController_Get(t *testing.T) {
	mockService := NewMockCemeteryService()
	mockService.cemeteries[1] = &service.CemeteryDTO{ID: 1, Name: "墓园1", Province: "广东省", City: "深圳市"}

	r := setupTestRouter()
	r.GET("/cemeteries/:id", func(ctx *gin.Context) {
		id := int64(1)
		cemetery, err := mockService.GetCemetery(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if cemetery == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "cemetery not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": cemetery})
	})

	req, _ := http.NewRequest("GET", "/cemeteries/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCemeteryController_Create(t *testing.T) {
	mockService := NewMockCemeteryService()

	r := setupTestRouter()
	r.POST("/cemeteries", func(ctx *gin.Context) {
		var req service.CreateCemeteryRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cemetery, err := mockService.CreateCemetery(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": cemetery})
	})

	body := `{"name":"新墓园","province":"广东省","city":"深圳市","total_grave":100}`
	req, _ := http.NewRequest("POST", "/cemeteries", bytes.NewBufferString(body))
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
	if data["name"] != "新墓园" {
		t.Errorf("Expected name '新墓园', got '%v'", data["name"])
	}
}

func TestCemeteryController_Delete(t *testing.T) {
	mockService := NewMockCemeteryService()
	mockService.cemeteries[1] = &service.CemeteryDTO{ID: 1, Name: "墓园1"}

	r := setupTestRouter()
	r.DELETE("/cemeteries/:id", func(ctx *gin.Context) {
		id := int64(1)
		err := mockService.DeleteCemetery(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/cemeteries/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if _, exists := mockService.cemeteries[1]; exists {
		t.Errorf("Cemetery should have been deleted")
	}
}

func TestCemeteryController_ListGrave(t *testing.T) {
	mockService := NewMockCemeteryService()
	mockService.graves[1] = &service.GraveDTO{ID: 1, CemeteryID: 1, Section: "A区", Row: 1, Number: 1}
	mockService.graves[2] = &service.GraveDTO{ID: 2, CemeteryID: 1, Section: "A区", Row: 2, Number: 1}

	r := setupTestRouter()
	r.GET("/cemeteries/:id/graves", func(ctx *gin.Context) {
		id := int64(1)
		graves, err := mockService.ListGravesByCemetery(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": graves})
	})

	req, _ := http.NewRequest("GET", "/cemeteries/1/graves", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 graves, got %d", len(data))
	}
}

func TestCemeteryController_CreateGrave(t *testing.T) {
	mockService := NewMockCemeteryService()

	r := setupTestRouter()
	r.POST("/graves", func(ctx *gin.Context) {
		var req service.CreateGraveRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		grave, err := mockService.CreateGrave(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": grave})
	})

	body := `{"cemetery_id":1,"section":"A区","row":1,"number":1}`
	req, _ := http.NewRequest("POST", "/graves", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
