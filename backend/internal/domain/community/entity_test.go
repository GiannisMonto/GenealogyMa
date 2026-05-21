package community

import (
	"testing"
	"time"
)

func TestUserProfileEntity_Validate(t *testing.T) {
	tests := []struct {
		name    string
		p       *UserProfile
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid profile",
			p: &UserProfile{
				UserID:   1,
				Nickname: "测试用户",
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			p: &UserProfile{
				UserID:   0,
				Nickname: "测试用户",
			},
			wantErr: true,
			errMsg:  "user_id is required",
		},
		{
			name: "empty nickname",
			p: &UserProfile{
				UserID:   1,
				Nickname: "",
			},
			wantErr: true,
			errMsg:  "nickname cannot be empty",
		},
		{
			name: "nickname too long",
			p: &UserProfile{
				UserID:   1,
				Nickname: string(make([]byte, 51)),
			},
			wantErr: true,
			errMsg:  "nickname too long (max 50)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.p.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUserProfileEntity_Fields(t *testing.T) {
	birthDate := "1990-01-01"
	now := time.Now()
	p := &UserProfile{
		ID:             1,
		UserID:         100,
		Nickname:       "测试用户",
		AvatarURL:      "https://example.com/avatar.jpg",
		Bio:            "这是一个简介",
		Gender:         "male",
		BirthDate:      &birthDate,
		Province:       "浙江",
		City:           "杭州",
		FollowerCount:  100,
		FollowingCount: 50,
		PostCount:      200,
		PrivacyLevel:   0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if p.ID != 1 {
		t.Errorf("UserProfile.ID = %v, want 1", p.ID)
	}
	if p.UserID != 100 {
		t.Errorf("UserProfile.UserID = %v, want 100", p.UserID)
	}
	if p.Nickname != "测试用户" {
		t.Errorf("UserProfile.Nickname = %v, want 测试用户", p.Nickname)
	}
	if p.Province != "浙江" {
		t.Errorf("UserProfile.Province = %v, want 浙江", p.Province)
	}
	if p.FollowerCount != 100 {
		t.Errorf("UserProfile.FollowerCount = %v, want 100", p.FollowerCount)
	}
	if p.PrivacyLevel != 0 {
		t.Errorf("UserProfile.PrivacyLevel = %v, want 0", p.PrivacyLevel)
	}
}

func TestPostEntity_Validate(t *testing.T) {
	tests := []struct {
		name    string
		p       *Post
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid post",
			p: &Post{
				UserID:  1,
				Content: "这是一条动态",
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			p: &Post{
				UserID:  0,
				Content: "这是一条动态",
			},
			wantErr: true,
			errMsg:  "user_id is required",
		},
		{
			name: "empty content",
			p: &Post{
				UserID:  1,
				Content: "",
			},
			wantErr: true,
			errMsg:  "content cannot be empty",
		},
		{
			name: "content too long",
			p: &Post{
				UserID:  1,
				Content: string(make([]byte, 5001)),
			},
			wantErr: true,
			errMsg:  "content too long (max 5000)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.p.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestPostEntity_Fields(t *testing.T) {
	now := time.Now()
	p := &Post{
		ID:           1,
		UserID:       100,
		Content:      "测试内容",
		Images:       []string{"img1.jpg", "img2.jpg"},
		Tag:          "日常",
		LikeCount:    50,
		CommentCount: 10,
		ShareCount:   5,
		IsDeleted:    false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if p.ID != 1 {
		t.Errorf("Post.ID = %v, want 1", p.ID)
	}
	if len(p.Images) != 2 {
		t.Errorf("Post.Images length = %v, want 2", len(p.Images))
	}
	if p.LikeCount != 50 {
		t.Errorf("Post.LikeCount = %v, want 50", p.LikeCount)
	}
	if p.IsDeleted != false {
		t.Errorf("Post.IsDeleted = %v, want false", p.IsDeleted)
	}
}

func TestCommentEntity_Validate(t *testing.T) {
	tests := []struct {
		name    string
		c       *Comment
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid comment",
			c: &Comment{
				PostID:  1,
				UserID:  1,
				Content: "这是一条评论",
			},
			wantErr: false,
		},
		{
			name: "missing post_id",
			c: &Comment{
				PostID:  0,
				UserID:  1,
				Content: "这是一条评论",
			},
			wantErr: true,
			errMsg:  "post_id is required",
		},
		{
			name: "missing user_id",
			c: &Comment{
				PostID:  1,
				UserID:  0,
				Content: "这是一条评论",
			},
			wantErr: true,
			errMsg:  "user_id is required",
		},
		{
			name: "empty content",
			c: &Comment{
				PostID:  1,
				UserID:  1,
				Content: "",
			},
			wantErr: true,
			errMsg:  "content cannot be empty",
		},
		{
			name: "content too long",
			c: &Comment{
				PostID:  1,
				UserID:  1,
				Content: string(make([]byte, 1001)),
			},
			wantErr: true,
			errMsg:  "content too long (max 1000)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestCommentEntity_Fields(t *testing.T) {
	parentID := int64(10)
	now := time.Now()
	c := &Comment{
		ID:        1,
		PostID:    100,
		UserID:    200,
		ParentID:  &parentID,
		Content:   "回复内容",
		LikeCount: 5,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if c.ID != 1 {
		t.Errorf("Comment.ID = %v, want 1", c.ID)
	}
	if c.ParentID == nil || *c.ParentID != 10 {
		t.Errorf("Comment.ParentID = %v, want 10", *c.ParentID)
	}
	if c.LikeCount != 5 {
		t.Errorf("Comment.LikeCount = %v, want 5", c.LikeCount)
	}
}

func TestMessageEntity_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       *Message
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid message",
			m: &Message{
				SenderID:   1,
				ReceiverID: 2,
				Content:    "你好",
			},
			wantErr: false,
		},
		{
			name: "missing sender_id",
			m: &Message{
				SenderID:   0,
				ReceiverID: 2,
				Content:    "你好",
			},
			wantErr: true,
			errMsg:  "sender_id is required",
		},
		{
			name: "missing receiver_id",
			m: &Message{
				SenderID:   1,
				ReceiverID: 0,
				Content:    "你好",
			},
			wantErr: true,
			errMsg:  "receiver_id is required",
		},
		{
			name: "self message",
			m: &Message{
				SenderID:   1,
				ReceiverID: 1,
				Content:    "你好",
			},
			wantErr: true,
			errMsg:  "cannot send message to yourself",
		},
		{
			name: "empty content",
			m: &Message{
				SenderID:   1,
				ReceiverID: 2,
				Content:    "",
			},
			wantErr: true,
			errMsg:  "content cannot be empty",
		},
		{
			name: "content too long",
			m: &Message{
				SenderID:   1,
				ReceiverID: 2,
				Content:    string(make([]byte, 2001)),
			},
			wantErr: true,
			errMsg:  "content too long (max 2000)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestMessageEntity_IsReadMessage(t *testing.T) {
	tests := []struct {
		name   string
		m      *Message
		isRead bool
	}{
		{
			name: "read message",
			m: &Message{
				IsRead: true,
			},
			isRead: true,
		},
		{
			name: "unread message",
			m: &Message{
				IsRead: false,
			},
			isRead: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.m.IsReadMessage() != tt.isRead {
				t.Errorf("IsReadMessage() = %v, want %v", tt.m.IsReadMessage(), tt.isRead)
			}
		})
	}
}

func TestMessageEntity_MarkAsRead(t *testing.T) {
	m := &Message{
		IsRead: false,
		ReadAt: nil,
	}

	m.MarkAsRead()

	if !m.IsRead {
		t.Errorf("Message.IsRead after MarkAsRead = %v, want true", m.IsRead)
	}
	if m.ReadAt == nil {
		t.Error("Message.ReadAt after MarkAsRead should not be nil")
	}
}

func TestMessageEntity_Fields(t *testing.T) {
	now := time.Now()
	m := &Message{
		ID:         1,
		SenderID:   100,
		ReceiverID: 200,
		Content:    "测试私信",
		IsRead:     false,
		ReadAt:     nil,
		CreatedAt:  now,
	}

	if m.ID != 1 {
		t.Errorf("Message.ID = %v, want 1", m.ID)
	}
	if m.SenderID != 100 {
		t.Errorf("Message.SenderID = %v, want 100", m.SenderID)
	}
	if m.ReceiverID != 200 {
		t.Errorf("Message.ReceiverID = %v, want 200", m.ReceiverID)
	}
	if m.Content != "测试私信" {
		t.Errorf("Message.Content = %v, want 测试私信", m.Content)
	}
}