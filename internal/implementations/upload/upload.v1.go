package upload

import (
	"bmt_product_service/dto/request"
	"bmt_product_service/global"
	"bmt_product_service/internal/services"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type UploadService struct {
	S3Client *s3.S3
}

const (
	film_image_base_key = "film-images/"
	film_video_base_key = "film-videos/"
	fab_image_base_ket  = "fab-images/"
)

func extractObjectKeyFromURL(objectURL string) (string, error) {
	parsedURL, err := url.Parse(objectURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	objectKey := strings.TrimPrefix(parsedURL.Path, "/")

	return objectKey, nil

}

// DeleteObject implements services.IUpload.
func (us *UploadService) DeleteObject(objectURL string) error {
	objectKey, err := extractObjectKeyFromURL(objectURL)
	if err != nil {
		return fmt.Errorf("an error occur when parsing object URL (%s): %v", objectURL, err)
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(global.Config.ServiceSetting.S3Setting.FilmBucketName),
		Key:    aws.String(objectKey),
	}

	_, err = us.S3Client.DeleteObject(input)
	if err != nil {
		return fmt.Errorf("failed to delete object '%s' from S3: %v", objectKey, err)
	}

	return nil
}

// UploadProductImageToS3 implements services.IS3.
func (us *UploadService) UploadProductImageToS3(message request.UploadImageReq, productType string) error {
	src, err := message.Image.Open()
	if err != nil {
		return fmt.Errorf("can't open file (image): %v", err)
	}
	defer src.Close()

	ext := filepath.Ext(message.Image.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	newFileName := strconv.Itoa(int(message.ProductId)) + "-" + uuid.New().String() + ext
	objectKey := ""

	switch productType {
	case global.FAB_TYPE:
		objectKey = fab_image_base_ket + newFileName
	case global.FILM_TYPE:
		objectKey = film_image_base_key + newFileName
	default:
		return fmt.Errorf("invalid product type; %s", productType)
	}

	_, err = us.S3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(global.Config.ServiceSetting.S3Setting.FilmBucketName),
		Key:         aws.String(objectKey),
		Body:        src,
		ContentType: aws.String(message.Image.Header.Get("Content-Type")),
	})
	if err != nil {
		return fmt.Errorf("upload to S3 (image) failure: %v", err)
	}

	return nil
}

// UploadFilmVideoToS3 implements services.IS3.
func (us *UploadService) UploadFilmVideoToS3(message request.UploadVideoReq) error {
	src, err := message.Video.Open()
	if err != nil {
		return fmt.Errorf("can't open file (video): %v", err)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("can't read file (video): %v", err)
	}

	// Get the file extension (eg ".mp4")
	ext := filepath.Ext(message.Video.Filename)
	contentType := message.Video.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(fileBytes)
	}

	newFileName := strconv.Itoa(int(message.ProductId)) + "-" + uuid.New().String() + ext
	objectKey := film_video_base_key + newFileName

	_, err = us.S3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(global.Config.ServiceSetting.S3Setting.FilmBucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("upload to S3 (video) failure: %v", err)
	}

	return nil
}

func NewUploadService() services.IUpload {
	s3Session, err := session.NewSession(&aws.Config{
		Region: aws.String(global.Config.ServiceSetting.S3Setting.AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			global.Config.ServiceSetting.S3Setting.AwsAccessKeyId,
			global.Config.ServiceSetting.S3Setting.AwsSercetAccessKeyId,
			""),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create AWS session: %v", err))
	}

	return &UploadService{
		S3Client: s3.New(s3Session),
	}
}
