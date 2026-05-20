package persistence

import (
	"context"
	"encoding/json"

	"github.com/genealogy-ma/platform/internal/domain/culture"
	"gorm.io/gorm"
)

// DocumentModel GORM模型对应documents表
type DocumentModel struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Title        string `gorm:"column:title;type:varchar(200);not null"`
	Content      string `gorm:"column:content;type:text;not null"`
	Category     string `gorm:"column:category;type:varchar(50);not null"`
	Author       string `gorm:"column:author;type:varchar(200)"`
	CreatedYear  *int   `gorm:"column:created_year"`
	Dynasty      string `gorm:"column:dynasty;type:varchar(100)"`
	Source       string `gorm:"column:source;type:varchar(500)"`
	ImageURLs    string `gorm:"column:image_urls;type:text"`
	ViewCount    int    `gorm:"column:view_count;default:0"`
	CollectCount int    `gorm:"column:collect_count;default:0"`
	CreatedAt    string `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    string `gorm:"column:updated_at;autoUpdateTime"`
}

func (DocumentModel) TableName() string {
	return "documents"
}

// StoryModel GORM模型对应stories表
type StoryModel struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Title     string `gorm:"column:title;type:varchar(200);not null"`
	Content   string `gorm:"column:content;type:text;not null"`
	Era       string `gorm:"column:era;type:varchar(100)"`
	Category  string `gorm:"column:category;type:varchar(50)"`
	Tags      string `gorm:"column:tags;type:text"`
	AudioURL  string `gorm:"column:audio_url;type:varchar(500)"`
	ImageURL  string `gorm:"column:image_url;type:varchar(500)"`
	ViewCount int    `gorm:"column:view_count;default:0"`
	CreatedAt string `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt string `gorm:"column:updated_at;autoUpdateTime"`
}

func (StoryModel) TableName() string {
	return "stories"
}

