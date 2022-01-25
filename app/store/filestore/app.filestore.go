package filestore

import (
	"bytes"
	"fmt"
	"orijinplus/app/models"
	"orijinplus/utils/faulterr"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FileStore struct {
	Session    *session.Session
	S3         *s3.S3
	BucketName string
}

func NewFilestore(sess *session.Session) *FileStore {
	return &FileStore{
		Session:    sess,
		S3:         s3.New(sess),
		BucketName: "aws-s3-bucket",
	}
}

func (fs *FileStore) ListBuckets() (*s3.ListBucketsOutput, *faulterr.FaultErr) {
	result, err := fs.S3.ListBuckets(nil)
	if err != nil {
		return nil, faulterr.NewInternalServerError(err.Error())
	}

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	return result, nil
}

func (fs *FileStore) ListBucketItems() *faulterr.FaultErr {
	resp, err := fs.S3.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(fs.BucketName)})
	if err != nil {
		return faulterr.NewInternalServerError(err.Error())
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	return nil
}

func (fs *FileStore) UploadFile(filename string, fileObj []byte) (*models.File, *faulterr.FaultErr) {
	timestamp := time.Now().UTC().Unix()
	key := strconv.Itoa(int(timestamp)) + "_" + filename

	uploader := s3manager.NewUploader(fs.Session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &fs.BucketName,
		Key:    &key,
		Body:   bytes.NewReader(fileObj),
	})
	if err != nil {
		return nil, faulterr.NewInternalServerError("file upload")
	}

	obj := &models.File{
		Name: key,
		URL:  fs.generateURL(key),
	}

	return obj, nil
}

func (fs *FileStore) DownloadFile(key string) (*models.File, *faulterr.FaultErr) {
	downloader := s3manager.NewDownloader(fs.Session)

	file, err := os.Create(key)
	if err != nil {
		return nil, faulterr.NewInternalServerError(err.Error())
	}

	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String("orijinplus-dev"),
			Key:    &key,
		})
	if err != nil {
		return nil, faulterr.NewInternalServerError(err.Error())
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	obj := &models.File{
		Name: key,
		URL:  fs.generateURL(key),
	}

	return obj, nil
}

func (fs *FileStore) generateURL(key string) string {
	return fs.BucketName + ".s3.amazonaws.com/" + key
}
