package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
)

// StorageProvider defines the interface for file storage operations
type StorageProvider interface {
	// UploadFile uploads a file to the storage provider
	UploadFile(file io.Reader, filename string, contentType string) (string, error)
	
	// GetFileURL returns the public URL for a file
	GetFileURL(filename string) (string, error)
	
	// DeleteFile deletes a file from the storage provider
	DeleteFile(filename string) error
	
	// GetFile retrieves a file from the storage provider
	GetFile(filename string) (io.ReadCloser, error)
}

// LocalStorageProvider implements file storage using the local filesystem
type LocalStorageProvider struct {
	basePath string
}

// NewLocalStorageProvider creates a new local storage provider
func NewLocalStorageProvider(basePath string) (*LocalStorageProvider, error) {
	// Ensure the base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	
	return &LocalStorageProvider{
		basePath: basePath,
	}, nil
}

// NewStorageProvider creates a new storage provider based on the configuration
func NewStorageProvider(cfg *config.Config) (StorageProvider, error) {
	providerType := cfg.StorageProvider
	if providerType == "" {
		providerType = "local" // Default to local storage
	}
	
	switch strings.ToLower(providerType) {
	case "local":
		// Use default path if not specified
		localPath := "./uploads"
		if cfg.OSSBucketName != "" {
			// Use OSS bucket name as a fallback for local path
			localPath = "./" + cfg.OSSBucketName
		}
		return NewLocalStorageProvider(localPath)
	case "oss":
		return NewOSSStorageProvider(cfg)
	default:
		// Default to local storage for unknown providers
		localPath := "./uploads"
		if cfg.OSSBucketName != "" {
			// Use OSS bucket name as a fallback for local path
			localPath = "./" + cfg.OSSBucketName
		}
		return NewLocalStorageProvider(localPath)
	}
}

// UploadFile uploads a file to local storage
func (p *LocalStorageProvider) UploadFile(file io.Reader, filename string, contentType string) (string, error) {
	// Create the full file path
	filePath := filepath.Join(p.basePath, filename)
	
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()
	
	// Copy the file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	
	return filePath, nil
}

// GetFileURL returns the file URL for local storage
func (p *LocalStorageProvider) GetFileURL(filename string) (string, error) {
	filePath := filepath.Join(p.basePath, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", filename)
	}
	
	// Return a file:// URL for local files
	return fmt.Sprintf("file://%s", filePath), nil
}

// DeleteFile deletes a file from local storage
func (p *LocalStorageProvider) DeleteFile(filename string) error {
	filePath := filepath.Join(p.basePath, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}
	
	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

// GetFile retrieves a file from local storage
func (p *LocalStorageProvider) GetFile(filename string) (io.ReadCloser, error) {
	filePath := filepath.Join(p.basePath, filename)
	
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	
	return file, nil
}

// OSSStorageProvider implements file storage using Alibaba Cloud OSS
type OSSStorageProvider struct {
	config *config.Config
}

// NewOSSStorageProvider creates a new OSS storage provider
func NewOSSStorageProvider(cfg *config.Config) (*OSSStorageProvider, error) {
	// Validate required OSS configuration
	if cfg.OSSAccessKeyID == "" || cfg.OSSAccessKeySecret == "" || cfg.OSSBucketName == "" {
		return nil, fmt.Errorf("OSS configuration incomplete: missing required credentials")
	}
	
	return &OSSStorageProvider{
		config: cfg,
	}, nil
}

// UploadFile uploads a file to OSS
func (p *OSSStorageProvider) UploadFile(file io.Reader, filename string, contentType string) (string, error) {
	// TODO: Implement OSS file upload
	// This would require the Alibaba Cloud OSS Go SDK
	// For now, return a placeholder implementation
	return "", fmt.Errorf("OSS upload implementation not yet implemented")
}

// GetFileURL returns the file URL for OSS
func (p *OSSStorageProvider) GetFileURL(filename string) (string, error) {
	// TODO: Implement OSS file URL generation
	// This would require the Alibaba Cloud OSS Go SDK
	// For now, return a placeholder implementation
	return "", fmt.Errorf("OSS URL generation implementation not yet implemented")
}

// DeleteFile deletes a file from OSS
func (p *OSSStorageProvider) DeleteFile(filename string) error {
	// TODO: Implement OSS file deletion
	// This would require the Alibaba Cloud OSS Go SDK
	// For now, return a placeholder implementation
	return fmt.Errorf("OSS delete implementation not yet implemented")
}

// GetFile retrieves a file from OSS
func (p *OSSStorageProvider) GetFile(filename string) (io.ReadCloser, error) {
	// TODO: Implement OSS file retrieval
	// This would require the Alibaba Cloud OSS Go SDK
	// For now, return a placeholder implementation
	return nil, fmt.Errorf("OSS get file implementation not yet implemented")
}