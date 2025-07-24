package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ktchat/backend/internal/encryption"
	"ktchat/backend/internal/file"
	"ktchat/backend/internal/models"
	"ktchat/backend/internal/websocket"
	"ktchat/backend/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service handles chat-related operations
type Service struct {
	db  *gorm.DB
	hub *websocket.Hub
	cfg *config.Config
}

// NewService creates a new chat service
func NewService(db *gorm.DB, hub *websocket.Hub) *Service {
	cfg := config.New()
	return &Service{
		db:  db,
		hub: hub,
		cfg: cfg,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// CreateRoomRequest represents a room creation request
type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

// MessageRequest represents a message request
type MessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// Login handles user authentication
func (s *Service) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, you would validate against AWS Cognito
	// For now, we'll use a simple mock authentication
	if req.Username == "admin" && req.Password == "password" {
		// Check if user exists, create if not
		var user models.User
		if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create user if not exists
				user = models.User{
					ID:        "user-123",
					Username:  req.Username,
					Email:     "admin@example.com",
					CognitoID: "mock-cognito-id",
				}
				if err := s.db.Create(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
					return
				}
				// Log successful user creation
				fmt.Printf("Created user: %+v\n", user)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
				return
			}
		} else {
			// Log existing user found
			fmt.Printf("Found existing user: %+v\n", user)
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

// Register handles user registration
func (s *Service) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, you would register with AWS Cognito
	// For now, we'll just return success
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"username": req.Username,
			"email":    req.Email,
		},
	})
}

// RefreshToken handles token refresh
func (s *Service) RefreshToken(c *gin.Context) {
	// Implementation for token refresh
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// GetRooms returns all available chat rooms
func (s *Service) GetRooms(c *gin.Context) {
	var rooms []models.ChatRoom
	if err := s.db.Preload("CreatedByUser").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// CreateRoom creates a new chat room
func (s *Service) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	room := models.ChatRoom{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
		IsPrivate:   req.IsPrivate,
	}

	if err := s.db.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"room": room})
}

// GetMessages returns messages for a specific room
func (s *Service) GetMessages(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	var messages []models.Message
	if err := s.db.Preload("User").Preload("Files").Where("room_id = ?", roomID).Order("created_at desc").Limit(50).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	// Decrypt messages
	for i := range messages {
		if messages[i].Encrypted {
			decrypted, err := encryption.Decrypt(messages[i].Content, s.cfg.Encryption.Key)
			if err == nil {
				messages[i].Content = decrypted
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// JoinRoom adds a user to a chat room
func (s *Service) JoinRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	userID := c.GetString("user_id")

	if roomID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID and user ID are required"})
		return
	}

	// Check if user is already in the room
	var existing models.RoomParticipant
	if err := s.db.Where("room_id = ? AND user_id = ? AND left_at IS NULL", roomID, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "User already in room"})
		return
	}

	participant := models.RoomParticipant{
		ID:       uuid.New().String(),
		RoomID:   roomID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}

	if err := s.db.Create(&participant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined room successfully"})
}

// LeaveRoom removes a user from a chat room
func (s *Service) LeaveRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	userID := c.GetString("user_id")

	if roomID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID and user ID are required"})
		return
	}

	now := time.Now()
	if err := s.db.Model(&models.RoomParticipant{}).Where("room_id = ? AND user_id = ? AND left_at IS NULL", roomID, userID).Update("left_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left room successfully"})
}

// SendMessage sends a message to a chat room
func (s *Service) SendMessage(c *gin.Context) {
	roomID := c.Param("roomID")
	userID := c.GetString("user_id")

	if roomID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID and user ID are required"})
		return
	}

	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Encrypt message content
	encryptedContent, err := encryption.Encrypt(req.Content, s.cfg.Encryption.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt message"})
		return
	}

	message := models.Message{
		ID:        uuid.New().String(),
		RoomID:    roomID,
		UserID:    userID,
		Content:   encryptedContent,
		Encrypted: true,
	}

	if err := s.db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// Broadcast message via WebSocket
	wsMessage := websocket.Message{
		Type:      "message",
		RoomID:    roomID,
		UserID:    userID,
		Content:   req.Content, // Send decrypted content to WebSocket
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if msgBytes, err := json.Marshal(wsMessage); err == nil {
		s.hub.BroadcastToRoom(roomID, msgBytes)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

// UploadFile handles file uploads
func (s *Service) UploadFile(c *gin.Context) {
	roomID := c.Param("roomID")
	userID := c.GetString("user_id")

	if roomID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID and user ID are required"})
		return
	}

	uploadedFile, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer uploadedFile.Close()

	// Validate file size
	if header.Size > s.cfg.FileUpload.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// Validate file type
	fileExt := strings.ToLower(header.Filename[strings.LastIndex(header.Filename, "."):])
	allowed := false
	for _, allowedType := range s.cfg.FileUpload.AllowedTypes {
		if fileExt == allowedType {
			allowed = true
			break
		}
	}
	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		return
	}

	// Save file
	filePath, err := file.SaveFile(uploadedFile, header.Filename, s.cfg.FileUpload.UploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create file attachment record
	attachment := models.FileAttachment{
		ID:        uuid.New().String(),
		FileName:  header.Filename,
		FilePath:  filePath,
		FileSize:  header.Size,
		MimeType:  header.Header.Get("Content-Type"),
		FileType:  file.GetFileType(fileExt),
	}

	if err := s.db.Create(&attachment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": attachment})
} 