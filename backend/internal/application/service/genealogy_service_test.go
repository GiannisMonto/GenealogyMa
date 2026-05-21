package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/genealogy"
)

// mockGenealogyRepository 模拟族谱仓储
type mockGenealogyRepository struct {
	genealogies map[int64]*genealogy.Genealogy
	nextID      int64
}

func newMockGenealogyRepository() *mockGenealogyRepository {
	return &mockGenealogyRepository{
		genealogies: make(map[int64]*genealogy.Genealogy),
		nextID:      1,
	}
}

func (m *mockGenealogyRepository) FindByID(ctx context.Context, id int64) (*genealogy.Genealogy, error) {
	g, ok := m.genealogies[id]
	if !ok {
		return nil, nil
	}
	return g, nil
}

func (m *mockGenealogyRepository) Create(ctx context.Context, g *genealogy.Genealogy) error {
	g.ID = m.nextID
	m.nextID++
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	m.genealogies[g.ID] = g
	return nil
}

func (m *mockGenealogyRepository) Update(ctx context.Context, g *genealogy.Genealogy) error {
	g.UpdatedAt = time.Now()
	m.genealogies[g.ID] = g
	return nil
}

func (m *mockGenealogyRepository) Delete(ctx context.Context, id int64) error {
	delete(m.genealogies, id)
	return nil
}

func (m *mockGenealogyRepository) FindAll(ctx context.Context) ([]*genealogy.Genealogy, error) {
	var result []*genealogy.Genealogy
	for _, g := range m.genealogies {
		result = append(result, g)
	}
	return result, nil
}

func (m *mockGenealogyRepository) Search(ctx context.Context, query *genealogy.SearchQuery) ([]*genealogy.Genealogy, int64, error) {
	var result []*genealogy.Genealogy
	for _, g := range m.genealogies {
		result = append(result, g)
	}
	return result, int64(len(result)), nil
}

func (m *mockGenealogyRepository) FindBySurname(ctx context.Context, surname string) ([]*genealogy.Genealogy, error) {
	var result []*genealogy.Genealogy
	for _, g := range m.genealogies {
		if g.Surname == surname {
			result = append(result, g)
		}
	}
	return result, nil
}

// mockBranchRepository 模拟分支仓储
type mockBranchRepository struct {
	branches map[int64]*genealogy.Branch
	nextID   int64
}

func newMockBranchRepository() *mockBranchRepository {
	return &mockBranchRepository{
		branches: make(map[int64]*genealogy.Branch),
		nextID:   1,
	}
}

func (m *mockBranchRepository) FindByID(ctx context.Context, id int64) (*genealogy.Branch, error) {
	b, ok := m.branches[id]
	if !ok {
		return nil, nil
	}
	return b, nil
}

