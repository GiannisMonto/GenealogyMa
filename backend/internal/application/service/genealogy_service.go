package service

import (
	"context"

	"github.com/genealogy-ma/platform/internal/domain/genealogy"
)

// GenealogyService 族谱应用服务
type GenealogyService struct {
	domainService *genealogy.Service
}

// NewGenealogyService 创建族谱应用服务
func NewGenealogyService(repo genealogy.Repository, branchRepo genealogy.BranchRepository, generationRepo genealogy.GenerationRepository) *GenealogyService {
	return &GenealogyService{
		domainService: genealogy.NewService(repo, branchRepo, generationRepo),
	}
}

// ===== DTO 定义 =====

// GenealogyDTO 族谱DTO
type GenealogyDTO struct {
	ID              int64            `json:"id"`
	Name            string           `json:"name"`
	Surname         string           `json:"surname"`
	Description     string           `json:"description"`
	OriginPlace     string           `json:"origin_place"`
	TotalGenerations int            `json:"total_generations"`
	TotalMembers    int              `json:"total_members"`
	Version         string           `json:"version"`
	ImageURL        string           `json:"image_url"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
	Branches        []*BranchDTO     `json:"branches,omitempty"`
	Generations     []*GenerationDTO `json:"generations,omitempty"`
}

// BranchDTO 分支DTO
type BranchDTO struct {
	ID              int64             `json:"id"`
	GenealogyID     int64             `json:"genealogy_id"`
	Name            string            `json:"name"`
	Code            string            `json:"code"`
	Description     string            `json:"description"`
	GenerationStart int               `json:"generation_start"`
	GenerationEnd   int               `json:"generation_end"`
	MemberCount     int               `json:"member_count"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
	Members         []*PersonSummaryDTO `json:"members,omitempty"`
}

// GenerationDTO 世代DTO
type GenerationDTO struct {
	ID           int64   `json:"id"`
	GenealogyID  int64   `json:"genealogy_id"`
	Generation   int     `json:"generation"`
	Name         string  `json:"name"`
	Sequence     int     `json:"sequence"`
	Description  string  `json:"description"`
	StartYear    *int    `json:"start_year,omitempty"`
	EndYear      *int    `json:"end_year,omitempty"`
	MemberCount  int     `json:"member_count"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	Active       bool    `json:"active"`
}

// PersonSummaryDTO 人物摘要DTO
type PersonSummaryDTO struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Gender     string  `json:"gender"`
	Generation int     `json:"generation"`
	BirthYear  *int    `json:"birth_year,omitempty"`
	DeathYear  *int    `json:"death_year,omitempty"`
	StyleName  string  `json:"style_name"`
}

// CreateGenealogyRequest 创建族谱请求
type CreateGenealogyRequest struct {
	Name            string  `json:"name" binding:"required"`
	Surname         string  `json:"surname" binding:"required"`
	Description     string  `json:"description"`
	OriginPlace     string  `json:"origin_place"`
	TotalGenerations int    `json:"total_generations"`
	TotalMembers    int     `json:"total_members"`
	Version         string  `json:"version"`
	ImageURL        string  `json:"image_url"`
}

// UpdateGenealogyRequest 更新族谱请求
type UpdateGenealogyRequest struct {
	Name            *string `json:"name"`
	Surname         *string `json:"surname"`
	Description     *string `json:"description"`
	OriginPlace     *string `json:"origin_place"`
	TotalGenerations *int    `json:"total_generations"`
	TotalMembers    *int    `json:"total_members"`
	Version         *string `json:"version"`
	ImageURL        *string `json:"image_url"`
}

// CreateBranchRequest 创建分支请求
type CreateBranchRequest struct {
	GenealogyID     int64  `json:"genealogy_id" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Code            string `json:"code"`
	Description     string `json:"description"`
	GenerationStart int    `json:"generation_start"`
	GenerationEnd   int    `json:"generation_end"`
}

// UpdateBranchRequest 更新分支请求
type UpdateBranchRequest struct {
	Name            *string `json:"name"`
	Code            *string `json:"code"`
	Description     *string `json:"description"`
	GenerationStart *int    `json:"generation_start"`
	GenerationEnd   *int    `json:"generation_end"`
}

