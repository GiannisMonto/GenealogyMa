package service

import (
	"context"

	"github.com/genealogy-ma/platform/internal/domain/person"
)

// PersonService 人物应用服务
type PersonService struct {
	domainService *person.Service
}

// NewPersonService 创建人物应用服务
func NewPersonService(repo person.Repository) *PersonService {
	return &PersonService{
		domainService: person.NewService(repo),
	}
}

// ===== DTO 定义 =====

// PersonDTO 人物详情DTO
type PersonDTO struct {
	ID               int64       `json:"id"`
	LegacyID         string      `json:"legacy_id"`
	Name             string      `json:"name"`
	StyleName        string      `json:"style_name"`
	Gender           string      `json:"gender"`
	Generation       int         `json:"generation"`
	BirthOrder       string      `json:"birth_order"`
	FatherID         *int64      `json:"father_id,omitempty"`
	LineagePath      string      `json:"lineage_path"`
	DetailText       string      `json:"detail_text"`
	BirthTimeText    string      `json:"birth_time_text"`
	DeathTimeText    string      `json:"death_time_text"`
	BirthPlace       string      `json:"birth_place"`
	BurialPlace      string      `json:"burial_place"`
	SonCount         int         `json:"son_count"`
	DaughterCount    int         `json:"daughter_count"`
	AdoptedHeirCount int         `json:"adopted_heir_count"`
	TotalChildren    int         `json:"total_children_count"`
	Age              int         `json:"age"`
	IsAlive          bool        `json:"is_alive"`
	FullName         string      `json:"full_name"`
	CreatedAt        string      `json:"created_at"`
	UpdatedAt        string      `json:"updated_at"`

	// 关联数据
	Spouses  []*SpouseDTO `json:"spouses,omitempty"`
	Children []*ChildDTO  `json:"children,omitempty"`
	Father   *PersonDTO   `json:"father,omitempty"`
}

// SpouseDTO 配偶DTO
type SpouseDTO struct {
	ID            int64  `json:"id"`
	SpouseType    string `json:"spouse_type"`
	Name          string `json:"name"`
	BirthTimeText string `json:"birth_time_text"`
	DeathTimeText string `json:"death_time_text"`
	BirthPlace    string `json:"birth_place"`
	BurialPlace   string `json:"burial_place"`
}

// ChildDTO 子女DTO
type ChildDTO struct {
	PersonDTO
	RelationType  string `json:"relation_type"`
	BirthOrderNum int    `json:"birth_order_num"`
	IsPrimary     bool   `json:"is_primary"`
}

// CreatePersonRequest 创建人物请求
type CreatePersonRequest struct {
	Name             string  `json:"name" binding:"required"`
	StyleName        string  `json:"style_name"`
	Gender           string  `json:"gender" binding:"required,oneof=男 女"`
	Generation       int     `json:"generation"`
	BirthOrder       string  `json:"birth_order"`
	FatherID         *int64  `json:"father_id"`
	DetailText       string  `json:"detail_text"`
	BirthTimeText    string  `json:"birth_time_text"`
	DeathTimeText    string  `json:"death_time_text"`
	BirthPlace       string  `json:"birth_place"`
	BurialPlace      string  `json:"burial_place"`
}

// UpdatePersonRequest 更新人物请求
type UpdatePersonRequest struct {
	Name             string  `json:"name"`
	StyleName        string  `json:"style_name"`
	Gender           string  `json:"gender"`
	Generation       *int    `json:"generation"`
	BirthOrder       string  `json:"birth_order"`
	FatherID         *int64  `json:"father_id"`
	DetailText       string  `json:"detail_text"`
	BirthTimeText    string  `json:"birth_time_text"`
	DeathTimeText    string  `json:"death_time_text"`
	BirthPlace       string  `json:"birth_place"`
	BurialPlace      string  `json:"burial_place"`
}

// SearchPersonRequest 搜索请求
type SearchPersonRequest struct {
	Keyword    string `form:"keyword"`
	Name       string `form:"name"`
	Gender     string `form:"gender"`
	Generation *int   `form:"generation"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=20"`
	SortBy     string `form:"sort_by"`
	SortDesc   bool   `form:"sort_desc,default=false"`
}

