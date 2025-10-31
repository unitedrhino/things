# AI 面向的数据库模型定义规范（GORM）

本规范面向 AI/脚本自动生成数据库模型结构体和操作方法，确保风格统一、可直接通过 GORM 操作数据库。

## 适用范围
- 目录：`internal/repo/relationDB` 或类似的数据访问层目录
- 文件：`*.go` 数据库模型和操作方法文件
- ORM：GORM v2
- 数据库：MySQL/PostgreSQL/SQLite 等

## 文件结构规范

### 基础文件结构
每个数据库模型文件应包含以下结构：

```go
package relationDB

import (
    "context"
    "database/sql"
    "time"
    
    "your-project/internal/types"
    "your-project/pkg/stores"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

// 1. 模型结构体定义
type Example struct {
    // 字段定义...
}

// 2. 表名方法
func (m *Example) TableName() string {
    return "examples"
}

// 3. Repository 结构体
type exampleRepo struct {
    db *gorm.DB
}

// 4. 过滤器结构体
type exampleFilter struct {
    // 过滤字段...
}

// 5. 构造函数
func newExampleRepo(in any) *exampleRepo {
    return &exampleRepo{db: stores.GetCommonConn(in)}
}

// 6. 过滤方法
func (p exampleRepo) fmtFilter(ctx context.Context, f exampleFilter) *gorm.DB {
    // 过滤逻辑...
}

// 7. CRUD 操作方法
// insert, findOneByFilter, findByFilter, countByFilter, update, delete, etc.
```

## 模型结构体定义规范

### 结构体命名规范
- 模型结构体：使用帕斯卡命名法，如 `User`、`Product`、`Order`
- 表名：使用复数形式的下划线命名法，如 `users`、`products`、`orders`
- 字段名：使用下划线命名法，如 `user_name`、`product_id`、`created_at`

### 基础字段规范
每个模型都应包含以下基础字段：

```go
type Example struct {
    ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // 主键ID
    
    // 租户相关字段（可选）
    TenantCode string `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"` // 租户编码
    
    // 项目相关字段（可选）
    ProjectID int64 `gorm:"column:project_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"` // 项目ID
    AreaID    int64 `gorm:"column:area_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"`    // 项目区域ID
    AreaIDPath string `gorm:"column:area_id_path;type:varchar(100);default:'';NOT NULL"`               // 项目区域ID路径
    
    // 部门相关字段（可选）
    DeptID          int64     `gorm:"column:dept_id;type:bigint;default:0;NOT NULL"`         // 部门ID
    DeptIDPath      string    `gorm:"column:dept_id_path;type:varchar(100);default:'';NOT NULL"` // 部门ID路径
    DeptUpdatedTime time.Time `gorm:"column:dept_updated_time;default:null"`                 // 部门更新时间
    
    // 业务字段
    // ... 具体业务字段
    
    // 通用时间字段
    CreatedTime time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP"`     // 创建时间
    UpdatedTime time.Time `gorm:"column:updated_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
    
    // 软删除字段
    DeletedTime time.Time `gorm:"column:deleted_time;default:null"` // 软删除时间
}
```

### 字段类型规范

#### 基础数据类型
```go
// 字符串字段
Name        string `gorm:"column:name;type:varchar(100);NOT NULL"`                    // 名称
Description string `gorm:"column:description;type:varchar(500);default:''"`          // 描述
Code        string `gorm:"column:code;type:varchar(50);uniqueIndex;NOT NULL"`        // 编码

// 数值字段
ID          int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`         // 主键ID
Status      int64  `gorm:"column:status;index;type:smallint;default:1;NOT NULL"`    // 状态
Sort        int64  `gorm:"column:sort;type:BIGINT;default:100"`                     // 排序
Count       int64  `gorm:"column:count;type:bigint;default:0"`                      // 数量

// 布尔字段
IsEnable    bool   `gorm:"column:is_enable;index;type:boolean;default:true"`        // 是否启用
IsOnline    bool   `gorm:"column:is_online;type:boolean;default:false;NOT NULL"`    // 是否在线

