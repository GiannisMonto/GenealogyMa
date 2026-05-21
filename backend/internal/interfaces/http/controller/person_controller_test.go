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
	"github.com/genealogy-ma/platform/internal/domain/person"
)

// MockPersonService 模拟人物服务
type MockPersonService struct {
	persons  map[int64]*service.PersonDTO
	nextID   int64
	trees    map[int64]*service.FamilyTreeResponse
	stats    *service.StatisticsResponse
}

func NewMockPersonService() *MockPersonService {
	return &MockPersonService{
		persons:  make(map[int64]*service.PersonDTO),
		nextID:   1,
		trees:    make(map[int64]*service.FamilyTreeResponse),
		stats:    &service.StatisticsResponse{TotalPersons: 0, ByGeneration: make(map[int]int64)},
	}
}

func (m *MockPersonService) GetPerson(ctx context.Context, id int64, withRelations bool) (*service.PersonDTO, error) {
	if p, ok := m.persons[id]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *MockPersonService) SearchPersons(ctx context.Context, req *service.SearchPersonRequest) ([]*service.PersonDTO, int64, error) {
	result := make([]*service.PersonDTO, 0, len(m.persons))
	for _, p := range m.persons {
		result = append(result, p)
	}
	return result, int64(len(result)), nil
}

func (m *MockPersonService) CreatePerson(ctx context.Context, req *service.CreatePersonRequest) (*service.PersonDTO, error) {
	p := &service.PersonDTO{
		ID:   m.nextID,
		Name: req.Name,
	}
	m.nextID++
	m.persons[p.ID] = p
	return p, nil
}

func (m *MockPersonService) UpdatePerson(ctx context.Context, id int64, req *service.UpdatePersonRequest) (*service.PersonDTO, error) {
	if p, ok := m.persons[id]; ok {
		p.Name = req.Name
		return p, nil
	}
	return nil, nil
}

func (m *MockPersonService) DeletePerson(ctx context.Context, id int64) error {
	delete(m.persons, id)
	return nil
}

func (m *MockPersonService) GetFamilyTree(ctx context.Context, personID int64, upDepth, downDepth int) (*service.FamilyTreeResponse, error) {
	if t, ok := m.trees[personID]; ok {
		return t, nil
	}
	return &service.FamilyTreeResponse{}, nil
}

func (m *MockPersonService) GetStatistics(ctx context.Context) (*service.StatisticsResponse, error) {
	return m.stats, nil
}

func (m *MockPersonService) BatchCreatePersons(ctx context.Context, req *service.BatchCreatePersonRequest) (*service.BatchResult, error) {
	result := &service.BatchResult{
		Results: make([]*service.BatchItemResult, 0, len(req.Persons)),
	}
	for _, p := range req.Persons {
		id := m.nextID
		m.nextID++
		m.persons[id] = &service.PersonDTO{ID: id, Name: p.Name, Generation: p.Generation}
		result.SuccessCount++
		result.Results = append(result.Results, &service.BatchItemResult{ID: id, Success: true})
	}
	return result, nil
}

func (m *MockPersonService) BatchUpdatePersons(ctx context.Context, req *service.BatchUpdatePersonRequest) (*service.BatchResult, error) {
	result := &service.BatchResult{
		Results: make([]*service.BatchItemResult, 0, len(req.Persons)),
	}
	for _, item := range req.Persons {
		if p, ok := m.persons[item.ID]; ok {
			if item.Name != "" {
				p.Name = item.Name
			}
			result.SuccessCount++
			result.Results = append(result.Results, &service.BatchItemResult{ID: item.ID, Success: true})
		} else {
			result.FailCount++
			result.Results = append(result.Results, &service.BatchItemResult{ID: item.ID, Success: false, Error: "not found"})
		}
	}
	return result, nil
}

