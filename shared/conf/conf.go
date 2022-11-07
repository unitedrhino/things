package conf

//反向代理配置
type ProxyConf struct {
	FrontDir         string            `json:",default=./dist"`                       //前端文件路径
	FrontDefaultPage string            `json:",default=front/iThingsCore/index.html"` //前端默认文件地址
	UrlProxy         map[string]string `json:",optional"`                             //反向http代理配置
}
