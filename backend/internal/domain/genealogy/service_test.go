package genealogy

import (
	"context"
	"errors"
	"testing"
)

// MockRepository 模拟族谱仓储
type MockRepository struct {
	genealogies map[int64]*Genealogy
	nextID      int64
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		genealogies: make(map[int64]*Genealogy),
		nextID:      1,
	}
}

func (m *MockRepository) FindByID(ctx context.Context, id int64) (*Genealogy, error) {
	if g, ok := m.genealogies[id]; ok {
		return g, nil
	}
	return nil, nil
}

func (m *MockRepository) Create(ctx context.Context, genealogy *Genealogy) error {
	genealogy.ID = m.nextID
	m.nextID++
	m.genealogies[genealogy.ID] = genealogy
	return nil
}

func (m *MockRepository) Update(ctx context.Context, genealogy *Genealogy) error {
	if _, ok := m.genealogies[genealogy.ID]; !ok {
		return errors.New("genealogy not found")
	}
	m.genealogies[genealogy.ID] = genealogy
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	delete(m.genealogies, id)
	return nil
}

func (m *MockRepository) FindAll(ctx context.Context) ([]*Genealogy, error) {
	result := make([]*Genealogy, 0, len(m.genealogies))
	for _, g := range m.genealogies {
		result = append(result, g)
	}
	return result, nil
}

func (m *MockRepository) Search(ctx context.Context, query *SearchQuery) ([]*Genealogy, int64, error) {
	result := make([]*Genealogy, 0, len(m.genealogies))
	for _, g := range m.genealogies {
		if query.Name != "" && g.Name != query.Name {
			continue
		}
		if query.Surname != "" && g.Surname != query.Surname {
			continue
		}
		result = append(result, g)
	}
	return result, int64(len(result)), nil
}

func (m *MockRepository) FindBySurname(ctx context.Context, surname string) ([]*Genealogy, error) {
	result := make([]*Genealogy, 0)
	for _, g := range m.genealogies {
		if g.Surname == surname {
			result = append(result, g)
		}
	}
	return result, nil
}

// MockBranchRepository 模拟分支仓储
type MockBranchRepository struct {
	branches map[int64]*Branch
	nextID   int64
}

func NewMockBranchRepository() *MockBranchRepository {
	return &MockBranchRepository{
		branches: make(map[int64]*Branch),
		nextID:   1,
	}
}

func (m *MockBranchRepository) FindByID(ctx context.Context, id int64) (*Branch, error) {
	if b, ok := m.branches[id]; ok {
		return b, nil
	}
	return nil, nil
}

