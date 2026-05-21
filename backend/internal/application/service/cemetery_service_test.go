package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/cemetery"
)

// mockCemeteryRepository 模拟墓园仓储
type mockCemeteryRepository struct {
	cemeteries map[int64]*cemetery.Cemetery
	nextID     int64
}

func newMockCemeteryRepository() *mockCemeteryRepository {
	return &mockCemeteryRepository{
		cemeteries: make(map[int64]*cemetery.Cemetery),
		nextID:     1,
	}
}

func (m *mockCemeteryRepository) Create(ctx context.Context, c *cemetery.Cemetery) error {
	c.ID = m.nextID
	m.nextID++
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	m.cemeteries[c.ID] = c
	return nil
}

func (m *mockCemeteryRepository) FindByID(ctx context.Context, id int64) (*cemetery.Cemetery, error) {
	c, ok := m.cemeteries[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *mockCemeteryRepository) FindAll(ctx context.Context) ([]*cemetery.Cemetery, error) {
	var result []*cemetery.Cemetery
	for _, c := range m.cemeteries {
		result = append(result, c)
	}
	return result, nil
}

func (m *mockCemeteryRepository) Update(ctx context.Context, c *cemetery.Cemetery) error {
	c.UpdatedAt = time.Now()
	m.cemeteries[c.ID] = c
	return nil
}

func (m *mockCemeteryRepository) Delete(ctx context.Context, id int64) error {
	delete(m.cemeteries, id)
	return nil
}

func (m *mockCemeteryRepository) Search(ctx context.Context, query *cemetery.SearchQuery) ([]*cemetery.Cemetery, int64, error) {
	var result []*cemetery.Cemetery
	for _, c := range m.cemeteries {
		result = append(result, c)
	}
	return result, int64(len(result)), nil
}

func (m *mockCemeteryRepository) FindByRegion(ctx context.Context, province, city string) ([]*cemetery.Cemetery, error) {
	var result []*cemetery.Cemetery
	for _, c := range m.cemeteries {
		if c.Province == province && c.City == city {
			result = append(result, c)
		}
	}
	return result, nil
}

// mockGraveRepository 模拟墓位仓储
type mockGraveRepository struct {
	graves   map[int64]*cemetery.Grave
	nextID   int64
}

func newMockGraveRepository() *mockGraveRepository {
	return &mockGraveRepository{
		graves: make(map[int64]*cemetery.Grave),
		nextID: 1,
	}
}

func (m *mockGraveRepository) Create(ctx context.Context, g *cemetery.Grave) error {
	g.ID = m.nextID
	m.nextID++
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	m.graves[g.ID] = g
	return nil
}

func (m *mockGraveRepository) FindByID(ctx context.Context, id int64) (*cemetery.Grave, error) {
	g, ok := m.graves[id]
	if !ok {
		return nil, nil
	}
	return g, nil
}

func (m *mockGraveRepository) FindByCemeteryID(ctx context.Context, cemeteryID int64) ([]*cemetery.Grave, error) {
	var result []*cemetery.Grave
	for _, g := range m.graves {
		if g.CemeteryID == cemeteryID {
			result = append(result, g)
		}
	}
	return result, nil
}

func (m *mockGraveRepository) FindByPersonID(ctx context.Context, personID int64) (*cemetery.Grave, error) {
	for _, g := range m.graves {
		if g.PersonID != nil && *g.PersonID == personID {
			return g, nil
		}
	}
	return nil, nil
}

func (m *mockGraveRepository) Update(ctx context.Context, g *cemetery.Grave) error {
	g.UpdatedAt = time.Now()
	m.graves[g.ID] = g
	return nil
}

func (m *mockGraveRepository) Delete(ctx context.Context, id int64) error {
	delete(m.graves, id)
	return nil
}

func (m *mockGraveRepository) DeleteByCemeteryID(ctx context.Context, cemeteryID int64) error {
	for id, g := range m.graves {
		if g.CemeteryID == cemeteryID {
			delete(m.graves, id)
		}
	}
	return nil
}

func (m *mockGraveRepository) Search(ctx context.Context, query *cemetery.GraveSearchQuery) ([]*cemetery.Grave, int64, error) {
	var result []*cemetery.Grave
	for _, g := range m.graves {
		result = append(result, g)
	}
	return result, int64(len(result)), nil
}

func newTestCemeteryService() (*CemeteryService, *mockCemeteryRepository, *mockGraveRepository) {
	cemeteryRepo := newMockCemeteryRepository()
	graveRepo := newMockGraveRepository()
	svc := NewCemeteryService(cemeteryRepo, graveRepo)
	return svc, cemeteryRepo, graveRepo
}

// ===== 测试用例 =====

func TestCreateCemetery(t *testing.T) {
	svc, repo, _ := newTestCemeteryService()
	ctx := context.Background()

	t.Run("create cemetery successfully", func(t *testing.T) {
		req := &CreateCemeteryRequest{
			Name:        "八宝山公墓",
			Description: "北京市最大的公墓",
			Province:    "北京市",
			City:        "北京市",
			District:    "石景山区",
			Address:     "八宝山革命公墓",
			Latitude:    39.9042,
			Longitude:   116.4074,
			TotalGrave:  10000,
		}

		result, err := svc.CreateCemetery(ctx, req)
		if err != nil {
			t.Fatalf("CreateCemetery() error = %v", err)
		}
		if result.Name != "八宝山公墓" {
			t.Errorf("expected name '八宝山公墓', got '%s'", result.Name)
		}
		if result.Province != "北京市" {
			t.Errorf("expected province '北京市', got '%s'", result.Province)
		}
		// 验证仓储中的记录
		if len(repo.cemeteries) != 1 {
			t.Errorf("expected 1 cemetery in repo, got %d", len(repo.cemeteries))
		}
	})

	t.Run("create cemetery with empty name should fail", func(t *testing.T) {
		req := &CreateCemeteryRequest{
			Name:     "",
			Province: "北京市",
			City:     "北京市",
		}

		_, err := svc.CreateCemetery(ctx, req)
		if err == nil {
			t.Error("expected error for empty name")
		}
	})
}

func TestGetCemetery(t *testing.T) {
	svc, repo, _ := newTestCemeteryService()
	ctx := context.Background()

	// 预先创建一个墓园
	repo.Create(ctx, &cemetery.Cemetery{
		Name:        "测试墓园",
		Province:    "四川省",
		City:        "成都市",
		TotalGrave:  5000,
	})

	t.Run("get existing cemetery", func(t *testing.T) {
		result, err := svc.GetCemetery(ctx, 1)
		if err != nil {
			t.Fatalf("GetCemetery() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected cemetery, got nil")
		}
		if result.Name != "测试墓园" {
			t.Errorf("expected name '测试墓园', got '%s'", result.Name)
		}
	})

	t.Run("get non-existent cemetery", func(t *testing.T) {
		result, err := svc.GetCemetery(ctx, 999)
		if err != nil {
			t.Fatalf("GetCemetery() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestListCemeteries(t *testing.T) {
	svc, repo, _ := newTestCemeteryService()
	ctx := context.Background()

	// 创建多个墓园
	for i := 0; i < 3; i++ {
		repo.Create(ctx, &cemetery.Cemetery{
			Name:       "墓园" + string(rune('A'+i)),
			Province:   "四川省",
			City:       "成都市",
			TotalGrave: 1000,
		})
	}

	t.Run("list all cemeteries", func(t *testing.T) {
		result, err := svc.ListCemeteries(ctx)
		if err != nil {
			t.Fatalf("ListCemeteries() error = %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 cemeteries, got %d", len(result))
		}
	})
}

func TestUpdateCemetery(t *testing.T) {
	svc, repo, _ := newTestCemeteryService()
	ctx := context.Background()

	// 预先创建一个墓园
	repo.Create(ctx, &cemetery.Cemetery{
		Name:       "原名",
		Province:   "四川省",
		City:       "成都市",
		TotalGrave: 1000,
	})

	t.Run("update cemetery successfully", func(t *testing.T) {
		newName := "新名称"
		newTotal := 2000
		req := &UpdateCemeteryRequest{
			Name:       &newName,
			TotalGrave: &newTotal,
		}

		result, err := svc.UpdateCemetery(ctx, 1, req)
		if err != nil {
			t.Fatalf("UpdateCemetery() error = %v", err)
		}
		if result.Name != "新名称" {
			t.Errorf("expected name '新名称', got '%s'", result.Name)
		}
		if result.TotalGrave != 2000 {
			t.Errorf("expected total grave 2000, got %d", result.TotalGrave)
		}
	})

	t.Run("update non-existent cemetery", func(t *testing.T) {
		newName := "新名称"
		req := &UpdateCemeteryRequest{
			Name: &newName,
		}

		result, err := svc.UpdateCemetery(ctx, 999, req)
		if err != nil {
			t.Fatalf("UpdateCemetery() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestDeleteCemetery(t *testing.T) {
	svc, repo, _ := newTestCemeteryService()
	ctx := context.Background()

	// 预先创建一个墓园
	repo.Create(ctx, &cemetery.Cemetery{
		Name:       "待删除",
		Province:   "四川省",
		City:       "成都市",
		TotalGrave: 1000,
	})

	t.Run("delete existing cemetery", func(t *testing.T) {
		err := svc.DeleteCemetery(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteCemetery() error = %v", err)
		}
		if len(repo.cemeteries) != 0 {
			t.Errorf("expected 0 cemeteries after delete, got %d", len(repo.cemeteries))
		}
	})

	t.Run("delete non-existent cemetery", func(t *testing.T) {
		err := svc.DeleteCemetery(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent cemetery")
		}
	})
}

func TestCreateGrave(t *testing.T) {
	svc, cemeteryRepo, graveRepo := newTestCemeteryService()
	ctx := context.Background()

	// 先创建一个墓园（使用服务返回的 cemeteryRepo）
	cemeteryRepo.Create(ctx, &cemetery.Cemetery{
		Name:       "测试墓园",
		Province:   "四川省",
		City:       "成都市",
		TotalGrave: 1000, // 设置足够的墓位数量
	})

	t.Run("create grave successfully", func(t *testing.T) {
		req := &CreateGraveRequest{
			CemeteryID: 1,
			Section:   "A区",
			Row:        1,
			Number:     1,
			Status:     "available",
		}

		result, err := svc.CreateGrave(ctx, req)
		if err != nil {
			t.Fatalf("CreateGrave() error = %v", err)
		}
		if result.Section != "A区" {
			t.Errorf("expected section 'A区', got '%s'", result.Section)
		}
		if result.Row != 1 {
			t.Errorf("expected row 1, got %d", result.Row)
		}
		if result.Status != "available" {
			t.Errorf("expected status 'available', got '%s'", result.Status)
		}
		// 验证仓储中的记录
		if len(graveRepo.graves) != 1 {
			t.Errorf("expected 1 grave in repo, got %d", len(graveRepo.graves))
		}
	})

	t.Run("create grave with invalid row", func(t *testing.T) {
		req := &CreateGraveRequest{
			CemeteryID: 1,
			Section:    "B区",
			Row:        0,
			Number:     1,
		}

		_, err := svc.CreateGrave(ctx, req)
		if err == nil {
			t.Error("expected error for invalid row")
		}
	})
}

func TestGetGrave(t *testing.T) {
	svc, _, graveRepo := newTestCemeteryService()
	ctx := context.Background()

	// 预先创建一个墓位
	graveRepo.Create(ctx, &cemetery.Grave{
		CemeteryID: 1,
		Section:    "A区",
		Row:        1,
		Number:     1,
		Status:     cemetery.GraveStatusAvailable,
	})

	t.Run("get existing grave", func(t *testing.T) {
		result, err := svc.GetGrave(ctx, 1)
		if err != nil {
			t.Fatalf("GetGrave() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected grave, got nil")
		}
		if result.Section != "A区" {
			t.Errorf("expected section 'A区', got '%s'", result.Section)
		}
	})

	t.Run("get non-existent grave", func(t *testing.T) {
		result, err := svc.GetGrave(ctx, 999)
		if err != nil {
			t.Fatalf("GetGrave() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestListGravesByCemetery(t *testing.T) {
	svc, _, graveRepo := newTestCemeteryService()
	ctx := context.Background()

	// 创建多个墓位
	for i := 0; i < 3; i++ {
		graveRepo.Create(ctx, &cemetery.Grave{
			CemeteryID: 1,
			Section:    "A区",
			Row:        i + 1,
			Number:     i + 1,
			Status:     cemetery.GraveStatusAvailable,
		})
	}
	// 创建一个属于其他墓园的墓位
	graveRepo.Create(ctx, &cemetery.Grave{
		CemeteryID: 2,
		Section:    "B区",
		Row:        1,
		Number:     1,
		Status:     cemetery.GraveStatusAvailable,
	})

	t.Run("list graves by cemetery", func(t *testing.T) {
		result, err := svc.ListGravesByCemetery(ctx, 1)
		if err != nil {
			t.Fatalf("ListGravesByCemetery() error = %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 graves, got %d", len(result))
		}
	})
}

func TestUpdateGrave(t *testing.T) {
	svc, _, graveRepo := newTestCemeteryService()
	ctx := context.Background()

	// 预先创建一个墓位
	graveRepo.Create(ctx, &cemetery.Grave{
		CemeteryID: 1,
		Section:    "A区",
		Row:        1,
		Number:     1,
		Status:     cemetery.GraveStatusAvailable,
	})

	t.Run("update grave successfully", func(t *testing.T) {
		newSection := "B区"
		newStatus := "occupied"
		req := &UpdateGraveRequest{
			Section: &newSection,
			Status:  &newStatus,
		}

		result, err := svc.UpdateGrave(ctx, 1, req)
		if err != nil {
			t.Fatalf("UpdateGrave() error = %v", err)
		}
		if result.Section != "B区" {
			t.Errorf("expected section 'B区', got '%s'", result.Section)
		}
		if result.Status != "occupied" {
			t.Errorf("expected status 'occupied', got '%s'", result.Status)
		}
	})

	t.Run("update non-existent grave", func(t *testing.T) {
		newSection := "B区"
		req := &UpdateGraveRequest{
			Section: &newSection,
		}

		result, err := svc.UpdateGrave(ctx, 999, req)
		if err != nil {
			t.Fatalf("UpdateGrave() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestDeleteGrave(t *testing.T) {
	svc, _, graveRepo := newTestCemeteryService()
	ctx := context.Background()

	// 预先创建一个墓位
	graveRepo.Create(ctx, &cemetery.Grave{
		CemeteryID: 1,
		Section:    "A区",
		Row:        1,
		Number:     1,
		Status:     cemetery.GraveStatusAvailable,
	})

	t.Run("delete existing grave", func(t *testing.T) {
		err := svc.DeleteGrave(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteGrave() error = %v", err)
		}
		if len(graveRepo.graves) != 0 {
			t.Errorf("expected 0 graves after delete, got %d", len(graveRepo.graves))
		}
	})

	t.Run("delete non-existent grave", func(t *testing.T) {
		err := svc.DeleteGrave(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent grave")
		}
	})
}

// 测试错误场景
func TestCemeteryServiceErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("repository error on create cemetery", func(t *testing.T) {
		badCemeteryRepo := &mockCemeteryRepositoryError{
			inner: newMockCemeteryRepository(),
		}
		badGraveRepo := newMockGraveRepository()
		badSvc := NewCemeteryService(badCemeteryRepo, badGraveRepo)

		req := &CreateCemeteryRequest{
			Name:     "测试",
			Province: "四川省",
			City:     "成都市",
		}

		_, err := badSvc.CreateCemetery(ctx, req)
		if err == nil {
			t.Error("expected error from repository")
		}
	})

	t.Run("repository error on create grave", func(t *testing.T) {
		goodCemeteryRepo := newMockCemeteryRepository()
		goodCemeteryRepo.Create(ctx, &cemetery.Cemetery{
			Name:     "测试墓园",
			Province: "四川省",
			City:     "成都市",
		})

		badGraveRepo := &mockGraveRepositoryError{
			inner: newMockGraveRepository(),
		}
		badSvc := NewCemeteryService(goodCemeteryRepo, badGraveRepo)

		req := &CreateGraveRequest{
			CemeteryID: 1,
			Section:    "A区",
			Row:        1,
			Number:     1,
		}

		_, err := badSvc.CreateGrave(ctx, req)
		if err == nil {
			t.Error("expected error from repository")
		}
	})
}

type mockCemeteryRepositoryError struct {
	inner *mockCemeteryRepository
}

func (m *mockCemeteryRepositoryError) Create(ctx context.Context, c *cemetery.Cemetery) error {
	return errors.New("repository error")
}
func (m *mockCemeteryRepositoryError) FindByID(ctx context.Context, id int64) (*cemetery.Cemetery, error) {
	return m.inner.FindByID(ctx, id)
}
func (m *mockCemeteryRepositoryError) FindAll(ctx context.Context) ([]*cemetery.Cemetery, error) {
	return m.inner.FindAll(ctx)
}
func (m *mockCemeteryRepositoryError) Update(ctx context.Context, c *cemetery.Cemetery) error {
	return m.inner.Update(ctx, c)
}
func (m *mockCemeteryRepositoryError) Delete(ctx context.Context, id int64) error {
	return m.inner.Delete(ctx, id)
}
func (m *mockCemeteryRepositoryError) Search(ctx context.Context, query *cemetery.SearchQuery) ([]*cemetery.Cemetery, int64, error) {
	return m.inner.Search(ctx, query)
}
func (m *mockCemeteryRepositoryError) FindByRegion(ctx context.Context, province, city string) ([]*cemetery.Cemetery, error) {
	return m.inner.FindByRegion(ctx, province, city)
}

type mockGraveRepositoryError struct {
	inner *mockGraveRepository
}

func (m *mockGraveRepositoryError) Create(ctx context.Context, g *cemetery.Grave) error {
	return errors.New("repository error")
}
func (m *mockGraveRepositoryError) FindByID(ctx context.Context, id int64) (*cemetery.Grave, error) {
	return m.inner.FindByID(ctx, id)
}
func (m *mockGraveRepositoryError) FindByCemeteryID(ctx context.Context, cemeteryID int64) ([]*cemetery.Grave, error) {
	return m.inner.FindByCemeteryID(ctx, cemeteryID)
}
func (m *mockGraveRepositoryError) FindByPersonID(ctx context.Context, personID int64) (*cemetery.Grave, error) {
	return m.inner.FindByPersonID(ctx, personID)
}
func (m *mockGraveRepositoryError) Update(ctx context.Context, g *cemetery.Grave) error {
	return m.inner.Update(ctx, g)
}
func (m *mockGraveRepositoryError) Delete(ctx context.Context, id int64) error {
	return m.inner.Delete(ctx, id)
}
func (m *mockGraveRepositoryError) DeleteByCemeteryID(ctx context.Context, cemeteryID int64) error {
	return m.inner.DeleteByCemeteryID(ctx, cemeteryID)
}
func (m *mockGraveRepositoryError) Search(ctx context.Context, query *cemetery.GraveSearchQuery) ([]*cemetery.Grave, int64, error) {
	return m.inner.Search(ctx, query)
}