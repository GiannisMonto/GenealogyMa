package person

import (
	"context"
	"testing"
	"time"
)

// mockRepository 模拟人物仓储
type mockRepository struct {
	persons   map[int64]*Person
	spouses   map[int64][]*Spouse
	children  map[int64][]*PersonChild
	nextID    int64
}

func newMockPersonRepository() *mockRepository {
	return &mockRepository{
		persons:  make(map[int64]*Person),
		spouses:  make(map[int64][]*Spouse),
		children: make(map[int64][]*PersonChild),
		nextID:   1,
	}
}

func (m *mockRepository) FindByID(ctx context.Context, id int64) (*Person, error) {
	p, ok := m.persons[id]
	if !ok {
		return nil, nil
	}
	// 深拷贝以避免修改原始数据
	result := *p
	result.Spouses = m.spouses[id]
	result.Children = m.children[id]
	return &result, nil
}

func (m *mockRepository) FindByLegacyID(ctx context.Context, legacyID string) (*Person, error) {
	for _, p := range m.persons {
		if p.LegacyID == legacyID {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) Create(ctx context.Context, p *Person) error {
	p.ID = m.nextID
	m.nextID++
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	m.persons[p.ID] = p
	return nil
}

func (m *mockRepository) Update(ctx context.Context, p *Person) error {
	p.UpdatedAt = time.Now()
	// 深拷贝以避免指针问题（模拟真实数据库行为）
	copy := *p
	m.persons[p.ID] = &copy
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id int64) error {
	delete(m.persons, id)
	return nil
}

func (m *mockRepository) Search(ctx context.Context, query *SearchQuery) ([]*Person, int64, error) {
	var result []*Person
	for _, p := range m.persons {
		result = append(result, p)
	}
	return result, int64(len(result)), nil
}

func (m *mockRepository) FindByGeneration(ctx context.Context, generation int) ([]*Person, error) {
	var result []*Person
	for _, p := range m.persons {
		if p.Generation == generation {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockRepository) FindByName(ctx context.Context, name string, fuzzy bool) ([]*Person, error) {
	var result []*Person
	for _, p := range m.persons {
		if p.Name == name {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockRepository) FindAncestors(ctx context.Context, personID int64, depth int) ([]*Person, error) {
	return nil, nil
}

func (m *mockRepository) FindDescendants(ctx context.Context, personID int64, depth int) ([]*Person, error) {
	return nil, nil
}

func (m *mockRepository) FindSiblings(ctx context.Context, personID int64) ([]*Person, error) {
	return nil, nil
}

func (m *mockRepository) GetFamilyTree(ctx context.Context, rootID int64, depth int) (*TreeResult, error) {
	return nil, nil
}

func (m *mockRepository) FindChildren(ctx context.Context, parentID int64) ([]*PersonChild, error) {
	children := m.children[parentID]
	if children == nil {
		return nil, nil
	}
	result := make([]*PersonChild, len(children))
	copy(result, children)
	return result, nil
}

func (m *mockRepository) FindSpouses(ctx context.Context, personID int64) ([]*Spouse, error) {
	spouses := m.spouses[personID]
	if spouses == nil {
		return nil, nil
	}
	result := make([]*Spouse, len(spouses))
	copy(result, spouses)
	return result, nil
}

func (m *mockRepository) FindParents(ctx context.Context, personID int64) ([]*Person, error) {
	return nil, nil
}

func (m *mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.persons)), nil
}

func (m *mockRepository) CountByGeneration(ctx context.Context) (map[int]int64, error) {
	result := make(map[int]int64)
	for _, p := range m.persons {
		result[p.Generation]++
	}
	return result, nil
}

func (m *mockRepository) BatchCreate(ctx context.Context, persons []*Person) error {
	return nil
}

func (m *mockRepository) BatchUpdate(ctx context.Context, persons []*Person) error {
	return nil
}

// ===== 测试用例 =====

func TestPersonValidate(t *testing.T) {
	tests := []struct {
		name    string
		person  *Person
		wantErr bool
	}{
		{
			name: "valid person",
			person: &Person{
				Name:       "测试",
				Gender:     GenderMale,
				Generation: 1,
				SonCount:   1,
				DaughterCount: 0,
				AdoptedHeirCount: 0,
				TotalChildren: 1,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			person: &Person{
				Name:       "",
				Gender:     GenderMale,
				Generation: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid gender",
			person: &Person{
				Name:       "测试",
				Gender:     "unknown",
				Generation: 1,
			},
			wantErr: true,
		},
		{
			name: "generation out of range",
			person: &Person{
				Name:       "测试",
				Gender:     GenderMale,
				Generation: 31,
			},
			wantErr: true,
		},
		{
			name: "children count mismatch",
			person: &Person{
				Name:             "测试",
				Gender:           GenderMale,
				Generation:       1,
				SonCount:         1,
				DaughterCount:    0,
				AdoptedHeirCount: 0,
				TotalChildren:    2, // 不匹配
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

func TestPersonIsAlive(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		person     *Person
		want       bool
	}{
		{
			name: "alive - no death info",
			person: &Person{
				DeathYear: nil,
			},
			want: true,
		},
		{
			name: "dead - has death year",
			person: &Person{
				DeathYear: intPtr(1980),
			},
			want: false,
		},
		{
			name: "alive - has death gregorian",
			person: &Person{
				DeathGregorian: &now,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.person.IsAlive(); got != tt.want {
				t.Errorf("Person.IsAlive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPersonAge(t *testing.T) {
	birthYear := 1980
	tests := []struct {
		name       string
		person     *Person
		want       int
	}{
		{
			name: "has birth year, alive",
			person: &Person{
				BirthYear: &birthYear,
				DeathYear: nil,
			},
			want: time.Now().Year() - 1980,
		},
		{
			name: "has birth year, dead",
			person: &Person{
				BirthYear: &birthYear,
				DeathYear: intPtr(2020),
			},
			want: 40,
		},
		{
			name: "no birth year",
			person: &Person{
				BirthYear: nil,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.person.Age(); got != tt.want {
				t.Errorf("Person.Age() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPersonFullName(t *testing.T) {
	tests := []struct {
		name     string
		person   *Person
		expected string
	}{
		{
			name:     "with style name",
			person:   &Person{Name: "张三", StyleName: "伯约"},
			expected: "张三（伯约）",
		},
		{
			name:     "without style name",
			person:   &Person{Name: "李四"},
			expected: "李四",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.person.FullName(); got != tt.expected {
				t.Errorf("Person.FullName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPersonHasChildren(t *testing.T) {
	tests := []struct {
		name   string
		person *Person
		want   bool
	}{
		{"has children", &Person{TotalChildren: 2}, true},
		{"no children", &Person{TotalChildren: 0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.person.HasChildren(); got != tt.want {
				t.Errorf("Person.HasChildren() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPersonHasSpouse(t *testing.T) {
	tests := []struct {
		name     string
		person   *Person
		expected bool
	}{
		{
			name:     "has spouses",
			person:   &Person{Spouses: []*Spouse{{Name: "王五"}}},
			expected: true,
		},
		{
			name:     "no spouses",
			person:   &Person{Spouses: []*Spouse{}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.person.HasSpouse(); got != tt.expected {
				t.Errorf("Person.HasSpouse() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestServiceCreatePerson(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewService(repo)
	ctx := context.Background()

	t.Run("create person successfully", func(t *testing.T) {
		p := &Person{
			Name:       "测试人物",
			Gender:     GenderMale,
			Generation: 1,
			SonCount:   1,
			DaughterCount: 0,
			AdoptedHeirCount: 0,
			TotalChildren: 1,
		}

		err := svc.CreatePerson(ctx, p)
		if err != nil {
			t.Fatalf("CreatePerson() error = %v", err)
		}
		if p.ID == 0 {
			t.Error("expected person ID to be set")
		}
	})

	t.Run("create person with invalid data", func(t *testing.T) {
		p := &Person{
			Name:   "", // 名称为空
			Gender: GenderMale,
		}

		err := svc.CreatePerson(ctx, p)
		if err == nil {
			t.Error("expected error for invalid data")
		}
	})
}

func TestServiceGetPerson(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewService(repo)
	ctx := context.Background()

	// 添加测试数据
	fatherID := int64(1)
	repo.persons[1] = &Person{
		ID:        1,
		Name:      "父亲",
		Gender:    GenderMale,
		Generation: 1,
	}
	repo.persons[2] = &Person{
		ID:        2,
		Name:      "儿子",
		Gender:    GenderMale,
		Generation: 2,
		FatherID: &fatherID,
	}
	repo.spouses[1] = []*Spouse{
		{ID: 1, Name: "配偶"},
	}
	repo.children[1] = []*PersonChild{
		{Person: repo.persons[2], RelationType: RelationBiological},
	}

	t.Run("get person without relations", func(t *testing.T) {
		// 创建一个没有关联数据的独立person用于测试
		repo.persons[99] = &Person{
			ID:        99,
			Name:      "无关联人物",
			Gender:    GenderMale,
			Generation: 1,
			Spouses:   nil,
			Children:  nil,
		}
		p, err := svc.GetPerson(ctx, 99, false)
		if err != nil {
			t.Fatalf("GetPerson() error = %v", err)
		}
		if p.Name != "无关联人物" {
			t.Errorf("expected name '无关联人物', got '%s'", p.Name)
		}
	})

	t.Run("get person with relations", func(t *testing.T) {
		p, err := svc.GetPerson(ctx, 1, true)
		if err != nil {
			t.Fatalf("GetPerson() error = %v", err)
		}
		if p.Name != "父亲" {
			t.Errorf("expected name '父亲', got '%s'", p.Name)
		}
		if len(p.Spouses) != 1 {
			t.Errorf("expected 1 spouse, got %d", len(p.Spouses))
		}
	})

	t.Run("get person with father", func(t *testing.T) {
		p, err := svc.GetPerson(ctx, 2, true)
		if err != nil {
			t.Fatalf("GetPerson() error = %v", err)
		}
		if p.Father == nil {
			t.Error("expected father to be loaded")
		}
		if p.Father.Name != "父亲" {
			t.Errorf("expected father name '父亲', got '%s'", p.Father.Name)
		}
	})

	t.Run("get non-existent person", func(t *testing.T) {
		_, err := svc.GetPerson(ctx, 999, false)
		if err != nil {
			t.Error("expected nil error for non-existent person")
		}
	})
}

func TestServiceUpdatePerson(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewService(repo)
	ctx := context.Background()

	repo.persons[1] = &Person{
		ID:        1,
		Name:      "原名",
		Gender:    GenderMale,
		Generation: 1,
		SonCount:   1,
		DaughterCount: 0,
		AdoptedHeirCount: 0,
		TotalChildren: 1,
	}

	t.Run("update person successfully", func(t *testing.T) {
		p, _ := repo.FindByID(ctx, 1)
		p.Name = "新名字"
		p.StyleName = "新字号"

		err := svc.UpdatePerson(ctx, p)
		if err != nil {
			t.Fatalf("UpdatePerson() error = %v", err)
		}

		updated, _ := repo.FindByID(ctx, 1)
		if updated.Name != "新名字" {
			t.Errorf("expected name '新名字', got '%s'", updated.Name)
		}
		if updated.StyleName != "新字号" {
			t.Errorf("expected style name '新字号', got '%s'", updated.StyleName)
		}
	})
}

func TestServiceDeletePerson(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewService(repo)
	ctx := context.Background()

	// 有子女的父亲
	repo.persons[1] = &Person{ID: 1, Name: "有子女的父亲", Generation: 1}
	repo.children[1] = []*PersonChild{
		{Person: &Person{ID: 2, Name: "子女"}},
	}

	// 没有子女的父亲
	repo.persons[3] = &Person{ID: 3, Name: "无子女的父亲", Generation: 1}

	t.Run("delete person with children should fail", func(t *testing.T) {
		err := svc.DeletePerson(ctx, 1)
		if err == nil {
			t.Error("expected error when deleting person with children")
		}
	})

	t.Run("delete person without children", func(t *testing.T) {
		err := svc.DeletePerson(ctx, 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestServiceSearchPersons(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewService(repo)
	ctx := context.Background()

	repo.persons[1] = &Person{Name: "张三", Gender: GenderMale, Generation: 1}
	repo.persons[2] = &Person{Name: "李四", Gender: GenderFemale, Generation: 2}
	repo.persons[3] = &Person{Name: "王五", Gender: GenderMale, Generation: 1}

	t.Run("search returns all persons", func(t *testing.T) {
		query := &SearchQuery{}
		persons, total, err := svc.SearchPersons(ctx, query)
		if err != nil {
			t.Fatalf("SearchPersons() error = %v", err)
		}
		if total != 3 {
			t.Errorf("expected total 3, got %d", total)
		}
		if len(persons) != 3 {
			t.Errorf("expected 3 persons, got %d", len(persons))
		}
	})

	t.Run("search respects page size limit", func(t *testing.T) {
		query := &SearchQuery{Page: 1, PageSize: 2}
		_, _, err := svc.SearchPersons(ctx, query)
		if err != nil {
			t.Fatalf("SearchPersons() error = %v", err)
		}
		// 注意：mock repository 不实现分页，pageSize 只影响领域服务
		// 领域服务会将 PageSize 限制在 100 以内
		if query.PageSize != 2 {
			t.Errorf("expected pageSize 2, got %d", query.PageSize)
		}
	})
}

func TestServiceGetStatistics(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewService(repo)
	ctx := context.Background()

	repo.persons[1] = &Person{Name: "第一代", Generation: 1}
	repo.persons[2] = &Person{Name: "第一代2", Generation: 1}
	repo.persons[3] = &Person{Name: "第二代", Generation: 2}

	t.Run("get statistics", func(t *testing.T) {
		stats, err := svc.GetStatistics(ctx)
		if err != nil {
			t.Fatalf("GetStatistics() error = %v", err)
		}
		if stats.TotalPersons != 3 {
			t.Errorf("expected total 3, got %d", stats.TotalPersons)
		}
		if stats.ByGeneration[1] != 2 {
			t.Errorf("expected generation 1 count 2, got %d", stats.ByGeneration[1])
		}
		if stats.ByGeneration[2] != 1 {
			t.Errorf("expected generation 2 count 1, got %d", stats.ByGeneration[2])
		}
	})
}

func TestServiceAddParentChildRelation(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewService(repo)
	ctx := context.Background()

	repo.persons[1] = &Person{ID: 1, Name: "父亲", Gender: GenderMale, Generation: 1}
	repo.persons[2] = &Person{ID: 2, Name: "儿子", Gender: GenderMale, Generation: 2, FatherID: nil}

	t.Run("add parent-child relation successfully", func(t *testing.T) {
		err := svc.AddParentChildRelation(ctx, 1, 2, RelationBiological, 1, true)
		if err != nil {
			t.Fatalf("AddParentChildRelation() error = %v", err)
		}

		// 验证儿子的父亲ID已更新
		child, _ := repo.FindByID(ctx, 2)
		if child.FatherID == nil || *child.FatherID != 1 {
			t.Errorf("expected FatherID to be 1, got %v", child.FatherID)
		}
		// 验证世代已更新
		if child.Generation != 2 {
			t.Errorf("expected Generation to be 2, got %d", child.Generation)
		}
	})

	t.Run("self-reference not allowed", func(t *testing.T) {
		err := svc.AddParentChildRelation(ctx, 1, 1, RelationBiological, 1, true)
		if err == nil {
			t.Error("expected error for self-reference")
		}
		if err.Error() != "self reference not allowed" {
			t.Errorf("expected 'self reference not allowed', got '%s'", err.Error())
		}
	})

	t.Run("parent not found", func(t *testing.T) {
		err := svc.AddParentChildRelation(ctx, 999, 2, RelationBiological, 1, true)
		if err == nil {
			t.Error("expected error when parent not found")
		}
		if err.Error() != "parent person not found: 999" {
			t.Errorf("expected 'parent person not found: 999', got '%s'", err.Error())
		}
	})

	t.Run("child not found", func(t *testing.T) {
		err := svc.AddParentChildRelation(ctx, 1, 999, RelationBiological, 1, true)
		if err == nil {
			t.Error("expected error when child not found")
		}
		if err.Error() != "child person not found: 999" {
			t.Errorf("expected 'child person not found: 999', got '%s'", err.Error())
		}
	})
}

func TestServiceCheckCycle(t *testing.T) {
	repo := newMockPersonRepositoryWithDescendants()
	svc := NewService(repo)
	ctx := context.Background()

	// 设置：1是2的父亲，2是3的父亲
	repo.persons[1] = &Person{ID: 1, Name: "祖父", Gender: GenderMale, Generation: 1}
	repo.persons[2] = &Person{ID: 2, Name: "父亲", Gender: GenderMale, Generation: 2, FatherID: int64Ptr(1)}
	repo.persons[3] = &Person{ID: 3, Name: "儿子", Gender: GenderMale, Generation: 3, FatherID: int64Ptr(2)}

	// 设置1的后代是2和3
	repo.SetDescendants(1, []*Person{repo.persons[2], repo.persons[3]})
	// 设置2的后代是3
	repo.SetDescendants(2, []*Person{repo.persons[3]})

	t.Run("cycle detected - ancestor as child", func(t *testing.T) {
		// 尝试让3成为1的父亲（形成循环：1->2->3->1）
		err := svc.AddParentChildRelation(ctx, 3, 1, RelationBiological, 1, true)
		if err == nil {
			t.Error("expected error when cycle detected")
		}
		if err.Error() != "cycle detected: person 3 is already a descendant of 1" {
			t.Errorf("expected 'cycle detected' error, got '%s'", err.Error())
		}
	})
}

// mockRepositoryWithDescendants 支持FindDescendants的mock
type mockRepositoryWithDescendants struct {
	*mockRepository
	descendants map[int64][]*Person
}

func newMockPersonRepositoryWithDescendants() *mockRepositoryWithDescendants {
	return &mockRepositoryWithDescendants{
		mockRepository: newMockPersonRepository(),
		descendants:    make(map[int64][]*Person),
	}
}

func (m *mockRepositoryWithDescendants) FindDescendants(ctx context.Context, personID int64, depth int) ([]*Person, error) {
	return m.descendants[personID], nil
}

func (m *mockRepositoryWithDescendants) SetDescendants(personID int64, descendants []*Person) {
	m.descendants[personID] = descendants
}

// 辅助函数
func intPtr(i int) *int {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}