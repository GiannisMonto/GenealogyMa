package service

import (
	"context"

	"github.com/genealogy-ma/platform/internal/domain/memorial"
)

// MemorialService 宗祠应用服务
type MemorialService struct {
	domainService *memorial.Service
}

// NewMemorialService 创建宗祠应用服务
func NewMemorialService(repo memorial.Repository, tabletRepo memorial.TabletRepository) *MemorialService {
	return &MemorialService{
		domainService: memorial.NewService(repo, tabletRepo),
	}
}

// ===== DTO 定义 =====

// MemorialHallDTO 宗祠DTO
type MemorialHallDTO struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Province     string         `json:"province"`
	City         string         `json:"city"`
	District     string         `json:"district"`
	Address      string         `json:"address"`
	Latitude     float64        `json:"latitude"`
	Longitude    float64        `json:"longitude"`
	BuildYear    *int           `json:"build_year,omitempty"`
	Style        string         `json:"style"`
	ImageURL     string         `json:"image_url"`
	TotalTablet  int            `json:"total_tablet"`
	UsedTablet   int            `json:"used_tablet"`
	UsageRate    float64        `json:"usage_rate"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
	Tables       []*TabletDTO   `json:"tablets,omitempty"`
}

// TabletDTO 牌位DTO
type TabletDTO struct {
	ID           int64   `json:"id"`
	HallID       int64   `json:"hall_id"`
	PersonID     *int64  `json:"person_id,omitempty"`
	PersonName   string  `json:"person_name"`
	Generation   int     `json:"generation"`
	TabletType   string  `json:"tablet_type"`
	Position     string  `json:"position"`
	Floor        int     `json:"floor"`
	Row          int     `json:"row"`
	Number       int     `json:"number"`
	Enthronement *string `json:"entronement,omitempty"`
	Note         string  `json:"note"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	Person       *PersonInfoDTO `json:"person,omitempty"`
}

// PersonInfoDTO 人物简要信息DTO
type PersonInfoDTO struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	StyleName string  `json:"style_name"`
	Gender    string  `json:"gender"`
	BirthYear *int    `json:"birth_year,omitempty"`
	DeathYear *int    `json:"death_year,omitempty"`
}

// CreateHallRequest 创建宗祠请求
type CreateHallRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Province    string  `json:"province" binding:"required"`
	City        string  `json:"city" binding:"required"`
	District    string  `json:"district"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	BuildYear   *int    `json:"build_year"`
	Style       string  `json:"style"`
	ImageURL    string  `json:"image_url"`
	TotalTablet int     `json:"total_tablet"`
}

// UpdateHallRequest 更新宗祠请求
type UpdateHallRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Province    *string `json:"province"`
	City        *string `json:"city"`
	District    *string `json:"district"`
	Address     *string `json:"address"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	BuildYear   *int    `json:"build_year"`
	Style       *string `json:"style"`
	ImageURL    *string `json:"image_url"`
	TotalTablet *int    `json:"total_tablet"`
}

// CreateTabletRequest 创建牌位请求
type CreateTabletRequest struct {
	HallID       int64  `json:"hall_id" binding:"required"`
	PersonID     *int64 `json:"person_id"`
	PersonName   string `json:"person_name" binding:"required"`
	Generation   int    `json:"generation"`
	TabletType   string `json:"tablet_type"`
	Position     string `json:"position"`
	Floor        int    `json:"floor"`
	Row          int    `json:"row"`
	Number       int    `json:"number" binding:"required,min=1"`
	Enthronement string `json:"entronement"`
	Note         string `json:"note"`
}

// UpdateTabletRequest 更新牌位请求
type UpdateTabletRequest struct {
	PersonID     *int64  `json:"person_id"`
	PersonName   *string `json:"person_name"`
	Generation   *int    `json:"generation"`
	TabletType   *string `json:"tablet_type"`
	Position     *string `json:"position"`
	Floor        *int    `json:"floor"`
	Row          *int    `json:"row"`
	Number       *int    `json:"number"`
	Enthronement *string `json:"entronement"`
	Note         *string `json:"note"`
}

