package main

import (
	"context"

	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type LogEntry struct {
	Timestamp string `json:"timestamp" binding:"required" parquet:"name=timestamp, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Level     string `json:"level" binding:"required" parquet:"name=level, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Message   string `json:"message" binding:"required" parquet:"name=message, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
	Service   string `json:"service" binding:"required" parquet:"name=service, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Host      string `json:"host" binding:"required" parquet:"name=host, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
}

var BufferQueue []LogEntry
var mu sync.Mutex

const N = 5

func connectToMinIO() *minio.Client {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return minioClient
}

var minioClient = connectToMinIO()

func main() {

	router := gin.Default()
	router.POST("/logs", handleLog)
	log.Println("🚀 Server running on http://localhost:8080")
	router.Run(":8080")
}

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

func WriteToParquet(buffer []LogEntry) {
	filename := "dummy.parquet"
	fw, err := local.NewLocalFileWriter(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	pw, err := writer.NewParquetWriter(fw, new(LogEntry), 4)
	if err != nil {
		fmt.Println(err)
		return
	}
	pw.RowGroupSize = 128 * 1024 * 1024 // 128MB
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	for _, log := range buffer {
		if err := pw.Write(log); err != nil {
			fmt.Println(err)
			return
		}
	}
	UploadToMinIO(minioClient, "auth-service", "ITsMe")
	if err := pw.WriteStop(); err != nil {
		fmt.Println("WriteStop error:", err)
	}
	fw.Close()

}

func handleLog(c *gin.Context) {
	var logEntry LogEntry
	if err := c.ShouldBindJSON(&logEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	BufferQueue = append(BufferQueue, logEntry)
	if len(BufferQueue) >= N {
		bufferCopy := make([]LogEntry, len(BufferQueue))
		copy(bufferCopy, BufferQueue)
		go WriteToParquet(bufferCopy)
		BufferQueue = []LogEntry{}
	}
	mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
