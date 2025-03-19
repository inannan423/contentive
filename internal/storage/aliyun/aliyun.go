package aliyun

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/google/uuid"
)

// AliyunOSSStorage implements the StorageProvider interface for Aliyun OSS
type AliyunOSSStorage struct {
	Client     *oss.Client
	BucketName string
	BaseURL    string
}

// NewAliyunOSSStorage creates a new Aliyun OSS storage provider
func NewAliyunOSSStorage(endpoint, accessKeyID, accessKeySecret, bucketName, baseURL string) (*AliyunOSSStorage, error) {
	// Extract region from endpoint (e.g., "oss-cn-beijing.aliyuncs.com" -> "cn-beijing")
	region := strings.TrimPrefix(strings.Split(endpoint, ".")[0], "oss-")

	// Create OSS client configuration
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")).
		WithRegion(region).    // Set the region
		WithEndpoint(endpoint) // Set the endpoint

	// Create OSS client
	client := oss.NewClient(cfg)

	// Check if bucket exists by trying to list objects (minimal operation)
	_, err := client.ListObjectsV2(context.Background(), &oss.ListObjectsV2Request{
		Bucket:  oss.Ptr(bucketName),
		MaxKeys: 1, // Just request one object to check if bucket exists
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %v", err)
	}

	return &AliyunOSSStorage{
		Client:     client,
		BucketName: bucketName,
		BaseURL:    baseURL,
	}, nil
}

// Upload uploads a file to Aliyun OSS
func (s *AliyunOSSStorage) Upload(file *multipart.FileHeader, directory string) (string, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Generate a unique filename
	ext := filepath.Ext(file.Filename)
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	filename := fmt.Sprintf("%s/%s/%s%s", directory, year, month, uuid.New().String()+ext)

	// Create upload request
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(s.BucketName),
		Key:    oss.Ptr(filename),
		Body:   src,
	}

	// Upload the file
	_, err = s.Client.PutObject(context.Background(), request)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Return the URL
	return fmt.Sprintf("%s/%s", s.BaseURL, filename), nil
}

// Delete deletes a file from Aliyun OSS
func (s *AliyunOSSStorage) Delete(path string) error {
	// Extract the object key from the URL
	objectKey := path
	if s.BaseURL != "" && len(path) > len(s.BaseURL)+1 {
		objectKey = path[len(s.BaseURL)+1:]
	}

	// Create delete request
	request := &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(s.BucketName),
		Key:    oss.Ptr(objectKey),
	}

	// Delete the object
	_, err := s.Client.DeleteObject(context.Background(), request)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// Get retrieves a file from Aliyun OSS
func (s *AliyunOSSStorage) Get(path string) (io.Reader, error) {
	// Extract the object key from the URL
	objectKey := path
	if s.BaseURL != "" && len(path) > len(s.BaseURL)+1 {
		objectKey = path[len(s.BaseURL)+1:]
	}

	// Create get request
	request := &oss.GetObjectRequest{
		Bucket: oss.Ptr(s.BucketName),
		Key:    oss.Ptr(objectKey),
	}

	// Get the object
	result, err := s.Client.GetObject(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	return result.Body, nil
}