// 时间字段
CreatedTime time.Time     `gorm:"column:created_time;default:CURRENT_TIMESTAMP"`     // 创建时间
UpdatedTime time.Time     `gorm:"column:updated_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
ExpTime     sql.NullTime  `gorm:"column:exp_time"`                                   // 过期时间,为0不限制
FirstLogin  sql.NullTime  `gorm:"column:first_login"`                                // 激活时间
LastLogin   sql.NullTime  `gorm:"column:last_login"`                                 // 最后上线时间

// 可空字符串字段
Phone       sql.NullString `gorm:"column:phone;type:varchar(20)"`                    // 手机号
Iccid       sql.NullString `gorm:"column:iccid;type:varchar(20)"`                    // SIM卡卡号
Address     string         `gorm:"column:address;type:varchar(512);default:''"`     // 所在地址
```

#### 复杂数据类型
```go
// JSON 字段
Tags        map[string]string      `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"`              // 标签
Config      map[string]string      `gorm:"column:config;type:json;serializer:json;NOT NULL;default:'{}'"`            // 配置
SchemaAlias map[string]string      `gorm:"column:schema_alias;type:json;serializer:json;NOT NULL;default:'{}'"`      // 物模型别名

// 地理位置字段
Position    stores.Point           `gorm:"column:position;"`                         // 设备的位置(默认百度坐标系BD09)

// 枚举类型字段
DeviceType  int64                  `gorm:"column:device_type;index;type:smallint;default:1"`                         // 设备类型:1:设备,2:网关,3:子设备
Status      def.DeviceStatus       `gorm:"column:status;index;type:smallint;default:1;NOT NULL"`                     // 设备状态 1-未激活，2-在线，3-离线 4-异常
MobileOperator def.MobileOperator  `gorm:"column:mobile_operator;type:smallint;default:10;NOT NULL"`                 // 移动运营商:1)移动 2)联通 3)电信 4)广电 10) 无

// 外键关联
ProductInfo *ProductInfo           `gorm:"foreignKey:ProductID;references:ProductID"` // 添加外键
```

### 索引规范
```go
// 主键索引
ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`

// 普通索引
TenantCode string `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`
Status     int64  `gorm:"column:status;index;type:smallint;default:1;NOT NULL"`

// 联合索引
ProjectID int64 `gorm:"column:project_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"`
AreaID    int64 `gorm:"column:area_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"`

// 唯一索引
ProductID  string `gorm:"column:product_id;type:varchar(100);uniqueIndex:idx_device_info_product_id_deviceName;NOT NULL"`
DeviceName string `gorm:"column:device_name;type:varchar(100);uniqueIndex:idx_device_info_product_id_deviceName;NOT NULL"`

// 复合索引
ProductID string `gorm:"column:product_id;type:varchar(100);index:idx_device_info_pd_dn;NOT NULL"`
DeviceName string `gorm:"column:device_name;type:varchar(100);index:idx_device_info_pd_dn;NOT NULL"`
```

## Repository 结构体规范

### Repository 命名规范
```go
type exampleRepo struct {
    db *gorm.DB
}
```

### 构造函数规范
```go
func newExampleRepo(in any) *exampleRepo {
    return &exampleRepo{db: stores.GetCommonConn(in)}
}
```

## 过滤器结构体规范

### 过滤器命名规范
```go
type exampleFilter struct {
    // 基础过滤字段
    ID          int64
    IDs         []int64
    TenantCode  string
    TenantCodes []string
    
    // 项目相关过滤
    ProjectID   int64
    ProjectIDs  []int64
    AreaID      int64
    AreaIDs     []int64
    AreaIDPath  string
    AreaIDPaths []string
    NotAreaID   int64
    NotAreaIDs  []int64
    
    // 业务字段过滤
    Name        string
    Names       []string
    Status      int64
    Statuses    []int64
    IsEnable    bool
    
    // 时间范围过滤
    CreatedTime *TimeRange
    UpdatedTime *TimeRange
    
    // 标签过滤
    Tags        map[string]string
    TagsLike    map[string]string
    
    // 关联查询
    WithProduct bool
    WithFiles   bool
    WithJobs    bool
    
    // 比较查询
    Count       *Cmp
    Price       *Cmp
    
    // 地理位置过滤
    Position    Point
    Range       int64
    
    // 分页相关
    Page        *PageInfo
    
    // 排序相关
    OrderBy     string
    OrderDesc   bool
}
```

