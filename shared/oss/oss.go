package oss

import (
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/minio/minio-go/v7"
)

type OSSer = *minio.Client

//bucket 列表
const (
	BucketFirmware = "firmware"
)

func NewOss(conf conf.OSSConf) (OSSer, error) {
	if conf.OSS != "minio" {
		return nil, fmt.Errorf("oss just support minio")
	}
	return NewMinio(conf.Minio)
}
