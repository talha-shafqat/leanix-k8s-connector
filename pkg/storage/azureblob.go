package storage

import (
	"errors"
	"fmt"

	azure "github.com/Azure/azure-sdk-for-go/storage"
)

// AzureBlobOpts options for azure blob storage
type AzureBlobOpts struct {
	AccountName string
	AccountKey  string
	Container   string
}

// AzureBlob is used to upload files to azure blob storage
type AzureBlob struct {
	Container *azure.Container
}

// NewAzureBlob creates a new AzureBlob
func NewAzureBlob(azureOpts *AzureBlobOpts) (*AzureBlob, error) {
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

	u := &AzureBlob{
		Container: containerRef,
	}

	return u, nil
}

// Upload uploads a file to azure blob storage
func (u AzureBlob) Upload(content []byte) error {
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
