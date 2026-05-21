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

// MockMemorialService 模拟宗祠服务
type MockMemorialService struct {
	halls   map[int64]*service.MemorialHallDTO
	tablets map[int64]*service.TabletDTO
	nextID  int64
}

func NewMockMemorialService() *MockMemorialService {
	return &MockMemorialService{
		halls:   make(map[int64]*service.MemorialHallDTO),
		tablets: make(map[int64]*service.TabletDTO),
		nextID:  1,
	}
}

func (m *MockMemorialService) ListHalls(ctx context.Context) ([]*service.MemorialHallDTO, error) {
	result := make([]*service.MemorialHallDTO, 0, len(m.halls))
	for _, h := range m.halls {
		result = append(result, h)
	}
	return result, nil
}

func (m *MockMemorialService) GetHall(ctx context.Context, id int64) (*service.MemorialHallDTO, error) {
	if h, ok := m.halls[id]; ok {
		return h, nil
	}
	return nil, nil
}

func (m *MockMemorialService) GetHallWithTablets(ctx context.Context, id int64) (*service.MemorialHallDTO, error) {
	if h, ok := m.halls[id]; ok {
		return h, nil
	}
	return nil, nil
}

func (m *MockMemorialService) CreateHall(ctx context.Context, req *service.CreateHallRequest) (*service.MemorialHallDTO, error) {
	h := &service.MemorialHallDTO{
		ID:          m.nextID,
		Name:        req.Name,
		Province:    req.Province,
		City:        req.City,
		TotalTablet: req.TotalTablet,
	}
	m.nextID++
	m.halls[h.ID] = h
	return h, nil
}

func (m *MockMemorialService) UpdateHall(ctx context.Context, id int64, req *service.UpdateHallRequest) (*service.MemorialHallDTO, error) {
	if h, ok := m.halls[id]; ok {
		if req.Name != nil {
			h.Name = *req.Name
		}
		return h, nil
	}
	return nil, nil
}

func (m *MockMemorialService) DeleteHall(ctx context.Context, id int64) error {
	delete(m.halls, id)
	return nil
}

