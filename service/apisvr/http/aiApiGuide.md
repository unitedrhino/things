## AI 面向的 API 定义规范（goctl api）

本规范面向 AI/脚本自动生成本项目的 `.api` 接口定义文件，确保风格统一、可直接通过 goctl 生成代码、Swagger 与 access 权限描述文件。

### 适用范围
- 目录：`service/apisvr/http` 及其子目录（如 `things/`）
- 文件：`.api` 接口定义文件、引用 `common.api` 的公共类型

### 基础结构
每个聚合入口 `.api` 顶部必须包含 `info` 与 `import`，子模块文件也应包含 `info`。

```startLine:endLine:service/apisvr/http/api.api
info (
	title:   "web端网关"
	desc:    "负责http协议的网关"
	author:  "杨磊"
	email:   "603685348@qq.com"
	version: "v1"
)

import "common.api" //公共结构体定义
import "things.api" //物模型管理
```

### 模块拆分与入口聚合
- 聚合入口：`api.api`、`things.api` 这类文件仅做模块聚合与说明，不直接定义路由。
- 业务模块：放置于子目录（如 `things/device/xxx.api`、`things/product/xxx.api`）。
- 规范示例：

```startLine:endLine:service/apisvr/http/things.api
info(
	title: "web端网关"
	desc: "负责http协议的网关"
	author: "杨磊"
	email: "603685348@qq.com"
	version: "v1"
)

import "things/protocol.api"  //协议相关功能接口
import "things/product.api"   //产品相关功能接口
import "things/device.api"    //设备相关功能接口
import "things/group.api"     //设备分组相关功能接口
import "things/ota.api"       //ota相关功能接口
import "things/schema/common.api"  //物模型相关接口
import "things/user/device.api"
import "things/slot.api"
```

### @server 规范（服务域、前缀、权限、中间件）
每个可实际提供路由的模块文件必须声明一个 `@server` 块，并使用统一前缀、权限前缀与中间件。

- `group`: 唯一的模块标识路径，使用目录风格，例如 `things/device/info`
- `prefix`: 统一使用 `/api/v1/<domain>/<sub>/...` 的 REST 前缀
- `accessCodePrefix`: 英文驼峰/帕斯卡前缀，用于权限编码前缀（唯一）
- `accessNamePrefix`: 中文前缀，用于权限名称
- `accessGroup`: 中文权限分组名（菜单/模块）
- `middleware`: 统一写法 `CheckTokenWare,InitCtxsWare`，如无需鉴权再在 `@doc` 中显式标注

示例：

```startLine:endLine:service/apisvr/http/things/device/info.api
@server(
	group: things/device/info
	prefix: /api/v1/things/device/info
	accessCodePrefix: "thingsDeviceInfo"
	accessNamePrefix: "设备信息"
	accessGroup: "设备管理"
	middleware:  CheckTokenWare,InitCtxsWare
)
```

### 路由与处理器定义
在 `service api { ... }` 内声明接口，按如下要求：

- 使用 `@doc` 标注元数据：
  - `summary`: 必填，中文一句话摘要
  - `isNeedAuth`: 是否需要鉴权（"true"/"false"，字符串）
  - `authType`: 鉴权类型（常用：`admin`）
  - `businessType`: 业务类型（可选：`find`、`modify`），用于审计/日志
  - `injectFormdataParam`: 表单上传字段名（如有文件上传，必填）
- `@handler`: 处理器名称使用小驼峰动词短语，如 `create`、`update`、`index`、`read`、`delete`
- HTTP 方法：统一使用 `post`（本项目约定）
- 路径：语义化短路径，使用连字符分词，如 `/multi-import`、`/bind/token/read`
- 请求体与返回体：使用下文类型规范中的结构体

示例：

```startLine:endLine:service/apisvr/http/things/device/info.api
service api {
	@doc(
		summary: "新增设备"
		isNeedAuth: "true"
		authType: "admin"
	)
	@handler create
	post /create (DeviceInfo) returns ()
}
```

### 公共类型与复用
统一从 `common.api` 引入通用结构体，如分页、排序、标签、地理位置、时间区间等。

```startLine:endLine:service/apisvr/http/common.api
type(
	PageInfo{ Page int64 `json:"page,optional" form:"page,optional"` ... }
	PageResp{ Page int64 `json:"page,optional"` PageSize int64 `json:"pageSize,optional"` Total int64 `json:"total"` }
	OrderBy{ Field string `json:"field,optional"` Sort int64 `json:"sort,optional"` }
	Tag{ Key string `json:"key"` Value string `json:"value"` }
	Point{ Longitude float64 `json:"longitude,range=[0:180]"` Latitude float64 `json:"latitude,range=[0:90]"` }
	DateRange{ Start string `json:"start,optional"` End string `json:"end,optional"` }
	TimeRange{ Start int64 `json:"start,optional"` End int64 `json:"end,optional"` }
)
```

