package oss

import (
	"context"
	"fmt"
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

func InitBuckets(ctx context.Context, ser OSSer) error {
	if exists, err := ser.BucketExists(ctx, BucketFirmware); err != nil {
		return err
	} else if exists {
		return nil
	}
	err := ser.MakeBucket(ctx, BucketFirmware, minio.MakeBucketOptions{})
	return err
}

func GetUploadUrl(conf conf.MinioConf) string {
	return fmt.Sprintf("%s/open/oss/upload", conf.GateWayHost)
}
