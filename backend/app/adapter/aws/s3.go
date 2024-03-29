package aws

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/gold-kou/ToeBeans/backend/app"
)

func UploadObject(bucket, filename string, file *os.File) (*s3manager.UploadOutput, error) {
	// TODO デバッグコードなので検証後に削除する
	log.Println(bucket)
	sess := session.Must(session.NewSession(generateS3Config()))
	uploader := s3manager.NewUploader(sess)
	return uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
}

// func UploadObject(bucket, filename string, file []byte) (*s3.PutObjectOutput, error) {
// 	sess := session.Must(session.NewSession(generateS3Config()))
// 	svc := s3.New(sess)
// 	return svc.PutObject(&s3.PutObjectInput{
// 		Body:   bytes.NewReader(file),
// 		Bucket: aws.String(bucket),
// 		Key:    aws.String(filename),
// 		ACL:    aws.String(s3.BucketCannedACLPublicRead),
// 	})
// }

func DeleteObject(bucket, filename string) (err error) {
	sess := session.Must(session.NewSession(generateS3Config()))
	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return
	}
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	return
}

func generateS3Config() *aws.Config {
	// use minio in local or test
	if app.IsLocal() || app.IsTest() {
		creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
		return &aws.Config{
			Credentials: creds,
			Region:      aws.String(region),
			// TODO コンテナ間でもlocalhostにできないか
			Endpoint:         aws.String("http://minio:9000"),
			S3ForcePathStyle: aws.Bool(true),
		}
	}

	// use S3 in prd
	// no need access key setting and secret key setting because of IAM Roll
	// no need region setting because of AWS_DEFAULT_REGION
	return &aws.Config{}
}
