package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/community"
)

// mockUserProfileRepository 模拟用户资料仓储
type mockUserProfileRepository struct {
	profiles map[int64]*community.UserProfile
	nextID   int64
}

func newMockUserProfileRepository() *mockUserProfileRepository {
	return &mockUserProfileRepository{
		profiles: make(map[int64]*community.UserProfile),
		nextID:   1,
	}
}

func (m *mockUserProfileRepository) FindByID(ctx context.Context, id int64) (*community.UserProfile, error) {
	p, ok := m.profiles[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockUserProfileRepository) FindByUserID(ctx context.Context, userID int64) (*community.UserProfile, error) {
	for _, p := range m.profiles {
		if p.UserID == userID {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockUserProfileRepository) Create(ctx context.Context, profile *community.UserProfile) error {
	profile.ID = m.nextID
	m.nextID++
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()
	m.profiles[profile.ID] = profile
	return nil
}

func (m *mockUserProfileRepository) Update(ctx context.Context, profile *community.UserProfile) error {
	profile.UpdatedAt = time.Now()
	m.profiles[profile.ID] = profile
	return nil
}

func (m *mockUserProfileRepository) Delete(ctx context.Context, id int64) error {
	delete(m.profiles, id)
	return nil
}

// mockPostRepository 模拟动态仓储
type mockPostRepository struct {
	posts   map[int64]*community.Post
	nextID  int64
}

func newMockPostRepository() *mockPostRepository {
	return &mockPostRepository{
		posts:  make(map[int64]*community.Post),
		nextID: 1,
	}
}

func (m *mockPostRepository) FindByID(ctx context.Context, id int64) (*community.Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockPostRepository) Create(ctx context.Context, post *community.Post) error {
	post.ID = m.nextID
	m.nextID++
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepository) Update(ctx context.Context, post *community.Post) error {
	post.UpdatedAt = time.Now()
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepository) Delete(ctx context.Context, id int64) error {
	delete(m.posts, id)
	return nil
}

func (m *mockPostRepository) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*community.Post, int64, error) {
	var result []*community.Post
	for _, p := range m.posts {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockPostRepository) FindByTag(ctx context.Context, tag string, page, pageSize int) ([]*community.Post, int64, error) {
	var result []*community.Post
	for _, p := range m.posts {
		if p.Tag == tag {
			result = append(result, p)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockPostRepository) Search(ctx context.Context, query *community.PostSearchQuery) ([]*community.Post, int64, error) {
	var result []*community.Post
	for _, p := range m.posts {
		result = append(result, p)
	}
	return result, int64(len(result)), nil
}

// mockCommentRepository 模拟评论仓储
type mockCommentRepository struct {
	comments map[int64]*community.Comment
	nextID   int64
}

func newMockCommentRepository() *mockCommentRepository {
	return &mockCommentRepository{
		comments: make(map[int64]*community.Comment),
		nextID:   1,
	}
}

func (m *mockCommentRepository) FindByID(ctx context.Context, id int64) (*community.Comment, error) {
	c, ok := m.comments[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *mockCommentRepository) Create(ctx context.Context, comment *community.Comment) error {
	comment.ID = m.nextID
	m.nextID++
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepository) Update(ctx context.Context, comment *community.Comment) error {
	comment.UpdatedAt = time.Now()
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepository) Delete(ctx context.Context, id int64) error {
	delete(m.comments, id)
	return nil
}

func (m *mockCommentRepository) FindByPostID(ctx context.Context, postID int64) ([]*community.Comment, error) {
	var result []*community.Comment
	for _, c := range m.comments {
		if c.PostID == postID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *mockCommentRepository) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*community.Comment, int64, error) {
	var result []*community.Comment
	for _, c := range m.comments {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, int64(len(result)), nil
}

// mockMessageRepository 模拟私信仓储
type mockMessageRepository struct {
	messages map[int64]*community.Message
	nextID   int64
}

func newMockMessageRepository() *mockMessageRepository {
	return &mockMessageRepository{
		messages: make(map[int64]*community.Message),
		nextID:   1,
	}
}

func (m *mockMessageRepository) FindByID(ctx context.Context, id int64) (*community.Message, error) {
	msg, ok := m.messages[id]
	if !ok {
		return nil, nil
	}
	return msg, nil
}

func (m *mockMessageRepository) Create(ctx context.Context, msg *community.Message) error {
	msg.ID = m.nextID
	m.nextID++
	msg.CreatedAt = time.Now()
	m.messages[msg.ID] = msg
	return nil
}

func (m *mockMessageRepository) Delete(ctx context.Context, id int64) error {
	delete(m.messages, id)
	return nil
}

func (m *mockMessageRepository) FindConversation(ctx context.Context, userID1, userID2 int64, page, pageSize int) ([]*community.Message, int64, error) {
	var result []*community.Message
	for _, msg := range m.messages {
		if (msg.SenderID == userID1 && msg.ReceiverID == userID2) ||
			(msg.SenderID == userID2 && msg.ReceiverID == userID1) {
			result = append(result, msg)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockMessageRepository) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*community.Message, int64, error) {
	var result []*community.Message
	for _, msg := range m.messages {
		if msg.SenderID == userID || msg.ReceiverID == userID {
			result = append(result, msg)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockMessageRepository) MarkAsRead(ctx context.Context, messageID int64) error {
	msg, ok := m.messages[messageID]
	if !ok {
		return nil
	}
	msg.IsRead = true
	now := time.Now()
	msg.ReadAt = &now
	return nil
}

func newTestCommunityService() (*CommunityService, *mockUserProfileRepository, *mockPostRepository, *mockCommentRepository, *mockMessageRepository) {
	profileRepo := newMockUserProfileRepository()
	postRepo := newMockPostRepository()
	commentRepo := newMockCommentRepository()
	messageRepo := newMockMessageRepository()
	svc := NewCommunityService(profileRepo, postRepo, commentRepo, messageRepo)
	return svc, profileRepo, postRepo, commentRepo, messageRepo
}

// ===== 测试用例 =====

func TestCreateProfile(t *testing.T) {
	svc, repo, _, _, _ := newTestCommunityService()
	ctx := context.Background()

	t.Run("create profile successfully", func(t *testing.T) {
		req := &CreateProfileRequest{
			UserID:   1,
			Nickname: "张三",
			Bio:      "你好，我是张三",
			Gender:   "男",
			Province: "四川省",
			City:     "成都市",
		}

		result, err := svc.CreateProfile(ctx, req)
		if err != nil {
			t.Fatalf("CreateProfile() error = %v", err)
		}
		if result.Nickname != "张三" {
			t.Errorf("expected nickname '张三', got '%s'", result.Nickname)
		}
		if len(repo.profiles) != 1 {
			t.Errorf("expected 1 profile in repo, got %d", len(repo.profiles))
		}
	})

	t.Run("create profile with empty nickname should fail", func(t *testing.T) {
		req := &CreateProfileRequest{
			UserID:   2,
			Nickname: "",
		}

		_, err := svc.CreateProfile(ctx, req)
		if err == nil {
			t.Error("expected error for empty nickname")
		}
	})
}

func TestGetProfile(t *testing.T) {
	svc, repo, _, _, _ := newTestCommunityService()
	ctx := context.Background()

	// 预先创建一个资料
	repo.Create(ctx, &community.UserProfile{
		UserID:   1,
		Nickname: "李四",
		Province: "四川省",
		City:     "成都市",
	})

	t.Run("get existing profile", func(t *testing.T) {
		result, err := svc.GetProfile(ctx, 1)
		if err != nil {
			t.Fatalf("GetProfile() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected profile, got nil")
		}
		if result.Nickname != "李四" {
			t.Errorf("expected nickname '李四', got '%s'", result.Nickname)
		}
	})

	t.Run("get non-existent profile", func(t *testing.T) {
		result, err := svc.GetProfile(ctx, 999)
		if err != nil {
			t.Fatalf("GetProfile() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestCreatePost(t *testing.T) {
	svc, _, postRepo, _, _ := newTestCommunityService()
	ctx := context.Background()

	t.Run("create post successfully", func(t *testing.T) {
		req := &CreatePostRequest{
			UserID:  1,
			Content: "这是一条测试动态",
			Tag:     "测试",
		}

		result, err := svc.CreatePost(ctx, req)
		if err != nil {
			t.Fatalf("CreatePost() error = %v", err)
		}
		if result.Content != "这是一条测试动态" {
			t.Errorf("expected content '这是一条测试动态', got '%s'", result.Content)
		}
		if result.Tag != "测试" {
			t.Errorf("expected tag '测试', got '%s'", result.Tag)
		}
		if len(postRepo.posts) != 1 {
			t.Errorf("expected 1 post in repo, got %d", len(postRepo.posts))
		}
	})

	t.Run("create post with empty content should fail", func(t *testing.T) {
		req := &CreatePostRequest{
			UserID:  1,
			Content: "",
		}

		_, err := svc.CreatePost(ctx, req)
		if err == nil {
			t.Error("expected error for empty content")
		}
	})
}

func TestGetPost(t *testing.T) {
	svc, _, postRepo, _, _ := newTestCommunityService()
	ctx := context.Background()

	// 预先创建一个动态
	postRepo.Create(ctx, &community.Post{
		UserID:  1,
		Content: "测试动态",
		Tag:     "测试",
	})

	t.Run("get existing post", func(t *testing.T) {
		result, err := svc.GetPost(ctx, 1)
		if err != nil {
			t.Fatalf("GetPost() error = %v", err)
		}
		if result == nil {
			t.Fatal("expected post, got nil")
		}
		if result.Content != "测试动态" {
			t.Errorf("expected content '测试动态', got '%s'", result.Content)
		}
	})

	t.Run("get non-existent post", func(t *testing.T) {
		result, err := svc.GetPost(ctx, 999)
		if err != nil {
			t.Fatalf("GetPost() error = %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestCreateComment(t *testing.T) {
	svc, _, postRepo, commentRepo, _ := newTestCommunityService()
	ctx := context.Background()

	// 先创建一个动态
	postRepo.Create(ctx, &community.Post{
		UserID:  1,
		Content: "测试动态",
	})

	t.Run("create comment successfully", func(t *testing.T) {
		req := &CreateCommentRequest{
			PostID:  1,
			UserID:  2,
			Content: "这是一条测试评论",
		}

		result, err := svc.CreateComment(ctx, req)
		if err != nil {
			t.Fatalf("CreateComment() error = %v", err)
		}
		if result.Content != "这是一条测试评论" {
			t.Errorf("expected content '这是一条测试评论', got '%s'", result.Content)
		}
		if len(commentRepo.comments) != 1 {
			t.Errorf("expected 1 comment in repo, got %d", len(commentRepo.comments))
		}
		// 验证动态的评论数是否增加
		post := postRepo.posts[1]
		if post.CommentCount != 1 {
			t.Errorf("expected post CommentCount 1, got %d", post.CommentCount)
		}
	})

	t.Run("create comment with empty content should fail", func(t *testing.T) {
		req := &CreateCommentRequest{
			PostID:  1,
			UserID:  2,
			Content: "",
		}

		_, err := svc.CreateComment(ctx, req)
		if err == nil {
			t.Error("expected error for empty content")
		}
	})
}

func TestSendMessage(t *testing.T) {
	svc, _, _, _, msgRepo := newTestCommunityService()
	ctx := context.Background()

	t.Run("send message successfully", func(t *testing.T) {
		req := &SendMessageRequest{
			SenderID:   1,
			ReceiverID: 2,
			Content:    "你好，这是测试消息",
		}

		result, err := svc.SendMessage(ctx, req)
		if err != nil {
			t.Fatalf("SendMessage() error = %v", err)
		}
		if result.Content != "你好，这是测试消息" {
			t.Errorf("expected content '你好，这是测试消息', got '%s'", result.Content)
		}
		if len(msgRepo.messages) != 1 {
			t.Errorf("expected 1 message in repo, got %d", len(msgRepo.messages))
		}
	})

	t.Run("send message to yourself should fail", func(t *testing.T) {
		req := &SendMessageRequest{
			SenderID:   1,
			ReceiverID: 1,
			Content:    "自言自语",
		}

		_, err := svc.SendMessage(ctx, req)
		if err == nil {
			t.Error("expected error when sending message to yourself")
		}
	})
}

func TestLikePost(t *testing.T) {
	svc, _, postRepo, _, _ := newTestCommunityService()
	ctx := context.Background()

	// 预先创建一个动态
	postRepo.Create(ctx, &community.Post{
		UserID:    1,
		Content:   "测试动态",
		LikeCount: 0,
	})

	t.Run("like post successfully", func(t *testing.T) {
		err := svc.LikePost(ctx, 1)
		if err != nil {
			t.Fatalf("LikePost() error = %v", err)
		}
		post := postRepo.posts[1]
		if post.LikeCount != 1 {
			t.Errorf("expected LikeCount 1, got %d", post.LikeCount)
		}
	})

	t.Run("like non-existent post should fail", func(t *testing.T) {
		err := svc.LikePost(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent post")
		}
	})
}

func TestDeletePost(t *testing.T) {
	svc, _, postRepo, _, _ := newTestCommunityService()
	ctx := context.Background()

	// 预先创建一个动态
	postRepo.Create(ctx, &community.Post{
		UserID:  1,
		Content: "待删除动态",
	})

	t.Run("delete existing post", func(t *testing.T) {
		err := svc.DeletePost(ctx, 1)
		if err != nil {
			t.Fatalf("DeletePost() error = %v", err)
		}
		post := postRepo.posts[1]
		if !post.IsDeleted {
			t.Error("expected post to be marked as deleted")
		}
	})

	t.Run("delete non-existent post should fail", func(t *testing.T) {
		err := svc.DeletePost(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent post")
		}
	})
}

// 测试错误场景
func TestCommunityServiceErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("repository error on create post", func(t *testing.T) {
		badProfileRepo := newMockUserProfileRepository()
		badPostRepo := &mockPostRepositoryError{inner: newMockPostRepository()}
		badCommentRepo := newMockCommentRepository()
		badMessageRepo := newMockMessageRepository()
		badSvc := NewCommunityService(badProfileRepo, badPostRepo, badCommentRepo, badMessageRepo)

		req := &CreatePostRequest{
			UserID:  1,
			Content: "测试",
		}

		_, err := badSvc.CreatePost(ctx, req)
		if err == nil {
			t.Error("expected error from repository")
		}
	})
}

type mockPostRepositoryError struct {
	inner *mockPostRepository
}

func (m *mockPostRepositoryError) FindByID(ctx context.Context, id int64) (*community.Post, error) {
	return m.inner.FindByID(ctx, id)
}
func (m *mockPostRepositoryError) Create(ctx context.Context, post *community.Post) error {
	return errors.New("repository error")
}
func (m *mockPostRepositoryError) Update(ctx context.Context, post *community.Post) error {
	return m.inner.Update(ctx, post)
}
func (m *mockPostRepositoryError) Delete(ctx context.Context, id int64) error {
	return m.inner.Delete(ctx, id)
}
func (m *mockPostRepositoryError) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*community.Post, int64, error) {
	return m.inner.FindByUserID(ctx, userID, page, pageSize)
}
func (m *mockPostRepositoryError) FindByTag(ctx context.Context, tag string, page, pageSize int) ([]*community.Post, int64, error) {
	return m.inner.FindByTag(ctx, tag, page, pageSize)
}
func (m *mockPostRepositoryError) Search(ctx context.Context, query *community.PostSearchQuery) ([]*community.Post, int64, error) {
	return m.inner.Search(ctx, query)
}