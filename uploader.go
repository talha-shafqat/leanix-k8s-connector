package main

import (
	"errors"
	"fmt"

	azure "github.com/Azure/azure-sdk-for-go/storage"
)

// LDIFUploader exposes a common interface for all upload mechanisms
type LDIFUploader interface {
	Upload(content []byte) error
}

// AzureStorageOpts options for azure blob storage
type AzureStorageOpts struct {
	AccountName string
	AccountKey  string
	Container   string
}

// NewLDIFUploader create a new uploader for the given storage engine
func NewLDIFUploader(storageEngine string, azureOpts *AzureStorageOpts) (LDIFUploader, error) {
	switch storageEngine {
	case "azure":
		if azureOpts == nil {
			return nil, errors.New("azure storage options must be set when using azure as storage target")
		}
		return NewAzureUploader(azureOpts)
	case "local":
		return &NoOpUploader{}, nil
	}
	return nil, fmt.Errorf("Unsupported storage engine %s", storageEngine)
}

// NoOpUploader implements the LDIFUploader interface without doing any upload
type NoOpUploader struct{}

// Upload does nothing (no op)
func (u *NoOpUploader) Upload(content []byte) error {
	return nil
}

// AzureUploader implements to LDIFUploader interface to upload files to azure blob storage
type AzureUploader struct {
	Container *azure.Container
}

// NewAzureUploader creates a new AzureUploader
func NewAzureUploader(azureOpts *AzureStorageOpts) (*AzureUploader, error) {
	if azureOpts == nil {
		return nil, errors.New("missing azure options")
	}

	client, err := azure.NewBasicClient(azureOpts.AccountName, azureOpts.AccountKey)
	if err != nil {
		panic(err)
	}

	blobClient := client.GetBlobService()
	containerRef := blobClient.GetContainerReference(azureOpts.Container)

	u := &AzureUploader{
		Container: containerRef,
	}

	return u, nil
}

// Upload uploads a file to azure blob storage
func (u AzureUploader) Upload(content []byte) error {
	if u.Container == nil {
		return errors.New("unable to obtain a container reference")
	}

	blobReference := u.Container.GetBlobReference("ldif.json")

	err := blobReference.PutAppendBlob(nil)
	if err == nil {
		err = blobReference.AppendBlock(content, nil)
	}

	return err
}
