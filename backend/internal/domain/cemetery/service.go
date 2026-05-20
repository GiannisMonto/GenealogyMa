package cemetery

import (
	"context"
	"fmt"
)

// Service 墓园领域服务
type Service struct {
	repo        Repository
	graveRepo   GraveRepository
}

// NewService 创建墓园服务
func NewService(repo Repository, graveRepo GraveRepository) *Service {
	return &Service{
		repo:      repo,
		graveRepo: graveRepo,
	}
}

// CreateCemetery 创建墓园
func (s *Service) CreateCemetery(ctx context.Context, cemetery *Cemetery) error {
	if err := cemetery.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, cemetery)
}

// UpdateCemetery 更新墓园
func (s *Service) UpdateCemetery(ctx context.Context, cemetery *Cemetery) error {
	if err := cemetery.Validate(); err != nil {
		return err
	}
	existing, err := s.repo.FindByID(ctx, cemetery.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("cemetery not found")
	}
	return s.repo.Update(ctx, cemetery)
}

// DeleteCemetery 删除墓园
func (s *Service) DeleteCemetery(ctx context.Context, id int64) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("cemetery not found")
	}
	return s.repo.Delete(ctx, id)
}

// GetCemetery 获取墓园详情
func (s *Service) GetCemetery(ctx context.Context, id int64) (*Cemetery, error) {
	return s.repo.FindByID(ctx, id)
}

// ListCemeteries 获取所有墓园
func (s *Service) ListCemeteries(ctx context.Context) ([]*Cemetery, error) {
	return s.repo.FindAll(ctx)
}

// SearchCemeteries 搜索墓园
func (s *Service) SearchCemeteries(ctx context.Context, query *SearchQuery) ([]*Cemetery, int64, error) {
	return s.repo.Search(ctx, query)
}

// GetCemeteryWithGraveCount 获取墓园（含墓位统计）
func (s *Service) GetCemeteryWithGraveCount(ctx context.Context, id int64) (*Cemetery, int, error) {
	cemetery, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	graves, err := s.graveRepo.FindByCemeteryID(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	return cemetery, len(graves), nil
}

// CreateGrave 创建墓位
func (s *Service) CreateGrave(ctx context.Context, grave *Grave) error {
	if err := grave.Validate(); err != nil {
		return err
	}
	// 检查墓园是否存在
	cemetery, err := s.repo.FindByID(ctx, grave.CemeteryID)
	if err != nil {
		return err
	}
	if cemetery == nil {
		return fmt.Errorf("cemetery not found")
	}
	// 检查墓位是否可用
	if !cemetery.HasAvailableSpace() {
		return fmt.Errorf("cemetery has no available space")
	}
	return s.graveRepo.Create(ctx, grave)
}

// UpdateGrave 更新墓位
func (s *Service) UpdateGrave(ctx context.Context, grave *Grave) error {
	if err := grave.Validate(); err != nil {
		return err
	}
	existing, err := s.graveRepo.FindByID(ctx, grave.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("grave not found")
	}
	return s.graveRepo.Update(ctx, grave)
}

// GetGrave 获取墓位详情
func (s *Service) GetGrave(ctx context.Context, id int64) (*Grave, error) {
	return s.graveRepo.FindByID(ctx, id)
}

// GetGraveByPersonID 根据人物ID获取墓位
func (s *Service) GetGraveByPersonID(ctx context.Context, personID int64) (*Grave, error) {
	return s.graveRepo.FindByPersonID(ctx, personID)
}

// ListGravesByCemetery 获取墓园下所有墓位
func (s *Service) ListGravesByCemetery(ctx context.Context, cemeteryID int64) ([]*Grave, error) {
	return s.graveRepo.FindByCemeteryID(ctx, cemeteryID)
}

// SearchGraves 搜索墓位
func (s *Service) SearchGraves(ctx context.Context, query *GraveSearchQuery) ([]*Grave, int64, error) {
	return s.graveRepo.Search(ctx, query)
}