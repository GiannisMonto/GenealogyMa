package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/memorial"
)

// mockMemorialRepository 模拟宗祠仓储
type mockMemorialRepository struct {
	halls map[int64]*memorial.MemorialHall
	nextID int64
}

func newMockMemorialRepository() *mockMemorialRepository {
	return &mockMemorialRepository{
		halls: make(map[int64]*memorial.MemorialHall),
		nextID: 1,
	}
}

func (m *mockMemorialRepository) FindByID(ctx context.Context, id int64) (*memorial.MemorialHall, error) {
	h, ok := m.halls[id]
	if !ok {
		return nil, nil
	}
	return h, nil
}

func (m *mockMemorialRepository) Create(ctx context.Context, h *memorial.MemorialHall) error {
	h.ID = m.nextID
	m.nextID++
	h.CreatedAt = time.Now()
	h.UpdatedAt = time.Now()
	m.halls[h.ID] = h
	return nil
}

func (m *mockMemorialRepository) Update(ctx context.Context, h *memorial.MemorialHall) error {
	h.UpdatedAt = time.Now()
	m.halls[h.ID] = h
	return nil
}

func (m *mockMemorialRepository) Delete(ctx context.Context, id int64) error {
	delete(m.halls, id)
	return nil
}

func (m *mockMemorialRepository) FindAll(ctx context.Context) ([]*memorial.MemorialHall, error) {
	var result []*memorial.MemorialHall
	for _, h := range m.halls {
		result = append(result, h)
	}
	return result, nil
}

func (m *mockMemorialRepository) Search(ctx context.Context, query *memorial.SearchQuery) ([]*memorial.MemorialHall, int64, error) {
	var result []*memorial.MemorialHall
	for _, h := range m.halls {
		result = append(result, h)
	}
	return result, int64(len(result)), nil
}

func (m *mockMemorialRepository) FindByRegion(ctx context.Context, province, city string) ([]*memorial.MemorialHall, error) {
	var result []*memorial.MemorialHall
	for _, h := range m.halls {
		if h.Province == province && h.City == city {
			result = append(result, h)
		}
	}
	return result, nil
}

// mockTabletRepository 模拟牌位仓储
type mockTabletRepository struct {
	tablets map[int64]*memorial.MemorialTablet
	nextID  int64
}

func newMockTabletRepository() *mockTabletRepository {
	return &mockTabletRepository{
		tablets: make(map[int64]*memorial.MemorialTablet),
		nextID:  1,
	}
}

