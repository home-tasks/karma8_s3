package storage

import (
	"context"
	"io"
	"os"
	"path"
)

type fileStorage struct {
	id      string
	baseDir string
}

func (fr *fileStorage) Download(ctx context.Context, w io.Writer, fileName string) error {
	f, err := os.Open(path.Join(fr.baseDir, fileName))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}

func (fr *fileStorage) Upload(ctx context.Context, r io.Reader, fileName string) (*FilePartInfo, error) {
	filePath := path.Join(fr.baseDir, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := f.Close(); err == nil {
			err = e
		}
	}()
	writtenBytes, err := io.Copy(f, r)
	if err != nil {
		return nil, err
	}

	return &FilePartInfo{
		StorageId: fr.GetId(),
		Size:      writtenBytes,
		Path:      filePath,
	}, nil
}

func (fr *fileStorage) GetId() string {
	return fr.id
}

func NewFileStorage(id string, baseDir string) Storage {
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		panic(err)
	}
	return &fileStorage{id, baseDir}
}
