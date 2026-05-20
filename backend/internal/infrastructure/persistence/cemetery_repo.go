package persistence

import (
	"context"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/cemetery"
	"gorm.io/gorm"
)

// CemeteryModel GORM模型对应cemeteries表
type CemeteryModel struct {
	ID          int64   `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string  `gorm:"column:name;type:varchar(200);not null"`
	Description string  `gorm:"column:description;type:text"`
	Province    string  `gorm:"column:province;type:varchar(100);not null"`
	City        string  `gorm:"column:city;type:varchar(100);not null"`
	District    string  `gorm:"column:district;type:varchar(100)"`
	Address     string  `gorm:"column:address;type:text"`
	Latitude    float64 `gorm:"column:latitude;type:decimal(10,7)"`
	Longitude   float64 `gorm:"column:longitude;type:decimal(10,7)"`
	TotalGrave  int     `gorm:"column:total_grave;default:0"`
	UsedGrave   int     `gorm:"column:used_grave;default:0"`
	ImageURL    string  `gorm:"column:image_url;type:varchar(500)"`
	CreatedAt   string  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   string  `gorm:"column:updated_at;autoUpdateTime"`
}

func (CemeteryModel) TableName() string {
	return "cemeteries"
}

// GraveModel GORM模型对应graves表
type GraveModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	CemeteryID int64  `gorm:"column:cemetery_id;not null"`
	PersonID   *int64 `gorm:"column:person_id"`
	Section    string `gorm:"column:section;type:varchar(50);not null"`
	Row        int    `gorm:"column:row;not null"`
	Number     int    `gorm:"column:number;not null"`
	Status     string `gorm:"column:status;type:varchar(20);default:'available'"`
	BuriedName string `gorm:"column:buried_name;type:varchar(200)"`
	BuriedDate *string `gorm:"column:buried_date;type:date"`
	BuriedYear *int    `gorm:"column:buried_year"`
	Note       string `gorm:"column:note;type:text"`
	CreatedAt  string `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  string `gorm:"column:updated_at;autoUpdateTime"`
}

func (GraveModel) TableName() string {
	return "graves"
}

// CemeteryRepositoryImpl 墓园仓储实现
type CemeteryRepositoryImpl struct {
	db *gorm.DB
}

// NewCemeteryRepository 创建墓园仓储
func NewCemeteryRepository(db *gorm.DB) cemetery.Repository {
	return &CemeteryRepositoryImpl{db: db}
}

// modelToCemetery converts a database model to domain entity
func modelToCemetery(m *CemeteryModel) *cemetery.Cemetery {
	return &cemetery.Cemetery{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Province:    m.Province,
		City:        m.City,
		District:    m.District,
		Address:     m.Address,
		Latitude:    m.Latitude,
		Longitude:   m.Longitude,
		TotalGrave:  m.TotalGrave,
		UsedGrave:   m.UsedGrave,
		ImageURL:    m.ImageURL,
	}
}

// cemeteryToModel converts a domain entity to database model
func cemeteryToModel(c *cemetery.Cemetery) *CemeteryModel {
	return &CemeteryModel{
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
	}
}

// FindByID implements cemetery.Repository
func (r *CemeteryRepositoryImpl) FindByID(ctx context.Context, id int64) (*cemetery.Cemetery, error) {
	var model CemeteryModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToCemetery(&model), nil
}

// Create implements cemetery.Repository
func (r *CemeteryRepositoryImpl) Create(ctx context.Context, c *cemetery.Cemetery) error {
	model := cemeteryToModel(c)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update implements cemetery.Repository
func (r *CemeteryRepositoryImpl) Update(ctx context.Context, c *cemetery.Cemetery) error {
	model := cemeteryToModel(c)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete implements cemetery.Repository
func (r *CemeteryRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&CemeteryModel{}, id).Error
}

// FindAll implements cemetery.Repository
func (r *CemeteryRepositoryImpl) FindAll(ctx context.Context) ([]*cemetery.Cemetery, error) {
	var models []CemeteryModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*cemetery.Cemetery, len(models))
	for i := range models {
		result[i] = modelToCemetery(&models[i])
	}
	return result, nil
}

// Search implements cemetery.Repository
func (r *CemeteryRepositoryImpl) Search(ctx context.Context, query *cemetery.SearchQuery) ([]*cemetery.Cemetery, int64, error) {
	var models []CemeteryModel
	var total int64

	db := r.db.WithContext(ctx).Model(&CemeteryModel{})

	// Apply filters
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Province != "" {
		db = db.Where("province = ?", query.Province)
	}
	if query.City != "" {
		db = db.Where("city = ?", query.City)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if query.SortBy != "" {
		order := query.SortBy
		if query.SortDesc {
			order += " DESC"
		} else {
			order += " ASC"
		}
		db = db.Order(order)
	}

	// Apply pagination
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*cemetery.Cemetery, len(models))
	for i := range models {
		result[i] = modelToCemetery(&models[i])
	}
	return result, total, nil
}