func (m *mockTabletRepository) FindByID(ctx context.Context, id int64) (*memorial.MemorialTablet, error) {
	t, ok := m.tablets[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (m *mockTabletRepository) FindByHallID(ctx context.Context, hallID int64) ([]*memorial.MemorialTablet, error) {
	var result []*memorial.MemorialTablet
	for _, t := range m.tablets {
		if t.HallID == hallID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *mockTabletRepository) FindByPersonID(ctx context.Context, personID int64) (*memorial.MemorialTablet, error) {
	for _, t := range m.tablets {
		if t.PersonID != nil && *t.PersonID == personID {
			return t, nil
		}
	}
	return nil, nil
}

func (m *mockTabletRepository) Create(ctx context.Context, t *memorial.MemorialTablet) error {
	t.ID = m.nextID
	m.nextID++
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	m.tablets[t.ID] = t
	return nil
}

func (m *mockTabletRepository) Update(ctx context.Context, t *memorial.MemorialTablet) error {
	t.UpdatedAt = time.Now()
	m.tablets[t.ID] = t
	return nil
}

func (m *mockTabletRepository) Delete(ctx context.Context, id int64) error {
	delete(m.tablets, id)
	return nil
}

func (m *mockTabletRepository) DeleteByHallID(ctx context.Context, hallID int64) error {
	for id, t := range m.tablets {
		if t.HallID == hallID {
			delete(m.tablets, id)
		}
	}
	return nil
}

func (m *mockTabletRepository) Search(ctx context.Context, query *memorial.TabletSearchQuery) ([]*memorial.MemorialTablet, int64, error) {
	var result []*memorial.MemorialTablet
	for _, t := range m.tablets {
		result = append(result, t)
	}
	return result, int64(len(result)), nil
}

// newTestMemorialService 创建测试用宗祠应用服务
func newTestMemorialService() (*MemorialService, *mockMemorialRepository, *mockTabletRepository) {
	hallRepo := newMockMemorialRepository()
	tabletRepo := newMockTabletRepository()
	svc := NewMemorialService(hallRepo, tabletRepo)
	return svc, hallRepo, tabletRepo
}

func TestCreateHall(t *testing.T) {
	svc, repo, _ := newTestMemorialService()
	ctx := context.Background()

	t.Run("create hall successfully", func(t *testing.T) {
		req := &CreateHallRequest{
			Name:        "王氏宗祠",
			Description: "王氏家族宗祠",
			Province:    "四川省",
			City:        "成都市",
			District:    "武侯区",
			Address:     "某街道123号",
			Latitude:    30.6587,
			Longitude:   104.0658,
			Style:       "中式传统",
			TotalTablet: 100,
		}

		result, err := svc.CreateHall(ctx, req)
		if err != nil {
			t.Fatalf("CreateHall() error = %v", err)
		}
		if result.Name != "王氏宗祠" {
			t.Errorf("expected name '王氏宗祠', got '%s'", result.Name)
		}
		if result.Province != "四川省" {
			t.Errorf("expected province '四川省', got '%s'", result.Province)
		}
		if result.TotalTablet != 100 {
			t.Errorf("expected total_tablet 100, got %d", result.TotalTablet)
		}
		if len(repo.halls) != 1 {
			t.Errorf("expected 1 hall in repo, got %d", len(repo.halls))
		}
	})
}

func TestGetHall(t *testing.T) {
	svc, repo, _ := newTestMemorialService()
	ctx := context.Background()

	// 预先创建一个宗祠
	repo.Create(ctx, &memorial.MemorialHall{
		Name:       "李氏宗祠",
		Province:   "四川省",
		City:       "成都市",
		TotalTablet: 50,
	})

	t.Run("get existing hall", func(t *testing.T) {
		result, err := svc.GetHall(ctx, 1)
		if err != nil {
			t.Fatalf("GetHall() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected hall, got nil")
		}
		if result.Name != "李氏宗祠" {
			t.Errorf("expected name '李氏宗祠', got '%s'", result.Name)
		}
	})

	t.Run("get non-existent hall", func(t *testing.T) {
		result, err := svc.GetHall(ctx, 999)
		if err != nil {
			t.Fatalf("GetHall() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestListHalls(t *testing.T) {
	svc, repo, _ := newTestMemorialService()
	ctx := context.Background()

	// 创建多个宗祠
	for i := 0; i < 3; i++ {
		repo.Create(ctx, &memorial.MemorialHall{
			Name:       "宗祠" + string(rune('A'+i)),
			Province:   "四川省",
			City:       "成都市",
		})
	}

	t.Run("list all halls", func(t *testing.T) {
		result, err := svc.ListHalls(ctx)
		if err != nil {
			t.Fatalf("ListHalls() error = %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 halls, got %d", len(result))
		}
	})
}

func TestUpdateHall(t *testing.T) {
	svc, repo, _ := newTestMemorialService()
	ctx := context.Background()

	// 预先创建一个宗祠
	repo.Create(ctx, &memorial.MemorialHall{
		Name:        "原名",
		Province:    "四川省",
		City:        "成都市",
		TotalTablet: 50,
	})

	t.Run("update hall successfully", func(t *testing.T) {
		newName := "新宗祠名"
		newTotal := 100
		req := &UpdateHallRequest{
			Name:       &newName,
			TotalTablet: &newTotal,
		}

		result, err := svc.UpdateHall(ctx, 1, req)
		if err != nil {
			t.Fatalf("UpdateHall() error = %v", err)
		}
		if result.Name != "新宗祠名" {
			t.Errorf("expected name '新宗祠名', got '%s'", result.Name)
		}
		if result.TotalTablet != 100 {
			t.Errorf("expected total_tablet 100, got %d", result.TotalTablet)
		}
	})

	t.Run("update non-existent hall", func(t *testing.T) {
		newName := "新名称"
		req := &UpdateHallRequest{
			Name: &newName,
		}

		result, err := svc.UpdateHall(ctx, 999, req)
		if err != nil {
			t.Fatalf("UpdateHall() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestDeleteHall(t *testing.T) {
	svc, repo, _ := newTestMemorialService()
	ctx := context.Background()

	// 预先创建一个宗祠
	repo.Create(ctx, &memorial.MemorialHall{
		Name:    "待删除",
		Province: "四川省",
		City:    "成都市",
	})

	t.Run("delete existing hall", func(t *testing.T) {
		err := svc.DeleteHall(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteHall() error = %v", err)
		}
		if len(repo.halls) != 0 {
			t.Errorf("expected 0 halls after delete, got %d", len(repo.halls))
		}
	})

	t.Run("delete non-existent hall", func(t *testing.T) {
		err := svc.DeleteHall(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent hall")
		}
	})
}

func TestCreateTablet(t *testing.T) {
	svc, hallRepo, tabletRepo := newTestMemorialService()
	ctx := context.Background()

	// 先创建一个宗祠
	hallRepo.Create(ctx, &memorial.MemorialHall{
		Name:       "测试宗祠",
		Province:   "四川省",
		City:       "成都市",
		TotalTablet: 100,
	})

	t.Run("create tablet successfully", func(t *testing.T) {
		req := &CreateTabletRequest{
			HallID:     1,
			PersonName: "王阳明",
			Generation: 5,
			TabletType: "ancestor",
			Position:   "A区",
			Floor:      1,
			Row:        1,
			Number:     1,
		}

		result, err := svc.CreateTablet(ctx, req)
		if err != nil {
			t.Fatalf("CreateTablet() error = %v", err)
		}
		if result.PersonName != "王阳明" {
			t.Errorf("expected person_name '王阳明', got '%s'", result.PersonName)
		}
		if result.Generation != 5 {
			t.Errorf("expected generation 5, got %d", result.Generation)
		}
		if result.Floor != 1 {
			t.Errorf("expected floor 1, got %d", result.Floor)
		}
		if len(tabletRepo.tablets) != 1 {
			t.Errorf("expected 1 tablet in repo, got %d", len(tabletRepo.tablets))
		}
	})
}

func TestGetTablet(t *testing.T) {
	svc, _, tabletRepo := newTestMemorialService()
	ctx := context.Background()

	// 预先创建一个牌位
	tabletRepo.Create(ctx, &memorial.MemorialTablet{
		HallID:     1,
		PersonName: "孔子",
		Generation: 1,
		TabletType: memorial.TabletTypeAncestor,
		Floor:      1,
		Row:        1,
		Number:     1,
	})

	t.Run("get existing tablet", func(t *testing.T) {
		result, err := svc.GetTablet(ctx, 1)
		if err != nil {
			t.Fatalf("GetTablet() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected tablet, got nil")
		}
		if result.PersonName != "孔子" {
			t.Errorf("expected person_name '孔子', got '%s'", result.PersonName)
		}
	})

	t.Run("get non-existent tablet", func(t *testing.T) {
		result, err := svc.GetTablet(ctx, 999)
		if err != nil {
			t.Fatalf("GetTablet() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestListTabletsByHall(t *testing.T) {
	svc, _, tabletRepo := newTestMemorialService()
	ctx := context.Background()

	// 创建多个牌位
	for i := 0; i < 3; i++ {
		tabletRepo.Create(ctx, &memorial.MemorialTablet{
			HallID:     1,
			PersonName: "人物" + string(rune('A'+i)),
			Generation: i + 1,
			Floor:      1,
			Row:        i + 1,
			Number:     i + 1,
		})
	}

	t.Run("list tablets by hall", func(t *testing.T) {
		result, err := svc.ListTablets(ctx, 1)
		if err != nil {
			t.Fatalf("ListTablets() error = %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 tablets, got %d", len(result))
		}
	})
}

func TestUpdateTablet(t *testing.T) {
	svc, _, tabletRepo := newTestMemorialService()
	ctx := context.Background()

	// 预先创建一个牌位
	tabletRepo.Create(ctx, &memorial.MemorialTablet{
		HallID:     1,
		PersonName: "原名",
		Generation: 1,
		Floor:      1,
		Row:        1,
		Number:     1,
	})

	t.Run("update tablet successfully", func(t *testing.T) {
		newName := "新名字"
		newFloor := 2
		req := &UpdateTabletRequest{
			PersonName: &newName,
			Floor:      &newFloor,
		}

		result, err := svc.UpdateTablet(ctx, 1, req)
		if err != nil {
			t.Fatalf("UpdateTablet() error = %v", err)
		}
		if result.PersonName != "新名字" {
			t.Errorf("expected person_name '新名字', got '%s'", result.PersonName)
		}
		if result.Floor != 2 {
			t.Errorf("expected floor 2, got %d", result.Floor)
		}
	})

	t.Run("update non-existent tablet", func(t *testing.T) {
		newName := "新名称"
		req := &UpdateTabletRequest{
			PersonName: &newName,
		}

		result, err := svc.UpdateTablet(ctx, 999, req)
		if err != nil {
			t.Fatalf("UpdateTablet() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestDeleteTablet(t *testing.T) {
	svc, _, tabletRepo := newTestMemorialService()
	ctx := context.Background()

	// 预先创建一个牌位
	tabletRepo.Create(ctx, &memorial.MemorialTablet{
		HallID:     1,
		PersonName: "待删除",
		Floor:      1,
		Row:        1,
		Number:     1,
	})

	t.Run("delete existing tablet", func(t *testing.T) {
		err := svc.DeleteTablet(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteTablet() error = %v", err)
		}
		if len(tabletRepo.tablets) != 0 {
			t.Errorf("expected 0 tablets after delete, got %d", len(tabletRepo.tablets))
		}
	})
}

func TestMemorialServiceErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("repository error on create hall", func(t *testing.T) {
		badHallRepo := &mockMemorialRepositoryError{
			inner: newMockMemorialRepository(),
		}
		badTabletRepo := newMockTabletRepository()
		badSvc := NewMemorialService(badHallRepo, badTabletRepo)

		req := &CreateHallRequest{
			Name:    "测试",
			Province: "四川省",
			City:    "成都市",
		}

		_, err := badSvc.CreateHall(ctx, req)
		if err == nil {
			t.Error("expected error from repository")
		}
	})

	t.Run("repository error on get hall", func(t *testing.T) {
		badHallRepo := &mockMemorialRepositoryError{
			inner: newMockMemorialRepository(),
		}
		badTabletRepo := newMockTabletRepository()
		badSvc := NewMemorialService(badHallRepo, badTabletRepo)

		_, err := badSvc.GetHall(ctx, 1)
		if err == nil {
			t.Error("expected error from repository")
		}
	})

	t.Run("repository error on create tablet", func(t *testing.T) {
		goodHallRepo := newMockMemorialRepository()
		goodHallRepo.Create(ctx, &memorial.MemorialHall{
			Name:    "测试宗祠",
			Province: "四川省",
			City:    "成都市",
		})

		badTabletRepo := &mockTabletRepositoryError{
			inner: newMockTabletRepository(),
		}
		badSvc := NewMemorialService(goodHallRepo, badTabletRepo)

		req := &CreateTabletRequest{
			HallID:     1,
			PersonName: "测试",
			Number:     1,
		}

		_, err := badSvc.CreateTablet(ctx, req)
		if err == nil {
			t.Error("expected error from repository")
		}
	})
}

type mockMemorialRepositoryError struct {
	inner *mockMemorialRepository
}

func (m *mockMemorialRepositoryError) FindByID(ctx context.Context, id int64) (*memorial.MemorialHall, error) {
	return nil, errors.New("repository error")
}
func (m *mockMemorialRepositoryError) Create(ctx context.Context, h *memorial.MemorialHall) error {
	return errors.New("repository error")
}
func (m *mockMemorialRepositoryError) Update(ctx context.Context, h *memorial.MemorialHall) error {
	return m.inner.Update(ctx, h)
}
func (m *mockMemorialRepositoryError) Delete(ctx context.Context, id int64) error {
	return m.inner.Delete(ctx, id)
}
func (m *mockMemorialRepositoryError) FindAll(ctx context.Context) ([]*memorial.MemorialHall, error) {
	return m.inner.FindAll(ctx)
}
func (m *mockMemorialRepositoryError) Search(ctx context.Context, query *memorial.SearchQuery) ([]*memorial.MemorialHall, int64, error) {
	return m.inner.Search(ctx, query)
}
func (m *mockMemorialRepositoryError) FindByRegion(ctx context.Context, province, city string) ([]*memorial.MemorialHall, error) {
	return m.inner.FindByRegion(ctx, province, city)
}

type mockTabletRepositoryError struct {
	inner *mockTabletRepository
}

func (m *mockTabletRepositoryError) FindByID(ctx context.Context, id int64) (*memorial.MemorialTablet, error) {
	return m.inner.FindByID(ctx, id)
}
func (m *mockTabletRepositoryError) FindByHallID(ctx context.Context, hallID int64) ([]*memorial.MemorialTablet, error) {
	return m.inner.FindByHallID(ctx, hallID)
}
func (m *mockTabletRepositoryError) FindByPersonID(ctx context.Context, personID int64) (*memorial.MemorialTablet, error) {
	return m.inner.FindByPersonID(ctx, personID)
}
func (m *mockTabletRepositoryError) Create(ctx context.Context, t *memorial.MemorialTablet) error {
	return errors.New("repository error")
}
func (m *mockTabletRepositoryError) Update(ctx context.Context, t *memorial.MemorialTablet) error {
	return m.inner.Update(ctx, t)
}
func (m *mockTabletRepositoryError) Delete(ctx context.Context, id int64) error {
	return m.inner.Delete(ctx, id)
}
func (m *mockTabletRepositoryError) DeleteByHallID(ctx context.Context, hallID int64) error {
	return m.inner.DeleteByHallID(ctx, hallID)
}
func (m *mockTabletRepositoryError) Search(ctx context.Context, query *memorial.TabletSearchQuery) ([]*memorial.MemorialTablet, int64, error) {
	return m.inner.Search(ctx, query)
}