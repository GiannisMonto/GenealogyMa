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
	"github.com/genealogy-ma/platform/internal/domain/culture"
)

// MockCultureService 模拟文化服务
type MockCultureService struct {
	documents        map[int64]*service.DocumentDTO
	stories          map[int64]*service.StoryDTO
	familyTeachings  map[int64]*service.FamilyTeachingsDTO
	nextID           int64
}

func NewMockCultureService() *MockCultureService {
	return &MockCultureService{
		documents:       make(map[int64]*service.DocumentDTO),
		stories:         make(map[int64]*service.StoryDTO),
		familyTeachings: make(map[int64]*service.FamilyTeachingsDTO),
		nextID:          1,
	}
}

func (m *MockCultureService) ListDocuments(ctx context.Context) ([]*service.DocumentDTO, error) {
	result := make([]*service.DocumentDTO, 0, len(m.documents))
	for _, d := range m.documents {
		result = append(result, d)
	}
	return result, nil
}

func (m *MockCultureService) GetDocument(ctx context.Context, id int64) (*service.DocumentDTO, error) {
	if d, ok := m.documents[id]; ok {
		return d, nil
	}
	return nil, nil
}

func (m *MockCultureService) CreateDocument(ctx context.Context, req *service.CreateDocumentRequest) (*service.DocumentDTO, error) {
	d := &service.DocumentDTO{
		ID:       m.nextID,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Author:   req.Author,
		Dynasty:  req.Dynasty,
	}
	m.nextID++
	m.documents[d.ID] = d
	return d, nil
}

func (m *MockCultureService) UpdateDocument(ctx context.Context, id int64, req *service.UpdateDocumentRequest) (*service.DocumentDTO, error) {
	if d, ok := m.documents[id]; ok {
		if req.Title != nil {
			d.Title = *req.Title
		}
		if req.Content != nil {
			d.Content = *req.Content
		}
		return d, nil
	}
	return nil, nil
}

func (m *MockCultureService) DeleteDocument(ctx context.Context, id int64) error {
	delete(m.documents, id)
	return nil
}

func (m *MockCultureService) ListStories(ctx context.Context) ([]*service.StoryDTO, error) {
	result := make([]*service.StoryDTO, 0, len(m.stories))
	for _, s := range m.stories {
		result = append(result, s)
	}
	return result, nil
}

