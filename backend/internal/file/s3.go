package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Service handles S3 file operations
type S3Service struct {
	client     *s3.Client
	bucketName string
	region     string
}

// NewS3Service creates a new S3 service instance
func NewS3Service(region, bucketName, accessKeyID, secretAccessKey string) (*S3Service, error) {
	// Create AWS config
	var cfg aws.Config
	var err error

	if accessKeyID != "" && secretAccessKey != "" {
		// Use provided credentials
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				accessKeyID,
				secretAccessKey,
				"",
			)),
		)
	} else {
		// Use default credentials (IAM role, environment variables, etc.)
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client:     client,
		bucketName: bucketName,
		region:     region,
	}, nil
}

// UploadFile uploads a file to S3
func (s *S3Service) UploadFile(file multipart.File, filename string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	timestamp := time.Now().Format("20060102150405")
	uniqueName := fmt.Sprintf("%s_%s%s", baseName, timestamp, ext)
	
	// Create S3 key (path in bucket)
	s3Key := fmt.Sprintf("uploads/%s", uniqueName)

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Upload to S3
	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(s3Key),
		Body:        strings.NewReader(string(content)),
		ContentType: aws.String(GetMimeType(filename)),
		ACL:         types.ObjectCannedACLPrivate, // Private by default
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return the S3 key as the file path
	return s3Key, nil
}

// DownloadFile downloads a file from S3
func (s *S3Service) DownloadFile(s3Key string) ([]byte, error) {
	result, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}
	defer result.Body.Close()

	content, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read S3 object body: %w", err)
	}

	return content, nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(s3Key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// FileExists checks if a file exists in S3
func (s *S3Service) FileExists(s3Key string) (bool, error) {
	_, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		// Check if it's a "not found" error
		var notFoundErr *types.NoSuchKey
		if errors.As(err, &notFoundErr) {
			return false, nil // File doesn't exist
		}
		return false, fmt.Errorf("failed to check file existence in S3: %w", err)
	}

	return true, nil
}

// GetFileSize returns the size of a file in S3
func (s *S3Service) GetFileSize(s3Key string) (int64, error) {
	result, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get file size from S3: %w", err)
	}

	if result.ContentLength == nil {
		return 0, fmt.Errorf("content length is nil")
	}

	return *result.ContentLength, nil
}

// GetSignedURL generates a signed URL for temporary access to a file
func (s *S3Service) GetSignedURL(s3Key string, expires time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	
	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s3Key),
	}, s3.WithPresignExpires(expires))

	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return request.URL, nil
} 