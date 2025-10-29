# AI 面向的 Logic 层增删改查规范（Go-Zero）

本规范面向 AI/脚本自动生成 Logic 层业务逻辑代码，确保风格统一、可直接通过 goctl 生成代码并实现完整的业务逻辑。

## 适用范围
- 目录：`internal/logic` 及其子目录
- 文件：`*Logic.go` 业务逻辑文件
- 框架：Go-Zero + GORM
- 数据库：MySQL/PostgreSQL/SQLite 等

## 文件结构规范

### 基础文件结构
每个 Logic 文件应包含以下结构：

```go
package examplelogic

import (
    "context"
    "time"
    
    "gitee.com/unitedrhino/share/ctxs"
    "gitee.com/unitedrhino/share/def"
    "gitee.com/unitedrhino/share/errors"
    "gitee.com/unitedrhino/share/utils"
    "gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
    "gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
    "gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
    "gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
    "github.com/zeromicro/go-zero/core/logx"
)

// 1. Logic 结构体定义
type ExampleCreateLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
    // Repository 依赖
    ExampleDB *relationDB.ExampleRepo
}

// 2. 构造函数
func NewExampleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExampleCreateLogic {
    return &ExampleCreateLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
        ExampleDB: relationDB.NewExampleRepo(ctx),
    }
}

// 3. 业务方法
func (l *ExampleCreateLogic) ExampleCreate(in *dm.ExampleCreateReq) (*dm.Empty, error) {
    // 业务逻辑实现
}
```

## Logic 结构体规范

### 结构体命名规范
- Logic 结构体：`<实体名><操作>Logic`，如 `DeviceInfoCreateLogic`、`ProductInfoIndexLogic`
- 包名：使用小写，如 `devicemanagelogic`、`productmanagelogic`
- 构造函数：`New<实体名><操作>Logic`

### 基础字段规范
每个 Logic 结构体都应包含以下基础字段：

```go
type ExampleCreateLogic struct {
    ctx    context.Context        // 上下文
    svcCtx *svc.ServiceContext   // 服务上下文
    logx.Logger                  // 日志记录器
    
    // Repository 依赖（根据业务需要添加）
    ExampleDB *relationDB.ExampleRepo
    UserDB    *relationDB.UserRepo
    Cache     *relationDB.CacheRepo
}
```

### 构造函数规范
```go
func NewExampleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExampleCreateLogic {
    return &ExampleCreateLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
        ExampleDB: relationDB.NewExampleRepo(ctx),
    }
}
```

## 业务方法规范

### 方法命名规范
- 方法名：与 gRPC 服务方法名一致，如 `DeviceInfoCreate`、`ProductInfoIndex`
- 参数：使用 protobuf 生成的请求类型
- 返回值：使用 protobuf 生成的响应类型或 `*dm.Empty`

### 方法签名规范
```go
// 创建操作
func (l *ExampleCreateLogic) ExampleCreate(in *dm.ExampleCreateReq) (*dm.Empty, error)

// 查询操作
func (l *ExampleIndexLogic) ExampleIndex(in *dm.ExampleIndexReq) (*dm.ExampleIndexResp, error)
func (l *ExampleReadLogic) ExampleRead(in *dm.ExampleReadReq) (*dm.Example, error)

// 更新操作
func (l *ExampleUpdateLogic) ExampleUpdate(in *dm.Example) (*dm.Empty, error)

// 删除操作
func (l *ExampleDeleteLogic) ExampleDelete(in *dm.ExampleDeleteReq) (*dm.Empty, error)

// 批量操作
func (l *ExampleMultiCreateLogic) ExampleMultiCreate(in *dm.ExampleMultiCreateReq) (*dm.ExampleMultiCreateResp, error)
func (l *ExampleMultiUpdateLogic) ExampleMultiUpdate(in *dm.ExampleMultiUpdateReq) (*dm.Empty, error)
func (l *ExampleMultiDeleteLogic) ExampleMultiDelete(in *dm.ExampleMultiDeleteReq) (*dm.Empty, error)
```

