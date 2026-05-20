package service

import (
	"context"

	"github.com/genealogy-ma/platform/internal/domain/culture"
)

// CultureService 文化应用服务
type CultureService struct {
	domainService *culture.Service
}

// NewCultureService 创建文化应用服务
func NewCultureService(docRepo culture.Repository, storyRepo culture.StoryRepository, teachingsRepo culture.FamilyTeachingsRepository) *CultureService {
	return &CultureService{
		domainService: culture.NewService(docRepo, storyRepo, teachingsRepo),
	}
}

// ===== DTO 定义 =====

// DocumentDTO 文献DTO
type DocumentDTO struct {
	ID           int64                 `json:"id"`
	Title        string                `json:"title"`
	Content      string                `json:"content"`
	Category     culture.DocumentCategory `json:"category"`
	Author       string                `json:"author"`
	CreatedYear  *int                  `json:"created_year"`
	Dynasty      string                `json:"dynasty"`
	Source       string                `json:"source"`
	ImageURLs    []string              `json:"image_urls"`
	ViewCount    int                   `json:"view_count"`
	CollectCount int                   `json:"collect_count"`
	CreatedAt    string                `json:"created_at"`
	UpdatedAt    string                `json:"updated_at"`
}

// StoryDTO 故事DTO
type StoryDTO struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Era       string   `json:"era"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	AudioURL  string   `json:"audio_url"`
	ImageURL  string   `json:"image_url"`
	ViewCount int      `json:"view_count"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// FamilyTeachingsDTO 家训DTO
type FamilyTeachingsDTO struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Generation int    `json:"generation"`
	OriginText string `json:"origin_text"`
	Meaning    string `json:"meaning"`
	UsageCount int    `json:"usage_count"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// CreateDocumentRequest 创建文献请求
type CreateDocumentRequest struct {
	Title       string                 `json:"title" binding:"required"`
	Content     string                 `json:"content" binding:"required"`
	Category    culture.DocumentCategory `json:"category" binding:"required"`
	Author      string                 `json:"author"`
	CreatedYear *int                   `json:"created_year"`
	Dynasty     string                 `json:"dynasty"`
	Source      string                 `json:"source"`
	ImageURLs   []string               `json:"image_urls"`
}

// UpdateDocumentRequest 更新文献请求
type UpdateDocumentRequest struct {
	Title       *string                `json:"title"`
	Content     *string                `json:"content"`
	Category    *culture.DocumentCategory `json:"category"`
	Author      *string                `json:"author"`
	CreatedYear *int                   `json:"created_year"`
	Dynasty     *string                `json:"dynasty"`
	Source      *string                `json:"source"`
	ImageURLs   []string               `json:"image_urls"`
}

// CreateStoryRequest 创建故事请求
type CreateStoryRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Era      string   `json:"era"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	AudioURL string   `json:"audio_url"`
	ImageURL string   `json:"image_url"`
}

// UpdateStoryRequest 更新故事请求
type UpdateStoryRequest struct {
	Title    *string  `json:"title"`
	Content  *string  `json:"content"`
	Era      *string  `json:"era"`
	Category *string  `json:"category"`
	Tags     []string `json:"tags"`
	AudioURL *string  `json:"audio_url"`
	ImageURL *string  `json:"image_url"`
}

// CreateFamilyTeachingsRequest 创建家训请求
type CreateFamilyTeachingsRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Generation int    `json:"generation"`
	OriginText string `json:"origin_text"`
	Meaning    string `json:"meaning"`
}

// UpdateFamilyTeachingsRequest 更新家训请求
type UpdateFamilyTeachingsRequest struct {
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	Generation *int    `json:"generation"`
	OriginText *string `json:"origin_text"`
	Meaning    *string `json:"meaning"`
}

// ===== 文献服务方法 =====

// ListDocuments 获取所有文献
func (s *CultureService) ListDocuments(ctx context.Context) ([]*DocumentDTO, error) {
	docs, err := s.domainService.ListDocuments(ctx)
	if err != nil {
		return nil, err
	}
	dtos := make([]*DocumentDTO, len(docs))
	for i, d := range docs {
		dtos[i] = documentToDTO(d)
	}
	return dtos, nil
}

