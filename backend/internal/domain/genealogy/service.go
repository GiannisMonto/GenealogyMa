package genealogy

import (
	"context"
	"fmt"
)

// Service 族谱领域服务
type Service struct {
	repo           Repository
	branchRepo     BranchRepository
	generationRepo GenerationRepository
}

// NewService 创建族谱服务
func NewService(repo Repository, branchRepo BranchRepository, generationRepo GenerationRepository) *Service {
	return &Service{
		repo:           repo,
		branchRepo:     branchRepo,
		generationRepo: generationRepo,
	}
}

// CreateGenealogy 创建族谱
func (s *Service) CreateGenealogy(ctx context.Context, genealogy *Genealogy) error {
	if err := genealogy.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, genealogy)
}

// UpdateGenealogy 更新族谱
func (s *Service) UpdateGenealogy(ctx context.Context, genealogy *Genealogy) error {
	if err := genealogy.Validate(); err != nil {
		return err
	}
	existing, err := s.repo.FindByID(ctx, genealogy.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("genealogy not found")
	}
	return s.repo.Update(ctx, genealogy)
}

// DeleteGenealogy 删除族谱
func (s *Service) DeleteGenealogy(ctx context.Context, id int64) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("genealogy not found")
	}
	// Delete all branches and generations first
	if err := s.branchRepo.DeleteByGenealogyID(ctx, id); err != nil {
		return err
	}
	if err := s.generationRepo.DeleteByGenealogyID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

// GetGenealogy 获取族谱详情
func (s *Service) GetGenealogy(ctx context.Context, id int64) (*Genealogy, error) {
	return s.repo.FindByID(ctx, id)
}

// ListGenealogies 获取所有族谱
func (s *Service) ListGenealogies(ctx context.Context) ([]*Genealogy, error) {
	return s.repo.FindAll(ctx)
}

// SearchGenealogies 搜索族谱
func (s *Service) SearchGenealogies(ctx context.Context, query *SearchQuery) ([]*Genealogy, int64, error) {
	return s.repo.Search(ctx, query)
}

// GetGenealogyWithDetails 获取族谱详情（含分支和世代）
func (s *Service) GetGenealogyWithDetails(ctx context.Context, id int64) (*Genealogy, error) {
	genealogy, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if genealogy == nil {
		return nil, fmt.Errorf("genealogy not found")
	}

	// Load branches
	branches, err := s.branchRepo.FindByGenealogyID(ctx, id)
	if err != nil {
		return nil, err
	}
	genealogy.Branches = branches

	// Load generations
	generations, err := s.generationRepo.FindByGenealogyID(ctx, id)
	if err != nil {
		return nil, err
	}
	genealogy.Generations = generations

	return genealogy, nil
}

// CreateBranch 创建分支
func (s *Service) CreateBranch(ctx context.Context, branch *Branch) error {
	if err := branch.Validate(); err != nil {
		return err
	}
	// Check genealogy exists
	genealogy, err := s.repo.FindByID(ctx, branch.GenealogyID)
	if err != nil {
		return err
	}
	if genealogy == nil {
		return fmt.Errorf("genealogy not found")
	}
	return s.branchRepo.Create(ctx, branch)
}

// UpdateBranch 更新分支
func (s *Service) UpdateBranch(ctx context.Context, branch *Branch) error {
	if err := branch.Validate(); err != nil {
		return err
	}
	existing, err := s.branchRepo.FindByID(ctx, branch.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("branch not found")
	}
	return s.branchRepo.Update(ctx, branch)
}

// DeleteBranch 删除分支
func (s *Service) DeleteBranch(ctx context.Context, id int64) error {
	existing, err := s.branchRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("branch not found")
	}
	return s.branchRepo.Delete(ctx, id)
}

// GetBranch 获取分支详情
func (s *Service) GetBranch(ctx context.Context, id int64) (*Branch, error) {
	return s.branchRepo.FindByID(ctx, id)
}

// ListBranches 获取族谱下所有分支
func (s *Service) ListBranches(ctx context.Context, genealogyID int64) ([]*Branch, error) {
	return s.branchRepo.FindByGenealogyID(ctx, genealogyID)
}

// SearchBranches 搜索分支
func (s *Service) SearchBranches(ctx context.Context, query *BranchSearchQuery) ([]*Branch, int64, error) {
	return s.branchRepo.Search(ctx, query)
}

// CreateGeneration 创建世代
func (s *Service) CreateGeneration(ctx context.Context, generation *Generation) error {
	if err := generation.Validate(); err != nil {
		return err
	}
	// Check genealogy exists
	genealogy, err := s.repo.FindByID(ctx, generation.GenealogyID)
	if err != nil {
		return err
	}
	if genealogy == nil {
		return fmt.Errorf("genealogy not found")
	}
	return s.generationRepo.Create(ctx, generation)
}

// UpdateGeneration 更新世代
func (s *Service) UpdateGeneration(ctx context.Context, generation *Generation) error {
	if err := generation.Validate(); err != nil {
		return err
	}
	existing, err := s.generationRepo.FindByID(ctx, generation.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("generation not found")
	}
	return s.generationRepo.Update(ctx, generation)
}

// DeleteGeneration 删除世代
func (s *Service) DeleteGeneration(ctx context.Context, id int64) error {
	existing, err := s.generationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("generation not found")
	}
	return s.generationRepo.Delete(ctx, id)
}

// GetGeneration 获取世代详情
func (s *Service) GetGeneration(ctx context.Context, id int64) (*Generation, error) {
	return s.generationRepo.FindByID(ctx, id)
}

// ListGenerations 获取族谱下所有世代
func (s *Service) ListGenerations(ctx context.Context, genealogyID int64) ([]*Generation, error) {
	return s.generationRepo.FindByGenealogyID(ctx, genealogyID)
}

// GetGenerationByNumber 根据世代序号获取世代
func (s *Service) GetGenerationByNumber(ctx context.Context, genealogyID int64, generation int) (*Generation, error) {
	return s.generationRepo.FindByGeneration(ctx, genealogyID, generation)
}