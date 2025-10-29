# AI 面向的 gRPC 接口定义规范（protobuf）

本规范面向 AI/脚本自动生成本项目的 `.proto` gRPC 接口定义文件，确保风格统一、可直接通过 goctl 生成代码。

## 适用范围
- 目录：`service/dmsvr/proto` 及其子目录
- 文件：`.proto` gRPC 接口定义文件
- 生成工具：goctl rpc protoc

## 基础结构

### 文件头部定义
每个 `.proto` 文件必须包含以下头部信息：

```protobuf
syntax = "proto3";
option go_package = "pb/dm";
package dm;

import "google/protobuf/wrappers.proto";
```

- `syntax = "proto3"`：使用 protobuf v3 语法
- `option go_package`：指定生成的 Go 包路径，格式为 `pb/<服务名>`
- `package`：包名，通常与目录名一致
- `import`：导入必要的依赖，常用 `google/protobuf/wrappers.proto` 用于可选字段

## 公共消息类型定义

### 基础类型
项目中使用的基础消息类型，应在文件开头统一定义：

```protobuf
message Empty {}

message TimeRange {
  int64 start = 1;
  int64 end = 2;
}

message PageInfo {
  int64 page = 1;
  int64 size = 2;
  repeated OrderBy orders = 3;
  message OrderBy {
    string field = 1;    // 排序的字段名
    int64 sort = 2;      // 排序方式：1 从小到大, 2 从大到小
  }
}

message CompareString {
  string CmpType = 1;    // "=":相等 "!=":不相等 ">":大于">=":大于等于"<":小于"<=":小于等于 "like":模糊查询
  string value = 2;      // 值
}

message CompareInt64 {
  string CmpType = 1;    // 比较类型，同上
  int64 value = 2;       // 值
}
```

### 通用标识符类型
```protobuf
message WithID {
  int64 id = 1;
}

message WithIDCode {
  int64 id = 1;
  string code = 2;
}

message WithIDChildren {
  int64 id = 1;
  bool withChildren = 2;
}

message IDPath {
  int64 id = 1;
  string idPath = 2;
  int64 noParentID = 3;
}

message IDsInfo {
  repeated int64 ids = 1;
  repeated string idPaths = 2;
}
```

### 核心实体类型
```protobuf
message DeviceCore {
  string productID = 1;    // 产品ID
  string deviceName = 2;   // 设备名称
}

message Point {
  double longitude = 1;    // 经度
  double latitude = 2;     // 纬度
}

message FileCore {
  string path = 1;         // 文件的路径
  bool isUpdate = 2;       // 如果要更新该文件,则需要将该参数置为true
}

message SendOption {
  int64 timeoutToFail = 1;     // 超时失败时间
  int64 requestTimeout = 2;    // 请求超时,超时后会进行重试
  int64 retryInterval = 3;     // 重试间隔
}
```

## 服务定义规范

### 服务命名规范
- 服务名使用帕斯卡命名法（PascalCase）
- 服务名应体现业务领域，如：`DeviceManage`、`ProductManage`、`SchemaManage`
- 服务名应简洁明了，避免过长的名称

### 服务方法定义规范
```protobuf
service DeviceManage {
  // 鉴定是否是root账号(提供给mqtt broker)
  rpc rootCheck(RootCheckReq) returns (Empty);
  
  // 新增设备
  rpc deviceInfoCreate(DeviceInfo) returns (Empty);
  
  // 更新设备
  rpc deviceInfoUpdate(DeviceInfo) returns (Empty);
  
  // 删除设备
  rpc deviceInfoDelete(DeviceInfoDeleteReq) returns (Empty);
  
  // 获取设备信息列表
  rpc deviceInfoIndex(DeviceInfoIndexReq) returns (DeviceInfoIndexResp);
  
  // 获取设备信息详情
  rpc deviceInfoRead(DeviceInfoReadReq) returns (DeviceInfo);
}
```

### 方法命名规范
- 方法名使用小驼峰命名法（camelCase）
- 方法名应体现操作类型和业务实体
- 标准操作命名：
  - `create`：新增
  - `update`：更新
  - `delete`：删除
  - `index`：列表查询
  - `read`：详情查询
  - `multiCreate`：批量新增
  - `multiUpdate`：批量更新
  - `multiDelete`：批量删除
  - `multiImport`：批量导入
  - `multiExport`：批量导出