// GetDocument 获取文献详情
func (s *CultureService) GetDocument(ctx context.Context, id int64) (*DocumentDTO, error) {
	d, err := s.domainService.GetDocument(ctx, id)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, nil
	}
	return documentToDTO(d), nil
}

// CreateDocument 创建文献
func (s *CultureService) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*DocumentDTO, error) {
	d := &culture.Document{
		Title:       req.Title,
		Content:     req.Content,
		Category:    req.Category,
		Author:      req.Author,
		CreatedYear: req.CreatedYear,
		Dynasty:     req.Dynasty,
		Source:      req.Source,
		ImageURLs:   req.ImageURLs,
	}
	if err := s.domainService.CreateDocument(ctx, d); err != nil {
		return nil, err
	}
	return documentToDTO(d), nil
}

// UpdateDocument 更新文献
func (s *CultureService) UpdateDocument(ctx context.Context, id int64, req *UpdateDocumentRequest) (*DocumentDTO, error) {
	existing, err := s.domainService.GetDocument(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Content != nil {
		existing.Content = *req.Content
	}
	if req.Category != nil {
		existing.Category = *req.Category
	}
	if req.Author != nil {
		existing.Author = *req.Author
	}
	if req.CreatedYear != nil {
		existing.CreatedYear = req.CreatedYear
	}
	if req.Dynasty != nil {
		existing.Dynasty = *req.Dynasty
	}
	if req.Source != nil {
		existing.Source = *req.Source
	}
	if req.ImageURLs != nil {
		existing.ImageURLs = req.ImageURLs
	}
	if err := s.domainService.UpdateDocument(ctx, existing); err != nil {
		return nil, err
	}
	return documentToDTO(existing), nil
}

// DeleteDocument 删除文献
func (s *CultureService) DeleteDocument(ctx context.Context, id int64) error {
	return s.domainService.DeleteDocument(ctx, id)
}

// ===== 故事服务方法 =====

// ListStories 获取所有故事
func (s *CultureService) ListStories(ctx context.Context) ([]*StoryDTO, error) {
	stories, err := s.domainService.ListStories(ctx)
	if err != nil {
		return nil, err
	}
	dtos := make([]*StoryDTO, len(stories))
	for i, st := range stories {
		dtos[i] = storyToDTO(st)
	}
	return dtos, nil
}

// GetStory 获取故事详情
func (s *CultureService) GetStory(ctx context.Context, id int64) (*StoryDTO, error) {
	st, err := s.domainService.GetStory(ctx, id)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, nil
	}
	return storyToDTO(st), nil
}

// CreateStory 创建故事
func (s *CultureService) CreateStory(ctx context.Context, req *CreateStoryRequest) (*StoryDTO, error) {
	st := &culture.Story{
		Title:    req.Title,
		Content:  req.Content,
		Era:      req.Era,
		Category: req.Category,
		Tags:     req.Tags,
		AudioURL: req.AudioURL,
		ImageURL: req.ImageURL,
	}
	if err := s.domainService.CreateStory(ctx, st); err != nil {
		return nil, err
	}
	return storyToDTO(st), nil
}

