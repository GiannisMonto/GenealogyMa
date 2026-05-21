package memorial

import (
	"context"
	"fmt"
)

// Service 宗祠领域服务
type Service struct {
	repo       Repository
	tabletRepo TabletRepository
}

// NewService 创建宗祠服务
func NewService(repo Repository, tabletRepo TabletRepository) *Service {
	return &Service{
		repo:       repo,
		tabletRepo: tabletRepo,
	}
}

// CreateHall 创建宗祠
func (s *Service) CreateHall(ctx context.Context, hall *MemorialHall) error {
	if err := hall.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, hall)
}

// UpdateHall 更新宗祠
func (s *Service) UpdateHall(ctx context.Context, hall *MemorialHall) error {
	if err := hall.Validate(); err != nil {
		return err
	}
	existing, err := s.repo.FindByID(ctx, hall.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("memorial hall not found")
	}
	return s.repo.Update(ctx, hall)
}

// DeleteHall 删除宗祠
func (s *Service) DeleteHall(ctx context.Context, id int64) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("memorial hall not found")
	}
	// Delete all tablets first
	if err := s.tabletRepo.DeleteByHallID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

// GetHall 获取宗祠详情
func (s *Service) GetHall(ctx context.Context, id int64) (*MemorialHall, error) {
	return s.repo.FindByID(ctx, id)
}

// ListHalls 获取所有宗祠
func (s *Service) ListHalls(ctx context.Context) ([]*MemorialHall, error) {
	return s.repo.FindAll(ctx)
}

// SearchHalls 搜索宗祠
func (s *Service) SearchHalls(ctx context.Context, query *SearchQuery) ([]*MemorialHall, int64, error) {
	return s.repo.Search(ctx, query)
}

// GetHallWithTablets 获取宗祠详情（含牌位）
func (s *Service) GetHallWithTablets(ctx context.Context, id int64) (*MemorialHall, error) {
	hall, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if hall == nil {
		return nil, fmt.Errorf("memorial hall not found")
	}

	tablets, err := s.tabletRepo.FindByHallID(ctx, id)
	if err != nil {
		return nil, err
	}
	hall.Tables = tablets

	return hall, nil
}

// CreateTablet 创建牌位
func (s *Service) CreateTablet(ctx context.Context, tablet *MemorialTablet) error {
	if err := tablet.Validate(); err != nil {
		return err
	}
	// Check hall exists
	hall, err := s.repo.FindByID(ctx, tablet.HallID)
	if err != nil {
		return err
	}
	if hall == nil {
		return fmt.Errorf("memorial hall not found")
	}
	// Check if hall has available space
	if !hall.HasAvailableSpace() {
		return fmt.Errorf("memorial hall has no available space")
	}
	return s.tabletRepo.Create(ctx, tablet)
}

// UpdateTablet 更新牌位
func (s *Service) UpdateTablet(ctx context.Context, tablet *MemorialTablet) error {
	if err := tablet.Validate(); err != nil {
		return err
	}
	existing, err := s.tabletRepo.FindByID(ctx, tablet.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("memorial tablet not found")
	}
	return s.tabletRepo.Update(ctx, tablet)
}

// DeleteTablet 删除牌位
func (s *Service) DeleteTablet(ctx context.Context, id int64) error {
	existing, err := s.tabletRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("memorial tablet not found")
	}
	return s.tabletRepo.Delete(ctx, id)
}

// GetTablet 获取牌位详情
func (s *Service) GetTablet(ctx context.Context, id int64) (*MemorialTablet, error) {
	return s.tabletRepo.FindByID(ctx, id)
}

// GetTabletByPersonID 根据人物ID获取牌位
func (s *Service) GetTabletByPersonID(ctx context.Context, personID int64) (*MemorialTablet, error) {
	return s.tabletRepo.FindByPersonID(ctx, personID)
}

// ListTablets 获取宗祠下所有牌位
func (s *Service) ListTablets(ctx context.Context, hallID int64) ([]*MemorialTablet, error) {
	return s.tabletRepo.FindByHallID(ctx, hallID)
}

// SearchTablets 搜索牌位
func (s *Service) SearchTablets(ctx context.Context, query *TabletSearchQuery) ([]*MemorialTablet, int64, error) {
	return s.tabletRepo.Search(ctx, query)
}