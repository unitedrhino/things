
# 产品概述
**联犀** 是一款基于 Go 语言开发的商业级 SaaS 云原生微服务物联网平台，致力于帮助企业快速构建自己的物联网应用，实现数字化转型。  
[文档](http://doc.ithings.net.cn/)
## 技术优势
- **高性能**：采用 Go 语言和高性能组件（如 emqx, nats, tdengine），结合 gozero 微服务框架，确保服务的稳定性与性能。
- **云原生可扩展**：一套代码支持 k8s, docker, 微服务及单体模式部署，便于业务架构的灵活切换。
- **易部署**：一键安装和运行，简化部署流程。
- **易集成**：支持通过 HTTP, gRPC 或直接集成，同时兼容 MySQL 和 PostgreSQL，实现无缝集成。
    
## 开源社区
- **GitHub**: [联犀 GitHub](https://github.com/unitedrhino/things)
- **Gitee**: [联犀 Gitee](https://gitee.com/unitedrhino/things)

# 产品架构
联犀物联网平台作为物联网架构中的关键中介，实现设备与应用层的高效联动。它不仅向下管理物联网设备，完成数据的收集与存储，而且向上为开发者和集成商提供统一的数据接口和工具，加速物联网解决方案的开发和部署。

通过 HTTP, gRPC 应用厂商可以快速将 联犀 集成到自己的系统中，实现轻量级且高效的物联网能力扩展。
<img style="width:300px;" src="./doc/assets/部署架构图.png">

## 产品价值

| 平台价值   | 描述                                        |
|--------|-------------------------------------------|
| 拓展能力强  | 支持单体和微服务架构，便于开发者在不同发展阶段灵活切换，无需维护两套代码。     |
| 高性能    | 使用 golang 开发，依赖的第三方服务少，适应多种性能要求，可以快速水平拓展。 |
| 数据价值   | 私有化部署，数据自主管理，无需担心公有云服务中断或成本问题。            |
| 解决方案底座 | 作为行业解决方案的数字底座，支持多行业共用物联网平台，沉淀行业经验和产品方案。   |

## 产品特性

- **设备接入**：支持 MQTT、CoAP 和 HTTP 等物联网协议，实现海量设备连接，同时支持协议网关，兼容任何协议。
- **远程控制**：通过 HTTP API 实现服务器对设备的精准控制和设备主动通知。
- **物模型**：支持标准物模型，有效管理设备属性、事件及行为。
- **RBAC权限**：采用基于角色的访问控制（RBAC），提供完善的用户、角色、菜单权限管理。
- **多租户多项目多应用**：支持低成本开发应用，便于多企业共享使用。
- **应用支撑**：提供 HTTP, gRPC 接口，简化物联网解决方案开发，缩短上市周期，节省研发时间和成本。
- **自主可控**：支持私有云、公有云、边缘部署等多种部署方式。
- **快速开发及维护**：联犀 通过简化的接入流程和模块化开发，优化了物联网平台的开发体验。它提供了商业级小程序和 App 模板，允许快速上线，同时支持多租户架构以降低维护成本，并具备灵活的扩展能力以应对设备数量增长。

### 物联网模块

## 技术栈

### 后端
1. 微服务框架：[go-zero](https://go-zero.dev/)
2. 高性能缓存：[redis](https://redis.io/)
3. 高性能消息队列：[nats](https://docs.nats.io/)
4. 关系型数据库：[mysql (推荐使用 MariaDB 或 MySQL 5.7)](https://mariadb.com/) 或 pgsql，未来将支持更多数据库
5. 微服务注册中心（单体可不使用）：etcd
6. 云原生轻量级对象存储：[minio](https://min.io/)
7. 开源、高性能、云原生时序数据处理平台：[tdengine](https://www.taosdata.com/)
8. 大规模可弹性伸缩的云原生分布式物联网 MQTT 消息服务器：[emqx](https://docs.emqx.com/zh/emqx/latest/)

### 前端
1. 渐进式 JavaScript 框架：[vue](https://cn.vuejs.org/)
2. 企业级设计组件：[ant design](https://antdv.com/docs/vue/introduce-cn/)

### 小程序
1. [uniapp vue3](https://uniapp.dcloud.net.cn/)

### app(安卓, iOS, 鸿蒙)
1. [uniapp x](https://doc.dcloud.net.cn/uni-app-x/)
## 贡献者

感谢所有已经做出贡献的人!

### 后端

<a href="https://github.com/unitedrhino/things/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=unitedrhino/things" />
</a>


## 社区

- 微信(加我拉微信群): `godLei6` (需备注“来自github”)
- [官网](http://doc.ithings.net.cn/)
- 微信二维码
- <img style="width: 300px;" src="./doc/assets/微信二维码.jpg">

## 收藏

<img src="https://starchart.cc/unitedrhino/things.svg">