### 过滤方法规范
```go
func (p exampleRepo) fmtFilter(ctx context.Context, f exampleFilter) *gorm.DB {
    db := p.db.WithContext(ctx)
    
    // 基础过滤
    if f.ID != 0 {
        db = db.Where("id = ?", f.ID)
    }
    if len(f.IDs) != 0 {
        db = db.Where("id in ?", f.IDs)
    }
    if f.TenantCode != "" {
        db = db.Where("tenant_code = ?", f.TenantCode)
    }
    if len(f.TenantCodes) != 0 {
        db = db.Where("tenant_code in ?", f.TenantCodes)
    }
    
    // 项目过滤
    if f.ProjectID != 0 {
        db = db.Where("project_id = ?", f.ProjectID)
    }
    if len(f.ProjectIDs) != 0 {
        db = db.Where("project_id in ?", f.ProjectIDs)
    }
    if f.AreaID != 0 {
        db = db.Where("area_id = ?", f.AreaID)
    }
    if len(f.AreaIDs) != 0 {
        db = db.Where("area_id in ?", f.AreaIDs)
    }
    if f.AreaIDPath != "" {
        db = db.Where("area_id_path like ?", f.AreaIDPath+"%")
    }
    if len(f.AreaIDPaths) != 0 {
        var conditions []string
        var args []interface{}
        for _, path := range f.AreaIDPaths {
            conditions = append(conditions, "area_id_path like ?")
            args = append(args, path+"%")
        }
        db = db.Where("("+strings.Join(conditions, " OR ")+")", args...)
    }
    
    // 业务字段过滤
    if f.Name != "" {
        db = db.Where("name like ?", "%"+f.Name+"%")
    }
    if len(f.Names) != 0 {
        db = db.Where("name in ?", f.Names)
    }
    if f.Status != 0 {
        db = db.Where("status = ?", f.Status)
    }
    if len(f.Statuses) != 0 {
        db = db.Where("status in ?", f.Statuses)
    }
    if f.IsEnable {
        db = db.Where("is_enable = ?", f.IsEnable)
    }
    
    // 时间范围过滤
    if f.CreatedTime != nil {
        if f.CreatedTime.Start != 0 {
            db = db.Where("created_time >= ?", time.Unix(f.CreatedTime.Start/1000, 0))
        }
        if f.CreatedTime.End != 0 {
            db = db.Where("created_time <= ?", time.Unix(f.CreatedTime.End/1000, 0))
        }
    }
    
    // 标签过滤
    if len(f.Tags) != 0 {
        for key, value := range f.Tags {
            db = db.Where("JSON_EXTRACT(tags, ?) = ?", "$."+key, value)
        }
    }
    if len(f.TagsLike) != 0 {
        for key, value := range f.TagsLike {
            db = db.Where("JSON_EXTRACT(tags, ?) like ?", "$."+key, "%"+value+"%")
        }
    }
    
    // 关联查询
    if f.WithProduct {
        db = db.Preload("ProductInfo")
    }
    if f.WithFiles {
        db = db.Preload("Files")
    }
    if f.WithJobs {
        db = db.Preload("Jobs")
    }
    
    // 比较查询
    if f.Count != nil {
        db = f.Count.Where(db, "count")
    }
    if f.Price != nil {
        db = f.Price.Where(db, "price")
    }
    
    // 地理位置过滤
    if f.Position.Longitude != 0 && f.Position.Latitude != 0 && f.Range > 0 {
        // 使用 Haversine 公式计算距离
        db = db.Where(`
            (6371 * acos(
                cos(radians(?)) * cos(radians(JSON_EXTRACT(position, '$.latitude'))) * 
                cos(radians(JSON_EXTRACT(position, '$.longitude')) - radians(?)) + 
                sin(radians(?)) * sin(radians(JSON_EXTRACT(position, '$.latitude')))
            )) <= ?
        `, f.Position.Latitude, f.Position.Longitude, f.Position.Latitude, f.Range/1000)
    }
    
    return db
}
```

## CRUD 操作方法规范

