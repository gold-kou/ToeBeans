package aws

import (
	"bytes"
	"flag"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func PutToS3(bucket, filename string, file []byte) (*s3manager.UploadOutput, error) {
	sess := session.Must(session.NewSession(generateS3Config()))
	uploader := s3manager.NewUploader(sess)
	return uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(file),
	})
}

func generateS3Config() *aws.Config {
	if os.Getenv("APP_ENV") == "development" || flag.Lookup("test.v") != nil {
		// use minio in local or test
		creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), "")
		return &aws.Config{
			Credentials:      creds,
			Region:           aws.String(os.Getenv("AWS_REGION")),
			Endpoint:         aws.String("http://minio:9000"),
			S3ForcePathStyle: aws.Bool(true),
		}
	}

	// use S3 in prd
	// no need access key setting and secret key setting because of IAM Roll
	// no need region setting because of AWS_DEFAULT_REGION
	return &aws.Config{}
}
