package common

import (
	"bytes"
	"flag"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func PutToS3(bucket, key string, target []byte) error {
	sess, err := session.NewSession(generateConfig())
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(target),
	})

	return err
}

func generateConfig() *aws.Config {
	if os.Getenv("RUNSERVER") == "LOCAL" || flag.Lookup("test.v") != nil {
		// use minio in local
		creds := credentials.NewStaticCredentials(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), "")
		return &aws.Config{
			Credentials:      creds,
			Region:           aws.String(os.Getenv("MINIO_REGION")),
			Endpoint:         aws.String("http://minio:9000"),
			S3ForcePathStyle: aws.Bool(true),
		}
	}

	// use S3 in prd
	// no need access key setting and secret key setting because of IAM Roll
	// no need region setting because of AWS_DEFAULT_REGION
	return &aws.Config{}
}