## CRUD 操作实现规范

### 创建操作规范
```go
func (l *ExampleCreateLogic) ExampleCreate(in *dm.Example) (*dm.Empty, error) {
    // 1. 权限检查
    if err := ctxs.IsAdmin(l.ctx); err != nil {
        return nil, err
    }
    
    // 2. 参数验证
    if err := l.validateCreateParams(in); err != nil {
        return nil, err
    }
    
    // 3. 业务规则检查
    if err := l.checkBusinessRules(in); err != nil {
        return nil, err
    }
    
    // 4. 数据转换
    data := l.toDBModel(in)
    
    // 5. 数据库操作
    err := l.ExampleDB.Insert(l.ctx, data)
    if err != nil {
        l.Errorf("ExampleCreate.Insert err=%+v", err)
        return nil, err
    }
    
    // 6. 缓存更新
    if err := l.updateCache(data); err != nil {
        l.Error(err)
    }
    
    // 7. 事件发布
    if err := l.publishEvent(data); err != nil {
        l.Error(err)
    }
    
    return &dm.Empty{}, nil
}
```

### 查询操作规范
```go
func (l *ExampleIndexLogic) ExampleIndex(in *dm.ExampleIndexReq) (*dm.ExampleIndexResp, error) {
    l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
    
    // 1. 权限检查
    l.ctx = ctxs.WithDefaultAllProject(l.ctx)
    
    // 2. 构建过滤器
    filter := l.buildFilter(in)
    
    // 3. 统计总数
    total, err := l.ExampleDB.CountByFilter(l.ctx, filter)
    if err != nil {
        return nil, err
    }
    
    // 4. 查询列表
    list, err := l.ExampleDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
    if err != nil {
        return nil, err
    }
    
    // 5. 数据转换
    respList := make([]*dm.Example, 0, len(list))
    for _, item := range list {
        respList = append(respList, logic.ToExampleInfo(l.ctx, l.svcCtx, item))
    }
    
    return &dm.ExampleIndexResp{
        List:  respList,
        Total: total,
    }, nil
}

func (l *ExampleReadLogic) ExampleRead(in *dm.ExampleReadReq) (*dm.Example, error) {
    // 1. 权限检查
    l.ctx = ctxs.WithDefaultAllProject(l.ctx)
    
    // 2. 查询数据
    data, err := l.ExampleDB.FindOneByFilter(l.ctx, relationDB.ExampleFilter{
        ID: in.Id,
    })
    if err != nil {
        if errors.Cmp(err, errors.NotFind) {
            return nil, errors.NotFind.AddDetailf("not find example id=%d", in.Id)
        }
        return nil, err
    }
    
    // 3. 数据转换
    return logic.ToExampleInfo(l.ctx, l.svcCtx, data), nil
}
```

### 更新操作规范
```go
func (l *ExampleUpdateLogic) ExampleUpdate(in *dm.Example) (*dm.Empty, error) {
    // 1. 权限检查
    l.ctx = ctxs.WithDefaultAllProject(l.ctx)
    
    // 2. 查询原数据
    oldData, err := l.ExampleDB.FindOneByFilter(l.ctx, relationDB.ExampleFilter{
        ID: in.Id,
    })
    if err != nil {
        if errors.Cmp(err, errors.NotFind) {
            return nil, errors.NotFind.AddDetailf("not find example id=%d", in.Id)
        }
        return nil, err
    }
    
    // 3. 业务规则检查
    if err := l.checkUpdateRules(oldData, in); err != nil {
        return nil, err
    }
    
    // 4. 数据更新
    if err := l.updateFields(oldData, in); err != nil {
        return nil, err
    }
    
    // 5. 数据库操作
    err = l.ExampleDB.Update(l.ctx, oldData)
    if err != nil {
        l.Errorf("ExampleUpdate.Update err=%+v", err)
        return nil, err
    }
    
    // 6. 缓存更新
    if err := l.updateCache(oldData); err != nil {
        l.Error(err)
    }
    
    // 7. 事件发布
    if err := l.publishEvent(oldData); err != nil {
        l.Error(err)
    }
    
    return &dm.Empty{}, nil
}
```

