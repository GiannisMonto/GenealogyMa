package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/person"
)

// intPtr int指针辅助函数
func intPtr(i int) *int {
	return &i
}

// mockExportPersonRepository 模拟人物仓储
type mockExportPersonRepository struct {
	persons map[int64]*person.Person
	nextID  int64
}

func newMockExportPersonRepository() *mockExportPersonRepository {
	return &mockExportPersonRepository{
		persons: make(map[int64]*person.Person),
		nextID:  1,
	}
}

func (m *mockExportPersonRepository) FindByID(ctx context.Context, id int64) (*person.Person, error) {
	p, ok := m.persons[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockExportPersonRepository) FindByLegacyID(ctx context.Context, legacyID string) (*person.Person, error) {
	for _, p := range m.persons {
		if p.LegacyID == legacyID {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockExportPersonRepository) Create(ctx context.Context, p *person.Person) error {
	p.ID = m.nextID
	m.nextID++
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	m.persons[p.ID] = p
	return nil
}

func (m *mockExportPersonRepository) Update(ctx context.Context, p *person.Person) error {
	p.UpdatedAt = time.Now()
	m.persons[p.ID] = p
	return nil
}

func (m *mockExportPersonRepository) Delete(ctx context.Context, id int64) error {
	delete(m.persons, id)
	return nil
}

func (m *mockExportPersonRepository) Search(ctx context.Context, query *person.SearchQuery) ([]*person.Person, int64, error) {
	var result []*person.Person
	for _, p := range m.persons {
		if query.Name != "" && p.Name != query.Name {
			continue
		}
		if query.Generation != nil && p.Generation != *query.Generation {
			continue
		}
		if query.Gender != nil && p.Gender != *query.Gender {
			continue
		}
		result = append(result, p)
	}
	return result, int64(len(result)), nil
}

func (m *mockExportPersonRepository) FindByGeneration(ctx context.Context, generation int) ([]*person.Person, error) {
	var result []*person.Person
	for _, p := range m.persons {
		if p.Generation == generation {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockExportPersonRepository) FindByName(ctx context.Context, name string, fuzzy bool) ([]*person.Person, error) {
	var result []*person.Person
	for _, p := range m.persons {
		if p.Name == name {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockExportPersonRepository) FindAncestors(ctx context.Context, personID int64, depth int) ([]*person.Person, error) {
	return nil, nil
}

func (m *mockExportPersonRepository) FindDescendants(ctx context.Context, personID int64, depth int) ([]*person.Person, error) {
	return nil, nil
}

func (m *mockExportPersonRepository) FindSiblings(ctx context.Context, personID int64) ([]*person.Person, error) {
	return nil, nil
}

func (m *mockExportPersonRepository) GetFamilyTree(ctx context.Context, rootID int64, depth int) (*person.TreeResult, error) {
	return nil, nil
}

func (m *mockExportPersonRepository) FindChildren(ctx context.Context, parentID int64) ([]*person.PersonChild, error) {
	return nil, nil
}

func (m *mockExportPersonRepository) FindSpouses(ctx context.Context, personID int64) ([]*person.Spouse, error) {
	return nil, nil
}

func (m *mockExportPersonRepository) FindParents(ctx context.Context, personID int64) ([]*person.Person, error) {
	return nil, nil
}

func (m *mockExportPersonRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.persons)), nil
}

func (m *mockExportPersonRepository) CountByGeneration(ctx context.Context) (map[int]int64, error) {
	result := make(map[int]int64)
	for _, p := range m.persons {
		result[p.Generation]++
	}
	return result, nil
}

func (m *mockExportPersonRepository) BatchCreate(ctx context.Context, persons []*person.Person) error {
	for _, p := range persons {
		if err := m.Create(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockExportPersonRepository) BatchUpdate(ctx context.Context, persons []*person.Person) error {
	for _, p := range persons {
		if err := m.Update(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

// newTestExportService 创建测试用导出服务
func newTestExportService() (*ExportService, *mockExportPersonRepository) {
	repo := newMockExportPersonRepository()
	svc := NewExportService(repo)
	return svc, repo
}

func TestExportPersonsToCSV(t *testing.T) {
	svc, repo := newTestExportService()
	ctx := context.Background()

	// 添加测试数据
	testPersons := []*person.Person{
		{
			ID:            1,
			Name:          "张三",
			StyleName:     "子曰",
			Gender:        person.GenderMale,
			Generation:    10,
			BirthPlace:    "四川省成都市",
			BurialPlace:   "四川省成都市",
			BirthTimeText: "1900年",
			DeathTimeText: "1970年",
			DeathYear:     intPtr(1970),
			LineagePath:   "1.2.1",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            2,
			Name:          "李四",
			StyleName:     "子文",
			Gender:        person.GenderFemale,
			Generation:    10,
			BirthPlace:    "四川省重庆市",
			BurialPlace:   "四川省重庆市",
			BirthTimeText: "1905年",
			// DeathYear not set - person is alive
			LineagePath:   "1.2.2",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}
	for _, p := range testPersons {
		repo.Create(ctx, p)
	}

	t.Run("export all persons to CSV", func(t *testing.T) {
		filter := &ExportFilter{
			Entity: ExportEntityPerson,
			Format: ExportFormatCSV,
		}

		result, err := svc.ExportPersonsToCSV(ctx, filter)
		if err != nil {
			t.Fatalf("ExportPersonsToCSV() error = %v", err)
		}
		if result.Count != 2 {
			t.Errorf("expected count 2, got %d", result.Count)
		}
		if result.Format != "csv" {
			t.Errorf("expected format 'csv', got '%s'", result.Format)
		}
		if len(result.Data) == 0 {
			t.Error("expected non-empty data")
		}

		// 验证CSV内容包含表头
		data := string(result.Data)
		if !strings.Contains(data, "ID") || !strings.Contains(data, "姓名") {
			t.Error("CSV should contain header row")
		}
		if !strings.Contains(data, "张三") {
			t.Error("CSV should contain person name")
		}
	})

	t.Run("export persons with generation filter", func(t *testing.T) {
		gen := 10
		filter := &ExportFilter{
			Entity:     ExportEntityPerson,
			Format:     ExportFormatCSV,
			Generation: &gen,
		}

		result, err := svc.ExportPersonsToCSV(ctx, filter)
		if err != nil {
			t.Fatalf("ExportPersonsToCSV() error = %v", err)
		}
		if result.Count != 2 {
			t.Errorf("expected count 2, got %d", result.Count)
		}
	})

	t.Run("export persons with gender filter", func(t *testing.T) {
		filter := &ExportFilter{
			Entity: ExportEntityPerson,
			Format: ExportFormatCSV,
			Gender: "男",
		}

		result, err := svc.ExportPersonsToCSV(ctx, filter)
		if err != nil {
			t.Fatalf("ExportPersonsToCSV() error = %v", err)
		}
		if result.Count != 1 {
			t.Errorf("expected count 1, got %d", result.Count)
		}

		data := string(result.Data)
		if !strings.Contains(data, "张三") {
			t.Error("CSV should contain male person")
		}
		if strings.Contains(data, "李四") {
			t.Error("CSV should not contain female person")
		}
	})

	t.Run("export with name filter", func(t *testing.T) {
		filter := &ExportFilter{
			Entity: ExportEntityPerson,
			Format: ExportFormatCSV,
			Name:   "张三",
		}

		result, err := svc.ExportPersonsToCSV(ctx, filter)
		if err != nil {
			t.Fatalf("ExportPersonsToCSV() error = %v", err)
		}
		if result.Count != 1 {
			t.Errorf("expected count 1, got %d", result.Count)
		}
	})
}

func TestExportPersonsToJSON(t *testing.T) {
	svc, repo := newTestExportService()
	ctx := context.Background()

	// 添加测试数据
	testPerson := &person.Person{
		ID:            1,
		Name:          "王五",
		StyleName:     "子武",
		Gender:        person.GenderMale,
		Generation:    11,
		BirthPlace:    "北京市",
		BurialPlace:   "北京市",
		BirthTimeText: "1920年",
		LineagePath:   "1.3.1",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	repo.Create(ctx, testPerson)

	t.Run("export all persons to JSON", func(t *testing.T) {
		filter := &ExportFilter{
			Entity: ExportEntityPerson,
			Format: ExportFormatJSON,
		}

		result, err := svc.ExportPersonsToJSON(ctx, filter)
		if err != nil {
			t.Fatalf("ExportPersonsToJSON() error = %v", err)
		}
		if result.Count != 1 {
			t.Errorf("expected count 1, got %d", result.Count)
		}
		if result.Format != "json" {
			t.Errorf("expected format 'json', got '%s'", result.Format)
		}

		data := string(result.Data)
		if !strings.Contains(data, "王五") {
			t.Error("JSON should contain person name")
		}
		if !strings.Contains(data, "子武") {
			t.Error("JSON should contain style name")
		}
		if !strings.Contains(data, "1920年") {
			t.Error("JSON should contain birth time")
		}
	})
}

func TestExport(t *testing.T) {
	svc, repo := newTestExportService()
	ctx := context.Background()

	// 添加测试数据
	testPerson := &person.Person{
		ID:            1,
		Name:          "赵六",
		StyleName:     "子禄",
		Gender:        person.GenderMale,
		Generation:    12,
		DeathYear:     intPtr(1980),
		LineagePath:   "1.4.1",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	repo.Create(ctx, testPerson)

	t.Run("export with default format", func(t *testing.T) {
		filter := &ExportFilter{
			Entity: ExportEntityPerson,
			// Format defaults to CSV
		}

		result, err := svc.Export(ctx, filter)
		if err != nil {
			t.Fatalf("Export() error = %v", err)
		}
		if result.Format != "csv" {
			t.Errorf("expected default format 'csv', got '%s'", result.Format)
		}
	})

	t.Run("export with default entity", func(t *testing.T) {
		filter := &ExportFilter{
			Format: ExportFormatCSV,
			// Entity defaults to person
		}

		result, err := svc.Export(ctx, filter)
		if err != nil {
			t.Fatalf("Export() error = %v", err)
		}
		if result.Count != 1 {
			t.Errorf("expected count 1, got %d", result.Count)
		}
	})

	t.Run("export with unsupported entity", func(t *testing.T) {
		filter := &ExportFilter{
			Entity: "unsupported",
			Format: ExportFormatCSV,
		}

		_, err := svc.Export(ctx, filter)
		if err == nil {
			t.Error("expected error for unsupported entity")
		}
	})
}

func TestGetExportableEntities(t *testing.T) {
	svc, _ := newTestExportService()

	entities := svc.GetExportableEntities()
	if len(entities) != 2 {
		t.Errorf("expected 2 entities, got %d", len(entities))
	}

	found := make(map[ExportEntity]bool)
	for _, e := range entities {
		found[e] = true
	}
	if !found[ExportEntityPerson] {
		t.Error("expected person entity")
	}
	if !found[ExportEntityGenealogy] {
		t.Error("expected genealogy entity")
	}
}

func TestGetExportFormats(t *testing.T) {
	svc, _ := newTestExportService()

	formats := svc.GetExportFormats()
	if len(formats) != 2 {
		t.Errorf("expected 2 formats, got %d", len(formats))
	}

	found := make(map[ExportFormat]bool)
	for _, f := range formats {
		found[f] = true
	}
	if !found[ExportFormatCSV] {
		t.Error("expected csv format")
	}
	if !found[ExportFormatJSON] {
		t.Error("expected json format")
	}
}

func TestEscapeJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`test"quote`, `test\"quote`},
		{"line\nbreak", "line\\nbreak"},
		{"tab\there", "tab\\there"},
		{"backslash\\", "backslash\\\\"},
	}

	for _, tt := range tests {
		result := escapeJSON(tt.input)
		if result != tt.expected {
			t.Errorf("escapeJSON(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExportEmptyResult(t *testing.T) {
	svc, _ := newTestExportService()
	ctx := context.Background()

	t.Run("export with no matching data", func(t *testing.T) {
		filter := &ExportFilter{
			Entity: ExportEntityPerson,
			Format: ExportFormatCSV,
			Name:   "不存在的姓名",
		}

		result, err := svc.ExportPersonsToCSV(ctx, filter)
		if err != nil {
			t.Fatalf("ExportPersonsToCSV() error = %v", err)
		}
		if result.Count != 0 {
			t.Errorf("expected count 0, got %d", result.Count)
		}

		// 验证CSV仍包含表头
		data := string(result.Data)
		if !strings.Contains(data, "ID") || !strings.Contains(data, "姓名") {
			t.Error("CSV should still contain header row")
		}
	})
}
