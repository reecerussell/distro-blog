package storage

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Service is a high level interface used for uploading
// and downloading data, to and from AWS S3 buckets.
type Service struct {
	bucketName string
	sess *session.Session
}

// New returns a new service instance for the given bucket name.
func New(bucketName string) (*Service, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &Service {
		bucketName: bucketName,
		sess: sess,
	}, nil
}

// Set uploads data to an S3 bucket with the service's bucket name, using the given key.
func (s *Service) Set(key string, data []byte) error {
	uploader := s3manager.NewUploader(s.sess)
	buf := bytes.NewBuffer(data)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucketName),
		Key: aws.String(key),
		Body: buf,
	})
	if err != nil {
		return fmt.Errorf("failed to upload to bucket '%s' with key '%s': %v", s.bucketName, key, err)
	}

	return nil
}

// basicWriter is a very basic implementation of io.WriterAt, to give
// the ability to download objects from S3 to a buffer/[]byte, instead
// of a file.
type basicWriter []byte

func (w *basicWriter) WriteAt(p []byte, off int64) (n int, err error) {
	buf := make([]byte, len(*w)+len(p))
	copy(buf[:], *w)
	copy(buf[off:], p)
	*w = buf

	return len(p), nil
}

// Get attempts to retrieve an item from the service's S3 bucket, with the given key.
func (s *Service) Get(key string) ([]byte, error) {
	downloader := s3manager.NewDownloader(s.sess)
	var buf basicWriter

	_, err := downloader.Download(&buf, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key: aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download item '%s' from bucket '%s': %v", key, s.bucketName, err)
	}

	return buf, nil
}

// Delete attempts to delete a specific object from S3.
func (s *Service) Delete(key string) error {
	svc := s3.New(s.sess)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Key: aws.String(key),
		Bucket: aws.String(s.bucketName),
	})
	if err != nil {
		return fmt.Errorf("unable to delete object from s3: %v", err)
	}

	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{
		Key: aws.String(key),
		Bucket: aws.String(s.bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from s3: %v", err)
	}

	return nil
}