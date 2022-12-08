package download_file_manager

import (
	"context"
	"log"
	"os"
	"wckd1/tg-youtube-podcasts-bot/internal/file_manager"
	"wckd1/tg-youtube-podcasts-bot/internal/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var _ file_manager.FileManager = (*DownloadFileManager)(nil)

type S3Uploader struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Uploader(ctx context.Context, awsCfg util.AWSConfig) Uploader {
	cfg := aws.Config{
		Region: awsCfg.Region,
		Credentials: credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     awsCfg.Key,
				SecretAccessKey: awsCfg.Secret,
			},
		},
	}

	return &S3Uploader{
		client: s3.NewFromConfig(cfg),
		bucket: awsCfg.Bucket,
		region: awsCfg.Region,
	}
}

func (u S3Uploader) Upload(ctx context.Context, file localFile) (url string, err error) {
	// Prepare file
	stat, err := os.Stat(file.path)
	if err != nil {
		log.Printf("[ERROR] failed to stat audio file: %w", err)
		return
	}

	fbody, err := os.Open(file.path)
	if err != nil {
		log.Printf("[ERROR] failed to open audio file: %w", err)
		return
	}
	defer fbody.Close()
	id := uuid.New().String()

	//upload to the s3 bucket
	_, err = u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(u.bucket),
		Key:           aws.String(file.path),
		Body:          fbody,
		ContentLength: stat.Size(),
	})
	if err != nil {
		log.Printf("[ERROR] failed to upload file: %w", err)
		return
	}

	url = "https://" + u.bucket + "." + "s3-" + u.region + ".amazonaws.com/" + id
	return
}
