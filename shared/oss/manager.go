package oss

import (
	"github.com/i-Things/things/shared/conf"
)

func newOssManager(setting conf.OssConf) (sm Handle, err error) {
	OssType := setting.OssType
	switch OssType {
	case "aliyun":
		sm, err = newAliYunOss(conf.AliYunConf{OssConf: setting})
	case "minio":
		sm, err = newMinio(conf.MinioConf{OssConf: setting})
	}
	return sm, err
}
