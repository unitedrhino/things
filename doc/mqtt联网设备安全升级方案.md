# MQTT 联网设备安全升级方案

## 1. 方案目标

当前设备已经可以通过 MQTT 接入云端，但存在几个安全问题：

- 设备连接使用普通 `mqtt://1883`，链路没有加密。
- 下行控制指令是明文 JSON，设备没有逐条验证指令签名。
- 场景联动、后台控制、App 控制最终都会下发 MQTT 指令，如果没有统一安全出口，后续容易出现绕过。
- 设备密钥不能明文保存到业务数据库，也不能出现在场景配置、日志、前端或普通接口参数中。

本方案目标：

- 设备连接升级为 `mqtts://8883`。
- 每台设备使用独立证书或独立密钥。
- MQTT Topic 做权限隔离，设备只能访问自己的 Topic。
- 所有下行控制消息都带签名和防重放字段。
- 场景联动不保存密码、不保存密钥，只保存业务动作。
- 设备控制密钥从第一版开始进入密钥托管服务，例如 KMS 或 OpenBao。
- 老设备通过 OTA 灰度迁移，不能一次切断旧通道。

## 2. 总体设计

整体链路如下：

```text
App 手动控制 / 场景联动 / 后台 API / 定时任务
        |
        v
云端权限校验
        |
        v
统一设备控制服务 DeviceInteract
        |
        v
安全签名服务生成签名消息
        |
        v
MQTT 服务通过 8883 下发
        |
        v
设备验签、防重放、执行动作
```

关键原则：

- 云端可以代用户执行场景联动，但必须经过权限校验和统一签名出口。
- 设备不直接相信 MQTT 消息，必须验签后才执行。
- 场景联动不是安全边界，安全边界在统一设备控制服务和设备端验签。
- 业务服务不直接保存或长期持有明文设备密钥。

## 3. 不采用设备密码方案

不采用以下设计：

- 场景联动动作里保存 `devicePassword`。
- 控制接口里传递 `devicePassword`。
- 使用默认密码 `000000`。
- 使用 `HMAC-MD5(JSON(params)+msgToken, devicePassword)`。
- 让云端通过用户输入密码生成设备控制签名。

原因：

- 密码会扩散到场景表、前端、接口、日志和调试工具里。
- 用户密码不是设备身份密钥，不适合作为设备控制签名根。
- 默认密码和 HMAC-MD5 不符合消费 IoT 安全基线。
- 密码变更后，旧场景、旧自动化、旧分享关系都会变得难维护。

正确做法：

```text
场景只保存业务动作
设备密钥由密钥托管服务保存
统一设备控制服务负责请求签名
设备负责验签
```

## 4. 传输层升级

### 4.1 当前状态

当前 MQTT 配置仍是普通 TCP：

```yaml
Mqtt:
  TCP: :1883
  WS: :1882
```

设备当前连接形式：

```text
mqtt://iot.ykhl.vip:1883
```

### 4.2 目标状态

目标连接形式：

```text
mqtts://iot.ykhl.vip:8883
```

升级后要求：

- MQTT 服务开启 8883 TLS 端口。
- 设备校验云端服务器证书。
- 设备逐步启用客户端证书。
- 云端校验设备客户端证书。
- 证书中绑定设备身份，例如：

```text
productID=15
deviceName=FC012CD085E8
```

### 4.3 老设备证书获取

老设备不能要求一开始就有客户端证书。

迁移路径：

```text
旧设备先通过 1883 在线
    |
    v
OTA 升级过渡固件
    |
    v
设备支持 8883 单向 TLS
    |
    v
设备本地生成私钥和 CSR
    |
    v
云端确认设备身份后签发客户端证书
    |
    v
设备保存证书
    |
    v
后续优先使用双向 TLS
```

要求：

- 私钥必须在设备端生成。
- 云端不能下发私钥。
- 证书可吊销、可轮换。

## 5. MQTT Topic 权限

设备只能访问自己的 Topic。

