package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Environment    string
	DatabaseURL    string
	JWTSecret      string
	AWSCognito     CognitoConfig
	FileUpload     FileUploadConfig
	Encryption     EncryptionConfig
}

// CognitoConfig holds AWS Cognito configuration
type CognitoConfig struct {
	Region          string
	UserPoolID      string
	ClientID        string
	ClientSecret    string
}

// FileUploadConfig holds file upload configuration
type FileUploadConfig struct {
	MaxFileSize     int64
	AllowedTypes    []string
	UploadPath      string
}

// EncryptionConfig holds encryption configuration
type EncryptionConfig struct {
	Key             string
	Algorithm       string
}

// New creates a new Config instance
func New() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseURL: getEnv("DATABASE_URL", "ktchat:password@tcp(localhost:3306)/ktchat?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		AWSCognito: CognitoConfig{
			Region:       getEnv("AWS_REGION", "ap-northeast-1"),
			UserPoolID:   getEnv("COGNITO_USER_POOL_ID", ""),
			ClientID:     getEnv("COGNITO_CLIENT_ID", ""),
			ClientSecret: getEnv("COGNITO_CLIENT_SECRET", ""),
		},
		FileUpload: FileUploadConfig{
			MaxFileSize:  50 * 1024 * 1024, // 50MB
			AllowedTypes: []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".mp4", ".avi", ".mov"},
			UploadPath:   getEnv("UPLOAD_PATH", "./uploads"),
		},
		Encryption: EncryptionConfig{
			Key:       getEnv("ENCRYPTION_KEY", "your-encryption-key-32-bytes-long"),
			Algorithm: "AES-256-GCM",
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 