### 方法参数和返回值规范
- 请求参数：使用 `Req` 后缀，如 `DeviceInfoCreateReq`
- 响应参数：使用 `Resp` 后缀，如 `DeviceInfoIndexResp`
- 无参数方法：使用 `Empty` 作为参数或返回值
- 单个实体返回：直接使用实体类型，如 `DeviceInfo`
- 列表返回：使用 `Resp` 类型包含列表和分页信息

## 消息类型定义规范

### 消息命名规范
- 消息名使用帕斯卡命名法（PascalCase）
- 请求消息：`<实体名><操作>Req`，如 `DeviceInfoCreateReq`
- 响应消息：`<实体名><操作>Resp`，如 `DeviceInfoIndexResp`
- 实体消息：直接使用实体名，如 `DeviceInfo`、`ProductInfo`

### 字段定义规范
```protobuf
message DeviceInfo {
  int64 id = 34;                                    // 主键ID
  string tenantCode = 26;                           // 租户号,只有default租户能查到这个字段
  string productID = 1;                             // 产品id 只读
  int64 projectID = 2;                              // 项目id 只读
  int64 areaID = 3;                                 // 项目区域id
  string areaIDPath = 42;                           // 区域id 路径
  string productName = 4;                           // 产品名称 只读
  string deviceName = 5;                            // 设备名称 读写
  int64 createdTime = 6;                            // 创建时间 只读
  string secret = 7;                                // 设备秘钥 只读
  string cert = 8;                                  // 设备证书 只读
  string imei = 9;                                  // IMEI号信息 只读
  string mac = 10;                                  // MAC号信息 只读
  google.protobuf.StringValue version = 11;         // 固件版本 读写
  string hardInfo = 12;                             // 模组硬件型号 只读
  string softInfo = 13;                             // 模组软件版本 只读
  Point Position = 14;                              // 设备定位,默认百度坐标系
  google.protobuf.StringValue address = 15;         // 所在地址 读写
  google.protobuf.StringValue adcode = 45;          // 地区编码 读写
  map<string, string> tags = 16;                    // 设备标签
  int64 isOnline = 17;                              // 在线状态 1离线 2在线 只读
  int64 firstLogin = 18;                            // 激活时间 只读
  int64 firstBind = 37;                             // 第一次绑定的时间
  int64 lastBind = 50;                              // 最后一次绑定时间
  int64 lastLogin = 19;                             // 最后上线时间 只读
  int64 lastOffline = 62;                           // 最后离线时间 只读
  int64 logLevel = 20;                              // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试 读写
  google.protobuf.StringValue deviceAlias = 21;     // 设备别名 读写
  int64 mobileOperator = 22;                        // 移动运营商:1)移动 2)联通 3)电信 4)广电
  google.protobuf.StringValue phone = 23;           // 手机号
  google.protobuf.StringValue iccid = 24;           // SIM卡卡号
  map<string, string> schemaAlias = 25;             // 设备物模型别名,如果是结构体类型则key为xxx.xxx
  google.protobuf.Int64Value rssi = 27;             // 设备信号（信号极好[-55— 0]，信号好[-70— -55]，信号一般[-85— -70]，信号差[-100— -85]）
  int64 ratedPower = 28;                            // 额定功率:单位w/h
  map<string, string> protocolConf = 29;            // 协议配置
  map<string, string> subProtocolConf = 51;         // 子协议配置
  int64 status = 30;                                // 设备状态 1-未激活，2-在线，3-离线 4-异常(频繁上下线,告警中) 5-设备到期
  int64 isEnable = 31;                              // 是否启用
  int64 deviceType = 33;                            // 设备类型:1:设备,2:网关,3:子设备
  int64 netType = 35;                               // 网络类型
  IDPathWithUpdate distributor = 36;                // 过滤分销商的数据
  google.protobuf.Int64Value expTime = 38;          // 到期时间,如果为0,则不限制
  int64 NeedConfirmJobID = 39;                      // 需要app确认升级的任务ID,为0是没有
  string NeedConfirmVersion = 40;                   // 待确认升级的版本
  int64 userID = 41;                                // 拥有者的用户ID
  string productImg = 43;                           // 产品图片
  int64 categoryID = 44;                            // 产品品类
  string lastIp = 49;                               // 最后登录的ip,只读
  string lastLocalIp = 63;                          // 最后的局域网ip,只读
  google.protobuf.StringValue desc = 47;            // 描述
  IDPathWithUpdate dept = 52;                       // 过滤分销商的数据
  DeviceInfo Gateway = 46;                          // 子设备绑定的网关
  int64 sort = 53;                                  // 排序
  string groupPurpose = 58;                         // 更新的时候使用,将设备该用途下的分组进行更新
  repeated int64 groupIDs = 59;                     // 如果要更新分组,必须填写groupPurpose
  map<string, IDsInfo> BelongGroup = 61;            // key是group的purpose, value是里面包含的分组id
  string deviceImg = 54;                            // 设备图片
  bool isUpdateDeviceImg = 55;                      // 只有这个参数为true的时候才会更新设备图片,传参为图片的file path
  string file = 56;                                 // 设备文件
  bool isUpdateFile = 57;                           // 只有这个参数为true的时候才会更新设备文件,传参为文件的file path
}
```

