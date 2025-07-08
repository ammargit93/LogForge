package db

import (
	"context"
	"fmt"
	"log"
	"path"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func UploadToMinIO(minioClient *minio.Client, service string, user string) error {
	bucketName := "logs"
	location := "us-east-1"

	// Create bucket if not exists
	err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", bucketName)
	}

	// Create folder path: username/service/yyyy-mm-dd/<uuid>.parquet
	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := fmt.Sprintf("%02d", now.Month())
	day := fmt.Sprintf("%02d", now.Day())
	logID := uuid.New().String()
	logFileName := logID + ".parquet"

	date := year + "-" + month + "-" + day
	// Use path.Join() instead of filepath.Join() to ensure forward slashes
	objectPath := path.Join(user, service, date, logFileName)

	// Upload
	filePath := "dummy.parquet"
	contentType := "application/octet-stream"

	_, err = minioClient.FPutObject(context.Background(), bucketName, objectPath, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload object: %v", err)
	}

	log.Printf("✅ Uploaded to MinIO: %s\n", objectPath)
	return nil
}
