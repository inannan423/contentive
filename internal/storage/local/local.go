package local

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LocalStorage struct {
	BasePath string
	BaseURL  string
}

// NewLocalStorage creates a new LocalStorage instance
func NewLocalStorage(basePath string, baseURL string) *LocalStorage {
	return &LocalStorage{
		BasePath: basePath,
		BaseURL:  baseURL,
	}
}

// Upload uploads a file to the local storage
func (s *LocalStorage) Upload(file *multipart.FileHeader, path string) (string, error) {
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	year := time.Now().Format("2006")
	month := time.Now().Format("01")

	fullPath := filepath.Join(s.BasePath, path, year, month, fileName)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	relativePath := strings.TrimPrefix(fullPath, s.BasePath)
	return fmt.Sprintf("%s%s", s.BaseURL, relativePath), nil
}

// Delete deletes a file from the local storage
func (s *LocalStorage) Delete(path string) error {
	relativePath := strings.TrimPrefix(path, s.BaseURL)
	fullPath := filepath.Join(s.BasePath, relativePath)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %v", err)
	}

	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// Get retrieves a file from the local storage
func (s *LocalStorage) Get(path string) (io.Reader, error) {
	relativePath := strings.TrimPrefix(path, s.BaseURL)
	fullPath := filepath.Join(s.BasePath, relativePath)

	// Open the file
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	return file, nil
}