### 字段类型规范
- **基础类型**：`int64`、`string`、`bool`、`double`、`bytes`
- **可选字段**：使用 `google.protobuf.StringValue`、`google.protobuf.Int64Value` 等
- **数组字段**：使用 `repeated` 关键字
- **映射字段**：使用 `map<string, string>` 等
- **嵌套消息**：直接使用消息类型名

### 字段编号规范
- 字段编号从 1 开始递增
- 已删除的字段编号不能重复使用
- 新增字段使用下一个可用编号
- 重要字段使用较小的编号（1-10）
- 可选字段使用较大的编号（20+）

### 字段注释规范
- 每个字段都应包含中文注释
- 注释应说明字段的用途和约束
- 枚举值应详细说明每个选项的含义
- 时间字段应说明时间戳格式（毫秒/秒）

## 分页和查询规范

### 分页请求规范
```protobuf
message DeviceInfoIndexReq {
  PageInfo page = 1;                    // 分页信息 只获取一个则不填
  string productID = 2;                 // 过滤条件: 产品id
  repeated string productIDs = 23;      // 过滤条件: 产品id
  string deviceName = 3;                // 过滤条件:模糊查询 设备名
  repeated string deviceNames = 4;      // 过滤条件:精准查询 设备名
  string deviceNameOrAlias = 58;        // 过滤条件:模糊查询 设备名或别名
  map<string, string> tags = 5;         // key tag过滤查询,非模糊查询 为tag的名,value为tag对应的值
  map<string, string> tagsLike = 46;    // key tag过滤查询,模糊查询 为tag的名,value为tag对应的值
  int64 range = 6;                      // 过滤条件:取距离坐标点固定范围内的设备
  Point Position = 7;                   // 设备定位,默认百度坐标系
  int64 areaID = 36;
  repeated int64 areaIDs = 8;           // 区域ids
  string areaIDPath = 32;               // 区域路径过滤
  repeated string areaIDPaths = 34;     // 区域路径过滤
  string deviceAlias = 9;               // 过滤条件:模糊查询 设备别名
  int64 isOnline = 10;                  // 在线状态过滤 1离线 2在线
  int64 productCategoryID = 11;         // 产品品类id
  repeated int64 productCategoryIDs = 33; // 产品品类id
  repeated DeviceCore devices = 12;
  int64 withShared = 13;                // 过滤分享的设备(这里只获取分享的设备) 1: 同时获取分享的设备 2:只获取分享的设备
  int64 withCollect = 24;               // 过滤收藏的设备(这里只获取收藏的设备) 1: 同时获取收藏的设备 2:只获取收藏的设备
  int64 netType = 25;                   // 通讯方式:1:其他,2:wi-fi,3:2G/3G/4G,4:5G,5:BLE,6:LoRaWAN
  string tenantCode = 14;               // 租户过滤
  repeated string versions = 15;        // 版本列表
  string notVersion = 39;               // 非版本
  int64 deviceType = 16;                // 过滤设备类型:0:全部,1:设备,2:网关,3:子设备
  repeated int64 deviceTypes = 21;
  DeviceCore gateway = 17;              // 获取网关下子设备列表
  int64 groupID = 18;
  repeated int64 groupIDs = 55;
  string groupIDPath = 56;
  repeated string groupIDPaths = 57;
  int64 notGroupID = 19;
  int64 parentGroupID = 42;
  string groupPurpose = 45;             // 设备分组用途 不填默认为default
  string groupName = 43;                // 模糊查询
  int64 notAreaID = 29;
  IDPath distributor = 20;              // 过滤分销商的数据
  int64 status = 22;
  repeated int64 statuses = 35;
  CompareInt64 ratedPower = 28;         // 额定功率:单位w/h
  int64 hasOwner = 30;                  // 是否被人拥有
  int64 userID = 31;                    // 用户id查询
  CompareInt64 expTime = 38;            // 到期时间
  string iccid = 40;                    // SIM卡卡号
  bool withGateway = 41;                // 同时返回子设备的网关
  string protocolCode = 51;             // 协议查询
  CompareInt64 rssi = 52;
  IDPath dept = 53;                     // 过滤分销商的数据
  map<string, CompareString> property = 54; // 设备最新属性过滤,key为属性的id,如果是结构体则key为 aaa.bbb 数组为aaa.1
}
```

