package culture

import (
	"testing"
	"time"
)

func TestDocumentCategoryConstants(t *testing.T) {
	categories := []DocumentCategory{
		CategoryClassic,
		CategoryGenealogy,
		CategoryMemorial,
		CategoryHistory,
	}
	expectedValues := []string{"classic", "genealogy", "memorial", "history"}

	for i, cat := range categories {
		if string(cat) != expectedValues[i] {
			t.Errorf("Expected category %s, got %s", expectedValues[i], cat)
		}
	}
}

func TestDocumentEntity_Validate(t *testing.T) {
	tests := []struct {
		name    string
		doc     Document
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid document",
			doc: Document{
				Title:    "测试文献",
				Content:  "文献内容",
				Category: CategoryClassic,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			doc: Document{
				Title:    "",
				Content:  "文献内容",
				Category: CategoryClassic,
			},
			wantErr: true,
			errMsg:  "title cannot be empty",
		},
		{
			name: "empty content",
			doc: Document{
				Title:    "测试文献",
				Content:  "",
				Category: CategoryClassic,
			},
			wantErr: true,
			errMsg:  "content cannot be empty",
		},
		{
			name: "empty category",
			doc: Document{
				Title:    "测试文献",
				Content:  "文献内容",
				Category: "",
			},
			wantErr: true,
			errMsg:  "category is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestDocument_Getters(t *testing.T) {
	year := 2024
	now := time.Now()
	doc := Document{
		ID:          1,
		Title:       "测试文献",
		Content:     "文献内容",
		Category:    CategoryClassic,
		Author:      "张三",
		CreatedYear: &year,
		Dynasty:     "清朝",
		Source:      "图书馆",
		ImageURLs:   []string{"http://example.com/1.jpg"},
		ViewCount:   100,
		CollectCount: 50,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if doc.ID != 1 {
		t.Errorf("Expected ID 1, got %d", doc.ID)
	}
	if doc.Title != "测试文献" {
		t.Errorf("Expected Title '测试文献', got %s", doc.Title)
	}
	if doc.Author != "张三" {
		t.Errorf("Expected Author '张三', got %s", doc.Author)
	}
	if doc.ViewCount != 100 {
		t.Errorf("Expected ViewCount 100, got %d", doc.ViewCount)
	}
	if doc.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, doc.CreatedAt)
	}
}

func TestStoryEntity_Validate(t *testing.T) {
	tests := []struct {
		name    string
		story   Story
		wantErr bool
		errMsg  string
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
			errMsg:  "title cannot be empty",
		},
		{
			name: "empty content",
			story: Story{
				Title:   "测试故事",
				Content: "",
			},
			wantErr: true,
			errMsg:  "content cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.story.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestStory_Getters(t *testing.T) {
	now := time.Now()
	story := Story{
		ID:        1,
		Title:     "测试故事",
		Content:   "故事内容",
		Era:       "清朝康熙年间",
		Category:  "历史",
		Tags:      []string{"tag1", "tag2"},
		AudioURL:  "http://example.com/audio.mp3",
		ImageURL:  "http://example.com/image.jpg",
		ViewCount: 200,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if story.ID != 1 {
		t.Errorf("Expected ID 1, got %d", story.ID)
	}
	if story.Title != "测试故事" {
		t.Errorf("Expected Title '测试故事', got %s", story.Title)
	}
	if story.Era != "清朝康熙年间" {
		t.Errorf("Expected Era '清朝康熙年间', got %s", story.Era)
	}
	if len(story.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(story.Tags))
	}
	if story.ViewCount != 200 {
		t.Errorf("Expected ViewCount 200, got %d", story.ViewCount)
	}
}

func TestFamilyTeachingsEntity_Validate(t *testing.T) {
	tests := []struct {
		name         string
		familyTeach  FamilyTeachings
		wantErr      bool
		errMsg       string
	}{
		{
			name: "valid family teachings",
			familyTeach: FamilyTeachings{
				Title:   "测试家训",
				Content: "家训内容",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			familyTeach: FamilyTeachings{
				Title:   "",
				Content: "家训内容",
			},
			wantErr: true,
			errMsg:  "title cannot be empty",
		},
		{
			name: "empty content",
			familyTeach: FamilyTeachings{
				Title:   "测试家训",
				Content: "",
			},
			wantErr: true,
			errMsg:  "content cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.familyTeach.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestFamilyTeachings_Getters(t *testing.T) {
	now := time.Now()
	ft := FamilyTeachings{
		ID:         1,
		Title:      "测试家训",
		Content:    "家训内容",
		Generation: 10,
		OriginText: "原文内容",
		Meaning:    "释义内容",
		UsageCount: 50,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if ft.ID != 1 {
		t.Errorf("Expected ID 1, got %d", ft.ID)
	}
	if ft.Title != "测试家训" {
		t.Errorf("Expected Title '测试家训', got %s", ft.Title)
	}
	if ft.Generation != 10 {
		t.Errorf("Expected Generation 10, got %d", ft.Generation)
	}
	if ft.UsageCount != 50 {
		t.Errorf("Expected UsageCount 50, got %d", ft.UsageCount)
	}
}

func TestDocumentCategoryIsValid(t *testing.T) {
	validCategories := []DocumentCategory{
		CategoryClassic,
		CategoryGenealogy,
		CategoryMemorial,
		CategoryHistory,
	}

	for _, cat := range validCategories {
		doc := Document{
			Title:    "Test",
			Content:  "Test",
			Category: cat,
		}
		if err := doc.Validate(); err != nil {
			t.Errorf("Valid category %s failed validation: %v", cat, err)
		}
	}
}

func TestDocument_WithAllFields(t *testing.T) {
	year := 2024
	now := time.Now()
	doc := Document{
		ID:           100,
		Title:        "完整文献",
		Content:      "这是完整内容",
		Category:     CategoryGenealogy,
		Author:       "李四",
		CreatedYear:  &year,
		Dynasty:      "明朝",
		Source:       "国家图书馆",
		ImageURLs:    []string{"img1.jpg", "img2.jpg"},
		ViewCount:    500,
		CollectCount: 100,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := doc.Validate(); err != nil {
		t.Errorf("Document with all fields should be valid: %v", err)
	}
	if doc.ID != 100 {
		t.Errorf("ID should be 100, got %d", doc.ID)
	}
	if doc.Dynasty != "明朝" {
		t.Errorf("Dynasty should be '明朝', got %s", doc.Dynasty)
	}
	if len(doc.ImageURLs) != 2 {
		t.Errorf("ImageURLs should have 2 elements, got %d", len(doc.ImageURLs))
	}
}