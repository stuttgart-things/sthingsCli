/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioEnvVars = []string{"MINIO_ACCESS_KEY", "MINIO_SECRET_KEY", "MINIO_ADDR", "MINIO_SECURE"}
)

func CreateMinioClient() (bool, *minio.Client) {

	if VerifyEnvVars(minioEnvVars) {

		accessKey, _ := os.LookupEnv("ACCESS_KEY")
		secretAccessKey, _ := os.LookupEnv("SECRET_ACCESS_KEY")
		minioServer, _ := os.LookupEnv("MINIO_ADDR")
		secure, _ := os.LookupEnv("MINIO_SECURE")

		s3Client, err := minio.New(minioServer, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
			Secure: sthingsBase.ConvertStringToBoolean(secure),
		})

		if err != nil {
			fmt.Println(err)
		}

		return true, s3Client

	} else {
		return false, nil
	}
}

func GetObjectsFromMinioBucket(minioClient *minio.Client, bucket string) []string {

	s3_objects := make([]string, 0, 4)

	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    "",
	}

	for object := range minioClient.ListObjects(context.Background(), bucket, opts) {
		if object.Err != nil {
			fmt.Println(object.Err)
		}
		s3_objects = append(s3_objects, object.Key)
	}

	if len(s3_objects) == 0 {
		fmt.Println("NOTHING FOUND!")
		os.Exit(0)
	}

	return s3_objects
}

func DownloadObjectFromMinioBucket(minioClient *minio.Client, bucket, objectname, destinationName string) bool {

	reader, err := minioClient.GetObject(context.Background(), bucket, objectname, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	localFile, err := os.Create(destinationName)
	if err != nil {
		log.Fatalln(err)
	}
	defer localFile.Close()

	stat, err := reader.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := io.CopyN(localFile, reader, stat.Size); err != nil {
		log.Fatalln(err)
		return false
	} else {
		return true
	}

}

func UploadObjectToMinioBucket(minioClient *minio.Client, bucket, sourcePath, objectName string) (bool, int64) {

	ctx := context.Background()
	contentType := "application/octet-stream"

	// UPLOAD THE FILE WITH FPUTOBJECT
	info, err := minioClient.FPutObject(ctx, bucket, objectName, sourcePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
		return false, 0
	}

	return true, info.Size
}
