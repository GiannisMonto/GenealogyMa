package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/genealogy-ma/platform/internal/application/service"
)

// MockCommunityService 模拟社区服务
type MockCommunityService struct {
	profiles map[int64]*service.UserProfileDTO
	posts    map[int64]*service.PostDTO
	comments map[int64]*service.CommentDTO
	messages map[int64]*service.MessageDTO
	nextID   int64
}

func NewMockCommunityService() *MockCommunityService {
	return &MockCommunityService{
		profiles: make(map[int64]*service.UserProfileDTO),
		posts:    make(map[int64]*service.PostDTO),
		comments: make(map[int64]*service.CommentDTO),
		messages: make(map[int64]*service.MessageDTO),
		nextID:   1,
	}

}

func (m *MockCommunityService) CreateProfile(ctx context.Context, req *service.CreateProfileRequest) (*service.UserProfileDTO, error) {
	p := &service.UserProfileDTO{
		ID:       m.nextID,
		UserID:   req.UserID,
		Nickname: req.Nickname,
		Bio:      req.Bio,
		Gender:   req.Gender,
		Province: req.Province,
		City:     req.City,
	}
	m.nextID++
	m.profiles[p.ID] = p
	return p, nil
}

func (m *MockCommunityService) GetProfile(ctx context.Context, id int64) (*service.UserProfileDTO, error) {
	if p, ok := m.profiles[id]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *MockCommunityService) GetProfileByUserID(ctx context.Context, userID int64) (*service.UserProfileDTO, error) {
	for _, p := range m.profiles {
		if p.UserID == userID {
			return p, nil
		}
	}
	return nil, nil
}

func (m *MockCommunityService) UpdateProfile(ctx context.Context, id int64, req *service.UpdateProfileRequest) (*service.UserProfileDTO, error) {
	if p, ok := m.profiles[id]; ok {
		if req.Nickname != nil {
			p.Nickname = *req.Nickname
		}
		if req.Bio != nil {
			p.Bio = *req.Bio
		}
		return p, nil
	}
	return nil, nil
}

func (m *MockCommunityService) DeleteProfile(ctx context.Context, id int64) error {
	delete(m.profiles, id)
	return nil
}

func (m *MockCommunityService) CreatePost(ctx context.Context, req *service.CreatePostRequest) (*service.PostDTO, error) {
	p := &service.PostDTO{
		ID:      m.nextID,
		UserID:  req.UserID,
		Content: req.Content,
		Tag:     req.Tag,
	}
	m.nextID++
	m.posts[p.ID] = p
	return p, nil
}

