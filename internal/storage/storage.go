package storage

import (
	"io"
	"mime/multipart"
)

type StorageProvider interface {
	Upload(file *multipart.FileHeader, path string) (string, error)
	Delete(path string) error
	Get(path string) (io.Reader, error)
}

var provider StorageProvider

func SetStorageProvider(p StorageProvider) {
	provider = p
}

func GetStorageProvider() StorageProvider {
	return provider
}
