package oss

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinio(conf conf.MinioConf) (OSSer, error) {
	// 初使化 minio client对象。
	minioClient, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}
