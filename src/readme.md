#本地安装goctl
1.本地将go-zero项目克隆下来：  git clone git@github.com:zeromicro/go-zero.git
2.到目录go-zero\tools\goctl 下 执行命令： go install    生成 goctl.exe ，再复制到Go的安装⽬录bin下
3.后续执行下面的各种goctl命令即可

# 环境初始化

protoc/protoc-gen-go/protoc-gen-grpc-go 依赖可以通过
`goctl env check -i -f` 一键安装

# api文件编译方法

```shell script
goctl api go -api http/api.api  -dir ./  --style=goZero
```

# 用户管理模块-usersvr

## 数据库文件生成

```shell script
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/things_user" -table="*" -dir ./internal/repo/mysql
```

## rpc文件编译方法
```shell script
goctl rpc protoc  proto/user.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --style=goZero
```

# 设备管理模块-dmsvr
##  rpc文件编译
```shell
goctl rpc protoc  proto/dm.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero
```

## model文件编译

```shell
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/things_dm" -table="*" -dir ./internal/repo/mysql -c 
```

# 设备交互模块-dcsvr

```shell
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/things_dc" -table="*" -dir ./internal/repo/mysql -c  
goctl rpc protoc  proto/dc.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --style=goZero

```

# 设备数据交互模块-ddsvr

```shell
goctl api go -api http/dd.api  -dir ./ --style=goZero
goctl rpc protoc  proto/dd.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --style=goZero
```