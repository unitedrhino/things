# sdklog topic

- ### 设备 主动查询 服务端 设置的日志等级

上行 topic: $log/up/operation/${productID}/${deviceName}

```json
{
    "method": "get_status",
    "timestamp": 1654844328734,
    "clientToken": "xxxxxx"
}
```

下行 topic: $log/down/operation/${productID}/${deviceName}

```json
{
    "method": "get_status_reply",
    "clientToken": "xxxxxxx",
    "timestamp": 1654844328734,
    "code": 0,
    "status": "成功",
    "data": {
        "logLevel": 2
    }
}
```

- ### 日志直传

上行 topic: $log/up/report/${productID}/${deviceName}

> - 日志级别 logLevel: 2)错误 3)告警 4)信息 5)调试, 不传默认为5
> - 日志内容 content：在后台 将以文本形式 直接展示
> - 日志时间戳 timestamp

```json
{
    "method": "report_info",
    "timestamp": 1654844328734,
    "clientToken": "xxxxxx",
    "params": [
        {
            "timestamp": 1654844328734,
            "logLevel": 5,
            "content": "long string,不要超过500k"
        },
        {
            "timestamp": 1654844328734,
            "logLevel": 4,
            "content": "long string"
        },
        {
            "timestamp": 1654844328734,
            "content": "sdsdadasfafasdf sdfasd sadfasdf sdfsdfs sdf4asdfsdf"
        }
    ]
}
```

下行 topic: $log/down/report/${productID}/${deviceName}

```json
{
    "method": "report_reply",
    "clientToken": "xxxxxx",
    "timestamp": 1656553866096,
    "code": 0,
    "status": "成功"
}
```

- ### 服务端 主动推送 修改日志等级

下行 topic: $log/down/update/${productID}/${deviceName}

```json
{
    "method": "get_status_reply",
    "clientToken": "xxxxxx",
    "code": 0,
    "status": "成功",
    "data": {
        "logLevel": 1
    }
}
```