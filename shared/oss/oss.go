package oss

import (
	"context"
	"io"

	"github.com/i-Things/things/shared/oss/common"
)

type Handle interface {
	SignedPutUrl(ctx context.Context, objectName string, expiredSec int64, opKv common.OptionKv) (string, error)
	SignedGetUrl(ctx context.Context, objectName string, expiredSec int64, opKv common.OptionKv) (string, error)
	Delete(ctx context.Context, objectName string, opKv common.OptionKv) error
	Upload(ctx context.Context, objectName string, reader io.Reader, opKv common.OptionKv) (string, error)
	GetObjectInfo(ctx context.Context, objectName string) (*common.StorageObjectInfo, error)
	PrivateBucket() Handle
	PublicBucket() Handle
	TemporaryBucket() Handle
	CopyFromTempBucket(tempPath, dstPath string) (string, error)
	GetUrl(path string) (string, error)
	//List(ctx context.Context)
}
