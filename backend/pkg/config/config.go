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
	S3             S3Config
}

// CognitoConfig holds AWS Cognito configuration
type CognitoConfig struct {
	Region          string
	UserPoolID      string
	ClientID        string
	ClientSecret    string
}

// S3Config holds AWS S3 configuration
type S3Config struct {
	Region          string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	UseS3           bool
}

// FileUploadConfig holds file upload configuration
type FileUploadConfig struct {
	MaxFileSize     int64
	AllowedTypes    []string
	UploadPath      string
	StorageType     string // "local" or "s3"
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
			StorageType:  getEnv("STORAGE_TYPE", "local"), // "local" or "s3"
		},
		Encryption: EncryptionConfig{
			Key:       getEnv("ENCRYPTION_KEY", "your-encryption-key-32-bytes-long"),
			Algorithm: "AES-256-GCM",
		},
		S3: S3Config{
			Region:          getEnv("AWS_REGION", "ap-northeast-1"),
			BucketName:      getEnv("S3_BUCKET_NAME", ""),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			UseS3:           getEnv("STORAGE_TYPE", "local") == "s3",
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