// FamilyTeachingsModel GORM模型对应family_teachings表
type FamilyTeachingsModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Title      string `gorm:"column:title;type:varchar(200);not null"`
	Content    string `gorm:"column:content;type:text;not null"`
	Generation int    `gorm:"column:generation;default:0"`
	OriginText string `gorm:"column:origin_text;type:text"`
	Meaning    string `gorm:"column:meaning;type:text"`
	UsageCount int    `gorm:"column:usage_count;default:0"`
	CreatedAt  string `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  string `gorm:"column:updated_at;autoUpdateTime"`
}

func (FamilyTeachingsModel) TableName() string {
	return "family_teachings"
}

// ===== 文献仓储实现 =====

// DocumentRepositoryImpl 文献仓储实现
type DocumentRepositoryImpl struct {
	db *gorm.DB
}

// NewDocumentRepository 创建文献仓储
func NewDocumentRepository(db *gorm.DB) culture.Repository {
	return &DocumentRepositoryImpl{db: db}
}

func modelToDocument(m *DocumentModel) *culture.Document {
	var imageURLs []string
	if m.ImageURLs != "" {
		_ = json.Unmarshal([]byte(m.ImageURLs), &imageURLs)
	}
	return &culture.Document{
		ID:           m.ID,
		Title:        m.Title,
		Content:      m.Content,
		Category:     culture.DocumentCategory(m.Category),
		Author:       m.Author,
		CreatedYear:  m.CreatedYear,
		Dynasty:      m.Dynasty,
		Source:       m.Source,
		ImageURLs:    imageURLs,
		ViewCount:    m.ViewCount,
		CollectCount: m.CollectCount,
	}
}

func documentToModel(d *culture.Document) *DocumentModel {
	var imageURLs string
	if len(d.ImageURLs) > 0 {
		b, _ := json.Marshal(d.ImageURLs)
		imageURLs = string(b)
	}
	return &DocumentModel{
		ID:           d.ID,
		Title:        d.Title,
		Content:      d.Content,
		Category:     string(d.Category),
		Author:       d.Author,
		CreatedYear:  d.CreatedYear,
		Dynasty:      d.Dynasty,
		Source:       d.Source,
		ImageURLs:    imageURLs,
		ViewCount:    d.ViewCount,
		CollectCount: d.CollectCount,
	}
}

func (r *DocumentRepositoryImpl) FindByID(ctx context.Context, id int64) (*culture.Document, error) {
	var model DocumentModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToDocument(&model), nil
}

func (r *DocumentRepositoryImpl) Create(ctx context.Context, d *culture.Document) error {
	model := documentToModel(d)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *DocumentRepositoryImpl) Update(ctx context.Context, d *culture.Document) error {
	model := documentToModel(d)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *DocumentRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&DocumentModel{}, id).Error
}

func (r *DocumentRepositoryImpl) FindAll(ctx context.Context) ([]*culture.Document, error) {
	var models []DocumentModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*culture.Document, len(models))
	for i := range models {
		result[i] = modelToDocument(&models[i])
	}
	return result, nil
}

func (r *DocumentRepositoryImpl) Search(ctx context.Context, query *culture.SearchQuery) ([]*culture.Document, int64, error) {
	var models []DocumentModel
	var total int64

	db := r.db.WithContext(ctx).Model(&DocumentModel{})

	if query.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+query.Keyword+"%")
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.Author != "" {
		db = db.Where("author LIKE ?", "%"+query.Author+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
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

	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*culture.Document, len(models))
	for i := range models {
		result[i] = modelToDocument(&models[i])
	}
	return result, total, nil
}

func (r *DocumentRepositoryImpl) FindByCategory(ctx context.Context, category culture.DocumentCategory) ([]*culture.Document, error) {
	var models []DocumentModel
	if err := r.db.WithContext(ctx).Where("category = ?", category).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*culture.Document, len(models))
	for i := range models {
		result[i] = modelToDocument(&models[i])
	}
	return result, nil
}

// ===== 故事仓储实现 =====

// StoryRepositoryImpl 故事仓储实现
type StoryRepositoryImpl struct {
	db *gorm.DB
}

// NewStoryRepository 创建故事仓储
func NewStoryRepository(db *gorm.DB) culture.StoryRepository {
	return &StoryRepositoryImpl{db: db}
}

func modelToStory(m *StoryModel) *culture.Story {
	var tags []string
	if m.Tags != "" {
		_ = json.Unmarshal([]byte(m.Tags), &tags)
	}
	return &culture.Story{
		ID:        m.ID,
		Title:     m.Title,
		Content:   m.Content,
		Era:       m.Era,
		Category:  m.Category,
		Tags:      tags,
		AudioURL:  m.AudioURL,
		ImageURL:  m.ImageURL,
		ViewCount: m.ViewCount,
	}
}

func storyToModel(s *culture.Story) *StoryModel {
	var tags string
	if len(s.Tags) > 0 {
		b, _ := json.Marshal(s.Tags)
		tags = string(b)
	}
	return &StoryModel{
		ID:        s.ID,
		Title:     s.Title,
		Content:   s.Content,
		Era:       s.Era,
		Category:  s.Category,
		Tags:      tags,
		AudioURL:  s.AudioURL,
		ImageURL:  s.ImageURL,
		ViewCount: s.ViewCount,
	}
}

func (r *StoryRepositoryImpl) FindByID(ctx context.Context, id int64) (*culture.Story, error) {
	var model StoryModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToStory(&model), nil
}

func (r *StoryRepositoryImpl) Create(ctx context.Context, s *culture.Story) error {
	model := storyToModel(s)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *StoryRepositoryImpl) Update(ctx context.Context, s *culture.Story) error {
	model := storyToModel(s)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *StoryRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&StoryModel{}, id).Error
}

func (r *StoryRepositoryImpl) FindAll(ctx context.Context) ([]*culture.Story, error) {
	var models []StoryModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*culture.Story, len(models))
	for i := range models {
		result[i] = modelToStory(&models[i])
	}
	return result, nil
}

func (r *StoryRepositoryImpl) Search(ctx context.Context, query *culture.StorySearchQuery) ([]*culture.Story, int64, error) {
	var models []StoryModel
	var total int64

	db := r.db.WithContext(ctx).Model(&StoryModel{})

	if query.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+query.Keyword+"%")
	}
	if query.Era != "" {
		db = db.Where("era = ?", query.Era)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
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

	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*culture.Story, len(models))
	for i := range models {
		result[i] = modelToStory(&models[i])
	}
	return result, total, nil
}

func (r *StoryRepositoryImpl) FindByEra(ctx context.Context, era string) ([]*culture.Story, error) {
	var models []StoryModel
	if err := r.db.WithContext(ctx).Where("era = ?", era).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*culture.Story, len(models))
	for i := range models {
		result[i] = modelToStory(&models[i])
	}
	return result, nil
}

// ===== 家训仓储实现 =====

// FamilyTeachingsRepositoryImpl 家训仓储实现
type FamilyTeachingsRepositoryImpl struct {
	db *gorm.DB
}

// NewFamilyTeachingsRepository 创建家训仓储
func NewFamilyTeachingsRepository(db *gorm.DB) culture.FamilyTeachingsRepository {
	return &FamilyTeachingsRepositoryImpl{db: db}
}

func modelToFamilyTeachings(m *FamilyTeachingsModel) *culture.FamilyTeachings {
	return &culture.FamilyTeachings{
		ID:         m.ID,
		Title:      m.Title,
		Content:    m.Content,
		Generation: m.Generation,
		OriginText: m.OriginText,
		Meaning:    m.Meaning,
		UsageCount: m.UsageCount,
	}
}

func familyTeachingsToModel(ft *culture.FamilyTeachings) *FamilyTeachingsModel {
	return &FamilyTeachingsModel{
		ID:         ft.ID,
		Title:      ft.Title,
		Content:    ft.Content,
		Generation: ft.Generation,
		OriginText: ft.OriginText,
		Meaning:    ft.Meaning,
		UsageCount: ft.UsageCount,
	}
}

func (r *FamilyTeachingsRepositoryImpl) FindByID(ctx context.Context, id int64) (*culture.FamilyTeachings, error) {
	var model FamilyTeachingsModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return modelToFamilyTeachings(&model), nil
}

func (r *FamilyTeachingsRepositoryImpl) Create(ctx context.Context, ft *culture.FamilyTeachings) error {
	model := familyTeachingsToModel(ft)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *FamilyTeachingsRepositoryImpl) Update(ctx context.Context, ft *culture.FamilyTeachings) error {
	model := familyTeachingsToModel(ft)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *FamilyTeachingsRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&FamilyTeachingsModel{}, id).Error
}

func (r *FamilyTeachingsRepositoryImpl) FindAll(ctx context.Context) ([]*culture.FamilyTeachings, error) {
	var models []FamilyTeachingsModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*culture.FamilyTeachings, len(models))
	for i := range models {
		result[i] = modelToFamilyTeachings(&models[i])
	}
	return result, nil
}

func (r *FamilyTeachingsRepositoryImpl) FindByGeneration(ctx context.Context, generation int) ([]*culture.FamilyTeachings, error) {
	var models []FamilyTeachingsModel
	if err := r.db.WithContext(ctx).Where("generation = ?", generation).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*culture.FamilyTeachings, len(models))
	for i := range models {
		result[i] = modelToFamilyTeachings(&models[i])
	}
	return result, nil
}