设备允许上报：

```text
$thing/up/property/{productID}/{deviceName}
$thing/up/event/{productID}/{deviceName}
$thing/up/action/{productID}/{deviceName}
$ota/up/report/{productID}/{deviceName}
```

设备允许订阅：

```text
$thing/down/property/{productID}/{deviceName}
$thing/down/action/{productID}/{deviceName}
$ota/down/upgrade/{productID}/{deviceName}
```

设备禁止：

```text
publish $thing/down/...
subscribe 其他设备的 topic
publish 其他设备的上行 topic
使用 # 通配符
使用 + 通配符
```

云端服务账号可以发布下行 Topic，但必须先经过业务权限校验和消息签名。

## 6. 消息签名协议

### 6.1 下行消息格式

所有下行控制消息统一使用安全消息格式：

```json
{
  "version": "secure-msg-v1",
  "method": "control",
  "msgToken": "019f7c7a7c9b",
  "seq": 1718170001,
  "ts": 1718170001,
  "params": {
    "power_switch": 1
  },
  "signAlg": "hmac-sha256",
  "sign": "base64url_no_padding"
}
```

行为调用示例：

```json
{
  "version": "secure-msg-v1",
  "method": "action",
  "msgToken": "019f7c7a7c9c",
  "seq": 1718170002,
  "ts": 1718170002,
  "actionID": "restart",
  "params": {},
  "signAlg": "hmac-sha256",
  "sign": "base64url_no_padding"
}
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `version` | 安全消息版本 |
| `method` | 控制类型，例如属性控制或行为调用 |
| `msgToken` | 本次请求编号，用于匹配响应 |
| `seq` | 单调递增序号，用于防重放 |
| `ts` | 云端生成消息的时间 |
| `params` | 控制参数 |
| `actionID` | 行为调用标识，仅行为调用需要 |
| `signAlg` | 签名算法 |
| `sign` | 消息签名 |

### 6.2 签名密钥

设备注册后会有设备密钥，例如 `device_secret`。

设备密钥不直接用于所有场景，而是派生出方向密钥：

```text
K_down = HKDF-SHA256(device_secret, "mqtt-down-sign-v1|" + productID + "|" + deviceName)
K_up   = HKDF-SHA256(device_secret, "mqtt-up-sign-v1|" + productID + "|" + deviceName)
```

说明：

- 下行控制使用 `K_down`。
- 上行状态和响应使用 `K_up`。
- 上下行密钥分开，避免互相复用。

### 6.3 签名内容

签名覆盖这些内容：

```text
productID
deviceName
version
method
msgToken
seq
ts
paramsHash
actionID（行为调用时）
```

其中：

```text
paramsHash = sha256(canonical_json(params))
```

这样可以保证：

- 参数被改，签名失效。
- 设备名被改，签名失效。
- 消息序号被改，签名失效。
- 旧消息重放，设备拒绝执行。

### 6.4 设备验签规则

设备收到下行消息后按顺序检查：

1. Topic 中的 `productID/deviceName` 必须等于本机身份。
2. `version` 必须是设备支持的版本。
3. `signAlg` 必须是设备支持的算法。
4. `seq` 必须大于设备保存的 `last_seq`。
5. `seq` 不能一次跳太大，避免恶意锁死设备。
6. 设备重新计算签名。
7. 签名一致才执行。
8. 验签通过后，先保存新的 `last_seq`，再执行动作。

失败处理：

- 不执行动作。
- 返回统一错误码。
- 不打印设备密钥、派生密钥、签名输入全量、MQTT 密码。

### 6.5 上行消息签名

建议上行状态、事件、控制响应也带签名。

上行示例：

```json
{
  "version": "secure-msg-v1",
  "method": "property.report",
  "msgToken": "019f7c7a7c9d",
  "seq": 1718171001,
  "ts": 1718171001,
  "params": {
    "power_switch": 1
  },
  "signAlg": "hmac-sha256",
  "sign": "base64url_no_padding"
}
```

云端验上行签名，并记录每台设备的上行序号，拒绝旧状态回放。

## 7. 场景联动安全机制

### 7.1 场景只保存业务动作

场景动作只保存类似内容：

```json
{
  "productID": "15",
  "deviceName": "FC012CD085E8",
  "type": "property",
  "dataID": "power_switch",
  "value": "1"
}
```

禁止保存：

```text
devicePassword
psk
deviceSecret
signKey
clientKey
```

### 7.2 场景执行流程

```text
场景触发
    |
    v
