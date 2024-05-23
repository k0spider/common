package tencent

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Config struct {
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	Bucket          string `yaml:"bucket"`
	Region          string `yaml:"region"`
	Folder          string `yaml:"folder"`
	ShowUrl         string `yaml:"showUrl"`
}

func PresignPutObject(ctx context.Context, conf Config, fileName string, lifetimeSecs int64) (*CosResponse, error) {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", conf.Bucket, conf.Region))
	b := &cos.BaseURL{BucketURL: u}
	ossClient := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.AccessKeyID,
			SecretKey: conf.SecretAccessKey,
			Expire:    time.Second * 10,
		},
	})
	preSignedUrl, err := ossClient.Object.GetPresignedURL2(ctx, http.MethodPut, path.Join(conf.Folder, fileName), time.Duration(lifetimeSecs*int64(time.Second)), nil)
	if err != nil {
		return nil, err
	}
	return &CosResponse{
		Url:       fmt.Sprintf("%s%s", conf.ShowUrl, preSignedUrl.Path),
		SignedUrl: preSignedUrl.String(),
	}, nil
}
