package db

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

func DownloadParquetFilesFromMinIO(minioClient *minio.Client, prefix string) error {
	bucketName := "logs"
	ctx := context.Background()

	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	err := os.MkdirAll("../uploads", 0755)
	if err != nil {
		return fmt.Errorf("failed to create uploads folder: %v", err)
	}

	for object := range objectCh {
		if object.Err != nil {
			log.Printf("❌ Failed to list: %v\n", object.Err)
			continue
		}
		if strings.HasSuffix(object.Key, ".parquet") {
			log.Printf("⬇️ Downloading %s", object.Key)
			err := downloadToUploads(minioClient, bucketName, object.Key)
			if err != nil {
				log.Printf("❌ Download failed: %v", err)
			}
		}
	}
	return nil
}

func downloadToUploads(minioClient *minio.Client, bucketName, objectKey string) error {
	ctx := context.Background()
	object, err := minioClient.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("get object failed: %v", err)
	}
	defer object.Close()

	localPath := filepath.Join("uploads", filepath.Base(objectKey))
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("create file failed: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, object)
	if err != nil {
		return fmt.Errorf("copy to local failed: %v", err)
	}

	return nil
}
