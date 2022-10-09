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
下面两种方式二选一
命令执行路径: ithings\src\syssvr\
```shell script
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/things_sys" -table="*" -dir ./internal/repo/mysql
goctl model mysql ddl -src="../../deploy/conf/mysql/sql/syssvr.sql"  -dir ./internal/repo/mysql 

```

## rpc文件编译方法
命令执行路径: ithings\src\syssvr\
```shell script
goctl rpc protoc  proto/sys.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --style=goZero -m
```

# 设备管理模块-dmsvr
##  rpc文件编译
```shell
protoc proto/* --go_out=. --go-grpc_out=.
goctl rpc protoc  proto/dm.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m
```

## model文件编译
下面两种方式二选一
```shell
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/things_dm" -table="*" -dir ./internal/repo/mysql 
goctl model mysql ddl -src="../../deploy/conf/mysql/sql/dmsvr.sql"  -dir ./internal/repo/mysql 

```

# 设备交互模块-disvr

```shell
goctl rpc protoc  proto/di.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --style=goZero -m

```

# 设备数据交互模块-ddsvr

```shell
goctl api go -api http/dd.api  -dir ./ --style=goZero
```