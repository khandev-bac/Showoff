package pkg

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	Client     *s3.Client
	Uploader   *manager.Uploader
	BucketName string
	Region     string
}

func NewUploader(bucketName string, region string) (*S3Uploader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	return &S3Uploader{
		Client:     client,
		Uploader:   uploader,
		BucketName: bucketName,
		Region:     region,
	}, nil
}
func (s *S3Uploader) UploadFile(file multipart.File, fileSize int64, filename string) (string, error) {
	key := fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), filename)
	upParams := &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   file,
		ACL:    "public-read",
	}
	_, err := s.Uploader.Upload(context.TODO(), upParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to s3: %w", err)
	}
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.BucketName, s.Region, key)
	return url, nil
}
