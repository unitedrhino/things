package conf

// 文件反向代理
type FileProxyConf struct {
	FrontDir string `json:",default=./dist"` //前端文件路径
	CoreDir  string `json:",default=app/core/"`
}

// http反向代理
type StaticProxyConf struct {
	Router       string //原路由
	Dest         string //目标路由
	DeletePrefix bool   `json:",optional"` //是否删除原路由路径
}

type ProxyConf struct {
	FileProxy   FileProxyConf      `json:""`
	StaticProxy []*StaticProxyConf `json:",optional"`
}