// FindByRegion implements cemetery.Repository
func (r *CemeteryRepositoryImpl) FindByRegion(ctx context.Context, province, city string) ([]*cemetery.Cemetery, error) {
	var models []CemeteryModel
	db := r.db.WithContext(ctx)
	if province != "" {
		db = db.Where("province = ?", province)
	}
	if city != "" {
		db = db.Where("city = ?", city)
	}
	if err := db.Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*cemetery.Cemetery, len(models))
	for i := range models {
		result[i] = modelToCemetery(&models[i])
	}
	return result, nil
}

// GraveRepositoryImpl 墓位仓储实现
type GraveRepositoryImpl struct {
	db *gorm.DB
}

// NewGraveRepository 创建墓位仓储
func NewGraveRepository(db *gorm.DB) cemetery.GraveRepository {
	return &GraveRepositoryImpl{db: db}
}

// modelToGrave converts a database model to domain entity
func modelToGrave(m *GraveModel) *cemetery.Grave {
	status := cemetery.GraveStatus(m.Status)
	g := &cemetery.Grave{
		ID:         m.ID,
		CemeteryID: m.CemeteryID,
		PersonID:   m.PersonID,
		Section:    m.Section,
		Row:        m.Row,
		Number:     m.Number,
		Status:     status,
		BuriedName: m.BuriedName,
		BuriedYear: m.BuriedYear,
		Note:       m.Note,
	}
	if m.BuriedDate != nil && *m.BuriedDate != "" {
		g.BuriedDate = toPtrTime(*m.BuriedDate)
	}
	return g
}

// graveToModel converts a domain entity to database model
func graveToModel(g *cemetery.Grave) *GraveModel {
	model := &GraveModel{
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
	}
	if g.BuriedDate != nil {
		dateStr := g.BuriedDate.Format("2006-01-02")
		model.BuriedDate = &dateStr
	}
	return model
}

// toPtrTime is a helper to convert string to time.Time pointer
func toPtrTime(s string) *time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

// FindByID implements cemetery.GraveRepository
func (r *GraveRepositoryImpl) FindByID(ctx context.Context, id int64) (*cemetery.Grave, error) {
	var model GraveModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToGrave(&model), nil
}

// FindByCemeteryID implements cemetery.GraveRepository
func (r *GraveRepositoryImpl) FindByCemeteryID(ctx context.Context, cemeteryID int64) ([]*cemetery.Grave, error) {
	var models []GraveModel
	if err := r.db.WithContext(ctx).Where("cemetery_id = ?", cemeteryID).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*cemetery.Grave, len(models))
	for i := range models {
		result[i] = modelToGrave(&models[i])
	}
	return result, nil
}

// FindByPersonID implements cemetery.GraveRepository
func (r *GraveRepositoryImpl) FindByPersonID(ctx context.Context, personID int64) (*cemetery.Grave, error) {
	var model GraveModel
	if err := r.db.WithContext(ctx).Where("person_id = ?", personID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToGrave(&model), nil
}

// Create implements cemetery.GraveRepository
func (r *GraveRepositoryImpl) Create(ctx context.Context, grave *cemetery.Grave) error {
	model := graveToModel(grave)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update implements cemetery.GraveRepository
func (r *GraveRepositoryImpl) Update(ctx context.Context, grave *cemetery.Grave) error {
	model := graveToModel(grave)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete implements cemetery.GraveRepository
func (r *GraveRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&GraveModel{}, id).Error
}

// DeleteByCemeteryID implements cemetery.GraveRepository
func (r *GraveRepositoryImpl) DeleteByCemeteryID(ctx context.Context, cemeteryID int64) error {
	return r.db.WithContext(ctx).Where("cemetery_id = ?", cemeteryID).Delete(&GraveModel{}).Error
}

// Search implements cemetery.GraveRepository
func (r *GraveRepositoryImpl) Search(ctx context.Context, query *cemetery.GraveSearchQuery) ([]*cemetery.Grave, int64, error) {
	var models []GraveModel
	var total int64

	db := r.db.WithContext(ctx).Model(&GraveModel{})

	// Apply filters
	if query.CemeteryID != nil {
		db = db.Where("cemetery_id = ?", *query.CemeteryID)
	}
	if query.PersonID != nil {
		db = db.Where("person_id = ?", *query.PersonID)
	}
	if query.Section != "" {
		db = db.Where("section = ?", query.Section)
	}
	if query.Row != nil {
		db = db.Where("row = ?", *query.Row)
	}
	if query.Number != nil {
		db = db.Where("number = ?", *query.Number)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*cemetery.Grave, len(models))
	for i := range models {
		result[i] = modelToGrave(&models[i])
	}
	return result, total, nil
}
