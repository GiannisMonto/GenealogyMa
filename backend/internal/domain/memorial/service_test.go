package memorial

import (
	"context"
	"errors"
	"testing"
)

// MockRepository 模拟宗祠仓储
type MockRepository struct {
	halls map[int64]*MemorialHall
	nextID int64
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		halls: make(map[int64]*MemorialHall),
		nextID: 1,
	}
}

func (m *MockRepository) FindByID(ctx context.Context, id int64) (*MemorialHall, error) {
	if h, ok := m.halls[id]; ok {
		return h, nil
	}
	return nil, nil
}

func (m *MockRepository) Create(ctx context.Context, hall *MemorialHall) error {
	hall.ID = m.nextID
	m.nextID++
	m.halls[hall.ID] = hall
	return nil
}

func (m *MockRepository) Update(ctx context.Context, hall *MemorialHall) error {
	if _, ok := m.halls[hall.ID]; !ok {
		return errors.New("memorial hall not found")
	}
	m.halls[hall.ID] = hall
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	delete(m.halls, id)
	return nil
}

func (m *MockRepository) FindAll(ctx context.Context) ([]*MemorialHall, error) {
	result := make([]*MemorialHall, 0, len(m.halls))
	for _, h := range m.halls {
		result = append(result, h)
	}
	return result, nil
}

func (m *MockRepository) Search(ctx context.Context, query *SearchQuery) ([]*MemorialHall, int64, error) {
	result := make([]*MemorialHall, 0, len(m.halls))
	for _, h := range m.halls {
		if query.Name != "" && h.Name != query.Name {
			continue
		}
		result = append(result, h)
	}
	return result, int64(len(result)), nil
}

func (m *MockRepository) FindByRegion(ctx context.Context, province, city string) ([]*MemorialHall, error) {
	result := make([]*MemorialHall, 0)
	for _, h := range m.halls {
		if h.Province == province && h.City == city {
			result = append(result, h)
		}
	}
	return result, nil
}

// MockTabletRepository 模拟牌位仓储
type MockTabletRepository struct {
	tablets map[int64]*MemorialTablet
	nextID  int64
}

func NewMockTabletRepository() *MockTabletRepository {
	return &MockTabletRepository{
		tablets: make(map[int64]*MemorialTablet),
		nextID:  1,
	}
}

