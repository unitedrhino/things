# 设备属性缓存管理器

这个包提供了设备属性缓存的管理功能，包括首次记录和最后记录的缓存操作。

## 功能特性

- **属性缓存管理**: 管理设备属性的首次记录和最后记录
- **变化检测**: 检测属性值是否发生变化，避免重复记录
- **批量操作**: 支持批量更新和获取属性记录
- **缓存清理**: 提供设备属性缓存的清理功能

## 主要组件

### PropertyCacheManager

属性缓存管理器，提供以下功能：

- `UpdatePropertyCache`: 更新属性缓存
- `GetPropertyFirstRecord`: 获取属性首次记录
- `GetPropertyLastRecord`: 获取属性最后记录
- `ClearPropertyCache`: 清除设备属性缓存
- `CheckIsChange`: 检查属性值是否发生变化
- `GetAllPropertyRecords`: 获取设备所有属性记录

### 工具函数

- `GenRedisPropertyFirstKey`: 生成设备属性首次记录的Redis键
- `GenRedisPropertyLastKey`: 生成设备属性最后记录的Redis键

## 使用方法

### 基本使用

```go
package main

import (
    "context"
    "time"
    
    "gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/cache"
    "gitee.com/unitedrhino/things/share/domain/schema"
    "github.com/zeromicro/go-zero/core/stores/kv"
)

func main() {
    // 创建KV存储实例
    var kvStore kv.Store
    
    // 创建缓存管理器
    cacheManager := cache.NewPropertyCacheManager(kvStore)
    
    // 更新属性缓存
    ctx := context.Background()
    productID := "product_001"
    deviceName := "device_001"
    
    property := &schema.Property{
        Identifier: "temperature",
        Define: schema.PropertyDefine{
            Type: schema.DataTypeFloat,
        },
        RecordMode: schema.RecordModeAll,
    }
    
    data := map[string]any{
        "temperature": 25.5,
    }
    
    err := cacheManager.UpdatePropertyCache(ctx, productID, deviceName, property, data, time.Now().Unix())
    if err != nil {
        // 处理错误
    }
}
```

### 检查属性变化

```go
// 检查属性值是否发生变化
dev := devices.Core{
    ProductID:  productID,
    DeviceName: deviceName,
}

data := msgThing.PropertyData{
    Identifier: "temperature",
    Param:      26.0,
    TimeStamp:  time.Now().Unix(),
}

hasChanged := cacheManager.CheckIsChange(ctx, dev, property, data)
if hasChanged {
    // 属性值发生了变化，需要更新
}
```

### 获取属性记录

```go
// 获取首次记录
firstRecord, err := cacheManager.GetPropertyFirstRecord(ctx, productID, deviceName, "temperature")
if err != nil {
    // 处理错误
}

// 获取最后记录
lastRecord, err := cacheManager.GetPropertyLastRecord(ctx, productID, deviceName, "temperature")
if err != nil {
    // 处理错误
}
```

### 清除缓存

```go
// 清除设备属性缓存
err := cacheManager.ClearPropertyCache(ctx, productID, deviceName)
if err != nil {
    // 处理错误
}
```

## Redis键格式

- 首次记录键: `device:thing:property:first:{productID}:{deviceName}`
- 最后记录键: `device:thing:property:last:{productID}:{deviceName}`

## 记录模式

支持以下记录模式：

- `RecordModeAll`: 记录所有值
- `RecordModeNone`: 不记录任何值
- `RecordModeChange`: 只记录变化的值

## 注意事项

1. 缓存管理器是线程安全的，可以在多个goroutine中并发使用
2. 属性数据会被序列化为JSON格式存储在Redis中
3. 首次记录只在值发生变化时更新
4. 最后记录每次都会更新
5. 建议在生产环境中监控缓存的使用情况和性能

## 测试

运行测试：

```bash
go test ./cache
```

运行测试并查看覆盖率：

```bash
go test -cover ./cache
```