检查场景是否启用
    |
    v
检查租户、项目、设备归属
    |
    v
检查场景是否仍有权限控制目标设备
    |
    v
调用 PropertyControlSend / ActionSend
    |
    v
统一设备控制服务校验物模型
    |
    v
安全签名服务生成签名消息
    |
    v
MQTT 下发
    |
    v
设备验签并执行
```

### 7.3 为什么场景联动安全

场景联动安全不是因为场景自己保存了密码，而是因为它不能绕过统一下发出口。

安全点：

- 场景不保存密钥，场景表泄露也不能伪造设备指令。
- 场景不生成签名，签名逻辑只有一份。
- App 控制、场景联动、后台控制都走同一个设备控制服务。
- 设备最终验签，不关心消息来自 App、场景还是后台。
- `seq` 防止别人抓包后重复执行旧指令。
- Topic 权限防止普通客户端伪造云端下行。

## 8. 设备密钥托管

### 8.1 基本要求

设备控制密钥必须进入密钥托管服务，例如 KMS 或 OpenBao。

要求：

- 业务数据库不保存明文 `device_secret`。
- 场景配置不保存明文 `device_secret`。
- API 层不返回明文 `device_secret`。
- 日志不打印明文 `device_secret`。
- 普通业务服务不长期持有明文 `device_secret`。

### 8.2 设备安全状态表

新增或扩展设备安全状态表：

```text
device_security_state
- product_id
- device_name
- key_ref
- key_version
- secret_version
- security_level
- last_seen_protocol
- cert_status
- cert_serial
- last_seq_down
- last_seq_up
- allow_1883_until
- migration_state
- updated_at
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `key_ref` | 密钥托管服务中的密钥引用 |
| `key_version` | 设备密钥版本 |
| `secret_version` | 设备侧密钥材料版本 |
| `security_level` | 当前安全等级 |
| `last_seen_protocol` | 最近一次连接协议 |
| `cert_status` | 设备证书状态 |
| `last_seq_down` | 最近下行序号 |
| `last_seq_up` | 最近上行序号 |
| `allow_1883_until` | 允许旧通道的截止时间 |
| `migration_state` | 迁移状态 |

密钥引用示例：

```text
device/15/FC012CD085E8/down/v1
device/15/FC012CD085E8/up/v1
```

### 8.3 安全签名服务

统一提供两个能力：

```text
SignDownlink：生成下行消息签名
VerifyUplink：校验上行消息签名
```

安全签名服务职责：

- 根据 `productID/deviceName` 找到设备密钥引用。
- 从密钥托管服务获取代签能力或短暂使用密钥。
- 派生 `K_down` / `K_up`。
- 生成或验证签名。
- 记录审计日志。

审计日志记录：

```text
调用来源
productID
deviceName
method
msgToken
结果
时间
```

审计日志不能记录：

```text
device_secret
K_down
K_up
完整签名输入
MQTT password
```

## 9. 老设备迁移计划

### 9.1 设备状态盘点

为每台设备记录：

```text
firmwareVersion
lastSeenProtocol: mqtt1883 | mqtts8883 | mtls8883
securityLevel: legacy | tls_psk | mtls | signed
certStatus: none | csr_created | issued | active | revoked
lastSeqDown
lastSeqUp
allow1883Until
migrationState
```

### 9.2 过渡固件

过渡固件必须支持：

