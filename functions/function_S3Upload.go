package functions

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"findsalon-backend/integrations"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func UploadToS3(file *multipart.FileHeader, folder string) (string, error) {
	if integrations.AwsS3Bucket == "" {
		return "", fmt.Errorf("S3 bucket not configured")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(integrations.AwsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			integrations.AwsAccessKeyId,
			integrations.AwsSecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	key := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	client := s3.NewFromConfig(cfg)
	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(integrations.AwsS3Bucket),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(contentType(ext)),
	})
	if err != nil {
		return "", err
	}

	return integrations.AwsS3BaseUrl + "/" + key, nil
}

func contentType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