func (m *MockBranchRepository) FindByGenealogyID(ctx context.Context, genealogyID int64) ([]*Branch, error) {
	result := make([]*Branch, 0)
	for _, b := range m.branches {
		if b.GenealogyID == genealogyID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (m *MockBranchRepository) Create(ctx context.Context, branch *Branch) error {
	branch.ID = m.nextID
	m.nextID++
	m.branches[branch.ID] = branch
	return nil
}

func (m *MockBranchRepository) Update(ctx context.Context, branch *Branch) error {
	if _, ok := m.branches[branch.ID]; !ok {
		return errors.New("branch not found")
	}
	m.branches[branch.ID] = branch
	return nil
}

func (m *MockBranchRepository) Delete(ctx context.Context, id int64) error {
	delete(m.branches, id)
	return nil
}

func (m *MockBranchRepository) DeleteByGenealogyID(ctx context.Context, genealogyID int64) error {
	for id, b := range m.branches {
		if b.GenealogyID == genealogyID {
			delete(m.branches, id)
		}
	}
	return nil
}

func (m *MockBranchRepository) Search(ctx context.Context, query *BranchSearchQuery) ([]*Branch, int64, error) {
	result := make([]*Branch, 0)
	for _, b := range m.branches {
		if query.GenealogyID != nil && b.GenealogyID != *query.GenealogyID {
			continue
		}
		if query.Name != "" && b.Name != query.Name {
			continue
		}
		result = append(result, b)
	}
	return result, int64(len(result)), nil
}

// MockGenerationRepository 模拟世代仓储
type MockGenerationRepository struct {
	generations map[int64]*Generation
	nextID      int64
}

func NewMockGenerationRepository() *MockGenerationRepository {
	return &MockGenerationRepository{
		generations: make(map[int64]*Generation),
		nextID:      1,
	}
}

func (m *MockGenerationRepository) FindByID(ctx context.Context, id int64) (*Generation, error) {
	if g, ok := m.generations[id]; ok {
		return g, nil
	}
	return nil, nil
}

func (m *MockGenerationRepository) FindByGenealogyID(ctx context.Context, genealogyID int64) ([]*Generation, error) {
	result := make([]*Generation, 0)
	for _, g := range m.generations {
		if g.GenealogyID == genealogyID {
			result = append(result, g)
		}
	}
	return result, nil
}

func (m *MockGenerationRepository) FindByGeneration(ctx context.Context, genealogyID int64, generation int) (*Generation, error) {
	for _, g := range m.generations {
		if g.GenealogyID == genealogyID && g.Generation == generation {
			return g, nil
		}
	}
	return nil, nil
}

func (m *MockGenerationRepository) Create(ctx context.Context, generation *Generation) error {
	generation.ID = m.nextID
	m.nextID++
	m.generations[generation.ID] = generation
	return nil
}

func (m *MockGenerationRepository) Update(ctx context.Context, generation *Generation) error {
	if _, ok := m.generations[generation.ID]; !ok {
		return errors.New("generation not found")
	}
	m.generations[generation.ID] = generation
	return nil
}

func (m *MockGenerationRepository) Delete(ctx context.Context, id int64) error {
	delete(m.generations, id)
	return nil
}

func (m *MockGenerationRepository) DeleteByGenealogyID(ctx context.Context, genealogyID int64) error {
	for id, g := range m.generations {
		if g.GenealogyID == genealogyID {
			delete(m.generations, id)
		}
	}
	return nil
}

// Tests

func TestCreateGenealogy(t *testing.T) {
	repo := NewMockRepository()
	branchRepo := NewMockBranchRepository()
	generationRepo := NewMockGenerationRepository()
	svc := NewService(repo, branchRepo, generationRepo)

	tests := []struct {
		name    string
		input   *Genealogy
		wantErr bool
	}{
		{
			name: "valid genealogy",
			input: &Genealogy{
				Name:             "王氏族谱",
				Surname:          "王",
				TotalGenerations: 22,
				TotalMembers:     5950,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: &Genealogy{
				Surname: "王",
			},
			wantErr: true,
		},
		{
			name: "empty surname",
			input: &Genealogy{
				Name: "王氏族谱",
			},
			wantErr: true,
		},
		{
			name: "negative generations",
			input: &Genealogy{
				Name:             "王氏族谱",
				Surname:          "王",
				TotalGenerations: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreateGenealogy(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGenealogy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetGenealogy(t *testing.T) {
	repo := NewMockRepository()
	branchRepo := NewMockBranchRepository()
	generationRepo := NewMockGenerationRepository()
	svc := NewService(repo, branchRepo, generationRepo)

	// Create a genealogy first
	genealogy := &Genealogy{
		Name:             "李氏族谱",
		Surname:          "李",
		TotalGenerations: 15,
		TotalMembers:     3000,
	}
	_ = svc.CreateGenealogy(context.Background(), genealogy)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "existing genealogy",
			id:      1,
			wantErr: false,
		},
		{
			name:    "non-existing genealogy",
			id:      999,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetGenealogy(context.Background(), tt.id)
			if err != nil {
				t.Errorf("GetGenealogy() error = %v", err)
				return
			}
			if tt.id == 1 && result == nil {
				t.Error("GetGenealogy() expected genealogy, got nil")
			}
			if tt.id == 999 && result != nil {
				t.Errorf("GetGenealogy() expected nil, got %v", result)
			}
		})
	}
}

func TestDeleteGenealogy(t *testing.T) {
	repo := NewMockRepository()
	branchRepo := NewMockBranchRepository()
	generationRepo := NewMockGenerationRepository()
	svc := NewService(repo, branchRepo, generationRepo)

	// Create genealogies
	genealogy1 := &Genealogy{Name: "赵氏族谱", Surname: "赵"}
	genealogy2 := &Genealogy{Name: "孙氏族谱", Surname: "孙"}
	_ = svc.CreateGenealogy(context.Background(), genealogy1)
	_ = svc.CreateGenealogy(context.Background(), genealogy2)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "delete existing",
			id:      1,
			wantErr: false,
		},
		{
			name:    "delete non-existing",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.DeleteGenealogy(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteGenealogy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateBranch(t *testing.T) {
	repo := NewMockRepository()
	branchRepo := NewMockBranchRepository()
	generationRepo := NewMockGenerationRepository()
	svc := NewService(repo, branchRepo, generationRepo)

	// Create a genealogy first
	genealogy := &Genealogy{Name: "周氏族谱", Surname: "周"}
	_ = svc.CreateGenealogy(context.Background(), genealogy)

	tests := []struct {
		name    string
		branch  *Branch
		wantErr bool
	}{
		{
			name: "valid branch",
			branch: &Branch{
				GenealogyID:      1,
				Name:             "长房",
				Code:             "A",
				GenerationStart:  1,
				GenerationEnd:    10,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			branch: &Branch{
				GenealogyID: 1,
			},
			wantErr: true,
		},
		{
			name: "non-existing genealogy",
			branch: &Branch{
				GenealogyID: 999,
				Name:        "测试分支",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreateBranch(context.Background(), tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateGeneration(t *testing.T) {
	repo := NewMockRepository()
	branchRepo := NewMockBranchRepository()
	generationRepo := NewMockGenerationRepository()
	svc := NewService(repo, branchRepo, generationRepo)

	// Create a genealogy first
	genealogy := &Genealogy{Name: "吴氏族谱", Surname: "吴"}
	_ = svc.CreateGenealogy(context.Background(), genealogy)

	startYear := 1900
	endYear := 1950

	tests := []struct {
		name       string
		generation *Generation
		wantErr    bool
	}{
		{
			name: "valid generation",
			generation: &Generation{
				GenealogyID: 1,
				Generation:  1,
				Name:        "道",
				Sequence:    1,
			},
			wantErr: false,
		},
		{
			name: "generation with year range",
			generation: &Generation{
				GenealogyID: 1,
				Generation:  2,
				Name:        "德",
				Sequence:    2,
				StartYear:   &startYear,
				EndYear:     &endYear,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			generation: &Generation{
				GenealogyID: 1,
				Generation:  3,
			},
			wantErr: true,
		},
		{
			name: "non-existing genealogy",
			generation: &Generation{
				GenealogyID: 999,
				Generation:  1,
				Name:        "测试",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreateGeneration(context.Background(), tt.generation)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGeneration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenealogyValidate(t *testing.T) {
	tests := []struct {
		name    string
		g       *Genealogy
		wantErr bool
	}{
		{
			name: "valid",
			g: &Genealogy{
				Name:             "陈氏族谱",
				Surname:          "陈",
				TotalGenerations: 20,
				TotalMembers:     5000,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			g: &Genealogy{
				Surname: "陈",
			},
			wantErr: true,
		},
		{
			name: "empty surname",
			g: &Genealogy{
				Name: "陈氏族谱",
			},
			wantErr: true,
		},
		{
			name: "negative generations",
			g: &Genealogy{
				Name:             "陈氏族谱",
				Surname:          "陈",
				TotalGenerations: -5,
			},
			wantErr: true,
		},
		{
			name: "negative members",
			g: &Genealogy{
				Name:          "陈氏族谱",
				Surname:       "陈",
				TotalMembers:  -100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.g.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Genealogy.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBranchValidate(t *testing.T) {
	tests := []struct {
		name    string
		b       *Branch
		wantErr bool
	}{
		{
			name: "valid branch",
			b: &Branch{
				GenealogyID:     1,
				Name:            "长房",
				GenerationStart: 1,
				GenerationEnd:   10,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			b: &Branch{
				GenealogyID: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid genealogy id",
			b: &Branch{
				Name: "测试",
			},
			wantErr: true,
		},
		{
			name: "invalid generation range",
			b: &Branch{
				GenealogyID:     1,
				Name:            "测试",
				GenerationStart: 10,
				GenerationEnd:   5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.b.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Branch.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerationValidate(t *testing.T) {
	tests := []struct {
		name    string
		g       *Generation
		wantErr bool
	}{
		{
			name: "valid generation",
			g: &Generation{
				GenealogyID: 1,
				Generation:  1,
				Name:        "道",
				Sequence:    1,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			g: &Generation{
				GenealogyID: 1,
				Generation:  1,
			},
			wantErr: true,
		},
		{
			name: "invalid genealogy id",
			g: &Generation{
				Generation: 1,
				Name:       "道",
			},
			wantErr: true,
		},
		{
			name: "negative generation",
			g: &Generation{
				GenealogyID: 1,
				Generation:  -1,
				Name:        "道",
			},
			wantErr: true,
		},
		{
			name: "negative sequence",
			g: &Generation{
				GenealogyID: 1,
				Generation:  1,
				Name:        "道",
				Sequence:    -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.g.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generation.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerationIsActive(t *testing.T) {
	endYear := 1950

	tests := []struct {
		name     string
		g        *Generation
		expected bool
	}{
		{
			name: "no end year",
			g: &Generation{
				Generation: 1,
				Name:       "道",
			},
			expected: true,
		},
		{
			name: "with end year",
			g: &Generation{
				Generation: 2,
				Name:       "德",
				EndYear:    &endYear,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.IsActive(); got != tt.expected {
				t.Errorf("Generation.IsActive() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetGenealogyWithDetails(t *testing.T) {
	repo := NewMockRepository()
	branchRepo := NewMockBranchRepository()
	generationRepo := NewMockGenerationRepository()
	svc := NewService(repo, branchRepo, generationRepo)

	// Create genealogy with branches and generations
	genealogy := &Genealogy{Name: "郑氏族谱", Surname: "郑"}
	_ = svc.CreateGenealogy(context.Background(), genealogy)

	branch := &Branch{GenealogyID: 1, Name: "长房", Code: "A"}
	_ = svc.CreateBranch(context.Background(), branch)

	generation := &Generation{GenealogyID: 1, Generation: 1, Name: "道", Sequence: 1}
	_ = svc.CreateGeneration(context.Background(), generation)

	result, err := svc.GetGenealogyWithDetails(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetGenealogyWithDetails() error = %v", err)
	}

	if result == nil {
		t.Fatal("GetGenealogyWithDetails() returned nil")
	}

	if len(result.Branches) != 1 {
		t.Errorf("expected 1 branch, got %d", len(result.Branches))
	}

	if len(result.Generations) != 1 {
		t.Errorf("expected 1 generation, got %d", len(result.Generations))
	}
}