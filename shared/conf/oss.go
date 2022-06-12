package conf

type OSSConf struct {
	OSS          string    `json:",default=minio,options=minio|ceph|aliyun|tenxunyun"`
	Minio        MinioConf `json:",optional"`
	AccessSecret string    `json:",default=password"` //jwt 认证秘钥
	AccessExpire int64     `json:",default=600"`      //jwt 过期时间 单位:秒
}

type MinioConf struct {
	Endpoint        string `json:",default=127.0.0.1:9000"`
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool   `json:",optional"`
	GateWayHost     string `json:",default=127.0.0.1:7777"` //api访问地址
}