// UpdateStory 更新故事
func (s *CultureService) UpdateStory(ctx context.Context, id int64, req *UpdateStoryRequest) (*StoryDTO, error) {
	existing, err := s.domainService.GetStory(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Content != nil {
		existing.Content = *req.Content
	}
	if req.Era != nil {
		existing.Era = *req.Era
	}
	if req.Category != nil {
		existing.Category = *req.Category
	}
	if req.Tags != nil {
		existing.Tags = req.Tags
	}
	if req.AudioURL != nil {
		existing.AudioURL = *req.AudioURL
	}
	if req.ImageURL != nil {
		existing.ImageURL = *req.ImageURL
	}
	if err := s.domainService.UpdateStory(ctx, existing); err != nil {
		return nil, err
	}
	return storyToDTO(existing), nil
}

// DeleteStory 删除故事
func (s *CultureService) DeleteStory(ctx context.Context, id int64) error {
	return s.domainService.DeleteStory(ctx, id)
}

// ===== 家训服务方法 =====

// ListFamilyTeachings 获取所有家训
func (s *CultureService) ListFamilyTeachings(ctx context.Context) ([]*FamilyTeachingsDTO, error) {
	items, err := s.domainService.ListFamilyTeachings(ctx)
	if err != nil {
		return nil, err
	}
	dtos := make([]*FamilyTeachingsDTO, len(items))
	for i, ft := range items {
		dtos[i] = familyTeachingsToDTO(ft)
	}
	return dtos, nil
}

// GetFamilyTeachings 获取家训详情
func (s *CultureService) GetFamilyTeachings(ctx context.Context, id int64) (*FamilyTeachingsDTO, error) {
	ft, err := s.domainService.GetFamilyTeachings(ctx, id)
	if err != nil {
		return nil, err
	}
	if ft == nil {
		return nil, nil
	}
	return familyTeachingsToDTO(ft), nil
}

// CreateFamilyTeachings 创建家训
func (s *CultureService) CreateFamilyTeachings(ctx context.Context, req *CreateFamilyTeachingsRequest) (*FamilyTeachingsDTO, error) {
	ft := &culture.FamilyTeachings{
		Title:      req.Title,
		Content:    req.Content,
		Generation: req.Generation,
		OriginText: req.OriginText,
		Meaning:    req.Meaning,
	}
	if err := s.domainService.CreateFamilyTeachings(ctx, ft); err != nil {
		return nil, err
	}
	return familyTeachingsToDTO(ft), nil
}

// UpdateFamilyTeachings 更新家训
func (s *CultureService) UpdateFamilyTeachings(ctx context.Context, id int64, req *UpdateFamilyTeachingsRequest) (*FamilyTeachingsDTO, error) {
	existing, err := s.domainService.GetFamilyTeachings(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Content != nil {
		existing.Content = *req.Content
	}
	if req.Generation != nil {
		existing.Generation = *req.Generation
	}
	if req.OriginText != nil {
		existing.OriginText = *req.OriginText
	}
	if req.Meaning != nil {
		existing.Meaning = *req.Meaning
	}
	if err := s.domainService.UpdateFamilyTeachings(ctx, existing); err != nil {
		return nil, err
	}
	return familyTeachingsToDTO(existing), nil
}

// DeleteFamilyTeachings 删除家训
func (s *CultureService) DeleteFamilyTeachings(ctx context.Context, id int64) error {
	return s.domainService.DeleteFamilyTeachings(ctx, id)
}

// ===== 辅助函数 =====

func documentToDTO(d *culture.Document) *DocumentDTO {
	return &DocumentDTO{
		ID:           d.ID,
		Title:        d.Title,
		Content:      d.Content,
		Category:     d.Category,
		Author:       d.Author,
		CreatedYear:  d.CreatedYear,
		Dynasty:      d.Dynasty,
		Source:       d.Source,
		ImageURLs:    d.ImageURLs,
		ViewCount:    d.ViewCount,
		CollectCount: d.CollectCount,
		CreatedAt:    d.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    d.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func storyToDTO(st *culture.Story) *StoryDTO {
	return &StoryDTO{
		ID:        st.ID,
		Title:     st.Title,
		Content:   st.Content,
		Era:       st.Era,
		Category:  st.Category,
		Tags:      st.Tags,
		AudioURL:  st.AudioURL,
		ImageURL:  st.ImageURL,
		ViewCount: st.ViewCount,
		CreatedAt: st.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: st.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func familyTeachingsToDTO(ft *culture.FamilyTeachings) *FamilyTeachingsDTO {
	return &FamilyTeachingsDTO{
		ID:         ft.ID,
		Title:      ft.Title,
		Content:    ft.Content,
		Generation: ft.Generation,
		OriginText: ft.OriginText,
		Meaning:    ft.Meaning,
		UsageCount: ft.UsageCount,
		CreatedAt:  ft.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  ft.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
