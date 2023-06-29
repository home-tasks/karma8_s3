package services

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/google/uuid"
	strategy "karma8_s3_hometask/src/server_choosing_strategy"
	"karma8_s3_hometask/src/storage"
	"karma8_s3_hometask/src/utils"
)

type FileService struct {
	Strategy strategy.ServerChoosingStrategy
	Storages []storage.Storage
}

type FileRecord struct {
	FileName string
	FileSize int64
	Storages [utils.FileParts]storage.Storage
}

func NewUploader(strategy strategy.ServerChoosingStrategy, storages []storage.Storage) *FileService {
	return &FileService{strategy, storages}
}

func (fs *FileService) upload(ctx context.Context, f *os.File) (*[utils.FileParts]*storage.FilePartInfo, *FileRecord, error) {
	stats, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}
	partSizes := utils.SplitSize(stats.Size())

	randomServers, err := fs.Strategy.Choose(ctx, partSizes)
	if err != nil {
		return nil, nil, err
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, nil, err
	}
	fileName := uuid.String() + "_" + stats.Name()

	offset := int64(0)
	multErrs := make([]error, 0)
	fileInfoResults := [utils.FileParts]*storage.FilePartInfo{}
	storages := [utils.FileParts]storage.Storage{}
	var wg sync.WaitGroup
	wg.Add(len(partSizes))

	for i, size := range partSizes {
		r := io.NewSectionReader(f, offset, size)
		go func(i int) {
			// rase condition can't happen because {i, elements of randomServers and storages} are unique
			storage := fs.Storages[randomServers[i]]
			storages[i] = storage
			info, err := storage.Upload(ctx, r, fileName)
			if err != nil {
				multErrs = append(multErrs, err)
			}
			fileInfoResults[i] = info
			wg.Done()
		}(i)
		offset += size
	}
	wg.Wait()

	fs.Strategy.CommitChoice(ctx, randomServers)

	if len(multErrs) != 0 {
		// TODO: with GO 1.20, change to errors.Join()
		return nil, nil, multErrs[0]
	}
	return &fileInfoResults, &FileRecord{fileName, stats.Size(), storages}, nil
}

func (fs *FileService) UploadFile(ctx context.Context, name string) (*[utils.FileParts]*storage.FilePartInfo, *FileRecord, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	fileInfos, fileRecord, err := fs.upload(ctx, f)
	if err != nil {
		return nil, nil, err
	}
	return fileInfos, fileRecord, f.Close()
}

func (fs *FileService) Download(ctx context.Context, fileRecord FileRecord, w io.Writer) (int64, error) {
	readers := [utils.FileParts]*io.PipeReader{}
	multErrs := make([]error, 0)
	for i := range fileRecord.Storages {
		var writer *io.PipeWriter
		readers[i], writer = io.Pipe()
		go func(i int) {
			err := fileRecord.Storages[i].Download(ctx, writer, fileRecord.FileName)
			if err != nil {
				multErrs = append(multErrs, err)
			}
			if err := writer.Close(); err != nil {
				multErrs = append(multErrs, err)
			}
		}(i)
	}

	if len(multErrs) != 0 {
		for _, r := range readers {
			if r != nil {
				r.Close()
			}
		}
		// TODO: with GO 1.20, change this and below to errors.Join()
		return 0, multErrs[0]
	}

	size := int64(0)
	for _, reader := range readers {
		n, err := io.Copy(w, reader)
		if err != nil {
			multErrs = append(multErrs, err)
		}
		size += n
	}

	if len(multErrs) != 0 {
		return 0, multErrs[0]
	}
	return size, nil
}
