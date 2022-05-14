# 环境初始化

protoc/protoc-gen-go/protoc-gen-grpc-go 依赖可以通过
`goctl env check -i -f` 一键安装

# api文件编译方法

```shell script
goctl api go -api zeroapi/api.api  -dir ./
```

# 用户管理模块-usersvr

## 数据库文件生成

```shell script
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/pet" -table="*" -dir ./internal/repo/mysql -c
```

## rpc文件编译方法
```shell script
goctl rpc protoc  proto/user.proto --go_out=./ --go-grpc_out=./ --zrpc_out=.
```

# 设备管理模块-dmsvr
##  rpc文件编译
```shell
goctl rpc protoc  proto/dm.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./
```

## model文件编译

```shell
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/dm" -table="*" -dir ./internal/repo/mysql -c 
```

# 设备交互模块-dcsvr

```shell
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/dc" -table="*" -dir ./internal/repo/mysql -c  
goctl rpc protoc  proto/dc.proto --go_out=./ --go-grpc_out=./ --zrpc_out=.

```

# 设备数据交互模块-ddsvr

```shell
goctl rpc protoc  proto/dd.proto --go_out=./ --go-grpc_out=./ --zrpc_out=.
```