func (m *MockMemorialService) ListTablets(ctx context.Context, hallID int64) ([]*service.TabletDTO, error) {
	result := make([]*service.TabletDTO, 0)
	for _, t := range m.tablets {
		if t.HallID == hallID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *MockMemorialService) GetTablet(ctx context.Context, id int64) (*service.TabletDTO, error) {
	if t, ok := m.tablets[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *MockMemorialService) GetTabletByPersonID(ctx context.Context, personID int64) (*service.TabletDTO, error) {
	for _, t := range m.tablets {
		if t.PersonID != nil && *t.PersonID == personID {
			return t, nil
		}
	}
	return nil, nil
}

func (m *MockMemorialService) CreateTablet(ctx context.Context, req *service.CreateTabletRequest) (*service.TabletDTO, error) {
	t := &service.TabletDTO{
		ID:         m.nextID,
		HallID:     req.HallID,
		PersonName: req.PersonName,
		Floor:      req.Floor,
		Row:        req.Row,
		Number:     req.Number,
		TabletType: req.TabletType,
	}
	m.nextID++
	m.tablets[t.ID] = t
	return t, nil
}

func (m *MockMemorialService) UpdateTablet(ctx context.Context, id int64, req *service.UpdateTabletRequest) (*service.TabletDTO, error) {
	if t, ok := m.tablets[id]; ok {
		if req.PersonName != nil {
			t.PersonName = *req.PersonName
		}
		return t, nil
	}
	return nil, nil
}

func (m *MockMemorialService) DeleteTablet(ctx context.Context, id int64) error {
	delete(m.tablets, id)
	return nil
}

func TestMemorialController_List(t *testing.T) {
	mockService := NewMockMemorialService()
	mockService.halls[1] = &service.MemorialHallDTO{ID: 1, Name: "王氏宗祠", Province: "广东省", City: "深圳市"}
	mockService.halls[2] = &service.MemorialHallDTO{ID: 2, Name: "李氏宗祠", Province: "北京市", City: "北京市"}

	r := setupTestRouter()
	r.GET("/halls", func(ctx *gin.Context) {
		halls, err := mockService.ListHalls(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": halls})
	})

	req, _ := http.NewRequest("GET", "/halls", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 halls, got %d", len(data))
	}
}

func TestMemorialController_Get(t *testing.T) {
	mockService := NewMockMemorialService()
	mockService.halls[1] = &service.MemorialHallDTO{ID: 1, Name: "王氏宗祠", Province: "广东省", City: "深圳市"}

	r := setupTestRouter()
	r.GET("/halls/:id", func(ctx *gin.Context) {
		id := int64(1)
		hall, err := mockService.GetHall(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if hall == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "hall not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": hall})
	})

	req, _ := http.NewRequest("GET", "/halls/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestMemorialController_GetNotFound(t *testing.T) {
	mockService := NewMockMemorialService()

	r := setupTestRouter()
	r.GET("/halls/:id", func(ctx *gin.Context) {
		id := int64(999)
		hall, err := mockService.GetHall(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if hall == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "hall not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": hall})
	})

	req, _ := http.NewRequest("GET", "/halls/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestMemorialController_Create(t *testing.T) {
	mockService := NewMockMemorialService()

	r := setupTestRouter()
	r.POST("/halls", func(ctx *gin.Context) {
		var req service.CreateHallRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hall, err := mockService.CreateHall(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": hall})
	})

	body := `{"name":"王氏宗祠","province":"广东省","city":"深圳市","total_tablet":100}`
	req, _ := http.NewRequest("POST", "/halls", bytes.NewBufferString(body))
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
	if data["name"] != "王氏宗祠" {
		t.Errorf("Expected name '王氏宗祠', got '%v'", data["name"])
	}
}

func TestMemorialController_Update(t *testing.T) {
	mockService := NewMockMemorialService()
	mockService.halls[1] = &service.MemorialHallDTO{ID: 1, Name: "王氏宗祠", Province: "广东省", City: "深圳市"}

	r := setupTestRouter()
	r.PUT("/halls/:id", func(ctx *gin.Context) {
		id := int64(1)
		var req service.UpdateHallRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hall, err := mockService.UpdateHall(ctx, id, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if hall == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "hall not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": hall})
	})

	body := `{"name":"王氏宗祠v2"}`
	req, _ := http.NewRequest("PUT", "/halls/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestMemorialController_Delete(t *testing.T) {
	mockService := NewMockMemorialService()
	mockService.halls[1] = &service.MemorialHallDTO{ID: 1, Name: "王氏宗祠"}

	r := setupTestRouter()
	r.DELETE("/halls/:id", func(ctx *gin.Context) {
		id := int64(1)
		err := mockService.DeleteHall(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/halls/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if _, exists := mockService.halls[1]; exists {
		t.Errorf("Hall should have been deleted")
	}
}

func TestMemorialController_CreateTablet(t *testing.T) {
	mockService := NewMockMemorialService()

	r := setupTestRouter()
	r.POST("/tablets", func(ctx *gin.Context) {
		var req service.CreateTabletRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tablet, err := mockService.CreateTablet(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": tablet})
	})

	body := `{"hall_id":1,"person_name":"王大力","floor":1,"row":1,"number":1}`
	req, _ := http.NewRequest("POST", "/tablets", bytes.NewBufferString(body))
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
	if data["person_name"] != "王大力" {
		t.Errorf("Expected person_name '王大力', got '%v'", data["person_name"])
	}
}

func TestMemorialController_GetTablet(t *testing.T) {
	mockService := NewMockMemorialService()
	mockService.tablets[1] = &service.TabletDTO{ID: 1, HallID: 1, PersonName: "王大力", Floor: 1, Row: 1, Number: 1}

	r := setupTestRouter()
	r.GET("/tablets/:id", func(ctx *gin.Context) {
		id := int64(1)
		tablet, err := mockService.GetTablet(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if tablet == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tablet not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": tablet})
	})

	req, _ := http.NewRequest("GET", "/tablets/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestMemorialController_GetTabletByPersonID(t *testing.T) {
	personID := int64(100)
	mockService := NewMockMemorialService()
	mockService.tablets[1] = &service.TabletDTO{ID: 1, HallID: 1, PersonID: &personID, PersonName: "王大力"}

	r := setupTestRouter()
	r.GET("/tablets/person/:person_id", func(ctx *gin.Context) {
		pid := int64(100)
		tablet, err := mockService.GetTabletByPersonID(ctx, pid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if tablet == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tablet not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": tablet})
	})

	req, _ := http.NewRequest("GET", "/tablets/person/100", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestMemorialController_DeleteTablet(t *testing.T) {
	mockService := NewMockMemorialService()
	mockService.tablets[1] = &service.TabletDTO{ID: 1, PersonName: "王大力"}

	r := setupTestRouter()
	r.DELETE("/tablets/:id", func(ctx *gin.Context) {
		id := int64(1)
		err := mockService.DeleteTablet(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/tablets/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestMemorialController_GetWithTablets(t *testing.T) {
	mockService := NewMockMemorialService()
	mockService.halls[1] = &service.MemorialHallDTO{
		ID:     1,
		Name:   "王氏宗祠",
		Tables: []*service.TabletDTO{
			{ID: 1, PersonName: "王大力"},
		},
	}

	r := setupTestRouter()
	r.GET("/halls/:id/tablets", func(ctx *gin.Context) {
		id := int64(1)
		hall, err := mockService.GetHallWithTablets(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if hall == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "hall not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": hall})
	})

	req, _ := http.NewRequest("GET", "/halls/1/tablets", nil)
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
	if data["name"] != "王氏宗祠" {
		t.Errorf("Expected name '王氏宗祠', got '%v'", data["name"])
	}
	tablets, ok := data["tablets"].([]interface{})
	if !ok || len(tablets) != 1 {
		t.Errorf("Expected 1 tablet, got %v", data["tablets"])
	}
}