// FamilyTreeResponse 族谱树响应
type FamilyTreeResponse struct {
	CenterPerson *PersonDTO   `json:"center_person"`
	Ancestors    []*PersonDTO `json:"ancestors"`
	Descendants  []*PersonDTO `json:"descendants"`
	Generations  int          `json:"generations"`
	TotalNodes   int          `json:"total_nodes"`
}

// StatisticsResponse 统计响应
type StatisticsResponse struct {
	TotalPersons int64            `json:"total_persons"`
	ByGeneration map[int]int64    `json:"by_generation"`
}

// ===== 转换函数 =====

func toPersonDTO(p *person.Person) *PersonDTO {
	if p == nil {
		return nil
	}
	return &PersonDTO{
		ID:               p.ID,
		LegacyID:         p.LegacyID,
		Name:             p.Name,
		StyleName:        p.StyleName,
		Gender:           string(p.Gender),
		Generation:       p.Generation,
		BirthOrder:       p.BirthOrder,
		FatherID:         p.FatherID,
		LineagePath:      p.LineagePath,
		DetailText:       p.DetailText,
		BirthTimeText:    p.BirthTimeText,
		DeathTimeText:    p.DeathTimeText,
		BirthPlace:       p.BirthPlace,
		BurialPlace:      p.BurialPlace,
		SonCount:         p.SonCount,
		DaughterCount:    p.DaughterCount,
		AdoptedHeirCount: p.AdoptedHeirCount,
		TotalChildren:    p.TotalChildren,
		Age:              p.Age(),
		IsAlive:          p.IsAlive(),
		FullName:         p.FullName(),
		CreatedAt:        p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toPersonDTOWithRelations(p *person.Person) *PersonDTO {
	dto := toPersonDTO(p)
	if dto == nil {
		return nil
	}

	// 配偶
	if len(p.Spouses) > 0 {
		dto.Spouses = make([]*SpouseDTO, len(p.Spouses))
		for i, s := range p.Spouses {
			dto.Spouses[i] = &SpouseDTO{
				ID:            s.ID,
				SpouseType:    string(s.SpouseType),
				Name:          s.Name,
				BirthTimeText: s.BirthTimeText,
				DeathTimeText: s.DeathTimeText,
				BirthPlace:    s.BirthPlace,
				BurialPlace:   s.BurialPlace,
			}
		}
	}

	// 子女
	if len(p.Children) > 0 {
		dto.Children = make([]*ChildDTO, len(p.Children))
		for i, c := range p.Children {
			dto.Children[i] = &ChildDTO{
				PersonDTO:     *toPersonDTO(c.Person),
				RelationType:  string(c.RelationType),
				BirthOrderNum: c.BirthOrderNum,
				IsPrimary:     c.IsPrimary,
			}
		}
	}

	// 父亲
	if p.Father != nil {
		dto.Father = toPersonDTO(p.Father)
	}

	return dto
}

// ===== 应用服务方法 =====

// GetPerson 获取人物详情
func (s *PersonService) GetPerson(ctx context.Context, id int64, withRelations bool) (*PersonDTO, error) {
	p, err := s.domainService.GetPerson(ctx, id, withRelations)
	if err != nil {
		return nil, err
	}
	return toPersonDTOWithRelations(p), nil
}

// CreatePerson 创建人物
func (s *PersonService) CreatePerson(ctx context.Context, req *CreatePersonRequest) (*PersonDTO, error) {
	p := &person.Person{
		Name:          req.Name,
		StyleName:     req.StyleName,
		Gender:        person.Gender(req.Gender),
		Generation:    req.Generation,
		BirthOrder:    req.BirthOrder,
		FatherID:      req.FatherID,
		DetailText:    req.DetailText,
		BirthTimeText: req.BirthTimeText,
		DeathTimeText: req.DeathTimeText,
		BirthPlace:    req.BirthPlace,
		BurialPlace:   req.BurialPlace,
	}

	if err := s.domainService.CreatePerson(ctx, p); err != nil {
		return nil, err
	}

	return toPersonDTO(p), nil
}

// UpdatePerson 更新人物
func (s *PersonService) UpdatePerson(ctx context.Context, id int64, req *UpdatePersonRequest) (*PersonDTO, error) {
	p, err := s.domainService.GetPerson(ctx, id, false)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		p.Name = req.Name
	}
	if req.StyleName != "" {
		p.StyleName = req.StyleName
	}
	if req.Gender != "" {
		p.Gender = person.Gender(req.Gender)
	}
	if req.Generation != nil {
		p.Generation = *req.Generation
	}
	if req.BirthOrder != "" {
		p.BirthOrder = req.BirthOrder
	}
	if req.FatherID != nil {
		p.FatherID = req.FatherID
	}
	p.DetailText = req.DetailText
	p.BirthTimeText = req.BirthTimeText
	p.DeathTimeText = req.DeathTimeText
	p.BirthPlace = req.BirthPlace
	p.BurialPlace = req.BurialPlace

	if err := s.domainService.UpdatePerson(ctx, p); err != nil {
		return nil, err
	}

	return toPersonDTO(p), nil
}

// DeletePerson 删除人物
func (s *PersonService) DeletePerson(ctx context.Context, id int64) error {
	return s.domainService.DeletePerson(ctx, id)
}

// SearchPersons 搜索人物
func (s *PersonService) SearchPersons(ctx context.Context, req *SearchPersonRequest) ([]*PersonDTO, int64, error) {
	var gender *person.Gender
	if req.Gender != "" {
		g := person.Gender(req.Gender)
		gender = &g
	}

	query := &person.SearchQuery{
		Keyword:    req.Keyword,
		Name:       req.Name,
		Gender:     gender,
		Generation: req.Generation,
		Page:       req.Page,
		PageSize:   req.PageSize,
		SortBy:     req.SortBy,
		SortDesc:   req.SortDesc,
	}

	persons, total, err := s.domainService.SearchPersons(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]*PersonDTO, len(persons))
	for i, p := range persons {
		dtos[i] = toPersonDTO(p)
	}

	return dtos, total, nil
}

