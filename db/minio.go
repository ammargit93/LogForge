package db

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func UploadToMinIO(minioClient *minio.Client, service string, user string) {
	bucketName := "logs"
	location := "us-east-1"

	err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", bucketName)
	}
	currYear := time.Now().Year()
	currMonth := time.Now().Month()
	currDay := time.Now().Day()
	logId := uuid.New()
	logFileName := service + "_" + strconv.Itoa(currYear) + "_" + currMonth.String() + "_" + strconv.Itoa(currDay) + "_" + logId.String() + ".parquet"

	path := []string{user, service, strconv.Itoa(currYear), currMonth.String(), strconv.Itoa(currDay), logFileName}
	objectName := strings.Join(path, "/")
	filePath := "dummy.parquet"
	contentType := "application/octet-stream"

	_, err = minioClient.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s\n", objectName)
}
