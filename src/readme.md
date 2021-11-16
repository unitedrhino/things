# api文件编译方法
```shell script
goctl api go -api webapi.api  -dir ./
```

# 数据库文件生成
```shell script
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3308)/pet" -table="*" -dir ./model -c
```

# rpc文件编译方法
```shell script
goctl rpc proto -src user.proto  -dir ./
```

# 设备管理模块
##  rpc文件编译
goctl rpc proto -src dm.proto  -dir ./
## model文件编译
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/things" -table="*" -dir ./model -c
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/things" -table="device_log" -dir ./model -c

#接口文档生成

# 设备交互模块
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/dcsvr" -table="*" -dir ./model -c  
goctl rpc proto -src dc.proto  -dir ./