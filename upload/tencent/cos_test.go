package tencent

import (
	"context"
	"fmt"
	"testing"
)

func InitConf() Config {
	return Config{
		AccessKeyID:     "xxxxxxxxxxx",
		SecretAccessKey: "xxxxxxxxxxxxxx",
		Bucket:          "xxxxxxxxxxxx",
		Region:          "xxxxxxx",
		Folder:          "",
		ShowUrl:         "https://xxxxxxx",
	}
}

func TestPresignPutObject(t *testing.T) {
	conf := InitConf()
	ctx := context.Background()
	resp, err := PresignPutObject(ctx, conf, "abc.png", 300)
	if err != nil {
		t.Errorf("PresignPutObject error %v", err)
	}
	fmt.Println(resp)
}
