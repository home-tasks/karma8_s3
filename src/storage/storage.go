package storage

import (
	"context"
	"io"
)

type Storage interface {
	Download(ctx context.Context, w io.Writer, fileName string) error
	Upload(ctx context.Context, r io.Reader, fileName string) (*FilePartInfo, error)

	GetId() string
}

type FilePartInfo struct {
	StorageId string
	Size      int64
	Path      string
}
