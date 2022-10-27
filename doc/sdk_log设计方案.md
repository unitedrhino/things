##sdklog topic

- ### 设备主动查询服务端设置的日志等级
上行topic: $log/up/operation/${productID}/${deviceName}
```json
{
  "method": "get_status",
  "timestamp": 1654844328734,
  "clientToken": "xxxxxx"
}
```
下行topic: $log/down/operation/${productID}/${deviceName}
```json
{
  "method": "get_status_reply",
  "clientToken": "xxxxxxx",
  "timestamp": 1654844328734,
  "code": 0,
  "status": "成功",
  "data": {
    "log_level": 2
  }
}
```
- ### 日志直传
上行topic: $log/up/report/${productID}/${deviceName}

> 日志级别log_level: 2)错误 3)告警 4)信息 5)调试  ,不传默认为5
> 日志内容content：在后台将以文本形式直接展示

```json
{                     
    "method": "report_info",    
    "timestamp": 1654844328734,
    "clientToken": "xxxxxx",   
    "params": [
      {
        "log_level":5,
        "content":"long string,不要超过500k"
      },
      {
        "log_level":4,
        "content":"long string"
      },
      {
        "content":"sdsdadasfafasdf sdfasd sadfasdf sdfsdfs sdf4asdfsdf"
      }
    ]
}
```
下行topic: $log/down/report/${productID}/${deviceName}
```json
{
    "method": "report_reply",
    "clientToken": "xxxxxx",
    "timestamp": 1656553866096,
    "code": 0,
    "status": "成功"
}
```
- ### 服务端主动推送修改日志等级
下行topic: $log/down/update/${productID}/${deviceName}
```json
{
  "method": "get_status_reply",
  "clientToken": "xxxxxx",
  "code": 0,
  "status": "成功",
  "data": {
    "log_level":1
  }
}
```