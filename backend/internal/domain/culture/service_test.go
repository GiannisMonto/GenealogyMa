package culture

import (
	"context"
	"testing"
)

// MockDocumentRepository 模拟文献仓储
type MockDocumentRepository struct {
	docs   map[int64]*Document
	nextID int64
}

func NewMockDocumentRepository() *MockDocumentRepository {
	return &MockDocumentRepository{
		docs:   make(map[int64]*Document),
		nextID: 1,
	}
}

func (m *MockDocumentRepository) FindByID(ctx context.Context, id int64) (*Document, error) {
	if d, ok := m.docs[id]; ok {
		return d, nil
	}
	return nil, nil
}

func (m *MockDocumentRepository) Create(ctx context.Context, doc *Document) error {
	doc.ID = m.nextID
	m.nextID++
	m.docs[doc.ID] = doc
	return nil
}

func (m *MockDocumentRepository) Update(ctx context.Context, doc *Document) error {
	m.docs[doc.ID] = doc
	return nil
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id int64) error {
	delete(m.docs, id)
	return nil
}

func (m *MockDocumentRepository) FindAll(ctx context.Context) ([]*Document, error) {
	result := make([]*Document, 0, len(m.docs))
	for _, d := range m.docs {
		result = append(result, d)
	}
	return result, nil
}

func (m *MockDocumentRepository) Search(ctx context.Context, query *SearchQuery) ([]*Document, int64, error) {
	result := make([]*Document, 0, len(m.docs))
	for _, d := range m.docs {
		if query.Keyword != "" && d.Title != query.Keyword {
			continue
		}
		result = append(result, d)
	}
	return result, int64(len(result)), nil
}

func (m *MockDocumentRepository) FindByCategory(ctx context.Context, category DocumentCategory) ([]*Document, error) {
	result := make([]*Document, 0)
	for _, d := range m.docs {
		if d.Category == category {
			result = append(result, d)
		}
	}
	return result, nil
}

// MockStoryRepository 模拟故事仓储
type MockStoryRepository struct {
	stories map[int64]*Story
	nextID  int64
}

func NewMockStoryRepository() *MockStoryRepository {
	return &MockStoryRepository{
		stories: make(map[int64]*Story),
		nextID:  1,
	}
}

