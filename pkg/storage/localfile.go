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

// Upload persists the ldif content and the log file content in a local files
func (u *LocalFile) Upload(ldif []byte, log []byte) error {
	err := u.uploadFile(path.Join(u.Path, LdifFileName), ldif)
	if err != nil {
		return err
	}
	err = u.uploadFile(path.Join(u.Path, LogFileName), log)
	if err != nil {
		return err
	}
	return nil
}

func (u *LocalFile) uploadFile(name string, content []byte) error {
	location := path.Join(u.Path, name)
	return ioutil.WriteFile(location, content, 0644)
}