### 删除操作规范
```go
func (l *ExampleDeleteLogic) ExampleDelete(in *dm.ExampleDeleteReq) (*dm.Empty, error) {
    // 1. 权限检查
    if err := ctxs.IsAdmin(l.ctx); err != nil {
        return nil, err
    }
    
    // 2. 查询数据
    data, err := l.ExampleDB.FindOneByFilter(l.ctx, relationDB.ExampleFilter{
        ID: in.Id,
    })
    if err != nil {
        if errors.Cmp(err, errors.NotFind) {
            return nil, errors.NotFind.AddDetailf("not find example id=%d", in.Id)
        }
        return nil, err
    }
    
    // 3. 业务规则检查
    if err := l.checkDeleteRules(data); err != nil {
        return nil, err
    }
    
    // 4. 数据库操作
    err = l.ExampleDB.Delete(l.ctx, in.Id)
    if err != nil {
        l.Errorf("ExampleDelete.Delete err=%+v", err)
        return nil, err
    }
    
    // 5. 缓存清理
    if err := l.clearCache(data); err != nil {
        l.Error(err)
    }
    
    // 6. 事件发布
    if err := l.publishEvent(data); err != nil {
        l.Error(err)
    }
    
    return &dm.Empty{}, nil
}
```

## 批量操作规范

### 批量创建规范
```go
func (l *ExampleMultiCreateLogic) ExampleMultiCreate(in *dm.ExampleMultiCreateReq) (*dm.ExampleMultiCreateResp, error) {
    // 1. 权限检查
    if err := ctxs.IsAdmin(l.ctx); err != nil {
        return nil, err
    }
    
    // 2. 参数验证
    if err := l.validateMultiCreateParams(in); err != nil {
        return nil, err
    }
    
    // 3. 数据转换
    dataList := make([]*relationDB.Example, 0, len(in.Examples))
    for _, item := range in.Examples {
        dataList = append(dataList, l.toDBModel(item))
    }
    
    // 4. 批量插入
    err := l.ExampleDB.MultiInsert(l.ctx, dataList)
    if err != nil {
        l.Errorf("ExampleMultiCreate.MultiInsert err=%+v", err)
        return nil, err
    }
    
    // 5. 事件发布
    for _, data := range dataList {
        if err := l.publishEvent(data); err != nil {
            l.Error(err)
        }
    }
    
    return &dm.ExampleMultiCreateResp{}, nil
}
```

### 批量更新规范
```go
func (l *ExampleMultiUpdateLogic) ExampleMultiUpdate(in *dm.ExampleMultiUpdateReq) (*dm.Empty, error) {
    // 1. 权限检查
    if err := ctxs.IsAdmin(l.ctx); err != nil {
        return nil, err
    }
    
    // 2. 构建过滤器
    filter := l.buildMultiUpdateFilter(in)
    
    // 3. 查询数据
    dataList, err := l.ExampleDB.FindByFilter(l.ctx, filter, nil)
    if err != nil {
        return nil, err
    }
    
    // 4. 批量更新
    updates := l.buildUpdates(in)
    err = l.ExampleDB.UpdateWithField(l.ctx, filter, updates)
    if err != nil {
        l.Errorf("ExampleMultiUpdate.UpdateWithField err=%+v", err)
        return nil, err
    }
    
    // 5. 事件发布
    for _, data := range dataList {
        if err := l.publishEvent(data); err != nil {
            l.Error(err)
        }
    }
    
    return &dm.Empty{}, nil
}
```