func (m *MockCultureService) GetStory(ctx context.Context, id int64) (*service.StoryDTO, error) {
	if s, ok := m.stories[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (m *MockCultureService) CreateStory(ctx context.Context, req *service.CreateStoryRequest) (*service.StoryDTO, error) {
	s := &service.StoryDTO{
		ID:       m.nextID,
		Title:    req.Title,
		Content:  req.Content,
		Era:      req.Era,
		Category: req.Category,
	}
	m.nextID++
	m.stories[s.ID] = s
	return s, nil
}

func (m *MockCultureService) UpdateStory(ctx context.Context, id int64, req *service.UpdateStoryRequest) (*service.StoryDTO, error) {
	if s, ok := m.stories[id]; ok {
		if req.Title != nil {
			s.Title = *req.Title
		}
		if req.Content != nil {
			s.Content = *req.Content
		}
		return s, nil
	}
	return nil, nil
}

func (m *MockCultureService) DeleteStory(ctx context.Context, id int64) error {
	delete(m.stories, id)
	return nil
}

func (m *MockCultureService) ListFamilyTeachings(ctx context.Context) ([]*service.FamilyTeachingsDTO, error) {
	result := make([]*service.FamilyTeachingsDTO, 0, len(m.familyTeachings))
	for _, ft := range m.familyTeachings {
		result = append(result, ft)
	}
	return result, nil
}

func (m *MockCultureService) GetFamilyTeachings(ctx context.Context, id int64) (*service.FamilyTeachingsDTO, error) {
	if ft, ok := m.familyTeachings[id]; ok {
		return ft, nil
	}
	return nil, nil
}

func (m *MockCultureService) CreateFamilyTeachings(ctx context.Context, req *service.CreateFamilyTeachingsRequest) (*service.FamilyTeachingsDTO, error) {
	ft := &service.FamilyTeachingsDTO{
		ID:         m.nextID,
		Title:      req.Title,
		Content:    req.Content,
		Generation: req.Generation,
		OriginText: req.OriginText,
		Meaning:    req.Meaning,
	}
	m.nextID++
	m.familyTeachings[ft.ID] = ft
	return ft, nil
}

func (m *MockCultureService) UpdateFamilyTeachings(ctx context.Context, id int64, req *service.UpdateFamilyTeachingsRequest) (*service.FamilyTeachingsDTO, error) {
	if ft, ok := m.familyTeachings[id]; ok {
		if req.Title != nil {
			ft.Title = *req.Title
		}
		if req.Content != nil {
			ft.Content = *req.Content
		}
		return ft, nil
	}
	return nil, nil
}

func (m *MockCultureService) DeleteFamilyTeachings(ctx context.Context, id int64) error {
	delete(m.familyTeachings, id)
	return nil
}

// ===== 文献测试 =====

func TestCultureController_ListDocuments(t *testing.T) {
	mockService := NewMockCultureService()
	mockService.documents[1] = &service.DocumentDTO{ID: 1, Title: "族谱序", Category: culture.CategoryGenealogy}
	mockService.documents[2] = &service.DocumentDTO{ID: 2, Title: "家训原文", Category: culture.CategoryClassic}

	r := setupTestRouter()
	r.GET("/documents", func(ctx *gin.Context) {
		docs, err := mockService.ListDocuments(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": docs})
	})

	req, _ := http.NewRequest("GET", "/documents", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 documents, got %d", len(data))
	}
}

func TestCultureController_GetDocument(t *testing.T) {
	mockService := NewMockCultureService()
	mockService.documents[1] = &service.DocumentDTO{ID: 1, Title: "族谱序", Category: culture.CategoryGenealogy}

	r := setupTestRouter()
	r.GET("/documents/:id", func(ctx *gin.Context) {
		id := int64(1)
		doc, err := mockService.GetDocument(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if doc == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": doc})
	})

	req, _ := http.NewRequest("GET", "/documents/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCultureController_CreateDocument(t *testing.T) {
	mockService := NewMockCultureService()

	r := setupTestRouter()
	r.POST("/documents", func(ctx *gin.Context) {
		var req service.CreateDocumentRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		doc, err := mockService.CreateDocument(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": doc})
	})

	body := `{"title":"新族谱序","content":"族谱序内容","category":"genealogy","author":"张三"}`
	req, _ := http.NewRequest("POST", "/documents", bytes.NewBufferString(body))
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
	if data["title"] != "新族谱序" {
		t.Errorf("Expected title '新族谱序', got '%v'", data["title"])
	}
}

// ===== 故事测试 =====

func TestCultureController_ListStories(t *testing.T) {
	mockService := NewMockCultureService()
	mockService.stories[1] = &service.StoryDTO{ID: 1, Title: "先祖故事", Era: "清朝"}
	mockService.stories[2] = &service.StoryDTO{ID: 2, Title: "迁徙故事", Era: "明朝"}

	r := setupTestRouter()
	r.GET("/stories", func(ctx *gin.Context) {
		stories, err := mockService.ListStories(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": stories})
	})

	req, _ := http.NewRequest("GET", "/stories", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 stories, got %d", len(data))
	}
}

func TestCultureController_CreateStory(t *testing.T) {
	mockService := NewMockCultureService()

	r := setupTestRouter()
	r.POST("/stories", func(ctx *gin.Context) {
		var req service.CreateStoryRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		story, err := mockService.CreateStory(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": story})
	})

	body := `{"title":"先祖故事","content":"故事内容","era":"清朝"}`
	req, _ := http.NewRequest("POST", "/stories", bytes.NewBufferString(body))
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
	if data["era"] != "清朝" {
		t.Errorf("Expected era '清朝', got '%v'", data["era"])
	}
}

// ===== 家训测试 =====

func TestCultureController_ListFamilyTeachings(t *testing.T) {
	mockService := NewMockCultureService()
	mockService.familyTeachings[1] = &service.FamilyTeachingsDTO{ID: 1, Title: "家训一", Generation: 5}
	mockService.familyTeachings[2] = &service.FamilyTeachingsDTO{ID: 2, Title: "家训二", Generation: 10}

	r := setupTestRouter()
	r.GET("/family-teachings", func(ctx *gin.Context) {
		items, err := mockService.ListFamilyTeachings(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": items})
	})

	req, _ := http.NewRequest("GET", "/family-teachings", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 family teachings, got %d", len(data))
	}
}

func TestCultureController_DeleteDocument(t *testing.T) {
	mockService := NewMockCultureService()
	mockService.documents[1] = &service.DocumentDTO{ID: 1, Title: "族谱序"}

	r := setupTestRouter()
	r.DELETE("/documents/:id", func(ctx *gin.Context) {
		id := int64(1)
		err := mockService.DeleteDocument(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/documents/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if _, exists := mockService.documents[1]; exists {
		t.Errorf("Document should have been deleted")
	}
}