### 插入操作
```go
// 单条插入
func (p exampleRepo) insert(ctx context.Context, data *Example) error {
    result := p.db.WithContext(ctx).Create(data)
    return stores.ErrFmt(result.Error)
}

// 批量插入（支持冲突更新）
func (p exampleRepo) multiInsert(ctx context.Context, data []*Example) error {
    err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&Example{}).Create(data).Error
    return stores.ErrFmt(err)
}
```

### 查询操作
```go
// 根据ID查询单条
func (p exampleRepo) findOne(ctx context.Context, id int64) (*Example, error) {
    var result Example
    err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
    if err != nil {
        return nil, stores.ErrFmt(err)
    }
    return &result, nil
}

// 根据过滤器查询单条
func (p exampleRepo) findOneByFilter(ctx context.Context, f exampleFilter) (*Example, error) {
    var result Example
    db := p.fmtFilter(ctx, f)
    err := db.First(&result).Error
    if err != nil {
        return nil, stores.ErrFmt(err)
    }
    return &result, nil
}

// 根据过滤器查询列表（支持分页）
func (p exampleRepo) findByFilter(ctx context.Context, f exampleFilter, page *PageInfo) ([]*Example, error) {
    var results []*Example
    db := p.fmtFilter(ctx, f).Model(&Example{})
    if page != nil {
        db = page.ToGorm(db)
    }
    err := db.Find(&results).Error
    if err != nil {
        return nil, stores.ErrFmt(err)
    }
    return results, nil
}

// 统计数量
func (p exampleRepo) countByFilter(ctx context.Context, f exampleFilter) (size int64, err error) {
    db := p.fmtFilter(ctx, f).Model(&Example{})
    err = db.Count(&size).Error
    return size, stores.ErrFmt(err)
}
```

### 更新操作
```go
// 根据ID更新
func (p exampleRepo) update(ctx context.Context, data *Example) error {
    err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
    return stores.ErrFmt(err)
}

// 根据过滤器更新指定字段
func (p exampleRepo) updateWithField(ctx context.Context, f exampleFilter, updates map[string]any) error {
    db := p.fmtFilter(ctx, f)
    err := db.Model(&Example{}).Updates(updates).Error
    return stores.ErrFmt(err)
}

// 批量更新
func (p exampleRepo) multiUpdate(ctx context.Context, data []*Example) error {
    err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        for _, item := range data {
            if err := tx.Where("id = ?", item.ID).Save(item).Error; err != nil {
                return err
            }
        }
        return nil
    })
    return stores.ErrFmt(err)
}
```

### 删除操作
```go
// 根据ID删除
func (p exampleRepo) delete(ctx context.Context, id int64) error {
    err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&Example{}).Error
    return stores.ErrFmt(err)
}

// 根据过滤器删除
func (p exampleRepo) deleteByFilter(ctx context.Context, f exampleFilter) error {
    db := p.fmtFilter(ctx, f)
    err := db.Delete(&Example{}).Error
    return stores.ErrFmt(err)
}

// 批量删除
func (p exampleRepo) multiDelete(ctx context.Context, ids []int64) error {
    err := p.db.WithContext(ctx).Where("id in ?", ids).Delete(&Example{}).Error
    return stores.ErrFmt(err)
}
```

## 特殊操作方法规范

### 存在性检查
```go
func (p exampleRepo) exists(ctx context.Context, f exampleFilter) (bool, error) {
    var count int64
    db := p.fmtFilter(ctx, f).Model(&Example{})
    err := db.Count(&count).Error
    if err != nil {
        return false, stores.ErrFmt(err)
    }
    return count > 0, nil
}
```

### 聚合查询
```go
func (p exampleRepo) aggregateByFilter(ctx context.Context, f exampleFilter, field string, fn string) (interface{}, error) {
    var result interface{}
    db := p.fmtFilter(ctx, f).Model(&Example{})
    err := db.Select(fmt.Sprintf("%s(%s) as result", fn, field)).Scan(&result).Error
    if err != nil {
        return nil, stores.ErrFmt(err)
    }
    return result, nil
}
```

### 分组查询
```go
func (p exampleRepo) groupByFilter(ctx context.Context, f exampleFilter, groupBy string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    db := p.fmtFilter(ctx, f).Model(&Example{})
    err := db.Select(fmt.Sprintf("%s, count(*) as count", groupBy)).Group(groupBy).Find(&results).Error
    if err != nil {
        return nil, stores.ErrFmt(err)
    }
    return results, nil
}
```

