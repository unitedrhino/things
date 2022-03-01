# api文件编译方法
```shell script
goctl api go -api webapi.api  -dir ./
```
# 用户管理模块
## 数据库文件生成
```shell script
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/pet" -table="*" -dir ./model -c
```
## rpc文件编译方法
```shell script
goctl rpc proto -src user.proto  -dir ./
```

# 设备管理模块
##  rpc文件编译
```shell
goctl rpc proto -src dm.proto  -dir ./
```

## model文件编译

```shell
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/dm" -table="*" -dir ./internal/repo/mysql -c 
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/dm" -table="device_log" -dir ./internal/repo/mysql
```

# 设备交互模块

```shell
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/dc" -table="*" -dir ./model -c  
goctl rpc proto -src dc.proto  -dir ./
```
