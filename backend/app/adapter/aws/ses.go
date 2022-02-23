package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gold-kou/ToeBeans/backend/app"
)

func SendEmail(to, title, body string) error {
	sess, err := session.NewSession(generateSESConfig())
	if err != nil {
		return err
	}
	svc := ses.New(sess)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(title),
			},
		},
		Source: aws.String(emailFrom),
	}
	_, err = svc.SendEmail(input)
	if err != nil {
		return err
	}
	return nil
}

func generateSESConfig() *aws.Config {
	if app.IsLocal() || app.IsTest() {
		return &aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		}
	}
	return &aws.Config{}
}
