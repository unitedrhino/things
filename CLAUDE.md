# Things 仓库

> IoT 开源模块，负责设备、产品、物模型、协议、OTA 等能力。

## 仓库信息

| 项目 | 值 |
|------|-----|
| 远程 | `git@codeup.aliyun.com:642f7dca8b47795dae985084/ee/things.git` |
| 当前分支 | `dev`（最新） |
| 模块 | `gitee.com/unitedrhino/things` |
| Go 版本 | 1.24.4 |

## 关键目录

```
things/
├── service/
│   ├── apisvr/          # IoT HTTP 网关
│   ├── dmsvr/           # 设备管理 RPC（产品、设备、物模型、OTA）
│   ├── dgsvr/           # 设备网关 RPC（鉴权、消息路由、MQTT）
│   └── mqttsvr/         # MQTT 服务
└── Makefile
```

## 常用操作

```bash
# 编译全部
cd things && make build

# 仅编译 API 网关
cd things && make build.api
```

## 依赖升级

升级 share / core 依赖后执行 `go mod tidy`，然后打 tag：
```bash
git tag v1.5.x
git push origin v1.5.x
git push gitee v1.5.x
```

## README 媒体兼容

- ⚠️ Gitee 会将仓库内 SVG 以 `text/plain` 和 `nosniff` 返回，README 中不得用 SVG 承载必须直接展示的动画；优先使用动画 WebP 或 GIF，SVG 只作为源文件保留
- 更新 README 媒体后，必须分别检查 GitHub、Gitee 的资源响应类型，并在仓库页面确认图片已完成加载且自然尺寸不为零