// ===== 服务方法 =====

// ListHalls 获取所有宗祠
func (s *MemorialService) ListHalls(ctx context.Context) ([]*MemorialHallDTO, error) {
	halls, err := s.domainService.ListHalls(ctx)
	if err != nil {
		return nil, err
	}
	dtos := make([]*MemorialHallDTO, len(halls))
	for i, h := range halls {
		dtos[i] = hallToDTO(h)
	}
	return dtos, nil
}

// GetHall 获取宗祠详情
func (s *MemorialService) GetHall(ctx context.Context, id int64) (*MemorialHallDTO, error) {
	h, err := s.domainService.GetHall(ctx, id)
	if err != nil {
		return nil, err
	}
	if h == nil {
		return nil, nil
	}
	return hallToDTO(h), nil
}

// GetHallWithTablets 获取宗祠详情（含牌位）
func (s *MemorialService) GetHallWithTablets(ctx context.Context, id int64) (*MemorialHallDTO, error) {
	h, err := s.domainService.GetHallWithTablets(ctx, id)
	if err != nil {
		return nil, err
	}
	if h == nil {
		return nil, nil
	}
	return hallToDTO(h), nil
}

// CreateHall 创建宗祠
func (s *MemorialService) CreateHall(ctx context.Context, req *CreateHallRequest) (*MemorialHallDTO, error) {
	h := &memorial.MemorialHall{
		Name:        req.Name,
		Description: req.Description,
		Province:    req.Province,
		City:        req.City,
		District:    req.District,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		BuildYear:   req.BuildYear,
		Style:       req.Style,
		ImageURL:    req.ImageURL,
		TotalTablet: req.TotalTablet,
	}

	if err := s.domainService.CreateHall(ctx, h); err != nil {
		return nil, err
	}

	return hallToDTO(h), nil
}

