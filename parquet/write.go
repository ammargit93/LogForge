package parquet

import (
	"fmt"
	"log"
	"log-engine/db"
	"log-engine/models"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

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

var MinioClient = connectToMinIO()

func WriteToParquet(buffer []models.LogEntry, servicename, username string) {
	filename := "dummy.parquet"
	uname, _ := GetCreds()
	if username != uname {
		return
	}
	fw, err := local.NewLocalFileWriter(filename)
	if err != nil {
		fmt.Println("WriteToParquet Error", err)
		return
	}

	pw, err := writer.NewParquetWriter(fw, new(models.LogEntry), 4)
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

	if err := pw.WriteStop(); err != nil {
		fmt.Println("WriteStop error:", err)
	}
	fw.Close()

	db.UploadToMinIO(MinioClient, servicename, username)
}
