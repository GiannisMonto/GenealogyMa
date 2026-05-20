package service

import (
	"context"

	"github.com/genealogy-ma/platform/internal/domain/cemetery"
)

// CemeteryService 墓园应用服务
type CemeteryService struct {
	domainService *cemetery.Service
}

// NewCemeteryService 创建墓园应用服务
func NewCemeteryService(repo cemetery.Repository, graveRepo cemetery.GraveRepository) *CemeteryService {
	return &CemeteryService{
		domainService: cemetery.NewService(repo, graveRepo),
	}
}

// ===== DTO 定义 =====

// CemeteryDTO 墓园DTO
type CemeteryDTO struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Province    string  `json:"province"`
	City        string  `json:"city"`
	District    string  `json:"district"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	TotalGrave  int     `json:"total_grave"`
	UsedGrave   int     `json:"used_grave"`
	ImageURL    string  `json:"image_url"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// GraveDTO 墓位DTO
type GraveDTO struct {
	ID          int64   `json:"id"`
	CemeteryID  int64   `json:"cemetery_id"`
	PersonID    *int64  `json:"person_id"`
	Section     string  `json:"section"`
	Row         int     `json:"row"`
	Number      int     `json:"number"`
	Status      string  `json:"status"`
	BuriedName  string  `json:"buried_name"`
	BuriedDate  *string `json:"buried_date,omitempty"`
	BuriedYear  *int    `json:"buried_year,omitempty"`
	Note        string  `json:"note"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// CreateCemeteryRequest 创建墓园请求
type CreateCemeteryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Province    string  `json:"province" binding:"required"`
	City        string  `json:"city" binding:"required"`
	District    string  `json:"district"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	TotalGrave  int     `json:"total_grave"`
	ImageURL    string  `json:"image_url"`
}

