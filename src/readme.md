# 本地安装goctl
1.本地将go-zero项目克隆下来：  `git clone git@github.com:i4de/go-zero.git`
2.到目录go-zero\tools\goctl 下 执行命令： `go install`  
3.后续执行下面的各种goctl命令即可

# 环境初始化

protoc/protoc-gen-go/protoc-gen-grpc-go 依赖可以通过
`goctl env check -i -f` 一键安装

# api文件编译方法
命令执行路径: ithings\src\apisvr\
```shell script
goctl api go -api http/api.api  -dir ./  --style=goZero
```

# 系统管理模块-syssvr

## 数据库文件生成

命令执行路径: ithings\src\syssvr\
```shell script
goctl model mysql ddl -src="../../deploy/conf/mysql/sql/syssvr.sql"  -dir ./internal/repo/mysql -icreatedTime,updatedTime,deletedTime

```

## rpc文件编译方法
命令执行路径: ithings\src\syssvr\
```shell script
goctl rpc protoc  proto/sys.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --style=goZero -m
```

# 设备管理模块-dmsvr
##  rpc文件编译
命令执行路径: ithings\src\dmsvr\
```shell
protoc proto/* --go_out=. --go-grpc_out=.
goctl rpc protoc  proto/dm.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m
```

## model文件编译
命令执行路径: ithings\src\dmsvr\
```shell
goctl model mysql ddl -src="../../deploy/conf/mysql/sql/dmsvr.sql"  -dir ./internal/repo/mysql -icreatedTime,updatedTime,deletedTime
```

# 设备交互模块-disvr
命令执行路径: ithings\src\disvr\
```shell
goctl rpc protoc  proto/di.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --style=goZero -m

```

# 设备数据交互模块-ddsvr
命令执行路径: ithings\src\ddsvr\
```shell
goctl api go -api http/dd.api  -dir ./ --style=goZero
```