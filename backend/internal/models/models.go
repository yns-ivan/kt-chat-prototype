package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username  string    `json:"username" gorm:"uniqueIndex;type:varchar(255);not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	CognitoID string    `json:"cognito_id" gorm:"uniqueIndex;type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// ChatRoom represents a chat room
type ChatRoom struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null;index"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedBy   string    `json:"created_by" gorm:"type:uuid;not null;index"`
	IsPrivate   bool      `json:"is_private" gorm:"default:false;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	CreatedByUser User           `json:"created_by_user" gorm:"foreignKey:CreatedBy"`
	Participants  []RoomParticipant `json:"participants" gorm:"foreignKey:RoomID"`
	Messages      []Message      `json:"messages" gorm:"foreignKey:RoomID"`
}

// RoomParticipant represents a user's participation in a chat room
type RoomParticipant struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RoomID    string    `json:"room_id" gorm:"type:uuid;not null;index"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
	JoinedAt  time.Time `json:"joined_at" gorm:"autoCreateTime"`
	LeftAt    *time.Time `json:"left_at"`
	
	// Relationships
	Room ChatRoom `json:"room" gorm:"foreignKey:RoomID"`
	User User     `json:"user" gorm:"foreignKey:UserID"`
}

// Message represents a chat message
type Message struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RoomID    string    `json:"room_id" gorm:"type:uuid;not null;index"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Encrypted bool      `json:"encrypted" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Room      ChatRoom     `json:"room" gorm:"foreignKey:RoomID"`
	User      User         `json:"user" gorm:"foreignKey:UserID"`
	Files     []FileAttachment `json:"files" gorm:"foreignKey:MessageID"`
}

// FileAttachment represents a file attached to a message
type FileAttachment struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	MessageID   string    `json:"message_id" gorm:"type:uuid;not null;index"`
	FileName    string    `json:"file_name" gorm:"type:varchar(255);not null"`
	FilePath    string    `json:"file_path" gorm:"type:varchar(500);not null"`
	FileSize    int64     `json:"file_size"`
	MimeType    string    `json:"mime_type" gorm:"type:varchar(100)"`
	FileType    string    `json:"file_type" gorm:"type:varchar(50);index"` // image, pdf, video
	ThumbnailPath string  `json:"thumbnail_path" gorm:"type:varchar(500)"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// Relationships
	Message Message `json:"message" gorm:"foreignKey:MessageID"`
}

// BeforeCreate is a GORM hook to set ID before creating records (fallback for older PostgreSQL versions)
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = generateUUID()
	}
	return nil
}

func (r *ChatRoom) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = generateUUID()
	}
	return nil
}

func (p *RoomParticipant) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = generateUUID()
	}
	return nil
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = generateUUID()
	}
	return nil
}

func (f *FileAttachment) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = generateUUID()
	}
	return nil
}

// generateUUID generates a UUID (simplified for now, should use proper UUID library)
func generateUUID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
} 