func (m *MockCommunityService) GetPost(ctx context.Context, id int64) (*service.PostDTO, error) {
	if p, ok := m.posts[id]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *MockCommunityService) GetUserPosts(ctx context.Context, userID int64, page, pageSize int) ([]*service.PostDTO, int64, error) {
	result := make([]*service.PostDTO, 0)
	for _, p := range m.posts {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockCommunityService) GetPostsByTag(ctx context.Context, tag string, page, pageSize int) ([]*service.PostDTO, int64, error) {
	result := make([]*service.PostDTO, 0)
	for _, p := range m.posts {
		if p.Tag == tag {
			result = append(result, p)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockCommunityService) UpdatePost(ctx context.Context, id int64, req *service.UpdatePostRequest) (*service.PostDTO, error) {
	if p, ok := m.posts[id]; ok {
		if req.Content != nil {
			p.Content = *req.Content
		}
		if req.Tag != nil {
			p.Tag = *req.Tag
		}
		return p, nil
	}
	return nil, nil
}

func (m *MockCommunityService) DeletePost(ctx context.Context, id int64) error {
	delete(m.posts, id)
	return nil
}

func (m *MockCommunityService) LikePost(ctx context.Context, id int64) error {
	if p, ok := m.posts[id]; ok {
		p.LikeCount++
	}
	return nil
}

func (m *MockCommunityService) CreateComment(ctx context.Context, req *service.CreateCommentRequest) (*service.CommentDTO, error) {
	c := &service.CommentDTO{
		ID:       m.nextID,
		PostID:   req.PostID,
		UserID:   req.UserID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}
	m.nextID++
	m.comments[c.ID] = c
	return c, nil
}

func (m *MockCommunityService) GetComment(ctx context.Context, id int64) (*service.CommentDTO, error) {
	if c, ok := m.comments[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *MockCommunityService) GetPostComments(ctx context.Context, postID int64) ([]*service.CommentDTO, error) {
	result := make([]*service.CommentDTO, 0)
	for _, c := range m.comments {
		if c.PostID == postID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *MockCommunityService) UpdateComment(ctx context.Context, id int64, req *service.UpdateCommentRequest) (*service.CommentDTO, error) {
	if c, ok := m.comments[id]; ok {
		if req.Content != nil {
			c.Content = *req.Content
		}
		return c, nil
	}
	return nil, nil
}

func (m *MockCommunityService) DeleteComment(ctx context.Context, id int64) error {
	delete(m.comments, id)
	return nil
}

func (m *MockCommunityService) SendMessage(ctx context.Context, req *service.SendMessageRequest) (*service.MessageDTO, error) {
	msg := &service.MessageDTO{
		ID:         m.nextID,
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}
	m.nextID++
	m.messages[msg.ID] = msg
	return msg, nil
}

func (m *MockCommunityService) GetConversation(ctx context.Context, userID1, userID2 int64, page, pageSize int) ([]*service.MessageDTO, int64, error) {
	result := make([]*service.MessageDTO, 0)
	for _, msg := range m.messages {
		if (msg.SenderID == userID1 && msg.ReceiverID == userID2) ||
			(msg.SenderID == userID2 && msg.ReceiverID == userID1) {
			result = append(result, msg)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockCommunityService) GetUserMessages(ctx context.Context, userID int64, page, pageSize int) ([]*service.MessageDTO, int64, error) {
	result := make([]*service.MessageDTO, 0)
	for _, msg := range m.messages {
		if msg.SenderID == userID || msg.ReceiverID == userID {
			result = append(result, msg)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockCommunityService) MarkMessageAsRead(ctx context.Context, id int64) error {
	if msg, ok := m.messages[id]; ok {
		msg.IsRead = true
	}
	return nil
}

func (m *MockCommunityService) DeleteMessage(ctx context.Context, id int64) error {
	delete(m.messages, id)
	return nil
}

// ===== 用户资料测试 =====

func TestCommunityController_GetProfile(t *testing.T) {
	mockService := NewMockCommunityService()
	mockService.profiles[1] = &service.UserProfileDTO{ID: 1, UserID: 100, Nickname: "测试用户"}

	r := setupTestRouter()
	r.GET("/profiles/:id", func(ctx *gin.Context) {
		id := int64(1)
		profile, err := mockService.GetProfile(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if profile == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": profile})
	})

	req, _ := http.NewRequest("GET", "/profiles/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected data in response")
	}
	if data["nickname"] != "测试用户" {
		t.Errorf("Expected nickname '测试用户', got '%v'", data["nickname"])
	}
}

func TestCommunityController_CreateProfile(t *testing.T) {
	mockService := NewMockCommunityService()

	r := setupTestRouter()
	r.POST("/profiles", func(ctx *gin.Context) {
		var req service.CreateProfileRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		profile, err := mockService.CreateProfile(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": profile})
	})

	body := `{"user_id":100,"nickname":"新用户","bio":"大家好","gender":"男","province":"广东","city":"深圳"}`
	req, _ := http.NewRequest("POST", "/profiles", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected data in response")
	}
	if data["nickname"] != "新用户" {
		t.Errorf("Expected nickname '新用户', got '%v'", data["nickname"])
	}
}

// ===== 动态测试 =====

func TestCommunityController_GetPost(t *testing.T) {
	mockService := NewMockCommunityService()
	mockService.posts[1] = &service.PostDTO{ID: 1, UserID: 100, Content: "测试动态", Tag: "家族"}

	r := setupTestRouter()
	r.GET("/posts/:id", func(ctx *gin.Context) {
		id := int64(1)
		post, err := mockService.GetPost(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if post == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": post})
	})

	req, _ := http.NewRequest("GET", "/posts/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected data in response")
	}
	if data["content"] != "测试动态" {
		t.Errorf("Expected content '测试动态', got '%v'", data["content"])
	}
}

func TestCommunityController_CreatePost(t *testing.T) {
	mockService := NewMockCommunityService()

	r := setupTestRouter()
	r.POST("/posts", func(ctx *gin.Context) {
		var req service.CreatePostRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		post, err := mockService.CreatePost(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": post})
	})

	body := `{"user_id":100,"content":"今天天气不错","tag":"日常"}`
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected data in response")
	}
	if data["tag"] != "日常" {
		t.Errorf("Expected tag '日常', got '%v'", data["tag"])
	}
}

func TestCommunityController_DeletePost(t *testing.T) {
	mockService := NewMockCommunityService()
	mockService.posts[1] = &service.PostDTO{ID: 1, UserID: 100, Content: "测试动态"}

	r := setupTestRouter()
	r.DELETE("/posts/:id", func(ctx *gin.Context) {
		id := int64(1)
		err := mockService.DeletePost(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if _, exists := mockService.posts[1]; exists {
		t.Errorf("Post should have been deleted")
	}
}

func TestCommunityController_LikePost(t *testing.T) {
	mockService := NewMockCommunityService()
	mockService.posts[1] = &service.PostDTO{ID: 1, UserID: 100, Content: "测试动态", LikeCount: 0}

	r := setupTestRouter()
	r.POST("/posts/:id/like", func(ctx *gin.Context) {
		id := int64(1)
		err := mockService.LikePost(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
	})

	req, _ := http.NewRequest("POST", "/posts/1/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if mockService.posts[1].LikeCount != 1 {
		t.Errorf("Expected like count 1, got %d", mockService.posts[1].LikeCount)
	}
}

// ===== 评论测试 =====

func TestCommunityController_CreateComment(t *testing.T) {
	mockService := NewMockCommunityService()

	r := setupTestRouter()
	r.POST("/comments", func(ctx *gin.Context) {
		var req service.CreateCommentRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		comment, err := mockService.CreateComment(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": comment})
	})

	body := `{"post_id":1,"user_id":100,"content":"说得好"}`
	req, _ := http.NewRequest("POST", "/comments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected data in response")
	}
	if data["content"] != "说得好" {
		t.Errorf("Expected content '说得好', got '%v'", data["content"])
	}
}

func TestCommunityController_GetPostComments(t *testing.T) {
	mockService := NewMockCommunityService()
	mockService.comments[1] = &service.CommentDTO{ID: 1, PostID: 1, UserID: 100, Content: "评论1"}
	mockService.comments[2] = &service.CommentDTO{ID: 2, PostID: 1, UserID: 101, Content: "评论2"}

	r := setupTestRouter()
	r.GET("/comments/post/:post_id", func(ctx *gin.Context) {
		postID := int64(1)
		comments, err := mockService.GetPostComments(ctx, postID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": comments})
	})

	req, _ := http.NewRequest("GET", "/comments/post/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(data))
	}
}

// ===== 私信测试 =====

func TestCommunityController_SendMessage(t *testing.T) {
	mockService := NewMockCommunityService()

	r := setupTestRouter()
	r.POST("/messages", func(ctx *gin.Context) {
		var req service.SendMessageRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		msg, err := mockService.SendMessage(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": msg})
	})

	body := `{"sender_id":1,"receiver_id":2,"content":"你好"}`
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected data in response")
	}
	if data["content"] != "你好" {
		t.Errorf("Expected content '你好', got '%v'", data["content"])
	}
}

func TestCommunityController_GetConversation(t *testing.T) {
	mockService := NewMockCommunityService()
	mockService.messages[1] = &service.MessageDTO{ID: 1, SenderID: 1, ReceiverID: 2, Content: "你好"}
	mockService.messages[2] = &service.MessageDTO{ID: 2, SenderID: 2, ReceiverID: 1, Content: "你好啊"}

	r := setupTestRouter()
	r.GET("/messages/conversation/:user_id", func(ctx *gin.Context) {
		userID := int64(2)
		messages, _, err := mockService.GetConversation(ctx, 1, userID, 1, 10)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": messages})
	})

	req, _ := http.NewRequest("GET", "/messages/conversation/2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(data))
	}
}

func TestCommunityController_MarkMessageAsRead(t *testing.T) {
	mockService := NewMockCommunityService()
	mockService.messages[1] = &service.MessageDTO{ID: 1, SenderID: 1, ReceiverID: 2, Content: "你好", IsRead: false}

	r := setupTestRouter()
	r.PUT("/messages/:id/read", func(ctx *gin.Context) {
		id := int64(1)
		err := mockService.MarkMessageAsRead(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "标记已读成功"})
	})

	req, _ := http.NewRequest("PUT", "/messages/1/read", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if !mockService.messages[1].IsRead {
		t.Error("Message should be marked as read")
	}
}
