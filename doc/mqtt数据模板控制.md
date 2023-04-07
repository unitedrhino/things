# mqtt 发送数据 服务器端对数据的控制

## 属性上报

### 以以下上报信息为例

```json
{
    "clientToken": "42251d61-b355-4f21-a04c8f8865c",
    "method": "report",
    "params": {
        "GPS_Info": {
            "longtitude": 0,
            "latitude": 0
        }
    }
}
```

如果 GPS_Info 修改为 错误的值，会报没有找到正确的属性

```json
{
    "method": "report_reply",
    "clientToken": "42251d61-b355-4f21-a04c8f8865c",
    "code": 406,
    "status": "Not find PropertySchema"
}
```

如果是 GPS_Info 里的 参数错误，则会替换为默认值0 并返回成功

```json
{
    "method": "report_reply",
    "clientToken": "42251d61-b355-4f21-a04c8f8865c",
    "code": 0,
    "status": "success"
}
```

### 记录 属性控制 请求及回复

请求

```json
{
    "method": "control",
    "clientToken": "clientToken-8446fc3b-6dd6-4028-a68a-cd485ea1dd26",
    "params": {
        "GPS_ExtInfo": {
            "satellites": 0,
            "collect_time": 1624377600,
            "latitude": 0,
            "longtitude": 0,
            "altitude": 0,
            "gps_speed": 0,
            "direction": 0,
            "location_state": 0
        },
        "ipaddr": "",
        "rssi": "",
        "imageUrl": "",
        "shuxing": "",
        "biashijigou": {
            "fwe": 0,
            "ase": 0
        },
        "serfa": 1,
        "GPS_Info": {
            "longtitude": 0,
            "latitude": 0
        }
    }
}
```

回复

```json
{
    "method": "control_reply",
    "clientToken": "123",
    "code": 0,
    "status": "some message where error"
}
```

## 事件上报

如果 上报事件中的 参数有错误

```json
{
    "method": "event_post",
    "clientToken": "99ab17d8-9475-46d0-82e0-1b92fe03940a",
    "version": "1.0",
    "eventId": "fesf",
    "timestamp": 1623322323740,
    "params": {
        "GPS_Info": {
            "longtitude": 0,
            "latitude": 0
        }
    }
}
```

则服务器会返回 以下信息

```json
{
    "method": "event_reply",
    "clientToken": "99ab17d8-9475-46d0-82e0-1b92fe03940a",
    "code": 406,
    "status": "Param count not match",
    "data": {}
}
```

如果是 正确的参数 如下:

```json
{
    "method": "event_post",
    "clientToken": "99ab17d8-9475-46d0-82e0-1b92fe03940a",
    "version": "1.0",
    "eventId": "fesf",
    "timestamp": 1623322323740,
    "params": {
        "se": 0,
        "dfa": 200
    }
}
```

返回参数为:

```json
{
    "method": "event_reply",
    "clientToken": "99ab17d8-9475-46d0-82e0-1b92fe03940a",
    "code": 0,
    "status": "",
    "data": {}
}
```

## 其他错误

如果是 bool类型, 参数应该为0和1, 如果错误 则返回

```json
{
    "method": "event_reply",
    "clientToken": "99ab17d8-9475-46d0-82e0-1b92fe03940a",
    "code": 406,
    "status": "se value true out of range:[0,1]",
    "data": {}
}
```

如果是 超过了最大值 则返回:
> 如果是 float型 步距小于设定值，日志会按照 真实的记录 不会四舍五入
```json
{
    "method": "event_reply",
    "clientToken": "99ab17d8-9475-46d0-82e0-1b92fe03940a",
    "code": 406,
    "status": "dfa value 9 out of range:[100,238]",
    "data": {}
}
```


## 操作行为

操作行为 有超时控制 估算为3秒-5秒之间, 在此 先记录请求报文,之后进行测试

```json
{
    "method": "action",
    "clientToken": "146761676::bf34c502-2470-44b7-bd54-5ab9d137c0f8",
    "actionId": "biaoshifu",
    "timestamp": 1623324528,
    "params": {
        "asdfwe": "fewf",
        "ee": 1
    }
}
```

# 备注

参考腾讯云文档: https://cloud.tencent.com/document/product/1081/34916

# 总结

腾讯云 会对每个请求的报文 进行格式解析及控制, 通讯日志 会记录原始的报文, 成功的请求 解析后的结果 会保存到 对应的日志中