// GetFamilyTree 获取族谱树
func (s *PersonService) GetFamilyTree(ctx context.Context, personID int64, upDepth, downDepth int) (*FamilyTreeResponse, error) {
	result, err := s.domainService.GetFamilyTree(ctx, personID, upDepth, downDepth)
	if err != nil {
		return nil, err
	}

	ancestors := make([]*PersonDTO, len(result.Ancestors))
	for i, a := range result.Ancestors {
		ancestors[i] = toPersonDTO(a)
	}

	descendants := make([]*PersonDTO, len(result.Descendants))
	for i, d := range result.Descendants {
		descendants[i] = toPersonDTO(d)
	}

	return &FamilyTreeResponse{
		CenterPerson: toPersonDTO(result.CenterPerson),
		Ancestors:    ancestors,
		Descendants:  descendants,
		Generations:  result.Generations,
		TotalNodes:   result.TotalNodes,
	}, nil
}

// GetStatistics 获取统计数据
func (s *PersonService) GetStatistics(ctx context.Context) (*StatisticsResponse, error) {
	stats, err := s.domainService.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	return &StatisticsResponse{
		TotalPersons: stats.TotalPersons,
		ByGeneration: stats.ByGeneration,
	}, nil
}

// ===== 批量操作 DTO =====

// BatchCreatePersonRequest 批量创建人物请求
type BatchCreatePersonRequest struct {
	Persons []*CreatePersonRequest `json:"persons" binding:"required,min=1,max=100"`
}

// BatchUpdatePersonRequest 批量更新人物请求
type BatchUpdatePersonRequest struct {
	Persons []*BatchUpdateItem `json:"persons" binding:"required,min=1,max=100"`
}

// BatchUpdateItem 批量更新单项
type BatchUpdateItem struct {
	ID               int64  `json:"id" binding:"required"`
	Name             string `json:"name"`
	StyleName        string `json:"style_name"`
	Gender           string `json:"gender"`
	Generation       *int   `json:"generation"`
	BirthOrder       string `json:"birth_order"`
	FatherID         *int64 `json:"father_id"`
	DetailText       string `json:"detail_text"`
	BirthTimeText    string `json:"birth_time_text"`
	DeathTimeText    string `json:"death_time_text"`
	BirthPlace       string `json:"birth_place"`
	BurialPlace      string `json:"burial_place"`
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1,max=100"`
}

