package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/vinicius73/thecollector/pkg/support"
)

type RemoteCredentials struct {
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
}

func (r RemoteCredentials) DigitalOcean(endpoint string) (*session.Session, error) {
	if endpoint == "" {
		endpoint = support.GetEnv("DO_SESSION_ENDOINT", "https://nyc3.digitaloceanspaces.com")
	}

	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(r.Key, r.Secret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String("us-east-1"),
	})
}
