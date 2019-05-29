package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// LocalFileOpts options for local file storage
type LocalFileOpts struct {
	Path string
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
