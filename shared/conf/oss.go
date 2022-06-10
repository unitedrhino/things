package conf

type OSSConf struct {
	OSS   string    `json:",default=minio,options=minio|ceph|aliyun|tenxunyun"`
	Minio MinioConf `json:",optional"`
}

type MinioConf struct {
	Endpoint        string `json:",default=127.0.0.1:9000"`
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool `json:",optional"`
}
