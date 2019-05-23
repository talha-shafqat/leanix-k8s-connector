package main

import "fmt"

// LDIFUploader exposes a common interface for all upload mechanisms
type LDIFUploader interface {
	Upload() error
}

// NewLDIFUploader create a new uploader for the given storage engine
func NewLDIFUploader(storageEngine string) (LDIFUploader, error) {
	switch storageEngine {
	case "azure":
		return NewAzureUploader(), nil
	case "local":
		return NoOpUploader{}, nil
	}
	return nil, fmt.Errorf("Unsupported storage engine %s", storageEngine)
}

// NoOpUploader implements the LDIFUploader interface without doing any upload
type NoOpUploader struct{}

// Upload does nothing (no op)
func (u NoOpUploader) Upload() error {
	return nil
}

// AzureUploader implements to LDIFUploader interface to upload files to azure blob storage
type AzureUploader struct {
}

// NewAzureUploader creates a new AzureUploader
func NewAzureUploader() LDIFUploader {
	return AzureUploader{}
}

// Upload uploads a file to azure blob storage
func (u AzureUploader) Upload() error {
	return nil
}