// CreateGenerationRequest 创建世代请求
type CreateGenerationRequest struct {
	GenealogyID int64  `json:"genealogy_id" binding:"required"`
	Generation  int     `json:"generation" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Sequence    int    `json:"sequence"`
	Description string `json:"description"`
	StartYear   *int   `json:"start_year"`
	EndYear     *int   `json:"end_year"`
}

// UpdateGenerationRequest 更新世代请求
type UpdateGenerationRequest struct {
	Name        *string `json:"name"`
	Sequence    *int    `json:"sequence"`
	Description *string `json:"description"`
	StartYear   *int    `json:"start_year"`
	EndYear     *int    `json:"end_year"`
}

// ===== 服务方法 =====

// ListGenealogies 获取所有族谱
func (s *GenealogyService) ListGenealogies(ctx context.Context) ([]*GenealogyDTO, error) {
	genealogies, err := s.domainService.ListGenealogies(ctx)
	if err != nil {
		return nil, err
	}
	dtos := make([]*GenealogyDTO, len(genealogies))
	for i, g := range genealogies {
		dtos[i] = genealogyToDTO(g)
	}
	return dtos, nil
}

// GetGenealogy 获取族谱详情
func (s *GenealogyService) GetGenealogy(ctx context.Context, id int64) (*GenealogyDTO, error) {
	g, err := s.domainService.GetGenealogy(ctx, id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, nil
	}
	return genealogyToDTO(g), nil
}

// GetGenealogyWithDetails 获取族谱详情（含分支和世代）
func (s *GenealogyService) GetGenealogyWithDetails(ctx context.Context, id int64) (*GenealogyDTO, error) {
	g, err := s.domainService.GetGenealogyWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, nil
	}
	return genealogyToDTO(g), nil
}

// CreateGenealogy 创建族谱
func (s *GenealogyService) CreateGenealogy(ctx context.Context, req *CreateGenealogyRequest) (*GenealogyDTO, error) {
	g := &genealogy.Genealogy{
		Name:            req.Name,
		Surname:         req.Surname,
		Description:     req.Description,
		OriginPlace:     req.OriginPlace,
		TotalGenerations: req.TotalGenerations,
		TotalMembers:    req.TotalMembers,
		Version:         req.Version,
		ImageURL:        req.ImageURL,
	}

	if err := s.domainService.CreateGenealogy(ctx, g); err != nil {
		return nil, err
	}

	return genealogyToDTO(g), nil
}

// UpdateGenealogy 更新族谱
func (s *GenealogyService) UpdateGenealogy(ctx context.Context, id int64, req *UpdateGenealogyRequest) (*GenealogyDTO, error) {
	existing, err := s.domainService.GetGenealogy(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Surname != nil {
		existing.Surname = *req.Surname
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.OriginPlace != nil {
		existing.OriginPlace = *req.OriginPlace
	}
	if req.TotalGenerations != nil {
		existing.TotalGenerations = *req.TotalGenerations
	}
	if req.TotalMembers != nil {
		existing.TotalMembers = *req.TotalMembers
	}
	if req.Version != nil {
		existing.Version = *req.Version
	}
	if req.ImageURL != nil {
		existing.ImageURL = *req.ImageURL
	}

	if err := s.domainService.UpdateGenealogy(ctx, existing); err != nil {
		return nil, err
	}

	return genealogyToDTO(existing), nil
}

// DeleteGenealogy 删除族谱
func (s *GenealogyService) DeleteGenealogy(ctx context.Context, id int64) error {
	return s.domainService.DeleteGenealogy(ctx, id)
}

// ListBranches 获取族谱下所有分支
func (s *GenealogyService) ListBranches(ctx context.Context, genealogyID int64) ([]*BranchDTO, error) {
	branches, err := s.domainService.ListBranches(ctx, genealogyID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*BranchDTO, len(branches))
	for i, b := range branches {
		dtos[i] = branchToDTO(b)
	}
	return dtos, nil
}

// GetBranch 获取分支详情
func (s *GenealogyService) GetBranch(ctx context.Context, id int64) (*BranchDTO, error) {
	b, err := s.domainService.GetBranch(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}
	return branchToDTO(b), nil
}

// CreateBranch 创建分支
func (s *GenealogyService) CreateBranch(ctx context.Context, req *CreateBranchRequest) (*BranchDTO, error) {
	b := &genealogy.Branch{
		GenealogyID:     req.GenealogyID,
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		GenerationStart: req.GenerationStart,
		GenerationEnd:   req.GenerationEnd,
	}

	if err := s.domainService.CreateBranch(ctx, b); err != nil {
		return nil, err
	}

	return branchToDTO(b), nil
}

// UpdateBranch 更新分支
func (s *GenealogyService) UpdateBranch(ctx context.Context, id int64, req *UpdateBranchRequest) (*BranchDTO, error) {
	existing, err := s.domainService.GetBranch(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Code != nil {
		existing.Code = *req.Code
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.GenerationStart != nil {
		existing.GenerationStart = *req.GenerationStart
	}
	if req.GenerationEnd != nil {
		existing.GenerationEnd = *req.GenerationEnd
	}

	if err := s.domainService.UpdateBranch(ctx, existing); err != nil {
		return nil, err
	}

	return branchToDTO(existing), nil
}

// DeleteBranch 删除分支
func (s *GenealogyService) DeleteBranch(ctx context.Context, id int64) error {
	return s.domainService.DeleteBranch(ctx, id)
}

// ListGenerations 获取族谱下所有世代
func (s *GenealogyService) ListGenerations(ctx context.Context, genealogyID int64) ([]*GenerationDTO, error) {
	generations, err := s.domainService.ListGenerations(ctx, genealogyID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*GenerationDTO, len(generations))
	for i, g := range generations {
		dtos[i] = generationToDTO(g)
	}
	return dtos, nil
}

// GetGeneration 获取世代详情
func (s *GenealogyService) GetGeneration(ctx context.Context, id int64) (*GenerationDTO, error) {
	g, err := s.domainService.GetGeneration(ctx, id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, nil
	}
	return generationToDTO(g), nil
}

// CreateGeneration 创建世代
func (s *GenealogyService) CreateGeneration(ctx context.Context, req *CreateGenerationRequest) (*GenerationDTO, error) {
	g := &genealogy.Generation{
		GenealogyID: req.GenealogyID,
		Generation:  req.Generation,
		Name:        req.Name,
		Sequence:    req.Sequence,
		Description: req.Description,
		StartYear:   req.StartYear,
		EndYear:     req.EndYear,
	}

	if err := s.domainService.CreateGeneration(ctx, g); err != nil {
		return nil, err
	}

	return generationToDTO(g), nil
}

// UpdateGeneration 更新世代
func (s *GenealogyService) UpdateGeneration(ctx context.Context, id int64, req *UpdateGenerationRequest) (*GenerationDTO, error) {
	existing, err := s.domainService.GetGeneration(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Sequence != nil {
		existing.Sequence = *req.Sequence
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.StartYear != nil {
		existing.StartYear = req.StartYear
	}
	if req.EndYear != nil {
		existing.EndYear = req.EndYear
	}

	if err := s.domainService.UpdateGeneration(ctx, existing); err != nil {
		return nil, err
	}

	return generationToDTO(existing), nil
}

// DeleteGeneration 删除世代
func (s *GenealogyService) DeleteGeneration(ctx context.Context, id int64) error {
	return s.domainService.DeleteGeneration(ctx, id)
}

// GetGenerationByNumber 根据世代序号获取世代
func (s *GenealogyService) GetGenerationByNumber(ctx context.Context, genealogyID int64, generation int) (*GenerationDTO, error) {
	g, err := s.domainService.GetGenerationByNumber(ctx, genealogyID, generation)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, nil
	}
	return generationToDTO(g), nil
}

// ===== 辅助函数 =====

func genealogyToDTO(g *genealogy.Genealogy) *GenealogyDTO {
	dto := &GenealogyDTO{
		ID:              g.ID,
		Name:            g.Name,
		Surname:         g.Surname,
		Description:     g.Description,
		OriginPlace:     g.OriginPlace,
		TotalGenerations: g.TotalGenerations,
		TotalMembers:    g.TotalMembers,
		Version:         g.Version,
		ImageURL:        g.ImageURL,
		CreatedAt:       g.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       g.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if len(g.Branches) > 0 {
		dto.Branches = make([]*BranchDTO, len(g.Branches))
		for i, b := range g.Branches {
			dto.Branches[i] = branchToDTO(b)
		}
	}

	if len(g.Generations) > 0 {
		dto.Generations = make([]*GenerationDTO, len(g.Generations))
		for i, gen := range g.Generations {
			dto.Generations[i] = generationToDTO(gen)
		}
	}

	return dto
}

func branchToDTO(b *genealogy.Branch) *BranchDTO {
	dto := &BranchDTO{
		ID:              b.ID,
		GenealogyID:     b.GenealogyID,
		Name:            b.Name,
		Code:            b.Code,
		Description:     b.Description,
		GenerationStart: b.GenerationStart,
		GenerationEnd:   b.GenerationEnd,
		MemberCount:     b.MemberCount,
		CreatedAt:       b.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       b.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if len(b.Members) > 0 {
		dto.Members = make([]*PersonSummaryDTO, len(b.Members))
		for i, m := range b.Members {
			dto.Members[i] = personSummaryToDTO(m)
		}
	}

	return dto
}

func generationToDTO(g *genealogy.Generation) *GenerationDTO {
	return &GenerationDTO{
		ID:          g.ID,
		GenealogyID: g.GenealogyID,
		Generation:  g.Generation,
		Name:        g.Name,
		Sequence:    g.Sequence,
		Description: g.Description,
		StartYear:   g.StartYear,
		EndYear:     g.EndYear,
		MemberCount: g.MemberCount,
		CreatedAt:   g.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   g.UpdatedAt.Format("2006-01-02 15:04:05"),
		Active:      g.IsActive(),
	}
}

func personSummaryToDTO(p *genealogy.PersonSummary) *PersonSummaryDTO {
	return &PersonSummaryDTO{
		ID:         p.ID,
		Name:       p.Name,
		Gender:     string(p.Gender),
		Generation: p.Generation,
		BirthYear:  p.BirthYear,
		DeathYear:  p.DeathYear,
		StyleName:  p.StyleName,
	}
}