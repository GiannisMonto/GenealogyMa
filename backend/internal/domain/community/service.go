package community

import (
	"context"
	"errors"
	"fmt"
)

// Service 社区领域服务
type Service struct {
	profileRepo   UserProfileRepository
	postRepo      PostRepository
	commentRepo   CommentRepository
	messageRepo   MessageRepository
}

// NewService 创建社区服务
func NewService(
	profileRepo UserProfileRepository,
	postRepo PostRepository,
	commentRepo CommentRepository,
	messageRepo MessageRepository,
) *Service {
	return &Service{
		profileRepo: profileRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		messageRepo: messageRepo,
	}
}

// UserProfile operations

func (s *Service) CreateProfile(ctx context.Context, profile *UserProfile) error {
	if err := profile.Validate(); err != nil {
		return err
	}
	return s.profileRepo.Create(ctx, profile)
}

func (s *Service) UpdateProfile(ctx context.Context, profile *UserProfile) error {
	if err := profile.Validate(); err != nil {
		return err
	}
	existing, err := s.profileRepo.FindByID(ctx, profile.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("profile not found")
	}
	return s.profileRepo.Update(ctx, profile)
}

func (s *Service) GetProfile(ctx context.Context, id int64) (*UserProfile, error) {
	return s.profileRepo.FindByID(ctx, id)
}

func (s *Service) GetProfileByUserID(ctx context.Context, userID int64) (*UserProfile, error) {
	return s.profileRepo.FindByUserID(ctx, userID)
}

func (s *Service) DeleteProfile(ctx context.Context, id int64) error {
	existing, err := s.profileRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("profile not found")
	}
	return s.profileRepo.Delete(ctx, id)
}

// Post operations

func (s *Service) CreatePost(ctx context.Context, post *Post) error {
	if err := post.Validate(); err != nil {
		return err
	}
	return s.postRepo.Create(ctx, post)
}

func (s *Service) UpdatePost(ctx context.Context, post *Post) error {
	if err := post.Validate(); err != nil {
		return err
	}
	existing, err := s.postRepo.FindByID(ctx, post.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("post not found")
	}
	return s.postRepo.Update(ctx, post)
}

func (s *Service) DeletePost(ctx context.Context, id int64) error {
	existing, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("post not found")
	}
	existing.IsDeleted = true
	return s.postRepo.Update(ctx, existing)
}

func (s *Service) GetPost(ctx context.Context, id int64) (*Post, error) {
	return s.postRepo.FindByID(ctx, id)
}

func (s *Service) GetUserPosts(ctx context.Context, userID int64, page, pageSize int) ([]*Post, int64, error) {
	return s.postRepo.FindByUserID(ctx, userID, page, pageSize)
}

func (s *Service) GetPostsByTag(ctx context.Context, tag string, page, pageSize int) ([]*Post, int64, error) {
	return s.postRepo.FindByTag(ctx, tag, page, pageSize)
}

func (s *Service) SearchPosts(ctx context.Context, query *PostSearchQuery) ([]*Post, int64, error) {
	return s.postRepo.Search(ctx, query)
}

func (s *Service) IncrementPostLikeCount(ctx context.Context, id int64) error {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if post == nil {
		return fmt.Errorf("post not found")
	}
	post.LikeCount++
	return s.postRepo.Update(ctx, post)
}

func (s *Service) IncrementPostCommentCount(ctx context.Context, id int64) error {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if post == nil {
		return fmt.Errorf("post not found")
	}
	post.CommentCount++
	return s.postRepo.Update(ctx, post)
}

// Comment operations

func (s *Service) CreateComment(ctx context.Context, comment *Comment) error {
	if err := comment.Validate(); err != nil {
		return err
	}
	err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return err
	}
	// 增加帖子评论数
	s.IncrementPostCommentCount(ctx, comment.PostID)
	return nil
}

func (s *Service) UpdateComment(ctx context.Context, comment *Comment) error {
	if err := comment.Validate(); err != nil {
		return err
	}
	existing, err := s.commentRepo.FindByID(ctx, comment.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("comment not found")
	}
	return s.commentRepo.Update(ctx, comment)
}

func (s *Service) DeleteComment(ctx context.Context, id int64) error {
	existing, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("comment not found")
	}
	return s.commentRepo.Delete(ctx, id)
}

func (s *Service) GetComment(ctx context.Context, id int64) (*Comment, error) {
	return s.commentRepo.FindByID(ctx, id)
}

func (s *Service) GetPostComments(ctx context.Context, postID int64) ([]*Comment, error) {
	return s.commentRepo.FindByPostID(ctx, postID)
}

func (s *Service) GetUserComments(ctx context.Context, userID int64, page, pageSize int) ([]*Comment, int64, error) {
	return s.commentRepo.FindByUserID(ctx, userID, page, pageSize)
}

// Message operations

func (s *Service) SendMessage(ctx context.Context, msg *Message) error {
	if err := msg.Validate(); err != nil {
		return err
	}
	return s.messageRepo.Create(ctx, msg)
}

func (s *Service) DeleteMessage(ctx context.Context, id int64) error {
	existing, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("message not found")
	}
	return s.messageRepo.Delete(ctx, id)
}

func (s *Service) GetConversation(ctx context.Context, userID1, userID2 int64, page, pageSize int) ([]*Message, int64, error) {
	return s.messageRepo.FindConversation(ctx, userID1, userID2, page, pageSize)
}

func (s *Service) GetUserMessages(ctx context.Context, userID int64, page, pageSize int) ([]*Message, int64, error) {
	return s.messageRepo.FindByUserID(ctx, userID, page, pageSize)
}

func (s *Service) MarkMessageAsRead(ctx context.Context, messageID int64) error {
	return s.messageRepo.MarkAsRead(ctx, messageID)
}