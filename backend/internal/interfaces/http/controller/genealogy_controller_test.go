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

// MockGenealogyService 模拟族谱服务
type MockGenealogyService struct {
	genealogies  map[int64]*service.GenealogyDTO
	branches     map[int64]*service.BranchDTO
	generations  map[int64]*service.GenerationDTO
	nextID       int64
}

func NewMockGenealogyService() *MockGenealogyService {
	return &MockGenealogyService{
		genealogies:  make(map[int64]*service.GenealogyDTO),
		branches:    make(map[int64]*service.BranchDTO),
		generations: make(map[int64]*service.GenerationDTO),
		nextID:      1,
	}
}

func (m *MockGenealogyService) ListGenealogies(ctx context.Context) ([]*service.GenealogyDTO, error) {
	result := make([]*service.GenealogyDTO, 0, len(m.genealogies))
	for _, g := range m.genealogies {
		result = append(result, g)
	}
	return result, nil
}

func (m *MockGenealogyService) GetGenealogy(ctx context.Context, id int64) (*service.GenealogyDTO, error) {
	if g, ok := m.genealogies[id]; ok {
		return g, nil
	}
	return nil, nil
}

func (m *MockGenealogyService) GetGenealogyWithDetails(ctx context.Context, id int64) (*service.GenealogyDTO, error) {
	if g, ok := m.genealogies[id]; ok {
		return g, nil
	}
	return nil, nil
}

func (m *MockGenealogyService) CreateGenealogy(ctx context.Context, req *service.CreateGenealogyRequest) (*service.GenealogyDTO, error) {
	g := &service.GenealogyDTO{
		ID:       m.nextID,
		Name:     req.Name,
		Surname:  req.Surname,
	}
	m.nextID++
	m.genealogies[g.ID] = g
	return g, nil
}

func (m *MockGenealogyService) UpdateGenealogy(ctx context.Context, id int64, req *service.UpdateGenealogyRequest) (*service.GenealogyDTO, error) {
	if g, ok := m.genealogies[id]; ok {
		if req.Name != nil {
			g.Name = *req.Name
		}
		return g, nil
	}
	return nil, nil
}

func (m *MockGenealogyService) DeleteGenealogy(ctx context.Context, id int64) error {
	delete(m.genealogies, id)
	return nil
}

