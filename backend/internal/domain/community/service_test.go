package community

import (
	"context"
	"testing"
)

// MockUserProfileRepository 模拟用户资料仓储
type MockUserProfileRepository struct {
	profiles map[int64]*UserProfile
	nextID   int64
}

func NewMockUserProfileRepository() *MockUserProfileRepository {
	return &MockUserProfileRepository{
		profiles: make(map[int64]*UserProfile),
		nextID:   1,
	}
}

func (m *MockUserProfileRepository) FindByID(ctx context.Context, id int64) (*UserProfile, error) {
	if p, ok := m.profiles[id]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *MockUserProfileRepository) FindByUserID(ctx context.Context, userID int64) (*UserProfile, error) {
	for _, p := range m.profiles {
		if p.UserID == userID {
			return p, nil
		}
	}
	return nil, nil
}

func (m *MockUserProfileRepository) Create(ctx context.Context, profile *UserProfile) error {
	profile.ID = m.nextID
	m.nextID++
	m.profiles[profile.ID] = profile
	return nil
}

func (m *MockUserProfileRepository) Update(ctx context.Context, profile *UserProfile) error {
	m.profiles[profile.ID] = profile
	return nil
}

func (m *MockUserProfileRepository) Delete(ctx context.Context, id int64) error {
	delete(m.profiles, id)
	return nil
}

// MockPostRepository 模拟动态仓储
type MockPostRepository struct {
	posts   map[int64]*Post
	nextID  int64
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		posts:  make(map[int64]*Post),
		nextID: 1,
	}
}

func (m *MockPostRepository) FindByID(ctx context.Context, id int64) (*Post, error) {
	if p, ok := m.posts[id]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *MockPostRepository) Create(ctx context.Context, post *Post) error {
	post.ID = m.nextID
	m.nextID++
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) Update(ctx context.Context, post *Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) Delete(ctx context.Context, id int64) error {
	delete(m.posts, id)
	return nil
}

