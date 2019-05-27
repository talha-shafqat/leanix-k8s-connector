package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	azure "github.com/Azure/azure-sdk-for-go/storage"
)

const (
	azureblobStorage string = "azureblob"
	fileStorage      string = "file"
)

// StorageBackend exposes a common interface for all storage mechanisms
type StorageBackend interface {
	Upload(content []byte) error
}

// AzureStorageOpts options for azure blob storage
type AzureStorageOpts struct {
	AccountName string
	AccountKey  string
	Container   string
}

// LocalFileOpts options for local file storage
type LocalFileOpts struct {
	Path string
}

// NewStorageBackend create a new storage backend for the given storage backend type
func NewStorageBackend(storageBackend string, azureOpts *AzureStorageOpts, localFileOpts *LocalFileOpts) (StorageBackend, error) {
	switch storageBackend {
	case azureblobStorage:
		if azureOpts == nil {
			return nil, errors.New("azure storage options must be set when using azure as storage target")
		}
		return NewAzureStorage(azureOpts)
	case fileStorage:
		return NewLocalFile(localFileOpts.Path)
	}
	return nil, fmt.Errorf("Unsupported storage backend type %s", storageBackend)
}

// LocalFile writes the content to disk
type LocalFile struct {
	Path string
}

// NewLocalFile create a LocalFile StorageBackend
func NewLocalFile(path string) (*LocalFile, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsDir() {
		return nil, fmt.Errorf("path %s is not a directory", path)
	}
	lf := &LocalFile{
		Path: path,
	}
	return lf, nil
}

// Upload persists the content in a local file
func (u *LocalFile) Upload(content []byte) error {
	location := path.Join(u.Path, "ldif.json")
	return ioutil.WriteFile(location, content, 0644)
}

// AzureStorage is used to upload files to azure blob storage
type AzureStorage struct {
	Container *azure.Container
}

// NewAzureStorage creates a new AzureStorage
func NewAzureStorage(azureOpts *AzureStorageOpts) (*AzureStorage, error) {
	if azureOpts == nil {
		return nil, errors.New("missing azure options")
	}

	client, err := azure.NewBasicClient(azureOpts.AccountName, azureOpts.AccountKey)
	if err != nil {
		panic(err)
	}

	blobClient := client.GetBlobService()
	containerRef := blobClient.GetContainerReference(azureOpts.Container)
	containerExists, err := containerRef.Exists()
	if err != nil {
		return nil, err
	}
	if !containerExists {
		return nil, fmt.Errorf("azure blob storage container %s does not exist", azureOpts.Container)
	}

	u := &AzureStorage{
		Container: containerRef,
	}

	return u, nil
}

// Upload uploads a file to azure blob storage
func (u AzureStorage) Upload(content []byte) error {
	if u.Container == nil {
		return errors.New("unable to obtain a container reference")
	}

	blobReference := u.Container.GetBlobReference("ldif.json")

	// create the blob if it does not exist
	err := blobReference.PutAppendBlob(nil)
	if err == nil {
		err = blobReference.AppendBlock(content, nil)
	}

	return err
}