// UpdateHall 更新宗祠
func (s *MemorialService) UpdateHall(ctx context.Context, id int64, req *UpdateHallRequest) (*MemorialHallDTO, error) {
	existing, err := s.domainService.GetHall(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

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
	if req.BuildYear != nil {
		existing.BuildYear = req.BuildYear
	}
	if req.Style != nil {
		existing.Style = *req.Style
	}
	if req.ImageURL != nil {
		existing.ImageURL = *req.ImageURL
	}
	if req.TotalTablet != nil {
		existing.TotalTablet = *req.TotalTablet
	}

	if err := s.domainService.UpdateHall(ctx, existing); err != nil {
		return nil, err
	}

	return hallToDTO(existing), nil
}

// DeleteHall 删除宗祠
func (s *MemorialService) DeleteHall(ctx context.Context, id int64) error {
	return s.domainService.DeleteHall(ctx, id)
}

// ListTablets 获取宗祠下所有牌位
func (s *MemorialService) ListTablets(ctx context.Context, hallID int64) ([]*TabletDTO, error) {
	tablets, err := s.domainService.ListTablets(ctx, hallID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*TabletDTO, len(tablets))
	for i, t := range tablets {
		dtos[i] = tabletToDTO(t)
	}
	return dtos, nil
}

// GetTablet 获取牌位详情
func (s *MemorialService) GetTablet(ctx context.Context, id int64) (*TabletDTO, error) {
	t, err := s.domainService.GetTablet(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	return tabletToDTO(t), nil
}

// GetTabletByPersonID 根据人物ID获取牌位
func (s *MemorialService) GetTabletByPersonID(ctx context.Context, personID int64) (*TabletDTO, error) {
	t, err := s.domainService.GetTabletByPersonID(ctx, personID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	return tabletToDTO(t), nil
}

// CreateTablet 创建牌位
func (s *MemorialService) CreateTablet(ctx context.Context, req *CreateTabletRequest) (*TabletDTO, error) {
	tabletType := memorial.TabletTypeAncestor
	if req.TabletType != "" {
		tabletType = memorial.TabletType(req.TabletType)
	}

	t := &memorial.MemorialTablet{
		HallID:     req.HallID,
		PersonID:   req.PersonID,
		PersonName: req.PersonName,
		Generation: req.Generation,
		TabletType: tabletType,
		Position:   req.Position,
		Floor:      req.Floor,
		Row:        req.Row,
		Number:     req.Number,
		Note:       req.Note,
	}

	if err := s.domainService.CreateTablet(ctx, t); err != nil {
		return nil, err
	}

	return tabletToDTO(t), nil
}

// UpdateTablet 更新牌位
func (s *MemorialService) UpdateTablet(ctx context.Context, id int64, req *UpdateTabletRequest) (*TabletDTO, error) {
	existing, err := s.domainService.GetTablet(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	if req.PersonID != nil {
		existing.PersonID = req.PersonID
	}
	if req.PersonName != nil {
		existing.PersonName = *req.PersonName
	}
	if req.Generation != nil {
		existing.Generation = *req.Generation
	}
	if req.TabletType != nil {
		existing.TabletType = memorial.TabletType(*req.TabletType)
	}
	if req.Position != nil {
		existing.Position = *req.Position
	}
	if req.Floor != nil {
		existing.Floor = *req.Floor
	}
	if req.Row != nil {
		existing.Row = *req.Row
	}
	if req.Number != nil {
		existing.Number = *req.Number
	}
	if req.Note != nil {
		existing.Note = *req.Note
	}

	if err := s.domainService.UpdateTablet(ctx, existing); err != nil {
		return nil, err
	}

	return tabletToDTO(existing), nil
}

// DeleteTablet 删除牌位
func (s *MemorialService) DeleteTablet(ctx context.Context, id int64) error {
	return s.domainService.DeleteTablet(ctx, id)
}

// ===== 辅助函数 =====

func hallToDTO(h *memorial.MemorialHall) *MemorialHallDTO {
	dto := &MemorialHallDTO{
		ID:          h.ID,
		Name:        h.Name,
		Description: h.Description,
		Province:    h.Province,
		City:        h.City,
		District:    h.District,
		Address:     h.Address,
		Latitude:    h.Latitude,
		Longitude:   h.Longitude,
		BuildYear:   h.BuildYear,
		Style:       h.Style,
		ImageURL:    h.ImageURL,
		TotalTablet: h.TotalTablet,
		UsedTablet:  h.UsedTablet,
		UsageRate:   h.UsageRate(),
		CreatedAt:   h.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   h.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if len(h.Tables) > 0 {
		dto.Tables = make([]*TabletDTO, len(h.Tables))
		for i, t := range h.Tables {
			dto.Tables[i] = tabletToDTO(t)
		}
	}

	return dto
}

func tabletToDTO(t *memorial.MemorialTablet) *TabletDTO {
	dto := &TabletDTO{
		ID:         t.ID,
		HallID:     t.HallID,
		PersonID:   t.PersonID,
		PersonName: t.PersonName,
		Generation: t.Generation,
		TabletType: string(t.TabletType),
		Position:   t.Position,
		Floor:      t.Floor,
		Row:        t.Row,
		Number:     t.Number,
		Note:       t.Note,
		CreatedAt:  t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if t.Enthronement != nil {
		dateStr := t.Enthronement.Format("2006-01-02")
		dto.Enthronement = &dateStr
	}

	if t.Person != nil {
		dto.Person = personInfoToDTO(t.Person)
	}

	return dto
}

func personInfoToDTO(p *memorial.PersonInfo) *PersonInfoDTO {
	return &PersonInfoDTO{
		ID:        p.ID,
		Name:      p.Name,
		StyleName: p.StyleName,
		Gender:    string(p.Gender),
		BirthYear: p.BirthYear,
		DeathYear: p.DeathYear,
	}
}