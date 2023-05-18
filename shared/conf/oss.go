package conf

import (
	"fmt"

	oss "github.com/aliyun/aliyun-oss-go-sdk/aliYunOss"
)

type OssConf struct {
	OssType             string `json:",options=minio|aliyun"`
	AccessKeyID         string
	AccessKeySecret     string
	PublicBucketName    string
	TemporaryBucketName string //临时桶,30分钟有效期
	PrivateBucketName   string
	Location            string
	UseSSL              bool   `json:",optional"`
	CustomHost          string `json:",default=127.0.0.1:7777"` //api访问地址
	ConnectTimeout      int64
	ReadWriteTimeout    int64
}

// minio本地存储配置
type MinioConf struct {
	OssConf
}

func (m MinioConf) GetEndPoint() string {
	return m.Location
}

// 阿里云oss配置
type AliYunConf struct {
	OssConf
}

func (a AliYunConf) GenClientOption() []oss.ClientOption {
	options := make([]oss.ClientOption, 0)
	options = append(options, oss.Timeout(a.ConnectTimeout, a.ReadWriteTimeout))
	if a.CustomHost != "" {
		options = append(options, oss.UseCname(true))
	}
	return options
}

func (a AliYunConf) GetEndPoint() string {
	scheme := "https"
	if !a.UseSSL {
		scheme = "http"
	}
	if a.CustomHost == "" {
		return fmt.Sprintf("%s://%s.aliyuncs.com", scheme, a.Location)
	}
	return fmt.Sprintf("%s://%s", scheme, a.CustomHost)
}