// UpdateCemeteryRequest 更新墓园请求
type UpdateCemeteryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Province    *string `json:"province"`
	City        *string `json:"city"`
	District    *string `json:"district"`
	Address     *string `json:"address"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	TotalGrave  *int    `json:"total_grave"`
	ImageURL    *string `json:"image_url"`
}

// CreateGraveRequest 创建墓位请求
type CreateGraveRequest struct {
	CemeteryID  int64  `json:"cemetery_id" binding:"required"`
	Section     string `json:"section" binding:"required"`
	Row         int    `json:"row" binding:"required,min=1"`
	Number      int    `json:"number" binding:"required,min=1"`
	Status      string `json:"status"`
	BuriedName  string `json:"buried_name"`
	BuriedDate  string `json:"buried_date"`
	BuriedYear  *int   `json:"buried_year"`
	Note        string `json:"note"`
}

// UpdateGraveRequest 更新墓位请求
type UpdateGraveRequest struct {
	Section     *string `json:"section"`
	Row         *int    `json:"row"`
	Number      *int    `json:"number"`
	Status      *string `json:"status"`
	BuriedName  *string `json:"buried_name"`
	BuriedDate  *string `json:"buried_date"`
	BuriedYear  *int    `json:"buried_year"`
	Note        *string `json:"note"`
}

// ===== 服务方法 =====

// ListCemeteries 获取所有墓园
func (s *CemeteryService) ListCemeteries(ctx context.Context) ([]*CemeteryDTO, error) {
	cemeteries, err := s.domainService.ListCemeteries(ctx)
	if err != nil {
		return nil, err
	}
	dtos := make([]*CemeteryDTO, len(cemeteries))
	for i, c := range cemeteries {
		dtos[i] = cemeteryToDTO(c)
	}
	return dtos, nil
}

// GetCemetery 获取墓园详情
func (s *CemeteryService) GetCemetery(ctx context.Context, id int64) (*CemeteryDTO, error) {
	c, err := s.domainService.GetCemetery(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	return cemeteryToDTO(c), nil
}

// CreateCemetery 创建墓园
func (s *CemeteryService) CreateCemetery(ctx context.Context, req *CreateCemeteryRequest) (*CemeteryDTO, error) {
	c := &cemetery.Cemetery{
		Name:        req.Name,
		Description: req.Description,
		Province:    req.Province,
		City:        req.City,
		District:    req.District,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		TotalGrave:  req.TotalGrave,
		ImageURL:    req.ImageURL,
	}

	if err := s.domainService.CreateCemetery(ctx, c); err != nil {
		return nil, err
	}

	return cemeteryToDTO(c), nil
}

// UpdateCemetery 更新墓园
func (s *CemeteryService) UpdateCemetery(ctx context.Context, id int64, req *UpdateCemeteryRequest) (*CemeteryDTO, error) {
	existing, err := s.domainService.GetCemetery(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	// 应用更新
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Province != nil {
		existing.Province = *req.Province
	}
	if req.City != nil {
		existing.City = *req.City
	}
	if req.District != nil {
		existing.District = *req.District
	}
	if req.Address != nil {
		existing.Address = *req.Address
	}
	if req.Latitude != nil {
		existing.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		existing.Longitude = *req.Longitude
	}
	if req.TotalGrave != nil {
		existing.TotalGrave = *req.TotalGrave
	}
	if req.ImageURL != nil {
		existing.ImageURL = *req.ImageURL
	}

	if err := s.domainService.UpdateCemetery(ctx, existing); err != nil {
		return nil, err
	}

	return cemeteryToDTO(existing), nil
}

// DeleteCemetery 删除墓园
func (s *CemeteryService) DeleteCemetery(ctx context.Context, id int64) error {
	return s.domainService.DeleteCemetery(ctx, id)
}

// ListGravesByCemetery 获取墓园下的所有墓位
func (s *CemeteryService) ListGravesByCemetery(ctx context.Context, cemeteryID int64) ([]*GraveDTO, error) {
	graves, err := s.domainService.ListGravesByCemetery(ctx, cemeteryID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*GraveDTO, len(graves))
	for i, g := range graves {
		dtos[i] = graveToDTO(g)
	}
	return dtos, nil
}

// GetGrave 获取墓位详情
func (s *CemeteryService) GetGrave(ctx context.Context, id int64) (*GraveDTO, error) {
	g, err := s.domainService.GetGrave(ctx, id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, nil
	}
	return graveToDTO(g), nil
}

// CreateGrave 创建墓位
func (s *CemeteryService) CreateGrave(ctx context.Context, req *CreateGraveRequest) (*GraveDTO, error) {
	g := &cemetery.Grave{
		CemeteryID: req.CemeteryID,
		Section:    req.Section,
		Row:        req.Row,
		Number:     req.Number,
		Status:     cemetery.GraveStatusAvailable,
		BuriedName: req.BuriedName,
		BuriedYear: req.BuriedYear,
		Note:       req.Note,
	}

	if req.Status != "" {
		g.Status = cemetery.GraveStatus(req.Status)
	}

	if err := s.domainService.CreateGrave(ctx, g); err != nil {
		return nil, err
	}

	return graveToDTO(g), nil
}

// UpdateGrave 更新墓位
func (s *CemeteryService) UpdateGrave(ctx context.Context, id int64, req *UpdateGraveRequest) (*GraveDTO, error) {
	existing, err := s.domainService.GetGrave(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	// 应用更新
	if req.Section != nil {
		existing.Section = *req.Section
	}
	if req.Row != nil {
		existing.Row = *req.Row
	}
	if req.Number != nil {
		existing.Number = *req.Number
	}
	if req.Status != nil {
		existing.Status = cemetery.GraveStatus(*req.Status)
	}
	if req.BuriedName != nil {
		existing.BuriedName = *req.BuriedName
	}
	if req.BuriedYear != nil {
		existing.BuriedYear = req.BuriedYear
	}
	if req.Note != nil {
		existing.Note = *req.Note
	}

	if err := s.domainService.UpdateGrave(ctx, existing); err != nil {
		return nil, err
	}

	return graveToDTO(existing), nil
}

// DeleteGrave 删除墓位
func (s *CemeteryService) DeleteGrave(ctx context.Context, id int64) error {
	return s.domainService.DeleteGrave(ctx, id)
}

// ===== 辅助函数 =====

func cemeteryToDTO(c *cemetery.Cemetery) *CemeteryDTO {
	return &CemeteryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Province:    c.Province,
		City:        c.City,
		District:    c.District,
		Address:     c.Address,
		Latitude:    c.Latitude,
		Longitude:   c.Longitude,
		TotalGrave:  c.TotalGrave,
		UsedGrave:   c.UsedGrave,
		ImageURL:    c.ImageURL,
		CreatedAt:   c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func graveToDTO(g *cemetery.Grave) *GraveDTO {
	dto := &GraveDTO{
		ID:         g.ID,
		CemeteryID: g.CemeteryID,
		PersonID:   g.PersonID,
		Section:    g.Section,
		Row:        g.Row,
		Number:     g.Number,
		Status:     string(g.Status),
		BuriedName: g.BuriedName,
		BuriedYear: g.BuriedYear,
		Note:       g.Note,
		CreatedAt:  g.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  g.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if g.BuriedDate != nil {
		dateStr := g.BuriedDate.Format("2006-01-02")
		dto.BuriedDate = &dateStr
	}
	return dto
}
