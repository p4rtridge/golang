package aws

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type AWSClient interface {
	SaveImage(context.Context, *[]byte) (string, error)
	DeleteImage(context.Context, string) error
}

const (
	BUCKET         = "order-service"
	IMAGE_BASE_URL = "https://storage.iloveyour.dad/images"
)

type awsClient struct {
	client *s3.Client
}

func NewAWSClient(client *s3.Client) AWSClient {
	return &awsClient{
		client: client,
	}
}

func (c *awsClient) SaveImage(ctx context.Context, data *[]byte) (string, error) {
	mimeType := http.DetectContentType(*data)

	imageFilename := fmt.Sprintf("%s.%s", uuid.New().String(), strings.Split(mimeType, "/")[1])
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(BUCKET),
		Key:         aws.String(fmt.Sprintf("images/%s", imageFilename)),
		Body:        bytes.NewReader(*data),
		ContentType: aws.String(mimeType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", IMAGE_BASE_URL, imageFilename), nil
}

func (c *awsClient) DeleteImage(ctx context.Context, imageUrl string) error {
	splittedImage := strings.Split(imageUrl, "/")
	imageName := splittedImage[len(splittedImage)-1]

	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(fmt.Sprintf("images/%s", imageName)),
	})
	if err != nil {
		return err
	}

	return nil
}
