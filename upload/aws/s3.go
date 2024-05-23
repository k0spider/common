package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"net/url"
	"path"
	"strings"
	"time"
)

type Config struct {
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	Bucket          string `yaml:"bucket"`
	Region          string `yaml:"region"`
	Folder          string `yaml:"folder"`
}

type Credential struct {
	AccessKeyID     string
	SecretAccessKey string
}

func (c *Credential) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID: c.AccessKeyID, SecretAccessKey: c.SecretAccessKey,
	}, nil
}

func PresignPutObject(conf Config, fileName string, lifetimeSecs int64) (*AwsS3Response, error) {
	presignClient := s3.NewPresignClient(s3.NewFromConfig(aws.Config{
		Region: conf.Region,
		Credentials: &Credential{
			AccessKeyID:     conf.AccessKeyID,
			SecretAccessKey: conf.SecretAccessKey,
		},
	}))
	putObjectInput := s3.PutObjectInput{
		Bucket: aws.String(conf.Bucket),
		Key:    aws.String(path.Join(conf.Folder, fileName)),
		ACL:    types.ObjectCannedACLPublicRead,
	}
	result, err := presignClient.PresignPutObject(context.TODO(), &putObjectInput, func(po *s3.PresignOptions) {
		po.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		return nil, err
	}
	ret := &AwsS3Response{Url: result.URL}
	urm, _ := url.ParseRequestURI(result.URL)
	ret.Name = fmt.Sprintf("%s://%s%s", urm.Scheme, urm.Host, urm.Path)
	ret.Header = make(map[string]string)
	for key, val := range result.SignedHeader {
		if len(val) == 0 {
			continue
		}
		ret.Header[key] = strings.Join(val, ",")
	}
	return ret, nil
}
