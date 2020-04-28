package storage

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// AzureBlobOpts options for azure blob storage
type AzureBlobOpts struct {
	AccountName string
	AccountKey  string
	Container   string
}

// AzureContainer is used to create containers and upload files to Azure blob storage
type AzureContainer struct {
	Container *azblob.ContainerURL
}

// NewAzureBlob creates a new AzureBlob
func NewAzureBlob(azureOpts *AzureBlobOpts) (*AzureContainer, error) {
	if azureOpts == nil {
		return nil, errors.New("missing azure options")
	}

	credential, err := azblob.NewSharedKeyCredential(azureOpts.AccountName, azureOpts.AccountKey)
	if err != nil {
		panic(err)
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", azureOpts.AccountName, azureOpts.Container))
	if err != nil {
		return nil, err
	}
	container := azblob.NewContainerURL(*URL, pipeline)

	ctx := context.Background()
	_, err = container.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	// err is ignored, as the container can already exists

	u := &AzureContainer{
		Container: &container,
	}

	return u, nil
}

// Upload uploads a file to azure blob storage
func (u *AzureContainer) Upload(ldif []byte, log []byte) error {
	err := u.uploadFile(LdifFileName, ldif)
	if err != nil {
		return err
	}
	err = u.uploadFile(LogFileName, log)
	if err != nil {
		return err
	}

	return nil
}

func (u *AzureContainer) uploadFile(name string, content []byte) error {
	filePath := fmt.Sprintf("/tmp/%s", name)
	err := ioutil.WriteFile(filePath, content, 0700)
	if err != nil {
		return err
	}

	blobURL := azblob.ContainerURL(*u.Container).NewBlockBlobURL(name)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{})

	return err
}