- 1883 旧连接。
- 8883 单向 TLS。
- 服务端证书校验。
- 安全消息解析和验签。
- A/B OTA。
- 固件签名校验。
- 升级失败回滚。
- 设备本地生成 CSR。
- 客户端证书保存。

### 9.3 灰度顺序

建议按以下顺序：

```text
10 台测试设备
1%
5%
20%
50%
100%
```

每批观察：

- OTA 成功率。
- 升级后在线率。
- TLS 连接成功率。
- 回退 1883 比例。
- 控制成功率。
- 验签失败率。
- 设备重启率。

### 9.4 证书迁移

流程：

```text
设备通过 8883 单向 TLS 上线
设备生成私钥和 CSR
云端校验原设备身份
云端签发客户端证书
设备保存证书
下次连接使用双向 TLS
云端标记证书状态为 active
```

### 9.5 强制安全

最终策略：

- 新固件设备必须使用安全消息格式。
- 双向 TLS 设备不允许回退 1883，除非后台临时白名单放开。
- 1883 只保留迁移白名单。
- 最终关闭公网 1883。

## 10. OTA 安全

OTA 不能只依赖 TLS。

必须具备：

- 固件 hash 校验。
- 固件签名校验。
- A/B 分区。
- 新固件启动后成功连云端才标记有效。
- 失败自动回滚。
- OTA 下发消息也必须带签名。

设备 OTA 流程：

```text
1. 验 OTA 下发消息签名。
2. 下载固件。
3. 校验 hash。
4. 校验固件签名。
5. 写入备用分区。
6. 重启进入新固件。
7. 新固件成功连接云端并上报版本。
8. 标记新固件有效。
```

## 11. 代码落点

当前 `/Users/ykhl/Desktop/oldThings` 中，重点落点：

```text
service/dmsvr/internal/logic/deviceinteract/propertyControlSendLogic.go
service/dmsvr/internal/logic/deviceinteract/actionSendLogic.go
service/dmsvr/proto/dm.proto
service/dgsvr/internal/repo/event/publish/pubDev/mqtt.go
share/clients/mqtt.go
service/apisvr/etc/mqtt.yaml
deploy/docker/conf/things/etc/mqtt.yaml
```

建议改造：

- 在 `PropertyControlSendLogic` 中，组好控制消息后调用安全签名服务。
- 在 `ActionSendLogic` 中，组好行为调用消息后调用安全签名服务。
- 不在 `PropertyControlSendReq` / `ActionSendReq` 中新增 `devicePassword`。
- 不改场景模型去保存密码。
- MQTT 客户端支持 TLS 配置。
- MQTT 发布日志不再打印完整 payload。
- 新增设备安全状态记录。
- 新增安全签名服务。

## 12. 验收标准

### 12.1 功能验收

- 旧设备在兼容模式下仍能控制。
- 新设备能收到安全消息并执行。
- 新设备能拒绝旧明文控制消息。
- 场景联动无需保存密码，仍能正常触发控制。
- App 手动控制、场景联动、后台 API 都进入同一签名出口。

### 12.2 安全验收

- 抓包重放旧控制消息，设备拒绝执行。
- 篡改 `params`，设备拒绝执行。
- 篡改 `seq`，设备拒绝执行。
- 篡改 Topic 中的 `deviceName`，设备拒绝执行。
- 设备不能 publish 到下行 Topic。
- 设备不能 subscribe 其他设备 Topic。
- 业务数据库看不到明文设备密钥。
- 日志中没有设备密钥、MQTT 密码、派生密钥、密钥托管服务访问凭据。

### 12.3 迁移验收

- 10 台测试设备 OTA 成功。
- 8883 单向 TLS 连接成功。
- 设备 CSR 申请证书成功。
- 双向 TLS 连接成功。
- 断网或证书失败时按策略回退，不造成批量离线。
- 新固件未成功连云端时自动回滚。

## 13. 后续增强

- 设备证书吊销。
- 设备密钥轮换。
- 高风险动作增加二次确认。
- 关闭公网 1883。
- 对敏感业务参数做应用层加密。