func (m *MockTabletRepository) FindByID(ctx context.Context, id int64) (*MemorialTablet, error) {
	if t, ok := m.tablets[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *MockTabletRepository) FindByHallID(ctx context.Context, hallID int64) ([]*MemorialTablet, error) {
	result := make([]*MemorialTablet, 0)
	for _, t := range m.tablets {
		if t.HallID == hallID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *MockTabletRepository) FindByPersonID(ctx context.Context, personID int64) (*MemorialTablet, error) {
	for _, t := range m.tablets {
		if t.PersonID != nil && *t.PersonID == personID {
			return t, nil
		}
	}
	return nil, nil
}

func (m *MockTabletRepository) Create(ctx context.Context, tablet *MemorialTablet) error {
	tablet.ID = m.nextID
	m.nextID++
	m.tablets[tablet.ID] = tablet
	return nil
}

func (m *MockTabletRepository) Update(ctx context.Context, tablet *MemorialTablet) error {
	if _, ok := m.tablets[tablet.ID]; !ok {
		return errors.New("memorial tablet not found")
	}
	m.tablets[tablet.ID] = tablet
	return nil
}

func (m *MockTabletRepository) Delete(ctx context.Context, id int64) error {
	delete(m.tablets, id)
	return nil
}

func (m *MockTabletRepository) DeleteByHallID(ctx context.Context, hallID int64) error {
	for id, t := range m.tablets {
		if t.HallID == hallID {
			delete(m.tablets, id)
		}
	}
	return nil
}

func (m *MockTabletRepository) Search(ctx context.Context, query *TabletSearchQuery) ([]*MemorialTablet, int64, error) {
	result := make([]*MemorialTablet, 0)
	for _, t := range m.tablets {
		if query.HallID != nil && t.HallID != *query.HallID {
			continue
		}
		if query.PersonID != nil && (t.PersonID == nil || *t.PersonID != *query.PersonID) {
			continue
		}
		if query.TabletType != nil && t.TabletType != *query.TabletType {
			continue
		}
		result = append(result, t)
	}
	return result, int64(len(result)), nil
}

// Tests

func TestCreateHall(t *testing.T) {
	repo := NewMockRepository()
	tabletRepo := NewMockTabletRepository()
	svc := NewService(repo, tabletRepo)

	tests := []struct {
		name    string
		input   *MemorialHall
		wantErr bool
	}{
		{
			name: "valid hall",
			input: &MemorialHall{
				Name:        "王氏宗祠",
				Province:    "四川",
				City:        "成都",
				TotalTablet: 100,
				UsedTablet:  0,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: &MemorialHall{
				Province: "四川",
				City:     "成都",
			},
			wantErr: true,
		},
		{
			name: "empty province",
			input: &MemorialHall{
				Name: "王氏宗祠",
				City: "成都",
			},
			wantErr: true,
		},
		{
			name: "negative total",
			input: &MemorialHall{
				Name:       "王氏宗祠",
				Province:   "四川",
				TotalTablet: -1,
			},
			wantErr: true,
		},
		{
			name: "used exceeds total",
			input: &MemorialHall{
				Name:        "王氏宗祠",
				Province:    "四川",
				TotalTablet: 10,
				UsedTablet:  20,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreateHall(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateHall() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetHall(t *testing.T) {
	repo := NewMockRepository()
	tabletRepo := NewMockTabletRepository()
	svc := NewService(repo, tabletRepo)

	hall := &MemorialHall{Name: "李氏宗祠", Province: "湖南", City: "长沙", TotalTablet: 50}
	_ = svc.CreateHall(context.Background(), hall)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "existing hall",
			id:      1,
			wantErr: false,
		},
		{
			name:    "non-existing hall",
			id:      999,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetHall(context.Background(), tt.id)
			if err != nil {
				t.Errorf("GetHall() error = %v", err)
				return
			}
			if tt.id == 1 && result == nil {
				t.Error("GetHall() expected hall, got nil")
			}
			if tt.id == 999 && result != nil {
				t.Errorf("GetHall() expected nil, got %v", result)
			}
		})
	}
}

func TestDeleteHall(t *testing.T) {
	repo := NewMockRepository()
	tabletRepo := NewMockTabletRepository()
	svc := NewService(repo, tabletRepo)

	hall1 := &MemorialHall{Name: "赵氏宗祠", Province: "广东", TotalTablet: 30}
	hall2 := &MemorialHall{Name: "孙氏宗祠", Province: "广西", TotalTablet: 40}
	_ = svc.CreateHall(context.Background(), hall1)
	_ = svc.CreateHall(context.Background(), hall2)

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
			err := svc.DeleteHall(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteHall() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateTablet(t *testing.T) {
	repo := NewMockRepository()
	tabletRepo := NewMockTabletRepository()
	svc := NewService(repo, tabletRepo)

	// Create a hall first
	hall := &MemorialHall{Name: "周氏宗祠", Province: "福建", TotalTablet: 100, UsedTablet: 0}
	_ = svc.CreateHall(context.Background(), hall)

	personID := int64(100)

	tests := []struct {
		name    string
		tablet  *MemorialTablet
		wantErr bool
	}{
		{
			name: "valid tablet with person_id",
			tablet: &MemorialTablet{
				HallID:     1,
				PersonID:   &personID,
				PersonName: "周总理",
				Floor:      1,
				Row:        1,
				Number:     1,
			},
			wantErr: false,
		},
		{
			name: "valid tablet with person_name only",
			tablet: &MemorialTablet{
				HallID:     1,
				PersonName: "周将军",
				Floor:      1,
				Row:        1,
				Number:     2,
			},
			wantErr: false,
		},
		{
			name: "empty name and no person_id",
			tablet: &MemorialTablet{
				HallID: 1,
				Floor:  1,
				Row:    1,
				Number: 3,
			},
			wantErr: true,
		},
		{
			name: "non-existing hall",
			tablet: &MemorialTablet{
				HallID:     999,
				PersonName: "测试",
				Floor:      1,
				Row:        1,
				Number:     1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreateTablet(context.Background(), tt.tablet)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTablet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTabletByPersonID(t *testing.T) {
	repo := NewMockRepository()
	tabletRepo := NewMockTabletRepository()
	svc := NewService(repo, tabletRepo)

	// Create hall and tablet
	hall := &MemorialHall{Name: "吴氏宗祠", Province: "浙江", TotalTablet: 50}
	_ = svc.CreateHall(context.Background(), hall)

	personID := int64(200)
	tablet := &MemorialTablet{
		HallID:     1,
		PersonID:   &personID,
		PersonName: "吴先生",
		Floor:      1,
		Row:        2,
		Number:     1,
	}
	_ = svc.CreateTablet(context.Background(), tablet)

	result, err := svc.GetTabletByPersonID(context.Background(), 200)
	if err != nil {
		t.Fatalf("GetTabletByPersonID() error = %v", err)
	}
	if result == nil {
		t.Fatal("GetTabletByPersonID() returned nil")
	}
	if result.PersonName != "吴先生" {
		t.Errorf("expected person name '吴先生', got '%s'", result.PersonName)
	}
}

func TestMemorialHallUsageRate(t *testing.T) {
	tests := []struct {
		name     string
		hall     *MemorialHall
		expected float64
	}{
		{
			name: "50% usage",
			hall: &MemorialHall{
				Name:        "测试宗祠",
				TotalTablet: 100,
				UsedTablet:  50,
			},
			expected: 50.0,
		},
		{
			name: "0% usage",
			hall: &MemorialHall{
				Name:        "空宗祠",
				TotalTablet: 100,
				UsedTablet:  0,
			},
			expected: 0.0,
		},
		{
			name: "100% usage",
			hall: &MemorialHall{
				Name:        "满宗祠",
				TotalTablet: 100,
				UsedTablet:  100,
			},
			expected: 100.0,
		},
		{
			name: "zero total",
			hall: &MemorialHall{
				Name:        "零容量",
				TotalTablet: 0,
				UsedTablet:  0,
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hall.UsageRate(); got != tt.expected {
				t.Errorf("MemorialHall.UsageRate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMemorialTabletGetPositionDesc(t *testing.T) {
	tablet := &MemorialTablet{
		Floor:  2,
		Row:    3,
		Number: 5,
	}

	expected := "2层 3排 5号"
	if got := tablet.GetPositionDesc(); got != expected {
		t.Errorf("GetPositionDesc() = %v, want %v", got, expected)
	}
}

func TestMemorialTabletIsOccupied(t *testing.T) {
	personID := int64(100)

	tests := []struct {
		name     string
		tablet   *MemorialTablet
		expected bool
	}{
		{
			name: "occupied with person_id",
			tablet: &MemorialTablet{
				PersonID: &personID,
			},
			expected: true,
		},
		{
			name: "not occupied",
			tablet: &MemorialTablet{
				PersonName: "无名氏",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tablet.IsOccupied(); got != tt.expected {
				t.Errorf("IsOccupied() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetHallWithTablets(t *testing.T) {
	repo := NewMockRepository()
	tabletRepo := NewMockTabletRepository()
	svc := NewService(repo, tabletRepo)

	// Create hall with tablets
	hall := &MemorialHall{Name: "郑氏宗祠", Province: "江西", TotalTablet: 100}
	_ = svc.CreateHall(context.Background(), hall)

	tablet1 := &MemorialTablet{HallID: 1, PersonName: "郑一", Floor: 1, Row: 1, Number: 1}
	tablet2 := &MemorialTablet{HallID: 1, PersonName: "郑二", Floor: 1, Row: 1, Number: 2}
	_ = svc.CreateTablet(context.Background(), tablet1)
	_ = svc.CreateTablet(context.Background(), tablet2)

	result, err := svc.GetHallWithTablets(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetHallWithTablets() error = %v", err)
	}

	if result == nil {
		t.Fatal("GetHallWithTablets() returned nil")
	}

	if len(result.Tables) != 2 {
		t.Errorf("expected 2 tablets, got %d", len(result.Tables))
	}
}

func TestSearchHalls(t *testing.T) {
	repo := NewMockRepository()
	tabletRepo := NewMockTabletRepository()
	svc := NewService(repo, tabletRepo)

	// Create halls - use exact name match for testing
	hall1 := &MemorialHall{Name: "王氏宗祠", Province: "四川", City: "成都"}
	hall2 := &MemorialHall{Name: "王氏宗祠", Province: "四川", City: "绵阳"}
	hall3 := &MemorialHall{Name: "李氏宗祠", Province: "湖南", City: "长沙"}
	_ = svc.CreateHall(context.Background(), hall1)
	_ = svc.CreateHall(context.Background(), hall2)
	_ = svc.CreateHall(context.Background(), hall3)

	result, total, err := svc.SearchHalls(context.Background(), &SearchQuery{Name: "王氏宗祠"})
	if err != nil {
		t.Fatalf("SearchHalls() error = %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
}

func TestSearchTablets(t *testing.T) {
	repo := NewMockRepository()
	tabletRepo := NewMockTabletRepository()
	svc := NewService(repo, tabletRepo)

	// Create hall and tablets
	hall := &MemorialHall{Name: "孙氏宗祠", Province: "山东", TotalTablet: 100}
	_ = svc.CreateHall(context.Background(), hall)

	personID1 := int64(1)
	personID2 := int64(2)
	tablet1 := &MemorialTablet{HallID: 1, PersonID: &personID1, PersonName: "孙武", Floor: 1, Row: 1, Number: 1, TabletType: TabletTypeAncestor}
	tablet2 := &MemorialTablet{HallID: 1, PersonID: &personID2, PersonName: "孙膑", Floor: 1, Row: 1, Number: 2, TabletType: TabletTypeAncestor}
	tablet3 := &MemorialTablet{HallID: 1, PersonName: "孙权", Floor: 2, Row: 1, Number: 1, TabletType: TabletTypeFounder}
	_ = svc.CreateTablet(context.Background(), tablet1)
	_ = svc.CreateTablet(context.Background(), tablet2)
	_ = svc.CreateTablet(context.Background(), tablet3)

	hallID := int64(1)
	ancestorType := TabletTypeAncestor
	results, total, err := svc.SearchTablets(context.Background(), &TabletSearchQuery{HallID: &hallID, TabletType: &ancestorType})
	if err != nil {
		t.Fatalf("SearchTablets() error = %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}