func (m *MockStoryRepository) FindByID(ctx context.Context, id int64) (*Story, error) {
	if s, ok := m.stories[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (m *MockStoryRepository) Create(ctx context.Context, story *Story) error {
	story.ID = m.nextID
	m.nextID++
	m.stories[story.ID] = story
	return nil
}

func (m *MockStoryRepository) Update(ctx context.Context, story *Story) error {
	m.stories[story.ID] = story
	return nil
}

func (m *MockStoryRepository) Delete(ctx context.Context, id int64) error {
	delete(m.stories, id)
	return nil
}

func (m *MockStoryRepository) FindAll(ctx context.Context) ([]*Story, error) {
	result := make([]*Story, 0, len(m.stories))
	for _, s := range m.stories {
		result = append(result, s)
	}
	return result, nil
}

func (m *MockStoryRepository) Search(ctx context.Context, query *StorySearchQuery) ([]*Story, int64, error) {
	result := make([]*Story, 0, len(m.stories))
	for _, s := range m.stories {
		if query.Keyword != "" && s.Title != query.Keyword {
			continue
		}
		result = append(result, s)
	}
	return result, int64(len(result)), nil
}

func (m *MockStoryRepository) FindByEra(ctx context.Context, era string) ([]*Story, error) {
	result := make([]*Story, 0)
	for _, s := range m.stories {
		if s.Era == era {
			result = append(result, s)
		}
	}
	return result, nil
}

// MockFamilyTeachingsRepository 模拟家训仓储
type MockFamilyTeachingsRepository struct {
	teachings map[int64]*FamilyTeachings
	nextID    int64
}

func NewMockFamilyTeachingsRepository() *MockFamilyTeachingsRepository {
	return &MockFamilyTeachingsRepository{
		teachings: make(map[int64]*FamilyTeachings),
		nextID:    1,
	}
}

func (m *MockFamilyTeachingsRepository) FindByID(ctx context.Context, id int64) (*FamilyTeachings, error) {
	if ft, ok := m.teachings[id]; ok {
		return ft, nil
	}
	return nil, nil
}

func (m *MockFamilyTeachingsRepository) Create(ctx context.Context, ft *FamilyTeachings) error {
	ft.ID = m.nextID
	m.nextID++
	m.teachings[ft.ID] = ft
	return nil
}

func (m *MockFamilyTeachingsRepository) Update(ctx context.Context, ft *FamilyTeachings) error {
	m.teachings[ft.ID] = ft
	return nil
}

func (m *MockFamilyTeachingsRepository) Delete(ctx context.Context, id int64) error {
	delete(m.teachings, id)
	return nil
}

func (m *MockFamilyTeachingsRepository) FindAll(ctx context.Context) ([]*FamilyTeachings, error) {
	result := make([]*FamilyTeachings, 0, len(m.teachings))
	for _, ft := range m.teachings {
		result = append(result, ft)
	}
	return result, nil
}

func (m *MockFamilyTeachingsRepository) FindByGeneration(ctx context.Context, generation int) ([]*FamilyTeachings, error) {
	result := make([]*FamilyTeachings, 0)
	for _, ft := range m.teachings {
		if ft.Generation == generation {
			result = append(result, ft)
		}
	}
	return result, nil
}

// Test Document entity
func TestDocument_Validate(t *testing.T) {
	tests := []struct {
		name    string
		doc     Document
		wantErr bool
	}{
		{
			name: "valid document",
			doc: Document{
				Title:    "测试文献",
				Content:  "这是内容",
				Category: CategoryClassic,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			doc: Document{
				Title:    "",
				Content:  "这是内容",
				Category: CategoryClassic,
			},
			wantErr: true,
		},
		{
			name: "empty content",
			doc: Document{
				Title:    "测试文献",
				Content:  "",
				Category: CategoryClassic,
			},
			wantErr: true,
		},
		{
			name: "empty category",
			doc: Document{
				Title:    "测试文献",
				Content:  "这是内容",
				Category: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Document.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test Story entity
func TestStory_Validate(t *testing.T) {
	tests := []struct {
		name    string
		story   Story
		wantErr bool
	}{
		{
			name: "valid story",
			story: Story{
				Title:   "测试故事",
				Content: "故事内容",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			story: Story{
				Title:   "",
				Content: "故事内容",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			story: Story{
				Title:   "测试故事",
				Content: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.story.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Story.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test FamilyTeachings entity
func TestFamilyTeachings_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ft      FamilyTeachings
		wantErr bool
	}{
		{
			name: "valid teachings",
			ft: FamilyTeachings{
				Title:   "测试家训",
				Content: "家训内容",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			ft: FamilyTeachings{
				Title:   "",
				Content: "家训内容",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			ft: FamilyTeachings{
				Title:   "测试家训",
				Content: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ft.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FamilyTeachings.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test Service - Document
func TestService_CreateDocument(t *testing.T) {
	repo := NewMockDocumentRepository()
	svc := NewService(repo, nil, nil)

	doc := &Document{
		Title:    "测试文献",
		Content:  "这是内容",
		Category: CategoryClassic,
	}

	err := svc.CreateDocument(context.Background(), doc)
	if err != nil {
		t.Errorf("CreateDocument() error = %v", err)
	}
	if doc.ID == 0 {
		t.Error("CreateDocument() did not set ID")
	}
}

func TestService_GetDocument(t *testing.T) {
	repo := NewMockDocumentRepository()
	svc := NewService(repo, nil, nil)

	doc := &Document{
		Title:    "测试文献",
		Content:  "这是内容",
		Category: CategoryClassic,
	}
	repo.Create(context.Background(), doc)

	got, err := svc.GetDocument(context.Background(), doc.ID)
	if err != nil {
		t.Errorf("GetDocument() error = %v", err)
	}
	if got == nil {
		t.Error("GetDocument() returned nil")
	}
	if got.Title != doc.Title {
		t.Errorf("GetDocument() title = %v, want %v", got.Title, doc.Title)
	}
}

func TestService_UpdateDocument(t *testing.T) {
	repo := NewMockDocumentRepository()
	svc := NewService(repo, nil, nil)

	doc := &Document{
		Title:    "测试文献",
		Content:  "这是内容",
		Category: CategoryClassic,
	}
	repo.Create(context.Background(), doc)

	doc.Title = "更新后的标题"
	err := svc.UpdateDocument(context.Background(), doc)
	if err != nil {
		t.Errorf("UpdateDocument() error = %v", err)
	}

	got, _ := svc.GetDocument(context.Background(), doc.ID)
	if got.Title != "更新后的标题" {
		t.Errorf("UpdateDocument() title = %v, want '更新后的标题'", got.Title)
	}
}

func TestService_DeleteDocument(t *testing.T) {
	repo := NewMockDocumentRepository()
	svc := NewService(repo, nil, nil)

	doc := &Document{
		Title:    "测试文献",
		Content:  "这是内容",
		Category: CategoryClassic,
	}
	repo.Create(context.Background(), doc)

	err := svc.DeleteDocument(context.Background(), doc.ID)
	if err != nil {
		t.Errorf("DeleteDocument() error = %v", err)
	}

	got, _ := svc.GetDocument(context.Background(), doc.ID)
	if got != nil {
		t.Error("DeleteDocument() document still exists")
	}
}

// Test Service - Story
func TestService_CreateStory(t *testing.T) {
	repo := NewMockStoryRepository()
	svc := NewService(nil, repo, nil)

	story := &Story{
		Title:   "测试故事",
		Content: "故事内容",
		Era:     "清朝",
	}

	err := svc.CreateStory(context.Background(), story)
	if err != nil {
		t.Errorf("CreateStory() error = %v", err)
	}
	if story.ID == 0 {
		t.Error("CreateStory() did not set ID")
	}
}

func TestService_ListStoriesByEra(t *testing.T) {
	repo := NewMockStoryRepository()
	svc := NewService(nil, repo, nil)

	repo.Create(context.Background(), &Story{
		Title:   "故事1",
		Content: "内容1",
		Era:     "清朝",
	})
	repo.Create(context.Background(), &Story{
		Title:   "故事2",
		Content: "内容2",
		Era:     "明朝",
	})

	stories, err := svc.ListStoriesByEra(context.Background(), "清朝")
	if err != nil {
		t.Errorf("ListStoriesByEra() error = %v", err)
	}
	if len(stories) != 1 {
		t.Errorf("ListStoriesByEra() len = %v, want 1", len(stories))
	}
}

// Test Service - FamilyTeachings
func TestService_CreateFamilyTeachings(t *testing.T) {
	repo := NewMockFamilyTeachingsRepository()
	svc := NewService(nil, nil, repo)

	ft := &FamilyTeachings{
		Title:      "测试家训",
		Content:   "家训内容",
		Generation: 1,
	}

	err := svc.CreateFamilyTeachings(context.Background(), ft)
	if err != nil {
		t.Errorf("CreateFamilyTeachings() error = %v", err)
	}
	if ft.ID == 0 {
		t.Error("CreateFamilyTeachings() did not set ID")
	}
}

func TestService_ListFamilyTeachingsByGeneration(t *testing.T) {
	repo := NewMockFamilyTeachingsRepository()
	svc := NewService(nil, nil, repo)

	repo.Create(context.Background(), &FamilyTeachings{
		Title:      "家训1",
		Content:   "内容1",
		Generation: 1,
	})
	repo.Create(context.Background(), &FamilyTeachings{
		Title:      "家训2",
		Content:   "内容2",
		Generation: 2,
	})

	teachings, err := svc.ListFamilyTeachingsByGeneration(context.Background(), 1)
	if err != nil {
		t.Errorf("ListFamilyTeachingsByGeneration() error = %v", err)
	}
	if len(teachings) != 1 {
		t.Errorf("ListFamilyTeachingsByGeneration() len = %v, want 1", len(teachings))
	}
}

// Test IncrementViewCount
func TestService_IncrementDocumentViewCount(t *testing.T) {
	repo := NewMockDocumentRepository()
	svc := NewService(repo, nil, nil)

	doc := &Document{
		Title:    "测试文献",
		Content:  "这是内容",
		Category: CategoryClassic,
		ViewCount: 0,
	}
	repo.Create(context.Background(), doc)

	err := svc.IncrementDocumentViewCount(context.Background(), doc.ID)
	if err != nil {
		t.Errorf("IncrementDocumentViewCount() error = %v", err)
	}

	got, _ := svc.GetDocument(context.Background(), doc.ID)
	if got.ViewCount != 1 {
		t.Errorf("IncrementDocumentViewCount() count = %v, want 1", got.ViewCount)
	}
}