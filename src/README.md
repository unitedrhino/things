[TOC]

# 本地安装goctl【非官方，请使用i4de/go-zero】

1. 本地将go-zero项目克隆下来：  `git clone git@github.com:i4de/go-zero.git`
2. 到目录go-zero\tools\goctl 下 执行命令： `go install`
3. 后续执行下面的各种goctl命令即可

# 环境初始化

protoc/protoc-gen-go/protoc-gen-grpc-go 依赖可以通过 `goctl env check -i -f` 一键安装

# api网关接口代理模块-apisvr

```shell
#命令执行路径: ithings\src\apisvr\
goctl api go -api http/api.api  -dir ./  --style=goZero
```

# 系统管理模块-syssvr

- 数据库文件生成

```shell
#命令执行路径: ithings\src\syssvr\
goctl model mysql ddl -src="../../deploy/conf/mysql/sql/syssvr.sql"  -dir ./internal/repo/mysql -i updatedTime,deletedTime,createdTime
```

- rpc文件编译方法

```shell
#命令执行路径: ithings\src\syssvr\
goctl rpc protoc  proto/sys.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --style=goZero -m
```

# 设备管理模块-dmsvr

- rpc文件编译方法

```shell
#命令执行路径: ithings\src\dmsvr\
#protoc proto/* --go_out=. --go-grpc_out=.
goctl rpc protoc  proto/dm.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m
```

- model文件编译

```shell
#命令执行路径: ithings\src\dmsvr\
goctl model mysql ddl -src="../../deploy/conf/mysql/sql/dmsvr.sql"  -dir ./internal/repo/mysql -i updatedTime,deletedTime,createdTime
```

# 设备交互模块-disvr

```shell
#命令执行路径: ithings\src\disvr\
goctl rpc protoc  proto/di.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --style=goZero -m
```

# 设备数据处理模块-ddsvr

```shell
#命令执行路径: ithings\src\ddsvr\
goctl api go -api http/dd.api  -dir ./ --style=goZero
```

# 协议规则引擎模块-rulesvr

- rpc文件编译

```shell
#命令执行路径: ithings\src\rulesvr\
#protoc  proto/* --go_out=. --go-grpc_out=.
goctl rpc protoc  proto/rule.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m
```

- model文件编译

```shell
#命令执行路径: ithings\src\rulesvr\
goctl model mysql ddl -src="../../deploy/conf/mysql/sql/rulesvr.sql"  -dir ./internal/repo/mysql -i updatedTime,deletedTime,createdTime
```