### 字段标签规范（极其重要）
- 使用反引号标签，优先 `json:"<name>[,string][,optional][,omitempty]"`
- `,optional`：字段为可选
- `,omitempty`：序列化时忽略零值（仅在需要减少响应体时使用）
- `,string`：以字符串传输数值（前端/网关统一 ID 字符串化，避免精度问题）。常用于 `id/userID/projectID/createdTime/...`
- 校验/约束：使用 `range=[min:max]`（如 `range=[0:2]`）
- 表单：如为表单参数可额外标注 `form:"<name>[,optional]"`
- 时间戳：统一毫秒/秒需在字段名或注释中说明，本仓库常用毫秒时间戳 `,string`

示例片段：

```startLine:endLine:service/apisvr/http/things/device/info.api
DeviceInfo{
	ProjectID int64 `json:"projectID,string,optional"`
	AreaID    int64 `json:"areaID,string,optional"`
	FirstLogin int64 `json:"firstLogin,optional,string"`
	IsOnline int64 `json:"isOnline,optional,range=[0:2],omitempty"`
}
```

### 分页与查询筛选约定
- 分页入参统一：`Page *PageInfo json:"page,optional"`，出参统一继承 `PageResp`
- ID/枚举批量筛选：数组使用复数名 `IDs/Types/Names`，必要时 `,string`
- 范围/比较：
  - 时间范围：`DateRange` 或 `TimeRange`
  - 数值比较：`*CompareInt64`，字符串比较：`*CompareString`
- 位置范围：`Position *Point` 与 `Range int64`

### 命名规范
- 类型名：`<领域><实体><动作/后缀>`，如 `DeviceInfoIndexReq`、`DeviceCountResp`
- 处理器：小驼峰动词短语：`create/update/delete/index/read/bind/unbind/multiImport/multiExport`
- 路径：全小写、连字符分词：`/multi-import`、`/bind/token/read`
- 访问控制前缀：`things<Device|Product|Schema><Sub>` 的帕斯卡命名，唯一

### 鉴权与中间件
- 默认模块需要鉴权：`middleware: CheckTokenWare,InitCtxsWare`
- 个别接口如无需鉴权，在 `@doc` 中明确 `isNeedAuth: "false"`（否则默认为需要）
- 管理员接口在 `@doc` 中添加 `authType: "admin"`

### 导入与导出（文件上传）
- 有文件上传时：
  - `@doc(injectFormdataParam: "file")`
  - 入参体按需定义占位（可空），或通过 `form:"file"` 注释说明
  - 路径使用 `/multi-import`、`/multi-export`

### 物模型相关约定
- 通用物模型接口置于 `things/schema/common.api`，并按需在 `things.api` 中聚合导入
- 物模型条目使用 `CommonSchemaInfo`，筛选采用 `Type/Types/Identifiers` 等字段

```startLine:endLine:service/apisvr/http/things/schema/common.api
@server(
	group: things/schema/common
	prefix: /api/v1/things/schema/common
	accessCodePrefix: "thingsSchemaCommon"
	accessNamePrefix: "通用物模型"
	accessGroup: "通用物模型"
	middleware:  CheckTokenWare,InitCtxsWare
)
```

### 设备领域约定（示例）
- 模块路径：`things/device/info`
- 典型接口：`/create`、`/update`、`/index`、`/read`、`/delete`、`/bind`、`/unbind`、`/multi-import`
- 统计接口：`/count`，使用 `DeviceCountReq/DeviceCountResp`
- 设备核心键：`ProductID`（string）、`DeviceName`（string），涉及 ID 类统一 `,string`

### 生成命令
在 `service/README.md` 已给出标准命令，保持不变：

```bash
cd apisvr && goctl api go -api http/api.api -dir ./ --style=goZero && goctl api swagger -filename swagger.json -api http/api.api -dir ./http && goctl api access -api http/api.api -dir ./http && cd ..

cd apisvr && goctl api swagger -filename swagger.json -api http/api.api -dir ./http && cd ..
cd apisvr && goctl api access  -api http/api.api -dir ./http && cd ..
```

### 校验清单（生成前自检）
- 已设置 `@server` 的 `group/prefix/accessCodePrefix/accessNamePrefix/accessGroup/middleware`
- 所有接口均含 `@doc(summary)`，并按需设置 `isNeedAuth`、`authType`、`businessType`、`injectFormdataParam`
- 请求/返回体字段均按标签规范：`json:"...[,string][,optional][,omitempty]"`
- ID/时间戳等需要字符串化的字段均添加了 `,string`
- 分页统一 `Page *PageInfo` 入参，响应包含 `PageResp`
- 路径统一小写、连字符、动宾短语；HTTP 方法统一 `post`

### 约束与不做事项
- 不在 `.api` 中编写业务说明文档，长描述仅放在 `summary` 或文件 `desc`
- 不在 `.api` 中混入具体实现或示例 JSON，仅保留类型与路由
- 不使用 GET/PUT/DELETE 等方法（本项目统一 `post`）

---

如需新增模块：
1) 在对应目录新建 `xxx.api`，按本规范填写 `info/@server/service/@doc/type`
2) 在上层聚合 `things/*.api` 或根聚合 `api.api` 中 `import` 该模块文件
3) 执行“生成命令”一节中的命令


