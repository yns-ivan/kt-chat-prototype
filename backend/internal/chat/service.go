package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"ktchat/backend/internal/auth"
	"ktchat/backend/internal/encryption"
	"ktchat/backend/internal/file"
	"ktchat/backend/internal/models"
	"ktchat/backend/internal/websocket"
	"ktchat/backend/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service handles chat-related operations
type Service struct {
	db           *gorm.DB
	hub          *websocket.Hub
	cfg          *config.Config
	cognitoAuth  *auth.CognitoService
}

// NewService creates a new chat service
func NewService(db *gorm.DB, hub *websocket.Hub, cognitoAuth *auth.CognitoService) *Service {
	cfg := config.New()
	return &Service{
		db:          db,
		hub:         hub,
		cfg:         cfg,
		cognitoAuth: cognitoAuth,
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

// SearchRequest represents a search request
type SearchRequest struct {
	Query     string `json:"query" binding:"required"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	UserID    string `json:"user_id"`
}

// Login handles user authentication with AWS Cognito
func (s *Service) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("DEBUG: Login request binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		}})
		return
	}

	// Debug: Log the received parameters
	fmt.Printf("DEBUG: Login request received - Username: %s, Password length: %d\n", req.Username, len(req.Password))

	// Authenticate with AWS Cognito
	cognitoReq := auth.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	cognitoResp, userInfo, err := s.cognitoAuth.Login(c.Request.Context(), cognitoReq)
	if err != nil {
		// Debug: Log the full error
		fmt.Printf("DEBUG: Cognito login error: %v\n", err)
		// Extract structured error information from AWS Cognito
		cognitoError := s.cognitoAuth.ExtractCognitoError(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": cognitoError})
		return
	}

	// Check if user exists in our database, create if not
	var user models.User
	if err := s.db.Where("cognito_id = ?", userInfo.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Check if user exists by username (for existing users before Cognito integration)
			if err := s.db.Where("username = ?", userInfo.Username).First(&user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Create new user if not exists
					user = models.User{
						ID:        uuid.New().String(),
						Username:  userInfo.Username,
						Email:     userInfo.Email,
						CognitoID: userInfo.ID,
					}
					if err := s.db.Create(&user).Error; err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
							"code":    "DATABASE_ERROR",
							"message": "Failed to create user",
						}})
						return
					}
					fmt.Printf("Created new user: %+v\n", user)
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
						"code":    "DATABASE_ERROR",
						"message": "Database error",
					}})
					return
				}
			} else {
				// Found existing user by username, update with Cognito ID
				user.CognitoID = userInfo.ID
				user.Email = userInfo.Email // Update email from Cognito
				if err := s.db.Save(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
						"code":    "DATABASE_ERROR",
						"message": "Failed to update user",
					}})
					return
				}
				fmt.Printf("Updated existing user with Cognito ID: %+v\n", user)
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
				"code":    "DATABASE_ERROR",
				"message": "Database error",
			}})
			return
		}
	} else {
		// Update user info if needed
		if user.Username != userInfo.Username || user.Email != userInfo.Email {
			user.Username = userInfo.Username
			user.Email = userInfo.Email
			if err := s.db.Save(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
					"code":    "DATABASE_ERROR",
					"message": "Failed to update user",
				}})
				return
			}
		}
		fmt.Printf("Found existing user by Cognito ID: %+v\n", user)
	}

	// Generate custom JWT token for internal use
	customToken, err := s.cognitoAuth.GenerateCustomToken(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "TOKEN_ERROR",
			"message": "Failed to generate token",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": customToken,
		"cognito_tokens": gin.H{
			"access_token":  cognitoResp.AccessToken,
			"refresh_token": cognitoResp.RefreshToken,
			"id_token":      cognitoResp.IDToken,
			"expires_in":    cognitoResp.ExpiresIn,
			"token_type":    cognitoResp.TokenType,
		},
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Register handles user registration with AWS Cognito
func (s *Service) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		}})
		return
	}

	// Register with AWS Cognito
	cognitoReq := auth.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	err := s.cognitoAuth.Register(c.Request.Context(), cognitoReq)
	if err != nil {
		// Extract structured error information from AWS Cognito
		cognitoError := s.cognitoAuth.ExtractCognitoError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": cognitoError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully. Please check your email for confirmation.",
		"user": gin.H{
			"username": req.Username,
			"email":    req.Email,
		},
	})
}

// RefreshToken handles token refresh with AWS Cognito
func (s *Service) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		}})
		return
	}

	// Refresh token with AWS Cognito
	cognitoReq := auth.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	cognitoResp, err := s.cognitoAuth.RefreshToken(c.Request.Context(), cognitoReq)
	if err != nil {
		// Extract structured error information from AWS Cognito
		cognitoError := s.cognitoAuth.ExtractCognitoError(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": cognitoError})
		return
	}

	// Extract user info from the new ID token
	userInfo, err := s.cognitoAuth.ValidateToken(c.Request.Context(), cognitoResp.IDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "TOKEN_VALIDATION_ERROR",
			"message": "Failed to validate token",
		}})
		return
	}

	// Generate new custom JWT token
	customToken, err := s.cognitoAuth.GenerateCustomToken(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "TOKEN_GENERATION_ERROR",
			"message": "Failed to generate token",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": customToken,
		"cognito_tokens": gin.H{
			"access_token": cognitoResp.AccessToken,
			"id_token":     cognitoResp.IDToken,
			"expires_in":   cognitoResp.ExpiresIn,
			"token_type":   cognitoResp.TokenType,
		},
	})
}

// GetRooms returns all available chat rooms
func (s *Service) GetRooms(c *gin.Context) {
	var rooms []models.ChatRoom
	if err := s.db.Preload("CreatedByUser").Preload("Participants.User").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
		return
	}

	// Filter out participants who have left the room (left_at IS NOT NULL)
	for i := range rooms {
		var activeParticipants []models.RoomParticipant
		for _, participant := range rooms[i].Participants {
			if participant.LeftAt == nil {
				activeParticipants = append(activeParticipants, participant)
			}
		}
		rooms[i].Participants = activeParticipants
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

// GetRoom returns a specific room with its participants
func (s *Service) GetRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	var room models.ChatRoom
	if err := s.db.Preload("CreatedByUser").Preload("Participants.User").Where("id = ?", roomID).First(&room).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// Filter out participants who have left the room
	var activeParticipants []models.RoomParticipant
	for _, participant := range room.Participants {
		if participant.LeftAt == nil {
			activeParticipants = append(activeParticipants, participant)
		}
	}
	room.Participants = activeParticipants

	c.JSON(http.StatusOK, gin.H{"room": room})
}

// GetMessages returns messages for a specific room
func (s *Service) GetMessages(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	var messages []models.Message
	if err := s.db.Preload("User").Preload("Files").Where("room_id = ?", roomID).Order("created_at asc").Limit(50).Find(&messages).Error; err != nil {
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

// SearchMessages searches messages in a specific room
func (s *Service) SearchMessages(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Build query
	query := s.db.Model(&models.Message{}).
		Preload("User").
		Preload("Files").
		Where("room_id = ?", roomID)

	// Add date range filter if provided
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			query = query.Where("created_at >= ?", startDate)
		}
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			// Add one day to include the end date
			endDate = endDate.Add(24 * time.Hour)
			query = query.Where("created_at < ?", endDate)
		}
	}

	// Add user filter if provided
	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}

	// Get total count for pagination
	var totalCount int64
	countQuery := query
	if err := countQuery.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count messages"})
		return
	}

	// Search in decrypted content
	var messages []models.Message
	if err := query.Order("created_at desc").Limit(req.Limit).Offset(req.Offset).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	// Decrypt and search in content
	var searchResults []models.Message
	for _, message := range messages {
		// Decrypt message content
		decryptedContent := message.Content
		if message.Encrypted {
			decrypted, err := encryption.Decrypt(message.Content, s.cfg.Encryption.Key)
			if err == nil {
				decryptedContent = decrypted
			}
		}

		// Check if query matches decrypted content (case-insensitive)
		if strings.Contains(strings.ToLower(decryptedContent), strings.ToLower(req.Query)) {
			// Update the message content with decrypted version
			message.Content = decryptedContent
			searchResults = append(searchResults, message)
		}
	}

	// Sort results by creation date (newest first)
	sort.Slice(searchResults, func(i, j int) bool {
		return searchResults[i].CreatedAt.After(searchResults[j].CreatedAt)
	})

	c.JSON(http.StatusOK, gin.H{
		"messages": searchResults,
		"total":    totalCount,
		"limit":    req.Limit,
		"offset":   req.Offset,
		"query":    req.Query,
	})
}

// JoinRoom adds a user to a chat room
func (s *Service) JoinRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	cognitoUserID := c.GetString("user_id")

	if roomID == "" || cognitoUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID and user ID are required"})
		return
	}

	// Look up the internal user ID using the Cognito ID
	var user models.User
	if err := s.db.Where("cognito_id = ?", cognitoUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Check if user is already in the room
	var existing models.RoomParticipant
	if err := s.db.Where("room_id = ? AND user_id = ? AND left_at IS NULL", roomID, user.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "User already in room"})
		return
	}

	participant := models.RoomParticipant{
		ID:       uuid.New().String(),
		RoomID:   roomID,
		UserID:   user.ID, // Use the internal user ID
		JoinedAt: time.Now(),
	}

	if err := s.db.Create(&participant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join room"})
		return
	}

	// Broadcast participant count update
	s.broadcastParticipantCount(roomID)

	c.JSON(http.StatusOK, gin.H{"message": "Joined room successfully"})
}

// LeaveRoom removes a user from a chat room
func (s *Service) LeaveRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	cognitoUserID := c.GetString("user_id")

	if roomID == "" || cognitoUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID and user ID are required"})
		return
	}

	// Look up the internal user ID using the Cognito ID
	var user models.User
	if err := s.db.Where("cognito_id = ?", cognitoUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	now := time.Now()
	if err := s.db.Model(&models.RoomParticipant{}).Where("room_id = ? AND user_id = ? AND left_at IS NULL", roomID, user.ID).Update("left_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave room"})
		return
	}

	// Broadcast participant count update
	s.broadcastParticipantCount(roomID)

	c.JSON(http.StatusOK, gin.H{"message": "Left room successfully"})
}

// SendMessage sends a message to a chat room
func (s *Service) SendMessage(c *gin.Context) {
	roomID := c.Param("roomID")
	cognitoUserID := c.GetString("user_id")

	if roomID == "" || cognitoUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID and user ID are required"})
		return
	}

	// Look up the internal user ID using the Cognito ID
	var user models.User
	if err := s.db.Where("cognito_id = ?", cognitoUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
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
		UserID:    user.ID, // Use the internal user ID
		Content:   encryptedContent,
		Encrypted: true,
	}

	if err := s.db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// Get current participant count
	var participantCount int64
	s.db.Model(&models.RoomParticipant{}).Where("room_id = ? AND left_at IS NULL", roomID).Count(&participantCount)

	// Broadcast message via WebSocket
	wsMessage := websocket.Message{
		Type:             "message",
		RoomID:           roomID,
		UserID:           user.ID, // Use the internal user ID
		Username:         user.Username, // Include the username
		Content:          req.Content, // Send decrypted content to WebSocket
		Timestamp:        message.CreatedAt.Format(time.RFC3339), // Use the same timestamp as the database record
		ParticipantCount: int(participantCount),
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

// ConfirmUser handles user account confirmation
func (s *Service) ConfirmUser(c *gin.Context) {
	var req struct {
		Username         string `json:"username" binding:"required"`
		ConfirmationCode string `json:"confirmation_code" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("DEBUG: ConfirmUser request binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		}})
		return
	}

	fmt.Printf("DEBUG: ConfirmUser request - Username: %s, Code: %s\n", req.Username, req.ConfirmationCode)

	// Check if Cognito service is available
	if s.cognitoAuth == nil {
		fmt.Printf("DEBUG: Cognito service is nil\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "SERVICE_UNAVAILABLE",
			"message": "Authentication service is not available",
		}})
		return
	}

	err := s.cognitoAuth.ConfirmUser(c.Request.Context(), req.Username, req.ConfirmationCode)
	if err != nil {
		fmt.Printf("DEBUG: ConfirmUser Cognito error: %v\n", err)
		cognitoError := s.cognitoAuth.ExtractCognitoError(err)
		fmt.Printf("DEBUG: Extracted Cognito error: %+v\n", cognitoError)
		c.JSON(http.StatusBadRequest, gin.H{"error": cognitoError})
		return
	}

	fmt.Printf("DEBUG: User confirmed successfully in service\n")
	c.JSON(http.StatusOK, gin.H{"message": "User confirmed successfully"})
}

// ResendConfirmationCode handles resending confirmation codes
func (s *Service) ResendConfirmationCode(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		}})
		return
	}

	// Check if Cognito service is available
	if s.cognitoAuth == nil {
		fmt.Printf("DEBUG: Cognito service is nil\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "SERVICE_UNAVAILABLE",
			"message": "Authentication service is not available",
		}})
		return
	}

	err := s.cognitoAuth.ResendConfirmationCode(c.Request.Context(), req.Username)
	if err != nil {
		cognitoError := s.cognitoAuth.ExtractCognitoError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": cognitoError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation code sent successfully"})
}

// broadcastParticipantCount broadcasts the current participant count to all clients in a room
func (s *Service) broadcastParticipantCount(roomID string) {
	var participantCount int64
	s.db.Model(&models.RoomParticipant{}).Where("room_id = ? AND left_at IS NULL", roomID).Count(&participantCount)

	wsMessage := websocket.Message{
		Type:             "participant_count",
		RoomID:           roomID,
		ParticipantCount: int(participantCount),
	}

	if msgBytes, err := json.Marshal(wsMessage); err == nil {
		s.hub.BroadcastToRoom(roomID, msgBytes)
	}
} 