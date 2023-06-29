package main

import (
	"context"
	"fmt"
	"os"

	strategy "karma8_s3_hometask/src/server_choosing_strategy"
	"karma8_s3_hometask/src/services"
	"karma8_s3_hometask/src/storage"
	"karma8_s3_hometask/src/utils"
)

type Config struct {
	Servers int
}

func NewConfig() *Config {
	// get from external resource
	return &Config{Servers: 9}
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}

type dbRecord struct {
	fileRecord services.FileRecord
	storageIds [utils.FileParts]string
}

func uploadAndPrint(ctx context.Context, inputFile string, fs *services.FileService) dbRecord {
	res, fileRecord, err := fs.UploadFile(ctx, inputFile)
	noErr(err)
	record := dbRecord{*fileRecord, [6]string{}}
	for i, r := range res {
		fmt.Printf("%d bytes of '%s' at storageId: %s | path: %s\n",
			r.Size, inputFile, r.StorageId, r.Path)
		record.storageIds[i] = r.StorageId
	}
	println()
	return record
}

func downloadAndPrint(ctx context.Context, dbRecord *dbRecord, fileService *services.FileService) {
	fmt.Printf("\n===== Content of '%s' is: =====\n", dbRecord.fileRecord.FileName)
	size, err := fileService.Download(ctx, dbRecord.fileRecord, os.Stdout)
	noErr(err)
	if size != dbRecord.fileRecord.FileSize {
		panic("The stored FileSize != downloaded FileSize")
	}
	println("\n=====")
}

func main() {
	ctx := context.Background()
	cfg := NewConfig()
	println("servers:", cfg.Servers)

	// set up
	// ======
	strategy := strategy.NewEqualUploadStrategy(strategy.ServerCount(cfg.Servers))

	storages := []storage.Storage{
		storage.NewFileStorage("1", "out/1"),
		storage.NewFileStorage("2", "out/2"),
		storage.NewFileStorage("3", "out/3"),
		storage.NewFileStorage("4", "out/4"),
		storage.NewFileStorage("5", "out/5"),
		storage.NewFileStorage("6", "out/6"),
		storage.NewFileStorage("7", "out/7"),
		storage.NewFileStorage("8", "out/8"),
		storage.NewFileStorage("9", "out/9"),
	}
	fileService := services.NewUploader(strategy, storages)
	// ======

	dbRecords := make([]dbRecord, 0, 2)

	// upload in.txt over 9 servers
	dbRecords = append(dbRecords, uploadAndPrint(ctx, "example_inputs/in.txt", fileService))

	// dynamically add a new server
	fileService.Storages = append(fileService.Storages, storage.NewFileStorage("10", "out/10"))
	fileService.Strategy.SetServerCount(fileService.Strategy.GetServerCount() + 1)
	println("A new server is added with the id:", fileService.Storages[9].GetId(), "\n")

	// upload in2.txt over 10 servers
	dbRecords = append(dbRecords, uploadAndPrint(ctx, "example_inputs/in2.txt", fileService))

	// download & print: in.txt
	downloadAndPrint(ctx, &dbRecords[0], fileService)
	// download & print: in2.txt
	downloadAndPrint(ctx, &dbRecords[1], fileService)

	println("\n(to remove file parts, execute: 'rm -r ./out')")
}
