package s3

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

// Client is the AWS S3 client object
type Client struct {
	Region string
	Bucket string
}

// Put uploads an object to S3
func (c Client) Put(folder, name string, body io.ReadSeeker) (string, error) {
	sess, err := c.session()
	if err != nil {
		return "", errors.Wrap(err, "error getting session")
	}

	key := fmt.Sprintf("%s/%s", folder, name)
	url := c.url(key)

	s3Client := s3.New(sess)
	result, err := s3Client.PutObject(&s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(url),
		Body:   body,
	})

	if err != nil {
		return "", errors.Wrap(err, "error uploading file to s3 using uploader")
	}
	log.Printf("File uploaded to S3 | Result: %#v\n", *result)
	return "", nil
}

func (c Client) session() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
	})
}

func (c Client) url(key string) string {
	bucket := strings.Trim(c.Bucket, "/")
	region := strings.Trim(c.Region, "/")
	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucket, region, key)
}