### 分页响应规范
```protobuf
message DeviceInfoIndexResp {
  repeated DeviceInfo list = 1;         // 设备信息
  int64 total = 2;                      // 总数(只有分页的时候会返回)
}
```

## 批量操作规范

### 批量请求规范
```protobuf
message DeviceInfoMultiUpdateReq {
  repeated DeviceCore devices = 1;      // 设备列表
  CompareInt64 FilterDistributorID = 2; // 过滤分销商ID
  int64 areaID = 4;                     // 项目区域id
  IDPath distributor = 20;              // 分销商的数据
  int64 ratedPower = 28;                // 额定功率:单位w/h
}
```

### 批量响应规范
```protobuf
message DeviceInfoMultiBindResp {
  repeated DeviceError errs = 1;        // 错误列表
}

message DeviceError {
  string productID = 1;                 // 产品id
  string deviceName = 2;                // 设备名称
  int64 code = 3;                       // 错误码
  string msg = 4;                       // 错误信息
}
```

## 导入导出规范

### 导入请求规范
```protobuf
message ProductInfoImportReq {
  string products = 2;                  // 产品导出的信息
}

message ImportResp {
  int64 total = 1;                      // 导入总接口数
  int64 errCount = 2;                   // 失败数
  int64 ignoreCount = 3;                // 忽略数
  int64 succCount = 4;                  // 成功数
}
```

### 导出请求规范
```protobuf
message ProductInfoExportReq {
  repeated string productIDs = 1;       // 产品ID列表
}

message ProductInfoExportResp {
  string products = 2;                  // 产品导出的信息
}
```

## 特殊业务类型规范

### 设备交互规范
```protobuf
message ActionSendReq {
  string productID = 1;                 // 产品id 获取产品id下的所有设备信息
  string deviceName = 2;                // 设备名
  string actionID = 3;                  // 产品数据模板中行为功能的标识符，由开发者自行根据设备的应用场景定义
  string inputParams = 4;               // 输入参数
  bool isAsync = 5;                     // 是否异步获取
  SendOption option = 6;                // 异步选项
}

message ActionSendResp {
  string msgToken = 1;                  // 调用id
  string outputParams = 2;              // 输出参数 注意：此字段可能返回 null，表示取不到有效值。
  string msg = 3;                       // 返回状态
  int64 code = 4;                       // 设备返回状态码
}
```

### 消息发布规范
```protobuf
message PublishMsg {
  string handle = 1;                    // 对应 mqtt topic的第一个 thing ota config 等等
  string type = 2;                      // 操作类型 从topic中提取 物模型下就是 property属性 event事件 action行为
  bytes payload = 3;                    // 消息内容
  int64 timestamp = 4;                  // 毫秒时间戳
  string productID = 5;                 // 产品ID
  string deviceName = 6;                // 设备名称
  string explain = 7;                   // 内部使用的拓展字段
  string protocolCode = 8;              // 如果有该字段则回复的时候也会带上该字段
}
```

## 生成命令

在 `service/README.md` 中已给出标准命令：

```bash
cd dmsvr && goctl rpc protoc proto/dm.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -m && cd ..
```

## 校验清单（生成前自检）

- 已设置正确的 `syntax`、`option go_package`、`package`
- 已导入必要的依赖包
- 服务名使用帕斯卡命名法
- 方法名使用小驼峰命名法
- 消息名使用帕斯卡命名法
- 字段编号连续且不重复
- 所有字段都有中文注释
- 分页请求包含 `PageInfo` 字段
- 分页响应包含 `total` 字段
- 批量操作使用 `repeated` 字段
- 可选字段使用 `google.protobuf.*Value` 类型
- 时间字段说明时间戳格式

## 约束与不做事项

- 不在 `.proto` 中编写业务说明文档，长描述仅放在字段注释中
- 不在 `.proto` 中混入具体实现或示例数据，仅保留类型与接口定义
- 不使用过时的 protobuf v2 语法
- 不重复使用已删除的字段编号

---

如需新增服务：
1) 在 `proto/` 目录新建 `xxx.proto`，按本规范填写服务定义和消息类型
2) 执行"生成命令"一节中的命令
3) 在对应的服务中实现业务逻辑