func (m *MockPersonService) BatchDeletePersons(ctx context.Context, req *service.BatchDeleteRequest) (*service.BatchResult, error) {
	result := &service.BatchResult{
		Results: make([]*service.BatchItemResult, 0, len(req.IDs)),
	}
	for _, id := range req.IDs {
		if _, ok := m.persons[id]; ok {
			delete(m.persons, id)
			result.SuccessCount++
			result.Results = append(result.Results, &service.BatchItemResult{ID: id, Success: true})
		} else {
			result.FailCount++
			result.Results = append(result.Results, &service.BatchItemResult{ID: id, Success: false, Error: "not found"})
		}
	}
	return result, nil
}

func setupPersonTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestPersonController_Get(t *testing.T) {
	mockService := NewMockPersonService()
	mockService.persons[1] = &service.PersonDTO{ID: 1, Name: "张三", Generation: 1}

	r := setupPersonTestRouter()
	r.GET("/persons/:id", func(ctx *gin.Context) {
		p, _ := mockService.GetPerson(ctx, 1, false)
		if p == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": p})
	})

	req, _ := http.NewRequest("GET", "/persons/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["data"] == nil {
		t.Error("Expected data in response")
	}
}

func TestPersonController_List(t *testing.T) {
	mockService := NewMockPersonService()
	mockService.persons[1] = &service.PersonDTO{ID: 1, Name: "张三", Generation: 1}
	mockService.persons[2] = &service.PersonDTO{ID: 2, Name: "李四", Generation: 2}

	r := setupPersonTestRouter()
	r.GET("/persons", func(ctx *gin.Context) {
		persons, total, _ := mockService.SearchPersons(ctx, &service.SearchPersonRequest{})
		ctx.JSON(http.StatusOK, gin.H{"data": persons, "total": total})
	})

	req, _ := http.NewRequest("GET", "/persons", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 persons, got %v", response["data"])
	}
}

