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
	bucket := strings.TrimSpace(integrations.AwsS3Bucket)
	region := strings.TrimSpace(integrations.AwsRegion)
	accessKeyID := strings.TrimSpace(integrations.AwsAccessKeyId)
	secretAccessKey := strings.TrimSpace(integrations.AwsSecretAccessKey)

	if bucket == "" {
		return "", fmt.Errorf("S3 bucket not configured")
	}
	if region == "" {
		return "", fmt.Errorf("AWS region not configured")
	}
	if accessKeyID == "" || secretAccessKey == "" {
		return "", fmt.Errorf("AWS credentials not configured")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
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
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(contentType(ext)),
	})
	if err != nil {
		return "", err
	}

	return s3ObjectURL(bucket, region, key), nil
}

func s3ObjectURL(bucket, region, key string) string {
	baseURL := strings.TrimSpace(integrations.AwsS3BaseUrl)
	if baseURL == "" || !isHTTPURL(baseURL) {
		baseURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucket, region)
	}
	return strings.TrimRight(baseURL, "/") + "/" + key
}

func isHTTPURL(value string) bool {
	return strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "http://")
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