### 批量删除规范
```go
func (l *ExampleMultiDeleteLogic) ExampleMultiDelete(in *dm.ExampleMultiDeleteReq) (*dm.Empty, error) {
    // 1. 权限检查
    if err := ctxs.IsAdmin(l.ctx); err != nil {
        return nil, err
    }
    
    // 2. 查询数据
    dataList, err := l.ExampleDB.FindByFilter(l.ctx, relationDB.ExampleFilter{
        IDs: in.Ids,
    }, nil)
    if err != nil {
        return nil, err
    }
    
    // 3. 业务规则检查
    if err := l.checkMultiDeleteRules(dataList); err != nil {
        return nil, err
    }
    
    // 4. 批量删除
    err = l.ExampleDB.MultiDelete(l.ctx, in.Ids)
    if err != nil {
        l.Errorf("ExampleMultiDelete.MultiDelete err=%+v", err)
        return nil, err
    }
    
    // 5. 事件发布
    for _, data := range dataList {
        if err := l.publishEvent(data); err != nil {
            l.Error(err)
        }
    }
    
    return &dm.Empty{}, nil
}
```

## 辅助方法规范

### 参数验证方法
```go
func (l *ExampleCreateLogic) validateCreateParams(in *dm.Example) error {
    if in.Name == "" {
        return errors.Parameter.AddMsg("名称不能为空")
    }
    if len(in.Name) > 100 {
        return errors.Parameter.AddMsg("名称长度不能超过100个字符")
    }
    return nil
}
```

### 业务规则检查方法
```go
func (l *ExampleCreateLogic) checkBusinessRules(in *dm.Example) error {
    // 检查名称是否重复
    exists, err := l.ExampleDB.Exists(l.ctx, relationDB.ExampleFilter{
        Name: in.Name,
    })
    if err != nil {
        return err
    }
    if exists {
        return errors.Duplicate.AddMsgf("名称重复: %s", in.Name)
    }
    return nil
}
```

### 数据转换方法
```go
func (l *ExampleCreateLogic) toDBModel(in *dm.Example) *relationDB.Example {
    uc := ctxs.GetUserCtx(l.ctx)
    return &relationDB.Example{
        Name:        in.Name,
        Description: in.Description,
        Status:      def.StatusActive,
        TenantCode:  uc.TenantCode,
        ProjectID:   uc.ProjectID,
        CreatedBy:   uc.UserID,
    }
}

func (l *ExampleUpdateLogic) updateFields(old *relationDB.Example, in *dm.Example) error {
    if in.Name != "" {
        old.Name = in.Name
    }
    if in.Description != nil {
        old.Description = in.Description.GetValue()
    }
    if in.Status != 0 {
        old.Status = in.Status
    }
    return nil
}
```

### 过滤器构建方法
```go
func (l *ExampleIndexLogic) buildFilter(in *dm.ExampleIndexReq) relationDB.ExampleFilter {
    return relationDB.ExampleFilter{
        Name:        in.Name,
        Status:      in.Status,
        TenantCode:  in.TenantCode,
        ProjectID:   in.ProjectID,
        CreatedTime: logic.ToTimeRange(in.CreatedTime),
    }
}
```

### 缓存操作方法
```go
func (l *ExampleCreateLogic) updateCache(data *relationDB.Example) error {
    return l.svcCtx.ExampleCache.SetData(l.ctx, data.ID, data)
}

func (l *ExampleDeleteLogic) clearCache(data *relationDB.Example) error {
    return l.svcCtx.ExampleCache.DelData(l.ctx, data.ID)
}
```

### 事件发布方法
```go
func (l *ExampleCreateLogic) publishEvent(data *relationDB.Example) error {
    return l.svcCtx.FastEvent.Publish(l.ctx, topics.ExampleCreate, data)
}
```

## 权限控制规范

### 管理员权限检查
```go
// 检查是否为管理员
if err := ctxs.IsAdmin(l.ctx); err != nil {
    return nil, err
}
```

