package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	Folder string
}

func (s LocalStorage) UploadFromURL(destPath string, url string) error {
	fullFilePath := filepath.Join(s.Folder, destPath)
	dirPath := filepath.Dir(fullFilePath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	res, err := downloadFile(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(filepath.Join(s.Folder, destPath))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, res.Body)
	return err
}