func (m *MockGenealogyService) ListBranches(ctx context.Context, genealogyID int64) ([]*service.BranchDTO, error) {
	result := make([]*service.BranchDTO, 0)
	for _, b := range m.branches {
		if b.GenealogyID == genealogyID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (m *MockGenealogyService) GetBranch(ctx context.Context, id int64) (*service.BranchDTO, error) {
	if b, ok := m.branches[id]; ok {
		return b, nil
	}
	return nil, nil
}

func (m *MockGenealogyService) CreateBranch(ctx context.Context, req *service.CreateBranchRequest) (*service.BranchDTO, error) {
	b := &service.BranchDTO{
		ID:          m.nextID,
		GenealogyID: req.GenealogyID,
		Name:        req.Name,
		Code:        req.Code,
	}
	m.nextID++
	m.branches[b.ID] = b
	return b, nil
}

func (m *MockGenealogyService) UpdateBranch(ctx context.Context, id int64, req *service.UpdateBranchRequest) (*service.BranchDTO, error) {
	if b, ok := m.branches[id]; ok {
		if req.Name != nil {
			b.Name = *req.Name
		}
		return b, nil
	}
	return nil, nil
}

func (m *MockGenealogyService) DeleteBranch(ctx context.Context, id int64) error {
	delete(m.branches, id)
	return nil
}

func (m *MockGenealogyService) ListGenerations(ctx context.Context, genealogyID int64) ([]*service.GenerationDTO, error) {
	result := make([]*service.GenerationDTO, 0)
	for _, g := range m.generations {
		if g.GenealogyID == genealogyID {
			result = append(result, g)
		}
	}
	return result, nil
}

func (m *MockGenealogyService) GetGeneration(ctx context.Context, id int64) (*service.GenerationDTO, error) {
	if g, ok := m.generations[id]; ok {
		return g, nil
	}
	return nil, nil
}

func (m *MockGenealogyService) CreateGeneration(ctx context.Context, req *service.CreateGenerationRequest) (*service.GenerationDTO, error) {
	g := &service.GenerationDTO{
		ID:          m.nextID,
		GenealogyID: req.GenealogyID,
		Generation:  req.Generation,
		Name:        req.Name,
	}
	m.nextID++
	m.generations[g.ID] = g
	return g, nil
}

func (m *MockGenealogyService) UpdateGeneration(ctx context.Context, id int64, req *service.UpdateGenerationRequest) (*service.GenerationDTO, error) {
	if g, ok := m.generations[id]; ok {
		if req.Name != nil {
			g.Name = *req.Name
		}
		return g, nil
	}
	return nil, nil
}

func (m *MockGenealogyService) DeleteGeneration(ctx context.Context, id int64) error {
	delete(m.generations, id)
	return nil
}

func (m *MockGenealogyService) GetGenerationByNumber(ctx context.Context, genealogyID int64, generation int) (*service.GenerationDTO, error) {
	for _, g := range m.generations {
		if g.GenealogyID == genealogyID && g.Generation == generation {
			return g, nil
		}
	}
	return nil, nil
}

func TestGenealogyController_List(t *testing.T) {
	mockService := NewMockGenealogyService()
	mockService.genealogies[1] = &service.GenealogyDTO{ID: 1, Name: "王氏族谱", Surname: "王"}
	mockService.genealogies[2] = &service.GenealogyDTO{ID: 2, Name: "李氏族谱", Surname: "李"}

	r := setupTestRouter()
	r.GET("/genealogies", func(ctx *gin.Context) {
		genealogies, err := mockService.ListGenealogies(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": genealogies})
	})

	req, _ := http.NewRequest("GET", "/genealogies", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 genealogies, got %d", len(data))
	}
}

func TestGenealogyController_Get(t *testing.T) {
	mockService := NewMockGenealogyService()
	mockService.genealogies[1] = &service.GenealogyDTO{ID: 1, Name: "王氏族谱", Surname: "王"}

	r := setupTestRouter()
	r.GET("/genealogies/:id", func(ctx *gin.Context) {
		id := int64(1)
		genealogy, err := mockService.GetGenealogy(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if genealogy == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "genealogy not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": genealogy})
	})

	req, _ := http.NewRequest("GET", "/genealogies/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGenealogyController_GetNotFound(t *testing.T) {
	mockService := NewMockGenealogyService()

	r := setupTestRouter()
	r.GET("/genealogies/:id", func(ctx *gin.Context) {
		id := int64(999)
		genealogy, err := mockService.GetGenealogy(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if genealogy == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "genealogy not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": genealogy})
	})

	req, _ := http.NewRequest("GET", "/genealogies/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGenealogyController_Create(t *testing.T) {
	mockService := NewMockGenealogyService()

	r := setupTestRouter()
	r.POST("/genealogies", func(ctx *gin.Context) {
		var req service.CreateGenealogyRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		genealogy, err := mockService.CreateGenealogy(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": genealogy})
	})

	body := `{"name":"王氏族谱","surname":"王","total_generations":22}`
	req, _ := http.NewRequest("POST", "/genealogies", bytes.NewBufferString(body))
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
	if data["name"] != "王氏族谱" {
		t.Errorf("Expected name '王氏族谱', got '%v'", data["name"])
	}
}

func TestGenealogyController_Update(t *testing.T) {
	mockService := NewMockGenealogyService()
	mockService.genealogies[1] = &service.GenealogyDTO{ID: 1, Name: "王氏族谱", Surname: "王"}

	r := setupTestRouter()
	r.PUT("/genealogies/:id", func(ctx *gin.Context) {
		id := int64(1)
		var req service.UpdateGenealogyRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		genealogy, err := mockService.UpdateGenealogy(ctx, id, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if genealogy == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "genealogy not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": genealogy})
	})

	body := `{"name":"王氏族谱v2"}`
	req, _ := http.NewRequest("PUT", "/genealogies/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGenealogyController_Delete(t *testing.T) {
	mockService := NewMockGenealogyService()
	mockService.genealogies[1] = &service.GenealogyDTO{ID: 1, Name: "王氏族谱"}

	r := setupTestRouter()
	r.DELETE("/genealogies/:id", func(ctx *gin.Context) {
		id := int64(1)
		err := mockService.DeleteGenealogy(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/genealogies/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if _, exists := mockService.genealogies[1]; exists {
		t.Errorf("Genealogy should have been deleted")
	}
}

func TestGenealogyController_CreateBranch(t *testing.T) {
	mockService := NewMockGenealogyService()

	r := setupTestRouter()
	r.POST("/branches", func(ctx *gin.Context) {
		var req service.CreateBranchRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		branch, err := mockService.CreateBranch(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": branch})
	})

	body := `{"genealogy_id":1,"name":"长房","code":"A"}`
	req, _ := http.NewRequest("POST", "/branches", bytes.NewBufferString(body))
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
	if data["name"] != "长房" {
		t.Errorf("Expected name '长房', got '%v'", data["name"])
	}
}

func TestGenealogyController_GetBranch(t *testing.T) {
	mockService := NewMockGenealogyService()
	mockService.branches[1] = &service.BranchDTO{ID: 1, GenealogyID: 1, Name: "长房", Code: "A"}

	r := setupTestRouter()
	r.GET("/branches/:id", func(ctx *gin.Context) {
		id := int64(1)
		branch, err := mockService.GetBranch(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if branch == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "branch not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": branch})
	})

	req, _ := http.NewRequest("GET", "/branches/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGenealogyController_DeleteBranch(t *testing.T) {
	mockService := NewMockGenealogyService()
	mockService.branches[1] = &service.BranchDTO{ID: 1, Name: "长房"}

	r := setupTestRouter()
	r.DELETE("/branches/:id", func(ctx *gin.Context) {
		id := int64(1)
		err := mockService.DeleteBranch(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/branches/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGenealogyController_CreateGeneration(t *testing.T) {
	mockService := NewMockGenealogyService()

	r := setupTestRouter()
	r.POST("/generations", func(ctx *gin.Context) {
		var req service.CreateGenerationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		generation, err := mockService.CreateGeneration(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": generation})
	})

	body := `{"genealogy_id":1,"generation":1,"name":"道","sequence":1}`
	req, _ := http.NewRequest("POST", "/generations", bytes.NewBufferString(body))
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
	if data["name"] != "道" {
		t.Errorf("Expected name '道', got '%v'", data["name"])
	}
}

func TestGenealogyController_GetGeneration(t *testing.T) {
	mockService := NewMockGenealogyService()
	mockService.generations[1] = &service.GenerationDTO{ID: 1, GenealogyID: 1, Generation: 1, Name: "道"}

	r := setupTestRouter()
	r.GET("/generations/:id", func(ctx *gin.Context) {
		id := int64(1)
		generation, err := mockService.GetGeneration(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if generation == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "generation not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": generation})
	})

	req, _ := http.NewRequest("GET", "/generations/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGenealogyController_GetGenerationByNumber(t *testing.T) {
	mockService := NewMockGenealogyService()
	mockService.generations[1] = &service.GenerationDTO{ID: 1, GenealogyID: 1, Generation: 1, Name: "道"}

	r := setupTestRouter()
	r.GET("/generations/by-number", func(ctx *gin.Context) {
		genealogyID := int64(1)
		generationNum := 1
		generation, err := mockService.GetGenerationByNumber(ctx, genealogyID, generationNum)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if generation == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "generation not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": generation})
	})

	req, _ := http.NewRequest("GET", "/generations/by-number?genealogy_id=1&generation=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGenealogyController_DeleteGeneration(t *testing.T) {
	mockService := NewMockGenealogyService()
	mockService.generations[1] = &service.GenerationDTO{ID: 1, Name: "道"}

	r := setupTestRouter()
	r.DELETE("/generations/:id", func(ctx *gin.Context) {
		id := int64(1)
		err := mockService.DeleteGeneration(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/generations/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGenealogyController_GetWithDetails(t *testing.T) {
	mockService := NewMockGenealogyService()
	mockService.genealogies[1] = &service.GenealogyDTO{
		ID:     1,
		Name:   "王氏族谱",
		Branches: []*service.BranchDTO{
			{ID: 1, Name: "长房"},
		},
	}

	r := setupTestRouter()
	r.GET("/genealogies/:id/details", func(ctx *gin.Context) {
		id := int64(1)
		genealogy, err := mockService.GetGenealogyWithDetails(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if genealogy == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "genealogy not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": genealogy})
	})

	req, _ := http.NewRequest("GET", "/genealogies/1/details", nil)
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
	if data["name"] != "王氏族谱" {
		t.Errorf("Expected name '王氏族谱', got '%v'", data["name"])
	}
	branches, ok := data["branches"].([]interface{})
	if !ok || len(branches) != 1 {
		t.Errorf("Expected 1 branch, got %v", data["branches"])
	}
}