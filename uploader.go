package main

// LDIFUploader exposes a common interface for all upload mechanisms
type LDIFUploader interface {
	Upload() error
}

// NoOpUploader implements the LDIFUploader interface without doing any upload
type NoOpUploader struct{}

// NewNoOpUploader creates a new NewNoOpUploader
func NewNoOpUploader() LDIFUploader {
	return NoOpUploader{}
}

// Upload does nothing (no op)
func (u NoOpUploader) Upload() error {
	return nil
}