func (m *MockPostRepository) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*Post, int64, error) {
	result := make([]*Post, 0)
	for _, p := range m.posts {
		if p.UserID == userID && !p.IsDeleted {
			result = append(result, p)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockPostRepository) FindByTag(ctx context.Context, tag string, page, pageSize int) ([]*Post, int64, error) {
	result := make([]*Post, 0)
	for _, p := range m.posts {
		if p.Tag == tag && !p.IsDeleted {
			result = append(result, p)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockPostRepository) Search(ctx context.Context, query *PostSearchQuery) ([]*Post, int64, error) {
	result := make([]*Post, 0)
	for _, p := range m.posts {
		if query.Keyword != "" && p.Content != query.Keyword {
			continue
		}
		if !p.IsDeleted {
			result = append(result, p)
		}
	}
	return result, int64(len(result)), nil
}

// MockCommentRepository 模拟评论仓储
type MockCommentRepository struct {
	comments map[int64]*Comment
	nextID   int64
}

func NewMockCommentRepository() *MockCommentRepository {
	return &MockCommentRepository{
		comments: make(map[int64]*Comment),
		nextID:   1,
	}
}

func (m *MockCommentRepository) FindByID(ctx context.Context, id int64) (*Comment, error) {
	if c, ok := m.comments[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *Comment) error {
	comment.ID = m.nextID
	m.nextID++
	m.comments[comment.ID] = comment
	return nil
}

func (m *MockCommentRepository) Update(ctx context.Context, comment *Comment) error {
	m.comments[comment.ID] = comment
	return nil
}

func (m *MockCommentRepository) Delete(ctx context.Context, id int64) error {
	delete(m.comments, id)
	return nil
}

func (m *MockCommentRepository) FindByPostID(ctx context.Context, postID int64) ([]*Comment, error) {
	result := make([]*Comment, 0)
	for _, c := range m.comments {
		if c.PostID == postID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *MockCommentRepository) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*Comment, int64, error) {
	result := make([]*Comment, 0)
	for _, c := range m.comments {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, int64(len(result)), nil
}

// MockMessageRepository 模拟私信仓储
type MockMessageRepository struct {
	messages map[int64]*Message
	nextID   int64
}

func NewMockMessageRepository() *MockMessageRepository {
	return &MockMessageRepository{
		messages: make(map[int64]*Message),
		nextID:   1,
	}
}

func (m *MockMessageRepository) FindByID(ctx context.Context, id int64) (*Message, error) {
	if msg, ok := m.messages[id]; ok {
		return msg, nil
	}
	return nil, nil
}

func (m *MockMessageRepository) Create(ctx context.Context, msg *Message) error {
	msg.ID = m.nextID
	m.nextID++
	m.messages[msg.ID] = msg
	return nil
}

func (m *MockMessageRepository) Delete(ctx context.Context, id int64) error {
	delete(m.messages, id)
	return nil
}

func (m *MockMessageRepository) FindConversation(ctx context.Context, userID1, userID2 int64, page, pageSize int) ([]*Message, int64, error) {
	result := make([]*Message, 0)
	for _, msg := range m.messages {
		if (msg.SenderID == userID1 && msg.ReceiverID == userID2) ||
			(msg.SenderID == userID2 && msg.ReceiverID == userID1) {
			result = append(result, msg)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockMessageRepository) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*Message, int64, error) {
	result := make([]*Message, 0)
	for _, msg := range m.messages {
		if msg.SenderID == userID || msg.ReceiverID == userID {
			result = append(result, msg)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockMessageRepository) MarkAsRead(ctx context.Context, messageID int64) error {
	if msg, ok := m.messages[messageID]; ok {
		msg.IsRead = true
	}
	return nil
}

// Test UserProfile entity
func TestUserProfile_Validate(t *testing.T) {
	tests := []struct {
		name    string
		profile UserProfile
		wantErr bool
	}{
		{
			name: "valid profile",
			profile: UserProfile{
				UserID:   1,
				Nickname: "测试用户",
			},
			wantErr: false,
		},
		{
			name: "zero user_id",
			profile: UserProfile{
				UserID:   0,
				Nickname: "测试用户",
			},
			wantErr: true,
		},
		{
			name: "empty nickname",
			profile: UserProfile{
				UserID:   1,
				Nickname: "",
			},
			wantErr: true,
		},
		{
			name: "nickname too long",
			profile: UserProfile{
				UserID:   1,
				Nickname: string(make([]byte, 51)),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.profile.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserProfile.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test Post entity
func TestPost_Validate(t *testing.T) {
	tests := []struct {
		name    string
		post    Post
		wantErr bool
	}{
		{
			name: "valid post",
			post: Post{
				UserID:  1,
				Content: "这是动态内容",
			},
			wantErr: false,
		},
		{
			name: "zero user_id",
			post: Post{
				UserID:  0,
				Content: "这是动态内容",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			post: Post{
				UserID:  1,
				Content: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.post.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Post.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test Comment entity
func TestComment_Validate(t *testing.T) {
	tests := []struct {
		name    string
		comment Comment
		wantErr bool
	}{
		{
			name: "valid comment",
			comment: Comment{
				PostID:  1,
				UserID:  1,
				Content: "这是评论内容",
			},
			wantErr: false,
		},
		{
			name: "zero post_id",
			comment: Comment{
				PostID:  0,
				UserID:  1,
				Content: "这是评论内容",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			comment: Comment{
				PostID:  1,
				UserID:  1,
				Content: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Comment.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test Message entity
func TestMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     Message
		wantErr bool
	}{
		{
			name: "valid message",
			msg: Message{
				SenderID:   1,
				ReceiverID: 2,
				Content:    "你好",
			},
			wantErr: false,
		},
		{
			name: "self message",
			msg: Message{
				SenderID:   1,
				ReceiverID: 1,
				Content:    "你好",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			msg: Message{
				SenderID:   1,
				ReceiverID: 2,
				Content:    "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Message.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test Service - Profile
func TestService_CreateProfile(t *testing.T) {
	repo := NewMockUserProfileRepository()
	svc := NewService(repo, nil, nil, nil)

	profile := &UserProfile{
		UserID:   1,
		Nickname: "测试用户",
	}

	err := svc.CreateProfile(context.Background(), profile)
	if err != nil {
		t.Errorf("CreateProfile() error = %v", err)
	}
	if profile.ID == 0 {
		t.Error("CreateProfile() did not set ID")
	}
}

func TestService_GetProfileByUserID(t *testing.T) {
	repo := NewMockUserProfileRepository()
	svc := NewService(repo, nil, nil, nil)

	profile := &UserProfile{
		UserID:   1,
		Nickname: "测试用户",
	}
	repo.Create(context.Background(), profile)

	got, err := svc.GetProfileByUserID(context.Background(), 1)
	if err != nil {
		t.Errorf("GetProfileByUserID() error = %v", err)
	}
	if got == nil {
		t.Error("GetProfileByUserID() returned nil")
	}
	if got.Nickname != profile.Nickname {
		t.Errorf("GetProfileByUserID() nickname = %v, want %v", got.Nickname, profile.Nickname)
	}
}

// Test Service - Post
func TestService_CreatePost(t *testing.T) {
	repo := NewMockPostRepository()
	svc := NewService(nil, repo, nil, nil)

	post := &Post{
		UserID:  1,
		Content: "这是动态内容",
	}

	err := svc.CreatePost(context.Background(), post)
	if err != nil {
		t.Errorf("CreatePost() error = %v", err)
	}
	if post.ID == 0 {
		t.Error("CreatePost() did not set ID")
	}
}

func TestService_DeletePost(t *testing.T) {
	repo := NewMockPostRepository()
	svc := NewService(nil, repo, nil, nil)

	post := &Post{
		UserID:  1,
		Content: "这是动态内容",
	}
	repo.Create(context.Background(), post)

	err := svc.DeletePost(context.Background(), post.ID)
	if err != nil {
		t.Errorf("DeletePost() error = %v", err)
	}

	got, _ := svc.GetPost(context.Background(), post.ID)
	if !got.IsDeleted {
		t.Error("DeletePost() post not marked as deleted")
	}
}

func TestService_IncrementPostLikeCount(t *testing.T) {
	repo := NewMockPostRepository()
	svc := NewService(nil, repo, nil, nil)

	post := &Post{
		UserID:    1,
		Content:   "这是动态内容",
		LikeCount: 0,
	}
	repo.Create(context.Background(), post)

	err := svc.IncrementPostLikeCount(context.Background(), post.ID)
	if err != nil {
		t.Errorf("IncrementPostLikeCount() error = %v", err)
	}

	got, _ := svc.GetPost(context.Background(), post.ID)
	if got.LikeCount != 1 {
		t.Errorf("IncrementPostLikeCount() count = %v, want 1", got.LikeCount)
	}
}

// Test Service - Message
func TestService_SendMessage(t *testing.T) {
	repo := NewMockMessageRepository()
	svc := NewService(nil, nil, nil, repo)

	msg := &Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "你好",
	}

	err := svc.SendMessage(context.Background(), msg)
	if err != nil {
		t.Errorf("SendMessage() error = %v", err)
	}
	if msg.ID == 0 {
		t.Error("SendMessage() did not set ID")
	}
}

func TestService_GetConversation(t *testing.T) {
	repo := NewMockMessageRepository()
	svc := NewService(nil, nil, nil, repo)

	repo.Create(context.Background(), &Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "你好",
	})
	repo.Create(context.Background(), &Message{
		SenderID:   2,
		ReceiverID: 1,
		Content:    "你好啊",
	})

	messages, _, err := svc.GetConversation(context.Background(), 1, 2, 1, 10)
	if err != nil {
		t.Errorf("GetConversation() error = %v", err)
	}
	if len(messages) != 2 {
		t.Errorf("GetConversation() len = %v, want 2", len(messages))
	}
}

func TestService_MarkMessageAsRead(t *testing.T) {
	repo := NewMockMessageRepository()
	svc := NewService(nil, nil, nil, repo)

	msg := &Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "你好",
		IsRead:     false,
	}
	repo.Create(context.Background(), msg)

	err := svc.MarkMessageAsRead(context.Background(), msg.ID)
	if err != nil {
		t.Errorf("MarkMessageAsRead() error = %v", err)
	}

	got, _ := repo.FindByID(context.Background(), msg.ID)
	if !got.IsRead {
		t.Error("MarkMessageAsRead() message not marked as read")
	}
}