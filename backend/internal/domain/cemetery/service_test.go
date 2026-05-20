package cemetery

import (
	"context"
	"errors"
	"testing"
)

// MockRepository 模拟墓园仓储
type MockRepository struct {
	cemeteries map[int64]*Cemetery
	nextID    int64
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		cemeteries: make(map[int64]*Cemetery),
		nextID:    1,
	}
}

func (m *MockRepository) FindByID(ctx context.Context, id int64) (*Cemetery, error) {
	if c, ok := m.cemeteries[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *MockRepository) Create(ctx context.Context, cemetery *Cemetery) error {
	cemetery.ID = m.nextID
	m.nextID++
	m.cemeteries[cemetery.ID] = cemetery
	return nil
}

func (m *MockRepository) Update(ctx context.Context, cemetery *Cemetery) error {
	if _, ok := m.cemeteries[cemetery.ID]; !ok {
		return errors.New("cemetery not found")
	}
	m.cemeteries[cemetery.ID] = cemetery
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	delete(m.cemeteries, id)
	return nil
}

func (m *MockRepository) FindAll(ctx context.Context) ([]*Cemetery, error) {
	result := make([]*Cemetery, 0, len(m.cemeteries))
	for _, c := range m.cemeteries {
		result = append(result, c)
	}
	return result, nil
}

func (m *MockRepository) Search(ctx context.Context, query *SearchQuery) ([]*Cemetery, int64, error) {
	result := make([]*Cemetery, 0, len(m.cemeteries))
	for _, c := range m.cemeteries {
		if query.Name != "" && c.Name != query.Name {
			continue
		}
		result = append(result, c)
	}
	return result, int64(len(result)), nil
}

func (m *MockRepository) FindByRegion(ctx context.Context, province, city string) ([]*Cemetery, error) {
	result := make([]*Cemetery, 0)
	for _, c := range m.cemeteries {
		if c.Province == province && c.City == city {
			result = append(result, c)
		}
	}
	return result, nil
}

// MockGraveRepository 模拟墓位仓储
type MockGraveRepository struct {
	graves  map[int64]*Grave
	nextID  int64
}

func NewMockGraveRepository() *MockGraveRepository {
	return &MockGraveRepository{
		graves: make(map[int64]*Grave),
		nextID: 1,
	}
}

func (m *MockGraveRepository) FindByID(ctx context.Context, id int64) (*Grave, error) {
	if g, ok := m.graves[id]; ok {
		return g, nil
	}
	return nil, nil
}

func (m *MockGraveRepository) FindByCemeteryID(ctx context.Context, cemeteryID int64) ([]*Grave, error) {
	result := make([]*Grave, 0)
	for _, g := range m.graves {
		if g.CemeteryID == cemeteryID {
			result = append(result, g)
		}
	}
	return result, nil
}

func (m *MockGraveRepository) FindByPersonID(ctx context.Context, personID int64) (*Grave, error) {
	for _, g := range m.graves {
		if g.PersonID != nil && *g.PersonID == personID {
			return g, nil
		}
	}
	return nil, nil
}

func (m *MockGraveRepository) Create(ctx context.Context, grave *Grave) error {
	grave.ID = m.nextID
	m.nextID++
	m.graves[grave.ID] = grave
	return nil
}

func (m *MockGraveRepository) Update(ctx context.Context, grave *Grave) error {
	if _, ok := m.graves[grave.ID]; !ok {
		return errors.New("grave not found")
	}
	m.graves[grave.ID] = grave
	return nil
}

func (m *MockGraveRepository) Delete(ctx context.Context, id int64) error {
	delete(m.graves, id)
	return nil
}

func (m *MockGraveRepository) Search(ctx context.Context, query *GraveSearchQuery) ([]*Grave, int64, error) {
	result := make([]*Grave, 0)
	for _, g := range m.graves {
		if query.CemeteryID != nil && g.CemeteryID != *query.CemeteryID {
			continue
		}
		if query.PersonID != nil && (g.PersonID == nil || *g.PersonID != *query.PersonID) {
			continue
		}
		if query.Section != "" && g.Section != query.Section {
			continue
		}
		result = append(result, g)
	}
	return result, int64(len(result)), nil
}

// Test Cemetery entity
func TestCemetery_Validate(t *testing.T) {
	tests := []struct {
		name    string
		c       Cemetery
		wantErr bool
	}{
		{
			name: "valid cemetery",
			c: Cemetery{
				Name:      "测试墓园",
				Province:  "广东省",
				TotalGrave: 100,
				UsedGrave:  10,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			c: Cemetery{
				Name:      "",
				Province:  "广东省",
				TotalGrave: 100,
				UsedGrave:  10,
			},
			wantErr: true,
		},
		{
			name: "empty province",
			c: Cemetery{
				Name:      "测试墓园",
				Province:  "",
				TotalGrave: 100,
				UsedGrave:  10,
			},
			wantErr: true,
		},
		{
			name: "negative total",
			c: Cemetery{
				Name:      "测试墓园",
				Province:  "广东省",
				TotalGrave: -1,
				UsedGrave:  10,
			},
			wantErr: true,
		},
		{
			name: "used greater than total",
			c: Cemetery{
				Name:      "测试墓园",
				Province:  "广东省",
				TotalGrave: 10,
				UsedGrave:  100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cemetery.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCemetery_UsageRate(t *testing.T) {
	tests := []struct {
		name     string
		total    int
		used     int
		expected float64
	}{
		{"0%", 100, 0, 0},
		{"50%", 100, 50, 50},
		{"100%", 100, 100, 100},
		{"zero total", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cemetery{TotalGrave: tt.total, UsedGrave: tt.used}
			if got := c.UsageRate(); got != tt.expected {
				t.Errorf("Cemetery.UsageRate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCemetery_HasAvailableSpace(t *testing.T) {
	tests := []struct {
		name     string
		total    int
		used     int
		expected bool
	}{
		{"has space", 100, 50, true},
		{"no space", 100, 100, false},
		{"zero total", 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cemetery{TotalGrave: tt.total, UsedGrave: tt.used}
			if got := c.HasAvailableSpace(); got != tt.expected {
				t.Errorf("Cemetery.HasAvailableSpace() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test Grave entity
func TestGrave_Validate(t *testing.T) {
	tests := []struct {
		name    string
		g       Grave
		wantErr bool
	}{
		{
			name: "valid grave",
			g: Grave{
				CemeteryID: 1,
				Section:    "A区",
				Row:        1,
				Number:     1,
			},
			wantErr: false,
		},
		{
			name: "zero cemetery_id",
			g: Grave{
				CemeteryID: 0,
				Section:    "A区",
				Row:        1,
				Number:     1,
			},
			wantErr: true,
		},
		{
			name: "empty section",
			g: Grave{
				CemeteryID: 1,
				Section:    "",
				Row:        1,
				Number:     1,
			},
			wantErr: true,
		},
		{
			name: "zero row",
			g: Grave{
				CemeteryID: 1,
				Section:    "A区",
				Row:        0,
				Number:     1,
			},
			wantErr: true,
		},
		{
			name: "zero number",
			g: Grave{
				CemeteryID: 1,
				Section:    "A区",
				Row:        1,
				Number:     0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.g.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Grave.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGrave_Availability(t *testing.T) {
	g := Grave{Status: GraveStatusAvailable}
	if !g.Availability() {
		t.Error("Grave.Availability() = false, want true")
	}
}

func TestGrave_IsOccupied(t *testing.T) {
	personID := int64(1)
	g := Grave{
		Status:   GraveStatusOccupied,
		PersonID: &personID,
	}
	if !g.IsOccupied() {
		t.Error("Grave.IsOccupied() = false, want true")
	}
}

// Test Service
func TestService_CreateCemetery(t *testing.T) {
	repo := NewMockRepository()
	svc := NewService(repo, nil)

	cemetery := &Cemetery{
		Name:      "测试墓园",
		Province:  "广东省",
		City:      "广州市",
		TotalGrave: 100,
		UsedGrave:  10,
	}

	err := svc.CreateCemetery(context.Background(), cemetery)
	if err != nil {
		t.Errorf("CreateCemetery() error = %v", err)
	}
	if cemetery.ID == 0 {
		t.Error("CreateCemetery() did not set ID")
	}
}

func TestService_GetCemetery(t *testing.T) {
	repo := NewMockRepository()
	svc := NewService(repo, nil)

	// 创建墓园
	cemetery := &Cemetery{
		Name:      "测试墓园",
		Province:  "广东省",
		TotalGrave: 100,
	}
	repo.Create(context.Background(), cemetery)

	// 获取墓园
	got, err := svc.GetCemetery(context.Background(), cemetery.ID)
	if err != nil {
		t.Errorf("GetCemetery() error = %v", err)
	}
	if got == nil {
		t.Error("GetCemetery() returned nil")
	}
	if got.Name != cemetery.Name {
		t.Errorf("GetCemetery() name = %v, want %v", got.Name, cemetery.Name)
	}
}

func TestService_UpdateCemetery(t *testing.T) {
	repo := NewMockRepository()
	svc := NewService(repo, nil)

	// 创建墓园
	cemetery := &Cemetery{
		Name:      "测试墓园",
		Province:  "广东省",
		TotalGrave: 100,
	}
	repo.Create(context.Background(), cemetery)

	// 更新墓园
	cemetery.Name = "更新后的墓园"
	err := svc.UpdateCemetery(context.Background(), cemetery)
	if err != nil {
		t.Errorf("UpdateCemetery() error = %v", err)
	}

	// 验证更新
	got, _ := svc.GetCemetery(context.Background(), cemetery.ID)
	if got.Name != "更新后的墓园" {
		t.Errorf("UpdateCemetery() name = %v, want '更新后的墓园'", got.Name)
	}
}

func TestService_DeleteCemetery(t *testing.T) {
	repo := NewMockRepository()
	svc := NewService(repo, nil)

	// 创建墓园
	cemetery := &Cemetery{
		Name:      "测试墓园",
		Province:  "广东省",
		TotalGrave: 100,
	}
	repo.Create(context.Background(), cemetery)

	// 删除墓园
	err := svc.DeleteCemetery(context.Background(), cemetery.ID)
	if err != nil {
		t.Errorf("DeleteCemetery() error = %v", err)
	}

	// 验证删除
	got, _ := svc.GetCemetery(context.Background(), cemetery.ID)
	if got != nil {
		t.Error("DeleteCemetery() cemetery still exists")
	}
}

func TestService_CreateGrave(t *testing.T) {
	repo := NewMockRepository()
	graveRepo := NewMockGraveRepository()
	svc := NewService(repo, graveRepo)

	// 创建墓园
	cemetery := &Cemetery{
		Name:      "测试墓园",
		Province:  "广东省",
		TotalGrave: 100,
		UsedGrave:  10,
	}
	repo.Create(context.Background(), cemetery)

	// 创建墓位
	grave := &Grave{
		CemeteryID: cemetery.ID,
		Section:    "A区",
		Row:        1,
		Number:     1,
		Status:     GraveStatusAvailable,
	}

	err := svc.CreateGrave(context.Background(), grave)
	if err != nil {
		t.Errorf("CreateGrave() error = %v", err)
	}
	if grave.ID == 0 {
		t.Error("CreateGrave() did not set ID")
	}
}

func TestService_GetGraveByPersonID(t *testing.T) {
	repo := NewMockRepository()
	graveRepo := NewMockGraveRepository()
	svc := NewService(repo, graveRepo)

	personID := int64(100)

	// 创建墓园和墓位
	cemetery := &Cemetery{
		Name:      "测试墓园",
		Province:  "广东省",
		TotalGrave: 100,
	}
	repo.Create(context.Background(), cemetery)

	grave := &Grave{
		CemeteryID: cemetery.ID,
		PersonID:   &personID,
		Section:    "A区",
		Row:        1,
		Number:     1,
		Status:     GraveStatusOccupied,
	}
	graveRepo.Create(context.Background(), grave)

	// 根据人物ID查询墓位
	got, err := svc.GetGraveByPersonID(context.Background(), personID)
	if err != nil {
		t.Errorf("GetGraveByPersonID() error = %v", err)
	}
	if got == nil {
		t.Error("GetGraveByPersonID() returned nil")
	}
	if got.PersonID == nil || *got.PersonID != personID {
		t.Error("GetGraveByPersonID() returned wrong person")
	}
}

func TestService_SearchGraves(t *testing.T) {
	repo := NewMockRepository()
	graveRepo := NewMockGraveRepository()
	svc := NewService(repo, graveRepo)

	// 创建墓园
	cemetery := &Cemetery{
		Name:      "测试墓园",
		Province:  "广东省",
		TotalGrave: 100,
	}
	repo.Create(context.Background(), cemetery)

	// 创建多个墓位
	for i := 1; i <= 5; i++ {
		graveRepo.Create(context.Background(), &Grave{
			CemeteryID: cemetery.ID,
			Section:    "A区",
			Row:        i,
			Number:     i,
		})
	}

	// 搜索墓位
	query := &GraveSearchQuery{
		CemeteryID: &cemetery.ID,
		Section:    "A区",
		Page:       1,
		PageSize:   10,
	}

	graves, total, err := svc.SearchGraves(context.Background(), query)
	if err != nil {
		t.Errorf("SearchGraves() error = %v", err)
	}
	if total != 5 {
		t.Errorf("SearchGraves() total = %v, want 5", total)
	}
	if len(graves) != 5 {
		t.Errorf("SearchGraves() len = %v, want 5", len(graves))
	}
}