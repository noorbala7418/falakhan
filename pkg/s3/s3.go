package s3

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/noorbala7418/falakhan/internal/model"
	"github.com/sirupsen/logrus"
)

const S3ACLPolicy = "private"

// ObjectS3Store gets file name and file body and pushes to object storage. then it returns filename again.
func ObjectS3Store(s3credentials model.S3Credential, filename string, filebody io.ReadCloser) (string, error) {
	// TODO: refactor it.
	defer filebody.Close()
	startUpload := time.Now()

	sess, err := session.NewSession(
		&aws.Config{
			Endpoint: aws.String(s3credentials.Endpoint),
			Region:   aws.String(s3credentials.Region),
			Credentials: credentials.NewStaticCredentials(
				s3credentials.AccessKey,
				s3credentials.SecretKey,
				"",
			),
		})

	if err != nil {
		logrus.Error("function ObjectS3Store. Error in create s3 session. err: ", err)
		return "", fmt.Errorf("function ObjectS3Store. Error in create s3 session. err: %w", err)
	}

	uploader := s3manager.NewUploader(sess)

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3credentials.Bucket),
		ACL:    aws.String(S3ACLPolicy),
		Key:    aws.String(filename),
		Body:   filebody,
	})

	if err != nil {
		logrus.Error("function ObjectS3Store. Error in upload object. err: ", err)
		return "", fmt.Errorf("function ObjectS3Store. Error in upload object. err: %w", err)
	}

	totalUploadTime := time.Now().Sub(startUpload)
	logrus.Info("function ObjectS3Store. totalUploadTime: ", totalUploadTime.Seconds(), " seconds")
	logrus.Info("function ObjectS3Store. Uploaded on location: " + up.Location)

	return up.Location, nil
}
