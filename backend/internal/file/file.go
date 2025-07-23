package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileHandler handles file operations
type FileHandler struct {
	UploadPath string
}

// NewFileHandler creates a new file handler
func NewFileHandler(uploadPath string) *FileHandler {
	return &FileHandler{
		UploadPath: uploadPath,
	}
}

// SaveFile saves an uploaded file to the filesystem
func (fh *FileHandler) SaveFile(file multipart.File, filename, uploadPath string) (string, error) {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	timestamp := time.Now().Format("20060102150405")
	uniqueName := fmt.Sprintf("%s_%s%s", baseName, timestamp, ext)
	filePath := filepath.Join(uploadPath, uniqueName)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return filePath, nil
}

// SaveFile is a standalone function for saving files
func SaveFile(file multipart.File, filename, uploadPath string) (string, error) {
	fh := NewFileHandler(uploadPath)
	return fh.SaveFile(file, filename, uploadPath)
}

// GetFileType determines the file type based on extension
func GetFileType(extension string) string {
	ext := strings.ToLower(extension)
	
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return "image"
	case ".pdf":
		return "pdf"
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm":
		return "video"
	case ".mp3", ".wav", ".ogg", ".aac":
		return "audio"
	case ".doc", ".docx", ".txt", ".rtf":
		return "document"
	default:
		return "other"
	}
}

// GenerateThumbnail generates a thumbnail for images and videos
// This is a placeholder - in production, you would use libraries like:
// - github.com/disintegration/imaging for images
// - ffmpeg for video thumbnails
func GenerateThumbnail(filePath, outputPath string) error {
	// For now, we'll just copy the file as a placeholder
	// In production, implement actual thumbnail generation
	
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	destFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// ValidateFileSize checks if the file size is within limits
func ValidateFileSize(size int64, maxSize int64) bool {
	return size <= maxSize
}

// ValidateFileType checks if the file type is allowed
func ValidateFileType(filename string, allowedTypes []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// DeleteFile deletes a file from the filesystem
func DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// GetFileSize returns the size of a file
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}
	return fileInfo.Size(), nil
}

// GetMimeType returns the MIME type of a file based on its extension
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".mp4":
		return "video/mp4"
	case ".avi":
		return "video/x-msvideo"
	case ".mov":
		return "video/quicktime"
	case ".txt":
		return "text/plain"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
} 