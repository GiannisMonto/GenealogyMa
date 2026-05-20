package culture

import (
	"context"
	"errors"
	"fmt"
)

// Service 文化领域服务
type Service struct {
	docRepo       Repository
	storyRepo     StoryRepository
	teachingsRepo FamilyTeachingsRepository
}

// NewService 创建文化服务
func NewService(docRepo Repository, storyRepo StoryRepository, teachingsRepo FamilyTeachingsRepository) *Service {
	return &Service{
		docRepo:       docRepo,
		storyRepo:     storyRepo,
		teachingsRepo: teachingsRepo,
	}
}

// Document operations

func (s *Service) CreateDocument(ctx context.Context, doc *Document) error {
	if err := doc.Validate(); err != nil {
		return err
	}
	return s.docRepo.Create(ctx, doc)
}

func (s *Service) UpdateDocument(ctx context.Context, doc *Document) error {
	if err := doc.Validate(); err != nil {
		return err
	}
	existing, err := s.docRepo.FindByID(ctx, doc.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("document not found")
	}
	return s.docRepo.Update(ctx, doc)
}

func (s *Service) DeleteDocument(ctx context.Context, id int64) error {
	existing, err := s.docRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("document not found")
	}
	return s.docRepo.Delete(ctx, id)
}

func (s *Service) GetDocument(ctx context.Context, id int64) (*Document, error) {
	return s.docRepo.FindByID(ctx, id)
}

func (s *Service) ListDocuments(ctx context.Context) ([]*Document, error) {
	return s.docRepo.FindAll(ctx)
}

func (s *Service) SearchDocuments(ctx context.Context, query *SearchQuery) ([]*Document, int64, error) {
	return s.docRepo.Search(ctx, query)
}

func (s *Service) ListDocumentsByCategory(ctx context.Context, category DocumentCategory) ([]*Document, error) {
	return s.docRepo.FindByCategory(ctx, category)
}

// Story operations

func (s *Service) CreateStory(ctx context.Context, story *Story) error {
	if err := story.Validate(); err != nil {
		return err
	}
	return s.storyRepo.Create(ctx, story)
}

func (s *Service) UpdateStory(ctx context.Context, story *Story) error {
	if err := story.Validate(); err != nil {
		return err
	}
	existing, err := s.storyRepo.FindByID(ctx, story.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("story not found")
	}
	return s.storyRepo.Update(ctx, story)
}

func (s *Service) DeleteStory(ctx context.Context, id int64) error {
	existing, err := s.storyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("story not found")
	}
	return s.storyRepo.Delete(ctx, id)
}

func (s *Service) GetStory(ctx context.Context, id int64) (*Story, error) {
	return s.storyRepo.FindByID(ctx, id)
}

func (s *Service) ListStories(ctx context.Context) ([]*Story, error) {
	return s.storyRepo.FindAll(ctx)
}

func (s *Service) SearchStories(ctx context.Context, query *StorySearchQuery) ([]*Story, int64, error) {
	return s.storyRepo.Search(ctx, query)
}

func (s *Service) ListStoriesByEra(ctx context.Context, era string) ([]*Story, error) {
	return s.storyRepo.FindByEra(ctx, era)
}

// FamilyTeachings operations

func (s *Service) CreateFamilyTeachings(ctx context.Context, ft *FamilyTeachings) error {
	if err := ft.Validate(); err != nil {
		return err
	}
	return s.teachingsRepo.Create(ctx, ft)
}

func (s *Service) UpdateFamilyTeachings(ctx context.Context, ft *FamilyTeachings) error {
	if err := ft.Validate(); err != nil {
		return err
	}
	existing, err := s.teachingsRepo.FindByID(ctx, ft.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("family teachings not found")
	}
	return s.teachingsRepo.Update(ctx, ft)
}

func (s *Service) DeleteFamilyTeachings(ctx context.Context, id int64) error {
	existing, err := s.teachingsRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("family teachings not found")
	}
	return s.teachingsRepo.Delete(ctx, id)
}

func (s *Service) GetFamilyTeachings(ctx context.Context, id int64) (*FamilyTeachings, error) {
	return s.teachingsRepo.FindByID(ctx, id)
}

func (s *Service) ListFamilyTeachings(ctx context.Context) ([]*FamilyTeachings, error) {
	return s.teachingsRepo.FindAll(ctx)
}

func (s *Service) ListFamilyTeachingsByGeneration(ctx context.Context, generation int) ([]*FamilyTeachings, error) {
	return s.teachingsRepo.FindByGeneration(ctx, generation)
}

// IncrementViewCount 增加浏览次数
func (s *Service) IncrementDocumentViewCount(ctx context.Context, id int64) error {
	doc, err := s.docRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("document not found")
	}
	doc.ViewCount++
	return s.docRepo.Update(ctx, doc)
}

func (s *Service) IncrementStoryViewCount(ctx context.Context, id int64) error {
	story, err := s.storyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if story == nil {
		return fmt.Errorf("story not found")
	}
	story.ViewCount++
	return s.storyRepo.Update(ctx, story)
}