## 关联关系规范

### 一对一关联
```go
type Example struct {
    ID        int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
    ProductID string `gorm:"column:product_id;type:varchar(100);NOT NULL"`
    
    // 一对一关联
    ProductInfo *ProductInfo `gorm:"foreignKey:ProductID;references:ProductID"`
}
```

### 一对多关联
```go
type Example struct {
    ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
    
    // 一对多关联
    Files []*ExampleFile `gorm:"foreignKey:ExampleID;references:ID"`
}
```

### 多对多关联
```go
type Example struct {
    ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
    
    // 多对多关联
    Tags []*Tag `gorm:"many2many:example_tags;foreignKey:ID;joinForeignKey:ExampleID;References:ID;joinReferences:TagID"`
}
```

## 事务操作规范

### 事务方法
```go
func (p exampleRepo) transaction(ctx context.Context, fn func(*gorm.DB) error) error {
    return p.db.WithContext(ctx).Transaction(fn)
}

// 使用示例
func (p exampleRepo) createWithFiles(ctx context.Context, data *Example, files []*ExampleFile) error {
    return p.transaction(ctx, func(tx *gorm.DB) error {
        // 创建主记录
        if err := tx.Create(data).Error; err != nil {
            return err
        }
        
        // 创建关联文件
        for _, file := range files {
            file.ExampleID = data.ID
            if err := tx.Create(file).Error; err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

## 错误处理规范

### 错误处理
```go
// 统一使用 stores.ErrFmt 处理错误
func (p exampleRepo) insert(ctx context.Context, data *Example) error {
    result := p.db.WithContext(ctx).Create(data)
    return stores.ErrFmt(result.Error)
}

// 查询错误处理
func (p exampleRepo) findOne(ctx context.Context, id int64) (*Example, error) {
    var result Example
    err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
    if err != nil {
        return nil, stores.ErrFmt(err)
    }
    return &result, nil
}
```

## 性能优化规范

### 预加载关联数据
```go
func (p exampleRepo) findWithRelations(ctx context.Context, id int64) (*Example, error) {
    var result Example
    err := p.db.WithContext(ctx).
        Preload("ProductInfo").
        Preload("Files").
        Preload("Tags").
        Where("id = ?", id).
        First(&result).Error
    if err != nil {
        return nil, stores.ErrFmt(err)
    }
    return &result, nil
}
```

### 批量操作优化
```go
// 批量插入优化
func (p exampleRepo) batchInsert(ctx context.Context, data []*Example, batchSize int) error {
    return p.db.WithContext(ctx).CreateInBatches(data, batchSize).Error
}

// 批量更新优化
func (p exampleRepo) batchUpdate(ctx context.Context, updates []map[string]interface{}) error {
    return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        for _, update := range updates {
            if err := tx.Model(&Example{}).Where("id = ?", update["id"]).Updates(update).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

## 校验清单（生成前自检）

- 已设置正确的包名
- 已导入必要的依赖包
- 模型结构体使用帕斯卡命名法
- 表名方法返回正确的表名
- Repository 结构体包含 `db *gorm.DB` 字段
- 过滤器结构体包含所有必要的过滤字段
- 构造函数使用小驼峰命名
- 过滤方法使用 `fmtFilter` 命名
- 所有 CRUD 方法都包含 `context.Context` 参数
- 所有方法都使用统一的错误处理
- 字段标签包含正确的 GORM 配置
- 索引配置合理且完整
- 关联关系配置正确
- 软删除字段配置正确

## 约束与不做事项

- 不在模型文件中编写业务逻辑，仅保留数据访问层代码
- 不在模型文件中混入具体实现或示例数据，仅保留结构体和方法定义
- 不使用过时的 GORM v1 语法
- 不直接使用 `gorm.ErrRecordNotFound`，统一使用错误处理函数

---

如需新增模型：
1) 在 `relationDB/` 目录新建 `xxx.go`，按本规范填写模型定义和操作方法
2) 在 `model.go` 中定义模型结构体
3) 实现完整的 CRUD 操作方法
4) 添加必要的过滤器和关联关系