// BatchResult 批量操作结果
type BatchResult struct {
	SuccessCount int               `json:"success_count"`
	FailCount    int               `json:"fail_count"`
	Results      []*BatchItemResult `json:"results"`
}

// BatchItemResult 批量操作单项结果
type BatchItemResult struct {
	ID    int64  `json:"id"`
	Success bool `json:"success"`
	Error  string `json:"error,omitempty"`
}

// BatchCreatePersons 批量创建人物
func (s *PersonService) BatchCreatePersons(ctx context.Context, req *BatchCreatePersonRequest) (*BatchResult, error) {
	result := &BatchResult{
		Results: make([]*BatchItemResult, 0, len(req.Persons)),
	}

	for _, p := range req.Persons {
		person := &person.Person{
			Name:          p.Name,
			StyleName:     p.StyleName,
			Gender:        person.Gender(p.Gender),
			Generation:    p.Generation,
			BirthOrder:    p.BirthOrder,
			FatherID:      p.FatherID,
			DetailText:    p.DetailText,
			BirthTimeText: p.BirthTimeText,
			DeathTimeText: p.DeathTimeText,
			BirthPlace:    p.BirthPlace,
			BurialPlace:   p.BurialPlace,
		}

		if err := s.domainService.CreatePerson(ctx, person); err != nil {
			result.FailCount++
			result.Results = append(result.Results, &BatchItemResult{
				ID:     0,
				Success: false,
				Error:  err.Error(),
			})
		} else {
			result.SuccessCount++
			result.Results = append(result.Results, &BatchItemResult{
				ID:     person.ID,
				Success: true,
			})
		}
	}

	return result, nil
}

// BatchUpdatePersons 批量更新人物
func (s *PersonService) BatchUpdatePersons(ctx context.Context, req *BatchUpdatePersonRequest) (*BatchResult, error) {
	result := &BatchResult{
		Results: make([]*BatchItemResult, 0, len(req.Persons)),
	}

	for _, item := range req.Persons {
		p, err := s.domainService.GetPerson(ctx, item.ID, false)
		if err != nil {
			result.FailCount++
			result.Results = append(result.Results, &BatchItemResult{
				ID:     item.ID,
				Success: false,
				Error:  "person not found",
			})
			continue
		}

		if item.Name != "" {
			p.Name = item.Name
		}
		if item.StyleName != "" {
			p.StyleName = item.StyleName
		}
		if item.Gender != "" {
			p.Gender = person.Gender(item.Gender)
		}
		if item.Generation != nil {
			p.Generation = *item.Generation
		}
		if item.BirthOrder != "" {
			p.BirthOrder = item.BirthOrder
		}
		if item.FatherID != nil {
			p.FatherID = item.FatherID
		}
		p.DetailText = item.DetailText
		p.BirthTimeText = item.BirthTimeText
		p.DeathTimeText = item.DeathTimeText
		p.BirthPlace = item.BirthPlace
		p.BurialPlace = item.BurialPlace

		if err := s.domainService.UpdatePerson(ctx, p); err != nil {
			result.FailCount++
			result.Results = append(result.Results, &BatchItemResult{
				ID:     item.ID,
				Success: false,
				Error:  err.Error(),
			})
		} else {
			result.SuccessCount++
			result.Results = append(result.Results, &BatchItemResult{
				ID:     item.ID,
				Success: true,
			})
		}
	}

	return result, nil
}

// BatchDeletePersons 批量删除人物
func (s *PersonService) BatchDeletePersons(ctx context.Context, req *BatchDeleteRequest) (*BatchResult, error) {
	result := &BatchResult{
		Results: make([]*BatchItemResult, 0, len(req.IDs)),
	}

	for _, id := range req.IDs {
		if err := s.domainService.DeletePerson(ctx, id); err != nil {
			result.FailCount++
			result.Results = append(result.Results, &BatchItemResult{
				ID:     id,
				Success: false,
				Error:  err.Error(),
			})
		} else {
			result.SuccessCount++
			result.Results = append(result.Results, &BatchItemResult{
				ID:     id,
				Success: true,
			})
		}
	}

	return result, nil
}