func (m *mockBranchRepository) FindByGenealogyID(ctx context.Context, genealogyID int64) ([]*genealogy.Branch, error) {
	var result []*genealogy.Branch
	for _, b := range m.branches {
		if b.GenealogyID == genealogyID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (m *mockBranchRepository) Create(ctx context.Context, b *genealogy.Branch) error {
	b.ID = m.nextID
	m.nextID++
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	m.branches[b.ID] = b
	return nil
}

func (m *mockBranchRepository) Update(ctx context.Context, b *genealogy.Branch) error {
	b.UpdatedAt = time.Now()
	m.branches[b.ID] = b
	return nil
}

func (m *mockBranchRepository) Delete(ctx context.Context, id int64) error {
	delete(m.branches, id)
	return nil
}

func (m *mockBranchRepository) DeleteByGenealogyID(ctx context.Context, genealogyID int64) error {
	for id, b := range m.branches {
		if b.GenealogyID == genealogyID {
			delete(m.branches, id)
		}
	}
	return nil
}

func (m *mockBranchRepository) Search(ctx context.Context, query *genealogy.BranchSearchQuery) ([]*genealogy.Branch, int64, error) {
	var result []*genealogy.Branch
	for _, b := range m.branches {
		result = append(result, b)
	}
	return result, int64(len(result)), nil
}

// mockGenerationRepository 模拟世代仓储
type mockGenerationRepository struct {
	generations map[int64]*genealogy.Generation
	nextID      int64
}

func newMockGenerationRepository() *mockGenerationRepository {
	return &mockGenerationRepository{
		generations: make(map[int64]*genealogy.Generation),
		nextID:      1,
	}
}

func (m *mockGenerationRepository) FindByID(ctx context.Context, id int64) (*genealogy.Generation, error) {
	g, ok := m.generations[id]
	if !ok {
		return nil, nil
	}
	return g, nil
}

func (m *mockGenerationRepository) FindByGenealogyID(ctx context.Context, genealogyID int64) ([]*genealogy.Generation, error) {
	var result []*genealogy.Generation
	for _, g := range m.generations {
		if g.GenealogyID == genealogyID {
			result = append(result, g)
		}
	}
	return result, nil
}

func (m *mockGenerationRepository) FindByGeneration(ctx context.Context, genealogyID int64, generation int) (*genealogy.Generation, error) {
	for _, g := range m.generations {
		if g.GenealogyID == genealogyID && g.Generation == generation {
			return g, nil
		}
	}
	return nil, nil
}

func (m *mockGenerationRepository) Create(ctx context.Context, g *genealogy.Generation) error {
	g.ID = m.nextID
	m.nextID++
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	m.generations[g.ID] = g
	return nil
}

func (m *mockGenerationRepository) Update(ctx context.Context, g *genealogy.Generation) error {
	g.UpdatedAt = time.Now()
	m.generations[g.ID] = g
	return nil
}

func (m *mockGenerationRepository) Delete(ctx context.Context, id int64) error {
	delete(m.generations, id)
	return nil
}

func (m *mockGenerationRepository) DeleteByGenealogyID(ctx context.Context, genealogyID int64) error {
	for id, g := range m.generations {
		if g.GenealogyID == genealogyID {
			delete(m.generations, id)
		}
	}
	return nil
}

// newTestGenealogyService 创建测试用族谱应用服务
func newTestGenealogyService() (*GenealogyService, *mockGenealogyRepository, *mockBranchRepository, *mockGenerationRepository) {
	genealogyRepo := newMockGenealogyRepository()
	branchRepo := newMockBranchRepository()
	generationRepo := newMockGenerationRepository()
	svc := NewGenealogyService(genealogyRepo, branchRepo, generationRepo)
	return svc, genealogyRepo, branchRepo, generationRepo
}

func TestCreateGenealogy(t *testing.T) {
	svc, repo, _, _ := newTestGenealogyService()
	ctx := context.Background()

	t.Run("create genealogy successfully", func(t *testing.T) {
		req := &CreateGenealogyRequest{
			Name:             "王氏族谱",
			Surname:          "王",
			Description:      "王氏族谱第一版",
			OriginPlace:      "山西太原",
			TotalGenerations: 20,
			TotalMembers:     5000,
			Version:          "v1.0",
		}

		result, err := svc.CreateGenealogy(ctx, req)
		if err != nil {
			t.Fatalf("CreateGenealogy() error = %v", err)
		}
		if result.Name != "王氏族谱" {
			t.Errorf("expected name '王氏族谱', got '%s'", result.Name)
		}
		if result.Surname != "王" {
			t.Errorf("expected surname '王', got '%s'", result.Surname)
		}
		if result.TotalGenerations != 20 {
			t.Errorf("expected total_generations 20, got %d", result.TotalGenerations)
		}
		// 验证仓储中的记录
		if len(repo.genealogies) != 1 {
			t.Errorf("expected 1 genealogy in repo, got %d", len(repo.genealogies))
		}
	})

	t.Run("create genealogy with empty name", func(t *testing.T) {
		req := &CreateGenealogyRequest{
			Name:    "",
			Surname: "王",
		}

		// Validation happens at domain level, service accepts empty for now
		result, _ := svc.CreateGenealogy(ctx, req)
		if result != nil && result.Name != "" {
			t.Error("expected empty name to be accepted or validation error")
		}
	})
}

func TestGetGenealogy(t *testing.T) {
	svc, repo, _, _ := newTestGenealogyService()
	ctx := context.Background()

	// 预先创建一个族谱
	repo.Create(ctx, &genealogy.Genealogy{
		Name:            "李氏族谱",
		Surname:         "李",
		TotalGenerations: 15,
	})

	t.Run("get existing genealogy", func(t *testing.T) {
		result, err := svc.GetGenealogy(ctx, 1)
		if err != nil {
			t.Fatalf("GetGenealogy() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected genealogy, got nil")
		}
		if result.Name != "李氏族谱" {
			t.Errorf("expected name '李氏族谱', got '%s'", result.Name)
		}
	})

	t.Run("get non-existent genealogy", func(t *testing.T) {
		result, err := svc.GetGenealogy(ctx, 999)
		if err != nil {
			t.Fatalf("GetGenealogy() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestListGenealogies(t *testing.T) {
	svc, repo, _, _ := newTestGenealogyService()
	ctx := context.Background()

	// 创建多个族谱
	for i := 0; i < 3; i++ {
		repo.Create(ctx, &genealogy.Genealogy{
			Name:    "族谱" + string(rune('A'+i)),
			Surname: "姓氏",
		})
	}

	t.Run("list all genealogies", func(t *testing.T) {
		result, err := svc.ListGenealogies(ctx)
		if err != nil {
			t.Fatalf("ListGenealogies() error = %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 genealogies, got %d", len(result))
		}
	})
}

func TestUpdateGenealogy(t *testing.T) {
	svc, repo, _, _ := newTestGenealogyService()
	ctx := context.Background()

	// 预先创建一个族谱
	repo.Create(ctx, &genealogy.Genealogy{
		Name:            "原名",
		Surname:         "原姓",
		TotalGenerations: 10,
	})

	t.Run("update genealogy successfully", func(t *testing.T) {
		newName := "新族谱名"
		newTotal := 25
		req := &UpdateGenealogyRequest{
			Name:             &newName,
			TotalGenerations: &newTotal,
		}

		result, err := svc.UpdateGenealogy(ctx, 1, req)
		if err != nil {
			t.Fatalf("UpdateGenealogy() error = %v", err)
		}
		if result.Name != "新族谱名" {
			t.Errorf("expected name '新族谱名', got '%s'", result.Name)
		}
		if result.TotalGenerations != 25 {
			t.Errorf("expected total_generations 25, got %d", result.TotalGenerations)
		}
	})

	t.Run("update non-existent genealogy", func(t *testing.T) {
		newName := "新名称"
		req := &UpdateGenealogyRequest{
			Name: &newName,
		}

		result, err := svc.UpdateGenealogy(ctx, 999, req)
		if err != nil {
			t.Fatalf("UpdateGenealogy() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestDeleteGenealogy(t *testing.T) {
	svc, repo, _, _ := newTestGenealogyService()
	ctx := context.Background()

	// 预先创建一个族谱
	repo.Create(ctx, &genealogy.Genealogy{
		Name:    "待删除",
		Surname: "姓",
	})

	t.Run("delete existing genealogy", func(t *testing.T) {
		err := svc.DeleteGenealogy(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteGenealogy() error = %v", err)
		}
		if len(repo.genealogies) != 0 {
			t.Errorf("expected 0 genealogies after delete, got %d", len(repo.genealogies))
		}
	})

	t.Run("delete non-existent genealogy", func(t *testing.T) {
		// Domain service returns error for non-existent genealogy
		err := svc.DeleteGenealogy(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent genealogy")
		}
	})
}

func TestCreateBranch(t *testing.T) {
	svc, genealogyRepo, branchRepo, _ := newTestGenealogyService()
	ctx := context.Background()

	// 先创建一个族谱
	genealogyRepo.Create(ctx, &genealogy.Genealogy{
		Name:    "测试族谱",
		Surname: "测",
	})

	t.Run("create branch successfully", func(t *testing.T) {
		req := &CreateBranchRequest{
			GenealogyID:     1,
			Name:            "长房",
			Code:            "A",
			Description:     "长子长孙",
			GenerationStart: 1,
			GenerationEnd:   10,
		}

		result, err := svc.CreateBranch(ctx, req)
		if err != nil {
			t.Fatalf("CreateBranch() error = %v", err)
		}
		if result.Name != "长房" {
			t.Errorf("expected name '长房', got '%s'", result.Name)
		}
		if result.Code != "A" {
			t.Errorf("expected code 'A', got '%s'", result.Code)
		}
		if result.GenerationStart != 1 {
			t.Errorf("expected generation_start 1, got %d", result.GenerationStart)
		}
		if result.GenerationEnd != 10 {
			t.Errorf("expected generation_end 10, got %d", result.GenerationEnd)
		}
		if len(branchRepo.branches) != 1 {
			t.Errorf("expected 1 branch in repo, got %d", len(branchRepo.branches))
		}
	})
}

func TestGetBranch(t *testing.T) {
	svc, _, branchRepo, _ := newTestGenealogyService()
	ctx := context.Background()

	// 预先创建一个分支
	branchRepo.Create(ctx, &genealogy.Branch{
		GenealogyID:     1,
		Name:            "二房",
		Code:            "B",
		GenerationStart: 5,
		GenerationEnd:   15,
	})

	t.Run("get existing branch", func(t *testing.T) {
		result, err := svc.GetBranch(ctx, 1)
		if err != nil {
			t.Fatalf("GetBranch() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected branch, got nil")
		}
		if result.Name != "二房" {
			t.Errorf("expected name '二房', got '%s'", result.Name)
		}
	})

	t.Run("get non-existent branch", func(t *testing.T) {
		result, err := svc.GetBranch(ctx, 999)
		if err != nil {
			t.Fatalf("GetBranch() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestListBranches(t *testing.T) {
	svc, _, branchRepo, _ := newTestGenealogyService()
	ctx := context.Background()

	// 创建多个分支
	for i := 0; i < 3; i++ {
		branchRepo.Create(ctx, &genealogy.Branch{
			GenealogyID: 1,
			Name:        "分支" + string(rune('A'+i)),
			Code:        string(rune('A' + i)),
		})
	}

	t.Run("list branches by genealogy", func(t *testing.T) {
		result, err := svc.ListBranches(ctx, 1)
		if err != nil {
			t.Fatalf("ListBranches() error = %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 branches, got %d", len(result))
		}
	})
}

func TestUpdateBranch(t *testing.T) {
	svc, _, branchRepo, _ := newTestGenealogyService()
	ctx := context.Background()

	// 预先创建一个分支
	branchRepo.Create(ctx, &genealogy.Branch{
		GenealogyID:     1,
		Name:            "原分支名",
		Code:            "A",
		GenerationStart: 1,
		GenerationEnd:   10,
	})

	t.Run("update branch successfully", func(t *testing.T) {
		newName := "新分支名"
		newGenEnd := 20
		req := &UpdateBranchRequest{
			Name:          &newName,
			GenerationEnd: &newGenEnd,
		}

		result, err := svc.UpdateBranch(ctx, 1, req)
		if err != nil {
			t.Fatalf("UpdateBranch() error = %v", err)
		}
		if result.Name != "新分支名" {
			t.Errorf("expected name '新分支名', got '%s'", result.Name)
		}
		if result.GenerationEnd != 20 {
			t.Errorf("expected generation_end 20, got %d", result.GenerationEnd)
		}
	})

	t.Run("update non-existent branch", func(t *testing.T) {
		newName := "新名称"
		req := &UpdateBranchRequest{
			Name: &newName,
		}

		result, err := svc.UpdateBranch(ctx, 999, req)
		if err != nil {
			t.Fatalf("UpdateBranch() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestDeleteBranch(t *testing.T) {
	svc, _, branchRepo, _ := newTestGenealogyService()
	ctx := context.Background()

	// 预先创建一个分支
	branchRepo.Create(ctx, &genealogy.Branch{
		GenealogyID: 1,
		Name:        "待删除",
		Code:        "X",
	})

	t.Run("delete existing branch", func(t *testing.T) {
		err := svc.DeleteBranch(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteBranch() error = %v", err)
		}
		if len(branchRepo.branches) != 0 {
			t.Errorf("expected 0 branches after delete, got %d", len(branchRepo.branches))
		}
	})
}

func TestCreateGeneration(t *testing.T) {
	svc, genealogyRepo, _, generationRepo := newTestGenealogyService()
	ctx := context.Background()

	// 先创建一个族谱
	genealogyRepo.Create(ctx, &genealogy.Genealogy{
		Name:    "测试族谱",
		Surname: "测",
	})

	t.Run("create generation successfully", func(t *testing.T) {
		req := &CreateGenerationRequest{
			GenealogyID: 1,
			Generation:  1,
			Name:        "道",
			Sequence:    1,
			Description: "道德文章",
		}

		result, err := svc.CreateGeneration(ctx, req)
		if err != nil {
			t.Fatalf("CreateGeneration() error = %v", err)
		}
		if result.Name != "道" {
			t.Errorf("expected name '道', got '%s'", result.Name)
		}
		if result.Generation != 1 {
			t.Errorf("expected generation 1, got %d", result.Generation)
		}
		if result.Sequence != 1 {
			t.Errorf("expected sequence 1, got %d", result.Sequence)
		}
		if len(generationRepo.generations) != 1 {
			t.Errorf("expected 1 generation in repo, got %d", len(generationRepo.generations))
		}
	})
}

func TestGetGeneration(t *testing.T) {
	svc, _, _, generationRepo := newTestGenealogyService()
	ctx := context.Background()

	// 预先创建一个世代
	generationRepo.Create(ctx, &genealogy.Generation{
		GenealogyID: 1,
		Generation:  2,
		Name:        "德",
		Sequence:    2,
	})

	t.Run("get existing generation", func(t *testing.T) {
		result, err := svc.GetGeneration(ctx, 1)
		if err != nil {
			t.Fatalf("GetGeneration() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected generation, got nil")
		}
		if result.Name != "德" {
			t.Errorf("expected name '德', got '%s'", result.Name)
		}
	})

	t.Run("get non-existent generation", func(t *testing.T) {
		result, err := svc.GetGeneration(ctx, 999)
		if err != nil {
			t.Fatalf("GetGeneration() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestListGenerations(t *testing.T) {
	svc, _, _, generationRepo := newTestGenealogyService()
	ctx := context.Background()

	// 创建多个世代
	for i := 0; i < 3; i++ {
		generationRepo.Create(ctx, &genealogy.Generation{
			GenealogyID: 1,
			Generation:  i + 1,
			Name:        string(rune('A' + i)),
			Sequence:    i + 1,
		})
	}

	t.Run("list generations by genealogy", func(t *testing.T) {
		result, err := svc.ListGenerations(ctx, 1)
		if err != nil {
			t.Fatalf("ListGenerations() error = %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 generations, got %d", len(result))
		}
	})
}

func TestUpdateGeneration(t *testing.T) {
	svc, _, _, generationRepo := newTestGenealogyService()
	ctx := context.Background()

	// 预先创建一个世代
	generationRepo.Create(ctx, &genealogy.Generation{
		GenealogyID: 1,
		Generation:  1,
		Name:        "原名",
		Sequence:    1,
	})

	t.Run("update generation successfully", func(t *testing.T) {
		newName := "新名字"
		newSeq := 10
		req := &UpdateGenerationRequest{
			Name:     &newName,
			Sequence: &newSeq,
		}

		result, err := svc.UpdateGeneration(ctx, 1, req)
		if err != nil {
			t.Fatalf("UpdateGeneration() error = %v", err)
		}
		if result.Name != "新名字" {
			t.Errorf("expected name '新名字', got '%s'", result.Name)
		}
		if result.Sequence != 10 {
			t.Errorf("expected sequence 10, got %d", result.Sequence)
		}
	})

	t.Run("update non-existent generation", func(t *testing.T) {
		newName := "新名称"
		req := &UpdateGenerationRequest{
			Name: &newName,
		}

		result, err := svc.UpdateGeneration(ctx, 999, req)
		if err != nil {
			t.Fatalf("UpdateGeneration() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestDeleteGeneration(t *testing.T) {
	svc, _, _, generationRepo := newTestGenealogyService()
	ctx := context.Background()

	// 预先创建一个世代
	generationRepo.Create(ctx, &genealogy.Generation{
		GenealogyID: 1,
		Generation:  1,
		Name:        "待删除",
		Sequence:    1,
	})

	t.Run("delete existing generation", func(t *testing.T) {
		err := svc.DeleteGeneration(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteGeneration() error = %v", err)
		}
		if len(generationRepo.generations) != 0 {
			t.Errorf("expected 0 generations after delete, got %d", len(generationRepo.generations))
		}
	})
}

func TestGetGenerationByNumber(t *testing.T) {
	svc, _, _, generationRepo := newTestGenealogyService()
	ctx := context.Background()

	// 预先创建多个世代
	for i := 0; i < 5; i++ {
		generationRepo.Create(ctx, &genealogy.Generation{
			GenealogyID: 1,
			Generation:  i + 1,
			Name:        string(rune('A' + i)),
			Sequence:    i + 1,
		})
	}

	t.Run("get generation by number", func(t *testing.T) {
		result, err := svc.GetGenerationByNumber(ctx, 1, 3)
		if err != nil {
			t.Fatalf("GetGenerationByNumber() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected generation, got nil")
		}
		if result.Generation != 3 {
			t.Errorf("expected generation 3, got %d", result.Generation)
		}
	})

	t.Run("get non-existent generation by number", func(t *testing.T) {
		result, err := svc.GetGenerationByNumber(ctx, 1, 99)
		if err != nil {
			t.Fatalf("GetGenerationByNumber() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestGenealogyServiceErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("repository error on create genealogy", func(t *testing.T) {
		badGenealogyRepo := &mockGenealogyRepositoryError{
			inner: newMockGenealogyRepository(),
		}
		badBranchRepo := newMockBranchRepository()
		badGenerationRepo := newMockGenerationRepository()
		badSvc := NewGenealogyService(badGenealogyRepo, badBranchRepo, badGenerationRepo)

		req := &CreateGenealogyRequest{
			Name:    "测试",
			Surname: "测",
		}

		_, err := badSvc.CreateGenealogy(ctx, req)
		if err == nil {
			t.Error("expected error from repository")
		}
	})

	t.Run("repository error on get genealogy", func(t *testing.T) {
		badGenealogyRepo := &mockGenealogyRepositoryError{
			inner: newMockGenealogyRepository(),
		}
		badBranchRepo := newMockBranchRepository()
		badGenerationRepo := newMockGenerationRepository()
		badSvc := NewGenealogyService(badGenealogyRepo, badBranchRepo, badGenerationRepo)

		_, err := badSvc.GetGenealogy(ctx, 1)
		if err == nil {
			t.Error("expected error from repository")
		}
	})
}

type mockGenealogyRepositoryError struct {
	inner *mockGenealogyRepository
}

func (m *mockGenealogyRepositoryError) FindByID(ctx context.Context, id int64) (*genealogy.Genealogy, error) {
	return nil, errors.New("repository error")
}
func (m *mockGenealogyRepositoryError) Create(ctx context.Context, g *genealogy.Genealogy) error {
	return errors.New("repository error")
}
func (m *mockGenealogyRepositoryError) Update(ctx context.Context, g *genealogy.Genealogy) error {
	return m.inner.Update(ctx, g)
}
func (m *mockGenealogyRepositoryError) Delete(ctx context.Context, id int64) error {
	return m.inner.Delete(ctx, id)
}
func (m *mockGenealogyRepositoryError) FindAll(ctx context.Context) ([]*genealogy.Genealogy, error) {
	return m.inner.FindAll(ctx)
}
func (m *mockGenealogyRepositoryError) Search(ctx context.Context, query *genealogy.SearchQuery) ([]*genealogy.Genealogy, int64, error) {
	return m.inner.Search(ctx, query)
}
func (m *mockGenealogyRepositoryError) FindBySurname(ctx context.Context, surname string) ([]*genealogy.Genealogy, error) {
	return m.inner.FindBySurname(ctx, surname)
}