### 项目权限检查
```go
// 设置默认项目权限
l.ctx = ctxs.WithDefaultAllProject(l.ctx)

// 检查根用户权限
if err := ctxs.IsRoot(l.ctx); err == nil {
    ctxs.GetUserCtx(l.ctx).AllTenant = true
    defer func() {
        ctxs.GetUserCtx(l.ctx).AllTenant = false
    }()
}
```

### 业务权限检查
```go
func (l *ExampleUpdateLogic) checkUpdatePermission(data *relationDB.Example) error {
    uc := ctxs.GetUserCtx(l.ctx)
    if !uc.IsAdmin && data.CreatedBy != uc.UserID {
        return errors.Permissions.AddMsg("无权限修改此记录")
    }
    return nil
}
```

## 错误处理规范

### 统一错误处理
```go
// 使用 errors 包进行错误处理
if err != nil {
    if errors.Cmp(err, errors.NotFind) {
        return nil, errors.NotFind.AddDetailf("not find example id=%d", in.Id)
    }
    return nil, errors.Database.AddDetail(err)
}

// 参数错误
if in.Name == "" {
    return nil, errors.Parameter.AddMsg("名称不能为空")
}

// 业务错误
if exists {
    return nil, errors.Duplicate.AddMsgf("名称重复: %s", in.Name)
}

// 权限错误
if !hasPermission {
    return nil, errors.Permissions.AddMsg("无权限执行此操作")
}
```

### 日志记录规范
```go
// 使用 logx 进行日志记录
l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
l.Errorf("ExampleCreate.Insert err=%+v", err)
l.Error(err)
```

## 事务处理规范

### 事务操作
```go
func (l *ExampleCreateLogic) createWithRelatedData(in *dm.Example) error {
    return stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
        // 创建主记录
        if err := relationDB.NewExampleRepo(tx).Insert(l.ctx, data); err != nil {
            return err
        }
        
        // 创建关联记录
        for _, item := range relatedData {
            if err := relationDB.NewRelatedRepo(tx).Insert(l.ctx, item); err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

## 文件处理规范

### 文件上传处理
```go
func (l *ExampleUpdateLogic) handleFileUpload(in *dm.Example) error {
    if in.File != "" && in.IsUpdateFile {
        // 删除旧文件
        if oldData.File != "" {
            if err := l.svcCtx.OssClient.PrivateBucket().Delete(l.ctx, oldData.File); err != nil {
                l.Errorf("Delete file err path:%v,err:%v", oldData.File, err)
            }
        }
        
        // 上传新文件
        newPath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessExample, oss.SceneFile, 
            fmt.Sprintf("%d/%s", oldData.ID, oss.GetFileNameWithPath(in.File)))
        path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(in.File, newPath)
        if err != nil {
            return errors.System.AddDetail(err)
        }
        oldData.File = path
    }
    return nil
}
```

## 校验清单（生成前自检）

- 已设置正确的包名和导入
- Logic 结构体包含必要的字段
- 构造函数正确初始化所有依赖
- 方法签名与 gRPC 服务一致
- 包含完整的权限检查
- 包含参数验证和业务规则检查
- 包含数据转换逻辑
- 包含数据库操作
- 包含缓存更新
- 包含事件发布
- 包含错误处理
- 包含日志记录
- 事务操作正确实现
- 文件处理逻辑完整

## 约束与不做事项

- 不在 Logic 层编写数据访问逻辑，仅调用 Repository 层
- 不在 Logic 层编写复杂的业务计算，可提取到 Domain 层
- 不在 Logic 层直接操作数据库，通过 Repository 层
- 不在 Logic 层编写 UI 相关逻辑
- 不在 Logic 层硬编码业务规则，使用配置或常量
- 不在 Logic 层直接处理 HTTP 请求，通过 gRPC 接口

---

如需新增 Logic：
1) 在对应的子目录新建 `xxxLogic.go`，按本规范填写业务逻辑
2) 实现完整的 CRUD 操作方法
3) 添加必要的辅助方法
4) 实现权限控制和错误处理
5) 添加缓存和事件处理逻辑
