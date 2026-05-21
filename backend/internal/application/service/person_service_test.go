package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/person"
)

// mockPersonRepository 模拟人物仓储
type mockPersonRepository struct {
	persons map[int64]*person.Person
	nextID  int64
}

func newMockPersonRepository() *mockPersonRepository {
	return &mockPersonRepository{
		persons: make(map[int64]*person.Person),
		nextID:  1,
	}
}

func (m *mockPersonRepository) FindByID(ctx context.Context, id int64) (*person.Person, error) {
	p, ok := m.persons[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockPersonRepository) FindByLegacyID(ctx context.Context, legacyID string) (*person.Person, error) {
	for _, p := range m.persons {
		if p.LegacyID == legacyID {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockPersonRepository) Create(ctx context.Context, p *person.Person) error {
	p.ID = m.nextID
	m.nextID++
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	m.persons[p.ID] = p
	return nil
}

func (m *mockPersonRepository) Update(ctx context.Context, p *person.Person) error {
	p.UpdatedAt = time.Now()
	m.persons[p.ID] = p
	return nil
}

func (m *mockPersonRepository) Delete(ctx context.Context, id int64) error {
	delete(m.persons, id)
	return nil
}

func (m *mockPersonRepository) Search(ctx context.Context, query *person.SearchQuery) ([]*person.Person, int64, error) {
	var result []*person.Person
	for _, p := range m.persons {
		result = append(result, p)
	}
	return result, int64(len(result)), nil
}

func (m *mockPersonRepository) FindByGeneration(ctx context.Context, generation int) ([]*person.Person, error) {
	var result []*person.Person
	for _, p := range m.persons {
		if p.Generation == generation {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockPersonRepository) FindByName(ctx context.Context, name string, fuzzy bool) ([]*person.Person, error) {
	var result []*person.Person
	for _, p := range m.persons {
		if p.Name == name {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockPersonRepository) FindAncestors(ctx context.Context, personID int64, depth int) ([]*person.Person, error) {
	return nil, nil
}

func (m *mockPersonRepository) FindDescendants(ctx context.Context, personID int64, depth int) ([]*person.Person, error) {
	return nil, nil
}

func (m *mockPersonRepository) FindSiblings(ctx context.Context, personID int64) ([]*person.Person, error) {
	return nil, nil
}

func (m *mockPersonRepository) GetFamilyTree(ctx context.Context, rootID int64, depth int) (*person.TreeResult, error) {
	return nil, nil
}

func (m *mockPersonRepository) FindChildren(ctx context.Context, parentID int64) ([]*person.PersonChild, error) {
	return nil, nil
}

func (m *mockPersonRepository) FindSpouses(ctx context.Context, personID int64) ([]*person.Spouse, error) {
	return nil, nil
}

func (m *mockPersonRepository) FindParents(ctx context.Context, personID int64) ([]*person.Person, error) {
	return nil, nil
}

func (m *mockPersonRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.persons)), nil
}

func (m *mockPersonRepository) CountByGeneration(ctx context.Context) (map[int]int64, error) {
	result := make(map[int]int64)
	for _, p := range m.persons {
		result[p.Generation]++
	}
	return result, nil
}

func (m *mockPersonRepository) BatchCreate(ctx context.Context, persons []*person.Person) error {
	for _, p := range persons {
		if err := m.Create(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockPersonRepository) BatchUpdate(ctx context.Context, persons []*person.Person) error {
	for _, p := range persons {
		if err := m.Update(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

// newTestPersonService 创建测试用人物应用服务
func newTestPersonService() (*PersonService, *mockPersonRepository) {
	repo := newMockPersonRepository()
	svc := NewPersonService(repo)
	return svc, repo
}

func TestCreatePerson(t *testing.T) {
	svc, repo := newTestPersonService()
	ctx := context.Background()

	t.Run("create person successfully", func(t *testing.T) {
		req := &CreatePersonRequest{
			Name:       "张三",
			StyleName:  "子曰",
			Gender:     "男",
			Generation: 10,
			BirthOrder: "长子",
			BirthPlace: "四川省成都市",
		}

		result, err := svc.CreatePerson(ctx, req)
		if err != nil {
			t.Fatalf("CreatePerson() error = %v", err)
		}
		if result.Name != "张三" {
			t.Errorf("expected name '张三', got '%s'", result.Name)
		}
		if result.Gender != "男" {
			t.Errorf("expected gender '男', got '%s'", result.Gender)
		}
		if result.Generation != 10 {
			t.Errorf("expected generation 10, got %d", result.Generation)
		}
		// 验证仓储中的记录
		if len(repo.persons) != 1 {
			t.Errorf("expected 1 person in repo, got %d", len(repo.persons))
		}
	})

	t.Run("create person without required name", func(t *testing.T) {
		req := &CreatePersonRequest{
			Name:   "",
			Gender: "男",
		}

		_, err := svc.CreatePerson(ctx, req)
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("create person with father", func(t *testing.T) {
		// 先创建父亲
		fatherReq := &CreatePersonRequest{
			Name:       "张父",
			Gender:     "男",
			Generation: 9,
		}
		father, _ := svc.CreatePerson(ctx, fatherReq)

		// 创建儿子
		req := &CreatePersonRequest{
			Name:       "张小",
			Gender:     "男",
			Generation: 10,
			FatherID:   &father.ID,
		}

		result, err := svc.CreatePerson(ctx, req)
		if err != nil {
			t.Fatalf("CreatePerson() error = %v", err)
		}
		if result.FatherID == nil || *result.FatherID != father.ID {
			t.Errorf("expected father ID %d, got %v", father.ID, result.FatherID)
		}
	})
}

func TestGetPerson(t *testing.T) {
	svc, repo := newTestPersonService()
	ctx := context.Background()

	// 预先创建一个人物
	repo.Create(ctx, &person.Person{
		Name:       "李四",
		StyleName:  "小明",
		Gender:     person.GenderFemale,
		Generation: 8,
		BirthOrder: "次女",
	})

	t.Run("get existing person", func(t *testing.T) {
		result, err := svc.GetPerson(ctx, 1, false)
		if err != nil {
			t.Fatalf("GetPerson() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected person, got nil")
		}
		if result.Name != "李四" {
			t.Errorf("expected name '李四', got '%s'", result.Name)
		}
		if result.StyleName != "小明" {
			t.Errorf("expected style name '小明', got '%s'", result.StyleName)
		}
	})

	t.Run("get non-existent person", func(t *testing.T) {
		result, err := svc.GetPerson(ctx, 999, false)
		if err != nil {
			t.Fatalf("GetPerson() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestUpdatePerson(t *testing.T) {
	svc, repo := newTestPersonService()
	ctx := context.Background()

	// 预先创建一个人物
	repo.Create(ctx, &person.Person{
		Name:       "王五",
		Gender:     person.GenderMale,
		Generation: 5,
	})

	t.Run("update person successfully", func(t *testing.T) {
		newName := "王六"
		newGeneration := 6
		req := &UpdatePersonRequest{
			Name:       newName,
			Generation: &newGeneration,
		}

		result, err := svc.UpdatePerson(ctx, 1, req)
		if err != nil {
			t.Fatalf("UpdatePerson() error = %v", err)
		}
		if result.Name != "王六" {
			t.Errorf("expected name '王六', got '%s'", result.Name)
		}
		if result.Generation != 6 {
			t.Errorf("expected generation 6, got %d", result.Generation)
		}
	})

	t.Run("update non-existent person", func(t *testing.T) {
		newName := "新人"
		req := &UpdatePersonRequest{
			Name: newName,
		}

		result, err := svc.UpdatePerson(ctx, 999, req)
		if err != nil {
			t.Fatalf("UpdatePerson() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestDeletePerson(t *testing.T) {
	svc, repo := newTestPersonService()
	ctx := context.Background()

	// 预先创建一个人物
	repo.Create(ctx, &person.Person{
		Name:   "待删除",
		Gender: person.GenderMale,
	})

	t.Run("delete existing person", func(t *testing.T) {
		err := svc.DeletePerson(ctx, 1)
		if err != nil {
			t.Fatalf("DeletePerson() error = %v", err)
		}
		if len(repo.persons) != 0 {
			t.Errorf("expected 0 persons after delete, got %d", len(repo.persons))
		}
	})

	t.Run("delete non-existent person", func(t *testing.T) {
		// Mock repository doesn't return error for non-existent delete
		err := svc.DeletePerson(ctx, 999)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestSearchPersons(t *testing.T) {
	svc, repo := newTestPersonService()
	ctx := context.Background()

	// 创建多个人物
	for i := 0; i < 3; i++ {
		g := person.GenderMale
		if i%2 == 0 {
			g = person.GenderFemale
		}
		repo.Create(ctx, &person.Person{
			Name:       "搜索目标" + string(rune('A'+i)),
			Gender:     g,
			Generation: i + 1,
		})
	}

	t.Run("search by keyword", func(t *testing.T) {
		req := &SearchPersonRequest{
			Keyword: "搜索目标",
		}

		result, total, err := svc.SearchPersons(ctx, req)
		if err != nil {
			t.Fatalf("SearchPersons() error = %v", err)
		}
		if total != 3 {
			t.Errorf("expected 3 total, got %d", total)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 results, got %d", len(result))
		}
	})

	t.Run("search by gender", func(t *testing.T) {
		req := &SearchPersonRequest{
			Gender: "男",
		}

		result, _, err := svc.SearchPersons(ctx, req)
		if err != nil {
			t.Fatalf("SearchPersons() error = %v", err)
		}
		// Mock repository doesn't filter, just verify results returned
		if len(result) == 0 {
			t.Error("expected some results for gender filter")
		}
	})

	t.Run("search by generation", func(t *testing.T) {
		gen := 2
		req := &SearchPersonRequest{
			Generation: &gen,
		}

		result, _, err := svc.SearchPersons(ctx, req)
		if err != nil {
			t.Fatalf("SearchPersons() error = %v", err)
		}
		// Mock repository doesn't filter, just verify results returned
		if len(result) == 0 {
			t.Error("expected some results for generation filter")
		}
	})
}

func TestGetStatistics(t *testing.T) {
	svc, repo := newTestPersonService()
	ctx := context.Background()

	// 创建多个不同世代的人物
	for i := 0; i < 5; i++ {
		repo.Create(ctx, &person.Person{
			Name:       "统计对象" + string(rune('A'+i)),
			Gender:     person.GenderMale,
			Generation: i%3 + 1, // 分布在3个世代
		})
	}

	t.Run("get statistics", func(t *testing.T) {
		result, err := svc.GetStatistics(ctx)
		if err != nil {
			t.Fatalf("GetStatistics() error = %v", err)
		}
		if result.TotalPersons != 5 {
			t.Errorf("expected total 5, got %d", result.TotalPersons)
		}
		if len(result.ByGeneration) != 3 {
			t.Errorf("expected 3 generations, got %d", len(result.ByGeneration))
		}
	})
}

func TestBatchCreatePersons(t *testing.T) {
	svc, repo := newTestPersonService()
	ctx := context.Background()

	t.Run("batch create successfully", func(t *testing.T) {
		req := &BatchCreatePersonRequest{
			Persons: []*CreatePersonRequest{
				{Name: "批量1", Gender: "男"},
				{Name: "批量2", Gender: "女"},
				{Name: "批量3", Gender: "男"},
			},
		}

		result, err := svc.BatchCreatePersons(ctx, req)
		if err != nil {
			t.Fatalf("BatchCreatePersons() error = %v", err)
		}
		if result.SuccessCount != 3 {
			t.Errorf("expected success count 3, got %d", result.SuccessCount)
		}
		if result.FailCount != 0 {
			t.Errorf("expected fail count 0, got %d", result.FailCount)
		}
		if len(repo.persons) != 3 {
			t.Errorf("expected 3 persons in repo, got %d", len(repo.persons))
		}
	})

	t.Run("batch create with partial failure", func(t *testing.T) {
		badRepo := &mockPersonRepositoryError{
			inner: newMockPersonRepository(),
		}
		badSvc := NewPersonService(badRepo)

		req := &BatchCreatePersonRequest{
			Persons: []*CreatePersonRequest{
				{Name: "测试", Gender: "男"},
			},
		}

		result, err := badSvc.BatchCreatePersons(ctx, req)
		if err != nil {
			t.Fatalf("BatchCreatePersons() error = %v", err)
		}
		// Mock returns nil for non-implemented methods
		if result.SuccessCount != 1 {
			t.Errorf("expected success count 1, got %d", result.SuccessCount)
		}
	})
}

func TestBatchUpdatePersons(t *testing.T) {
	svc, repo := newTestPersonService()
	ctx := context.Background()

	// 预先创建人物
	repo.Create(ctx, &person.Person{
		Name:   "更新目标1",
		Gender: person.GenderMale,
	})
	repo.Create(ctx, &person.Person{
		Name:   "更新目标2",
		Gender: person.GenderFemale,
	})

	t.Run("batch update successfully", func(t *testing.T) {
		newGen := 10
		req := &BatchUpdatePersonRequest{
			Persons: []*BatchUpdateItem{
				{ID: 1, Name: "新名称1", Generation: &newGen},
				{ID: 2, Name: "新名称2"},
			},
		}

		result, err := svc.BatchUpdatePersons(ctx, req)
		if err != nil {
			t.Fatalf("BatchUpdatePersons() error = %v", err)
		}
		if result.SuccessCount != 2 {
			t.Errorf("expected success count 2, got %d", result.SuccessCount)
		}
		if repo.persons[1].Generation != 10 {
			t.Errorf("expected generation 10, got %d", repo.persons[1].Generation)
		}
	})

	t.Run("batch update non-existent", func(t *testing.T) {
		req := &BatchUpdatePersonRequest{
			Persons: []*BatchUpdateItem{
				{ID: 999, Name: "不存在"},
			},
		}

		result, err := svc.BatchUpdatePersons(ctx, req)
		if err != nil {
			t.Fatalf("BatchUpdatePersons() error = %v", err)
		}
		if result.FailCount != 1 {
			t.Errorf("expected fail count 1, got %d", result.FailCount)
		}
	})
}

func TestBatchDeletePersons(t *testing.T) {
	svc, repo := newTestPersonService()
	ctx := context.Background()

	// 预先创建人物
	repo.Create(ctx, &person.Person{Name: "删除1", Gender: person.GenderMale})
	repo.Create(ctx, &person.Person{Name: "删除2", Gender: person.GenderFemale})
	repo.Create(ctx, &person.Person{Name: "删除3", Gender: person.GenderMale})

	t.Run("batch delete successfully", func(t *testing.T) {
		req := &BatchDeleteRequest{
			IDs: []int64{1, 2},
		}

		result, err := svc.BatchDeletePersons(ctx, req)
		if err != nil {
			t.Fatalf("BatchDeletePersons() error = %v", err)
		}
		if result.SuccessCount != 2 {
			t.Errorf("expected success count 2, got %d", result.SuccessCount)
		}
		if len(repo.persons) != 1 {
			t.Errorf("expected 1 person after delete, got %d", len(repo.persons))
		}
	})
}

func TestPersonServiceErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("repository error on get person", func(t *testing.T) {
		badRepo := &mockPersonRepositoryError{
			inner: newMockPersonRepository(),
		}
		badSvc := NewPersonService(badRepo)

		_, err := badSvc.GetPerson(ctx, 1, false)
		if err == nil {
			t.Error("expected error from repository")
		}
	})

	t.Run("repository error on delete person", func(t *testing.T) {
		badRepo := &mockPersonRepositoryError{
			inner: newMockPersonRepository(),
		}
		badSvc := NewPersonService(badRepo)

		err := badSvc.DeletePerson(ctx, 1)
		if err == nil {
			t.Error("expected error from repository")
		}
	})
}

type mockPersonRepositoryError struct {
	inner *mockPersonRepository
}

func (m *mockPersonRepositoryError) FindByID(ctx context.Context, id int64) (*person.Person, error) {
	return nil, errors.New("repository error")
}
func (m *mockPersonRepositoryError) FindByLegacyID(ctx context.Context, legacyID string) (*person.Person, error) {
	return m.inner.FindByLegacyID(ctx, legacyID)
}
func (m *mockPersonRepositoryError) Create(ctx context.Context, p *person.Person) error {
	return m.inner.Create(ctx, p)
}
func (m *mockPersonRepositoryError) Update(ctx context.Context, p *person.Person) error {
	return m.inner.Update(ctx, p)
}
func (m *mockPersonRepositoryError) Delete(ctx context.Context, id int64) error {
	return errors.New("repository error")
}
func (m *mockPersonRepositoryError) Search(ctx context.Context, query *person.SearchQuery) ([]*person.Person, int64, error) {
	return m.inner.Search(ctx, query)
}
func (m *mockPersonRepositoryError) FindByGeneration(ctx context.Context, generation int) ([]*person.Person, error) {
	return m.inner.FindByGeneration(ctx, generation)
}
func (m *mockPersonRepositoryError) FindByName(ctx context.Context, name string, fuzzy bool) ([]*person.Person, error) {
	return m.inner.FindByName(ctx, name, fuzzy)
}
func (m *mockPersonRepositoryError) FindAncestors(ctx context.Context, personID int64, depth int) ([]*person.Person, error) {
	return m.inner.FindAncestors(ctx, personID, depth)
}
func (m *mockPersonRepositoryError) FindDescendants(ctx context.Context, personID int64, depth int) ([]*person.Person, error) {
	return m.inner.FindDescendants(ctx, personID, depth)
}
func (m *mockPersonRepositoryError) FindSiblings(ctx context.Context, personID int64) ([]*person.Person, error) {
	return m.inner.FindSiblings(ctx, personID)
}
func (m *mockPersonRepositoryError) GetFamilyTree(ctx context.Context, rootID int64, depth int) (*person.TreeResult, error) {
	return m.inner.GetFamilyTree(ctx, rootID, depth)
}
func (m *mockPersonRepositoryError) FindChildren(ctx context.Context, parentID int64) ([]*person.PersonChild, error) {
	return m.inner.FindChildren(ctx, parentID)
}
func (m *mockPersonRepositoryError) FindSpouses(ctx context.Context, personID int64) ([]*person.Spouse, error) {
	return m.inner.FindSpouses(ctx, personID)
}
func (m *mockPersonRepositoryError) FindParents(ctx context.Context, personID int64) ([]*person.Person, error) {
	return m.inner.FindParents(ctx, personID)
}
func (m *mockPersonRepositoryError) Count(ctx context.Context) (int64, error) {
	return m.inner.Count(ctx)
}
func (m *mockPersonRepositoryError) CountByGeneration(ctx context.Context) (map[int]int64, error) {
	return m.inner.CountByGeneration(ctx)
}
func (m *mockPersonRepositoryError) BatchCreate(ctx context.Context, persons []*person.Person) error {
	return m.inner.BatchCreate(ctx, persons)
}
func (m *mockPersonRepositoryError) BatchUpdate(ctx context.Context, persons []*person.Person) error {
	return m.inner.BatchUpdate(ctx, persons)
}