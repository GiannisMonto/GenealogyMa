package persistence

import (
	"context"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/genealogy"
	"gorm.io/gorm"
)

// GenealogyModel GORM模型对应genealogies表
type GenealogyModel struct {
	ID              int64   `gorm:"column:id;primaryKey;autoIncrement"`
	Name            string  `gorm:"column:name;type:varchar(200);not null"`
	Surname         string  `gorm:"column:surname;type:varchar(100);not null"`
	Description     string  `gorm:"column:description;type:text"`
	OriginPlace     string  `gorm:"column:origin_place;type:varchar(200)"`
	TotalGenerations int    `gorm:"column:total_generations;default:0"`
	TotalMembers    int     `gorm:"column:total_members;default:0"`
	Version         string  `gorm:"column:version;type:varchar(50)"`
	ImageURL        string  `gorm:"column:image_url;type:varchar(500)"`
	CreatedAt       string  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       string  `gorm:"column:updated_at;autoUpdateTime"`
}

func (GenealogyModel) TableName() string {
	return "genealogies"
}

// BranchModel GORM模型对应branches表
type BranchModel struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement"`
	GenealogyID     int64     `gorm:"column:genealogy_id;not null"`
	Name            string    `gorm:"column:name;type:varchar(100);not null"`
	Code            string    `gorm:"column:code;type:varchar(20)"`
	Description     string    `gorm:"column:description;type:text"`
	GenerationStart int       `gorm:"column:generation_start;default:1"`
	GenerationEnd   int       `gorm:"column:generation_end;default:0"`
	MemberCount     int       `gorm:"column:member_count;default:0"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (BranchModel) TableName() string {
	return "branches"
}

// GenerationModel GORM模型对应generations表
type GenerationModel struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement"`
	GenealogyID  int64  `gorm:"column:genealogy_id;not null"`
	Generation   int    `gorm:"column:generation;not null"`
	Name         string `gorm:"column:name;type:varchar(50);not null"`
	Sequence     int    `gorm:"column:sequence;default:0"`
	Description  string `gorm:"column:description;type:text"`
	StartYear    *int   `gorm:"column:start_year"`
	EndYear      *int   `gorm:"column:end_year"`
	MemberCount  int    `gorm:"column:member_count;default:0"`
	CreatedAt    string `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    string `gorm:"column:updated_at;autoUpdateTime"`
}

func (GenerationModel) TableName() string {
	return "generations"
}

// GenealogyRepositoryImpl 族谱仓储实现
type GenealogyRepositoryImpl struct {
	db *gorm.DB
}

// NewGenealogyRepository 创建族谱仓储
func NewGenealogyRepository(db *gorm.DB) genealogy.Repository {
	return &GenealogyRepositoryImpl{db: db}
}

func modelToGenealogy(m *GenealogyModel) *genealogy.Genealogy {
	return &genealogy.Genealogy{
		ID:              m.ID,
		Name:            m.Name,
		Surname:         m.Surname,
		Description:     m.Description,
		OriginPlace:     m.OriginPlace,
		TotalGenerations: m.TotalGenerations,
		TotalMembers:    m.TotalMembers,
		Version:         m.Version,
		ImageURL:        m.ImageURL,
	}
}

func genealogyToModel(g *genealogy.Genealogy) *GenealogyModel {
	return &GenealogyModel{
		ID:              g.ID,
		Name:            g.Name,
		Surname:         g.Surname,
		Description:     g.Description,
		OriginPlace:     g.OriginPlace,
		TotalGenerations: g.TotalGenerations,
		TotalMembers:    g.TotalMembers,
		Version:         g.Version,
		ImageURL:        g.ImageURL,
	}
}

func modelToBranch(m *BranchModel) *genealogy.Branch {
	return &genealogy.Branch{
		ID:              m.ID,
		GenealogyID:     m.GenealogyID,
		Name:            m.Name,
		Code:            m.Code,
		Description:     m.Description,
		GenerationStart: m.GenerationStart,
		GenerationEnd:   m.GenerationEnd,
		MemberCount:     m.MemberCount,
	}
}

func branchToModel(b *genealogy.Branch) *BranchModel {
	return &BranchModel{
		ID:              b.ID,
		GenealogyID:     b.GenealogyID,
		Name:            b.Name,
		Code:            b.Code,
		Description:     b.Description,
		GenerationStart: b.GenerationStart,
		GenerationEnd:   b.GenerationEnd,
		MemberCount:     b.MemberCount,
	}
}

func modelToGeneration(m *GenerationModel) *genealogy.Generation {
	return &genealogy.Generation{
		ID:          m.ID,
		GenealogyID: m.GenealogyID,
		Generation:  m.Generation,
		Name:        m.Name,
		Sequence:    m.Sequence,
		Description: m.Description,
		StartYear:   m.StartYear,
		EndYear:     m.EndYear,
		MemberCount: m.MemberCount,
	}
}

func generationToModel(g *genealogy.Generation) *GenerationModel {
	return &GenerationModel{
		ID:          g.ID,
		GenealogyID: g.GenealogyID,
		Generation:  g.Generation,
		Name:        g.Name,
		Sequence:    g.Sequence,
		Description: g.Description,
		StartYear:   g.StartYear,
		EndYear:     g.EndYear,
		MemberCount: g.MemberCount,
	}
}

// FindByID implements genealogy.Repository
func (r *GenealogyRepositoryImpl) FindByID(ctx context.Context, id int64) (*genealogy.Genealogy, error) {
	var model GenealogyModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToGenealogy(&model), nil
}

// Create implements genealogy.Repository
func (r *GenealogyRepositoryImpl) Create(ctx context.Context, g *genealogy.Genealogy) error {
	model := genealogyToModel(g)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update implements genealogy.Repository
func (r *GenealogyRepositoryImpl) Update(ctx context.Context, g *genealogy.Genealogy) error {
	model := genealogyToModel(g)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete implements genealogy.Repository
func (r *GenealogyRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&GenealogyModel{}, id).Error
}

// FindAll implements genealogy.Repository
func (r *GenealogyRepositoryImpl) FindAll(ctx context.Context) ([]*genealogy.Genealogy, error) {
	var models []GenealogyModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*genealogy.Genealogy, len(models))
	for i, m := range models {
		result[i] = modelToGenealogy(&m)
	}
	return result, nil
}

// Search implements genealogy.Repository
func (r *GenealogyRepositoryImpl) Search(ctx context.Context, query *genealogy.SearchQuery) ([]*genealogy.Genealogy, int64, error) {
	var models []GenealogyModel
	db := r.db.WithContext(ctx)

	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Surname != "" {
		db = db.Where("surname = ?", query.Surname)
	}

	var total int64
	db.Model(&GenealogyModel{}).Count(&total)

	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	if query.SortBy != "" {
		order := query.SortBy
		if query.SortDesc {
			order += " DESC"
		} else {
			order += " ASC"
		}
		db = db.Order(order)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*genealogy.Genealogy, len(models))
	for i, m := range models {
		result[i] = modelToGenealogy(&m)
	}
	return result, total, nil
}

// FindBySurname implements genealogy.Repository
func (r *GenealogyRepositoryImpl) FindBySurname(ctx context.Context, surname string) ([]*genealogy.Genealogy, error) {
	var models []GenealogyModel
	if err := r.db.WithContext(ctx).Where("surname = ?", surname).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*genealogy.Genealogy, len(models))
	for i, m := range models {
		result[i] = modelToGenealogy(&m)
	}
	return result, nil
}

// BranchRepositoryImpl 分支仓储实现
type BranchRepositoryImpl struct {
	db *gorm.DB
}

// NewBranchRepository 创建分支仓储
func NewBranchRepository(db *gorm.DB) genealogy.BranchRepository {
	return &BranchRepositoryImpl{db: db}
}

// FindByID implements genealogy.BranchRepository
func (r *BranchRepositoryImpl) FindByID(ctx context.Context, id int64) (*genealogy.Branch, error) {
	var model BranchModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToBranch(&model), nil
}

// FindByGenealogyID implements genealogy.BranchRepository
func (r *BranchRepositoryImpl) FindByGenealogyID(ctx context.Context, genealogyID int64) ([]*genealogy.Branch, error) {
	var models []BranchModel
	if err := r.db.WithContext(ctx).Where("genealogy_id = ?", genealogyID).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*genealogy.Branch, len(models))
	for i, m := range models {
		result[i] = modelToBranch(&m)
	}
	return result, nil
}

// Create implements genealogy.BranchRepository
func (r *BranchRepositoryImpl) Create(ctx context.Context, b *genealogy.Branch) error {
	model := branchToModel(b)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update implements genealogy.BranchRepository
func (r *BranchRepositoryImpl) Update(ctx context.Context, b *genealogy.Branch) error {
	model := branchToModel(b)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete implements genealogy.BranchRepository
func (r *BranchRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&BranchModel{}, id).Error
}

// DeleteByGenealogyID implements genealogy.BranchRepository
func (r *BranchRepositoryImpl) DeleteByGenealogyID(ctx context.Context, genealogyID int64) error {
	return r.db.WithContext(ctx).Where("genealogy_id = ?", genealogyID).Delete(&BranchModel{}).Error
}

// Search implements genealogy.BranchRepository
func (r *BranchRepositoryImpl) Search(ctx context.Context, query *genealogy.BranchSearchQuery) ([]*genealogy.Branch, int64, error) {
	var models []BranchModel
	db := r.db.WithContext(ctx)

	if query.GenealogyID != nil {
		db = db.Where("genealogy_id = ?", *query.GenealogyID)
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Code != "" {
		db = db.Where("code = ?", query.Code)
	}

	var total int64
	db.Model(&BranchModel{}).Count(&total)

	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*genealogy.Branch, len(models))
	for i, m := range models {
		result[i] = modelToBranch(&m)
	}
	return result, total, nil
}

// GenerationRepositoryImpl 世代仓储实现
type GenerationRepositoryImpl struct {
	db *gorm.DB
}

// NewGenerationRepository 创建世代仓储
func NewGenerationRepository(db *gorm.DB) genealogy.GenerationRepository {
	return &GenerationRepositoryImpl{db: db}
}

// FindByID implements genealogy.GenerationRepository
func (r *GenerationRepositoryImpl) FindByID(ctx context.Context, id int64) (*genealogy.Generation, error) {
	var model GenerationModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToGeneration(&model), nil
}

// FindByGenealogyID implements genealogy.GenerationRepository
func (r *GenerationRepositoryImpl) FindByGenealogyID(ctx context.Context, genealogyID int64) ([]*genealogy.Generation, error) {
	var models []GenerationModel
	if err := r.db.WithContext(ctx).Where("genealogy_id = ?", genealogyID).Order("generation ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*genealogy.Generation, len(models))
	for i, m := range models {
		result[i] = modelToGeneration(&m)
	}
	return result, nil
}

// FindByGeneration implements genealogy.GenerationRepository
func (r *GenerationRepositoryImpl) FindByGeneration(ctx context.Context, genealogyID int64, generation int) (*genealogy.Generation, error) {
	var model GenerationModel
	if err := r.db.WithContext(ctx).Where("genealogy_id = ? AND generation = ?", genealogyID, generation).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToGeneration(&model), nil
}

// Create implements genealogy.GenerationRepository
func (r *GenerationRepositoryImpl) Create(ctx context.Context, g *genealogy.Generation) error {
	model := generationToModel(g)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update implements genealogy.GenerationRepository
func (r *GenerationRepositoryImpl) Update(ctx context.Context, g *genealogy.Generation) error {
	model := generationToModel(g)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete implements genealogy.GenerationRepository
func (r *GenerationRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&GenerationModel{}, id).Error
}

// DeleteByGenealogyID implements genealogy.GenerationRepository
func (r *GenerationRepositoryImpl) DeleteByGenealogyID(ctx context.Context, genealogyID int64) error {
	return r.db.WithContext(ctx).Where("genealogy_id = ?", genealogyID).Delete(&GenerationModel{}).Error
}