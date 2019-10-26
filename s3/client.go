package s3

import (
	"io"
	"log"

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
func (c Client) Put(key string, body io.ReadSeeker) error {
	sess, err := c.session()
	if err != nil {
		return errors.Wrap(err, "error getting session")
	}

	s3Client := s3.New(sess)
	result, err := s3Client.PutObject(&s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
		Body:   body,
	})

	if err != nil {
		return errors.Wrap(err, "error uploading file to s3 using uploader")
	}
	log.Printf("File uploaded to S3 | Result: %#v\n", *result)
	return nil
}

func (c Client) session() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
	})
}
