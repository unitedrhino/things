[TOC]

# 本地安装goctl【非官方，请使用i-Things/go-zero】

1. 本地将 `go-zero 项目克隆下来：  `git clone git@github.com:i-Things/go-zero.git`
2. 到目录 `go-zero\tools\goctl 下 执行命令： `go install`
3. 后续执行下面的各种goctl命令即可

# 环境初始化

`protoc/protoc-gen-go/protoc-gen-grpc-go` 依赖可以通过下列命令 一键安装

```shell
$ goctl env check --install --verbose --force
```

# 服务新增方案

## rpc服务
```
goctl rpc new rpcsvr  --style=goZero -m
```
## api服务
```
goctl api new viewsvr  --style=goZero 
```

# 库表新增方案

在每个服务的 `internal/repo/relationDB` 目录下有example.go 
1. 借助 `https://sql2gorm.mccode.info/` 生成对应的模型 放到 `internal/repo/relationDB/modle.go` 中
2. 复制 `internal/repo/relationDB/example.go` 到对应目录下,并修改表名
3. 将example.go中的Example替换为表名
4. 定制修改对应函数即可

# api网关接口代理模块-apisvr

```shell
#cd apisvr && goctl api go -api http/api.api  -dir ./  --style=goZero && cd ..
cd apisvr && goctl api go -api http/api.api  -dir ./  --style=goZero -ws && cd ..
```

# view大屏goview模块-viewsvr

```shell
#cd apisvr && goctl api go -api http/api.api  -dir ./  --style=goZero && cd ..
cd viewsvr && goctl api go -api http/view.api  -dir ./  --style=goZero  && cd ..
```

# 系统管理模块-syssvr

- rpc文件编译方法

```shell
cd syssvr && goctl rpc protoc  proto/sys.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --style=goZero -m && cd ..
```


# 设备管理模块-dmsvr

- rpc文件编译方法

```shell
#protoc proto/* --go_out=. --go-grpc_out=.
cd dmsvr && goctl rpc protoc  proto/dm.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m && cd ..
```


# 设备网关模块-dgsvr

```shell
cd dgsvr && goctl rpc protoc  proto/dg.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m && cd ..
```

# 用户设备模块-udsvr

- rpc文件编译

```shell
#protoc  proto/* --go_out=. --go-grpc_out=.
cd udsvr && goctl rpc protoc  proto/ud.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m && cd ..
```


# 视频管理引擎模块-vidsvr

- rpc文件编译

```shell
#protoc  proto/* --go_out=. --go-grpc_out=.
cd vidsvr && goctl rpc protoc  proto/vid.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m && cd ..
#mediaserver http-api
cd vidsvr && goctl api go -api http/vid.api  -dir ./  --style=goZero && cd ..
``` 


# 定时生产者引擎模块-timedschedulersvr

- rpc文件编译

```shell
#protoc  proto/* --go_out=. --go-grpc_out=.
cd timed/timedschedulersvr && goctl rpc protoc  proto/timedscheduler.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m && cd ../..
```

# 定时消费者者引擎模块-timedjobsvr

- rpc文件编译

```shell
#protoc  proto/* --go_out=. --go-grpc_out=.
cd timed/timedjobsvr && goctl rpc protoc  proto/timedjob.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m && cd ../..
```