func TestPersonController_List_WithPagination(t *testing.T) {
	mockService := NewMockPersonService()
	for i := int64(1); i <= 25; i++ {
		mockService.persons[i] = &service.PersonDTO{ID: i, Name: "张三", Generation: 1}
	}

	r := setupPersonTestRouter()
	r.GET("/persons", func(ctx *gin.Context) {
		persons, total, _ := mockService.SearchPersons(ctx, &service.SearchPersonRequest{Page: 1, PageSize: 10})
		_ = ctx.Query("page_size")
		ctx.JSON(http.StatusOK, gin.H{"data": persons, "total": total})
	})

	req, _ := http.NewRequest("GET", "/persons?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	total, ok := response["total"].(float64)
	if !ok || total != 25 {
		t.Errorf("Expected total 25, got %v", response["total"])
	}
}

func TestPersonController_Create(t *testing.T) {
	mockService := NewMockPersonService()

	r := setupPersonTestRouter()
	r.POST("/persons", func(ctx *gin.Context) {
		var req service.CreatePersonRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p, _ := mockService.CreatePerson(ctx, &req)
		ctx.JSON(http.StatusOK, gin.H{"data": p})
	})

	body := bytes.NewBufferString(`{"name":"王五","gender":"男","generation":1}`)
	req, _ := http.NewRequest("POST", "/persons", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok || data["name"] != "王五" {
		t.Errorf("Expected name '王五', got %v", response["data"])
	}
}

func TestPersonController_Create_ValidationError(t *testing.T) {
	r := setupPersonTestRouter()
	r.POST("/persons", func(ctx *gin.Context) {
		var req service.CreatePersonRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	})

	body := bytes.NewBufferString(`{"name":"","gender":"male"}`)
	req, _ := http.NewRequest("POST", "/persons", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPersonController_Update(t *testing.T) {
	mockService := NewMockPersonService()
	mockService.persons[1] = &service.PersonDTO{ID: 1, Name: "张三", Generation: 1}

	r := setupPersonTestRouter()
	r.PUT("/persons/:id", func(ctx *gin.Context) {
		var req service.UpdatePersonRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p, _ := mockService.UpdatePerson(ctx, 1, &req)
		if p == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": p})
	})

	body := bytes.NewBufferString(`{"name":"新名字"}`)
	req, _ := http.NewRequest("PUT", "/persons/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok || data["name"] != "新名字" {
		t.Errorf("Expected name '新名字', got %v", response["data"])
	}
}

func TestPersonController_Delete(t *testing.T) {
	mockService := NewMockPersonService()
	mockService.persons[1] = &service.PersonDTO{ID: 1, Name: "张三", Generation: 1}

	r := setupPersonTestRouter()
	r.DELETE("/persons/:id", func(ctx *gin.Context) {
		_ = mockService.DeletePerson(ctx, 1)
		ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
	})

	req, _ := http.NewRequest("DELETE", "/persons/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if mockService.persons[1] != nil {
		t.Error("Expected person to be deleted")
	}
}

func TestPersonController_GetFamilyTree(t *testing.T) {
	mockService := NewMockPersonService()
	mockService.trees[1] = &service.FamilyTreeResponse{
		CenterPerson: &service.PersonDTO{ID: 1, Name: "张三"},
		Ancestors:    []*service.PersonDTO{},
		Descendants:  []*service.PersonDTO{},
		Generations:  1,
		TotalNodes:   1,
	}

	r := setupPersonTestRouter()
	r.GET("/persons/:id/tree", func(ctx *gin.Context) {
		tree, _ := mockService.GetFamilyTree(ctx, 1, 5, 5)
		ctx.JSON(http.StatusOK, gin.H{"data": tree})
	})

	req, _ := http.NewRequest("GET", "/persons/1/tree", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["data"] == nil {
		t.Error("Expected tree data in response")
	}
}

func TestPersonController_GetStatistics(t *testing.T) {
	mockService := NewMockPersonService()
	mockService.stats = &service.StatisticsResponse{
		TotalPersons:  100,
		ByGeneration: map[int]int64{1: 10, 2: 20, 3: 30},
	}

	r := setupPersonTestRouter()
	r.GET("/persons/statistics", func(ctx *gin.Context) {
		stats, _ := mockService.GetStatistics(ctx)
		ctx.JSON(http.StatusOK, gin.H{"data": stats})
	})

	req, _ := http.NewRequest("GET", "/persons/statistics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected data to be map, got %v", response["data"])
		return
	}
	if data["total_persons"] != float64(100) {
		t.Errorf("Expected total_persons 100, got %v", data["total_persons"])
	}
}

func TestPersonController_Get_NotFound(t *testing.T) {
	r := setupPersonTestRouter()
	r.GET("/persons/:id", func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	req, _ := http.NewRequest("GET", "/persons/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestPersonController_InvalidID(t *testing.T) {
	r := setupPersonTestRouter()
	r.GET("/persons/:id", func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
	})

	req, _ := http.NewRequest("GET", "/persons/invalid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPersonEntity_Validation(t *testing.T) {
	tests := []struct {
		name    string
		person  *person.Person
		wantErr bool
	}{
		{
			name: "valid person",
			person: &person.Person{
				Name:       "测试",
				Gender:     person.GenderMale,
				Generation: 1,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			person: &person.Person{
				Name:       "",
				Gender:     person.GenderMale,
				Generation: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.person.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Person.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPersonController_BatchCreate(t *testing.T) {
	mockService := NewMockPersonService()

	r := setupPersonTestRouter()
	r.POST("/persons/batch", func(ctx *gin.Context) {
		var req service.BatchCreatePersonRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, _ := mockService.BatchCreatePersons(ctx, &req)
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	body := bytes.NewBufferString(`{"persons":[{"name":"张三","gender":"男","generation":1},{"name":"李四","gender":"女","generation":2}]}`)
	req, _ := http.NewRequest("POST", "/persons/batch", body)
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
		return
	}
	if data["success_count"] != float64(2) {
		t.Errorf("Expected success_count 2, got %v", data["success_count"])
	}
}

func TestPersonController_BatchCreate_EmptyList(t *testing.T) {
	r := setupPersonTestRouter()
	r.POST("/persons/batch", func(ctx *gin.Context) {
		var req service.BatchCreatePersonRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(req.Persons) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty list"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	})

	body := bytes.NewBufferString(`{"persons":[]}`)
	req, _ := http.NewRequest("POST", "/persons/batch", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPersonController_BatchUpdate(t *testing.T) {
	mockService := NewMockPersonService()
	mockService.persons[1] = &service.PersonDTO{ID: 1, Name: "张三", Generation: 1}
	mockService.persons[2] = &service.PersonDTO{ID: 2, Name: "李四", Generation: 2}

	r := setupPersonTestRouter()
	r.PUT("/persons/batch", func(ctx *gin.Context) {
		var req service.BatchUpdatePersonRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, _ := mockService.BatchUpdatePersons(ctx, &req)
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	body := bytes.NewBufferString(`{"persons":[{"id":1,"name":"新张三"},{"id":2,"name":"新李四"}]}`)
	req, _ := http.NewRequest("PUT", "/persons/batch", body)
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
		return
	}
	if data["success_count"] != float64(2) {
		t.Errorf("Expected success_count 2, got %v", data["success_count"])
	}
}

func TestPersonController_BatchUpdate_PartialFailure(t *testing.T) {
	mockService := NewMockPersonService()
	mockService.persons[1] = &service.PersonDTO{ID: 1, Name: "张三", Generation: 1}

	r := setupPersonTestRouter()
	r.PUT("/persons/batch", func(ctx *gin.Context) {
		var req service.BatchUpdatePersonRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, _ := mockService.BatchUpdatePersons(ctx, &req)
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	body := bytes.NewBufferString(`{"persons":[{"id":1,"name":"新张三"},{"id":999,"name":"不存在"}]}`)
	req, _ := http.NewRequest("PUT", "/persons/batch", body)
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
		return
	}
	if data["success_count"] != float64(1) {
		t.Errorf("Expected success_count 1, got %v", data["success_count"])
	}
	if data["fail_count"] != float64(1) {
		t.Errorf("Expected fail_count 1, got %v", data["fail_count"])
	}
}

func TestPersonController_BatchDelete(t *testing.T) {
	mockService := NewMockPersonService()
	mockService.persons[1] = &service.PersonDTO{ID: 1, Name: "张三", Generation: 1}
	mockService.persons[2] = &service.PersonDTO{ID: 2, Name: "李四", Generation: 2}

	r := setupPersonTestRouter()
	r.DELETE("/persons/batch", func(ctx *gin.Context) {
		var req service.BatchDeleteRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, _ := mockService.BatchDeletePersons(ctx, &req)
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	body := bytes.NewBufferString(`{"ids":[1,2]}`)
	req, _ := http.NewRequest("DELETE", "/persons/batch", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if mockService.persons[1] != nil {
		t.Error("Expected person 1 to be deleted")
	}
	if mockService.persons[2] != nil {
		t.Error("Expected person 2 to be deleted")
	}
}

func TestPersonController_BatchDelete_PartialFailure(t *testing.T) {
	mockService := NewMockPersonService()
	mockService.persons[1] = &service.PersonDTO{ID: 1, Name: "张三", Generation: 1}

	r := setupPersonTestRouter()
	r.DELETE("/persons/batch", func(ctx *gin.Context) {
		var req service.BatchDeleteRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, _ := mockService.BatchDeletePersons(ctx, &req)
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	body := bytes.NewBufferString(`{"ids":[1,999]}`)
	req, _ := http.NewRequest("DELETE", "/persons/batch", body)
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
		return
	}
	if data["success_count"] != float64(1) {
		t.Errorf("Expected success_count 1, got %v", data["success_count"])
	}
	if data["fail_count"] != float64(1) {
		t.Errorf("Expected fail_count 1, got %v", data["fail_count"])
	}
}

func TestPersonController_BatchDelete_EmptyList(t *testing.T) {
	r := setupPersonTestRouter()
	r.DELETE("/persons/batch", func(ctx *gin.Context) {
		var req service.BatchDeleteRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(req.IDs) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty list"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	})

	body := bytes.NewBufferString(`{"ids":[]}`)
	req, _ := http.NewRequest("DELETE", "/persons/batch", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}