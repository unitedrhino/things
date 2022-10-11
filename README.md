# 物联网云平台ithings

iThings是一个基于golang开发的轻量级云原生微服务物联网平台.

定位于:
* 高性能 -- 使用golang编写,基于gozero微服务框架高性能的同时保证服务的稳定性
* 云原生易拓展 -- 支持k8s,docker微服务及单体模式部署,便于业务从单体到微服务的切换
* 易部署 -- 一键安装所有依赖,一键运行iThings

git地址:
* 后端github地址:[https://github.com/i4de/ithings](https://github.com/i4de/ithings)
* 后端gitee地址:[https://gitee.com/i-Things/things](https://gitee.com/i-Things/things)
* 前端github地址:[https://github.com/i4de/iThings-admin-react](https://github.com/i4de/iThings-admin-react)

## 产品架构
iThings物联网平台提供跨不同设备和数据源的通用接入及管理平台,在整个物联网架构中起到承上启下的中介作用，
联动感知层及应用层之间的所有交互——向下连接、管理物联网设备端并完成感知数据的归集与存储，
向上为应用开发商与系统集成商提供应用开发的统一数据接口及共性模块工具。  
应用厂商可以通过http,grpc,及golang包引入的方式快速集成进自己的系统中,迅速获取物联网平台的能力,轻量而高效.

<img src="./doc/assets/iThings架构图.png">

## 产品价值

| 平台价值   | 描述                                                                                  |
|--------|-------------------------------------------------------------------------------------|
| 拓展能力强  | 公司发展有不同的阶段,应用场景也有不同的需要,iThings同时支持单体和微服务便于开发者不需要维护两套代码同时单体不足以满足公司需要时也可以很便捷的切换为微服务架构 |
| 高性能    | 使用golang开发,依赖的第三方服务很少,对于性能要求不高的场景甚至可以部署到低性能的arm上                                    |
| 数据价值   | 私有化部署,所有数据都可以自己去管理,不用担心公有云停服及收费高的问题                                                 |
| 解决方案底座 | 可做为行业解决方案的数字底座，支持多行业解决方案共用同一个物联网平台底座，为企业沉淀更多行业经验和产品方案。                              |

## 产品特性

* 设备接入：通过行业标准物联网协议（MQTT、CoAP 和 HTTP）实现海量设备连接,也可以通过规则引擎接入自定义协议。
* 远程控制：使用http api可以实现服务器对设备的精准控制和设备主动向服务器发送请求通知。
* 物模型：支持标准物模型,有效管理设备的属性,事件及行为。
* RBAC权限：权限控制采用 RBAC，基于角色的访问控制。封装完善的用户、角色、菜单等基础权限功能。
* 应用支撑：提供http,grpc或golang直接导入包的形式将数据流转至业务层面，屏蔽了物联网层面的复杂特性，让企业专注业务应用开发，缩短物联网解决方案的上市周期，为企业节省大量研发时间和成本。
* 自主可控：私有云、公有云、边缘部署.
* 快速开发及维护: 提供设备云端调试日志,设备本地日志,在线设备调试 快速开发及调试问题

## 产品功能

| 功能   | 描述                                        |
|------|-------------------------------------------|
| 后台管理 | 独立后台,提供用户管理,角色管理,菜单管理                     |
| 设备管理 | 包含设备信息维护、数据收集等基础功能，支持设备参数配置、模型配置、在线调试等功能。 |
| 产品管理 | 支持产品管理,物模型管理,配置管理                         |
| 设备分组 | 支持设备最多三级分组                                |

## 技术栈
### 后端
1. 微服务框架:[go-zero](https://go-zero.dev/)
2. 高性能缓存:[redis](https://redis.io/)
3. 高性能消息队列:[nats](https://docs.nats.io/)
4. 数据库:mysql
5. 微服务注册中心(单体可不使用):etcd
6. 云原生轻量级对象存储:[minio](https://min.io/)
7. 开源、高性能、云原生,极简的时序数据处理平台:[tdengine](https://www.taosdata.com/)
8. 大规模可弹性伸缩的云原生分布式物联网MQTT消息服务器:[emqx](https://docs.emqx.com/zh/enterprise/v4.4/getting-started/install-ee.html)
### 前端
1. 用于构建用户界面的 JavaScript 库:[react](https://reactjs.org/)
2. 企业级设计组件:[ant design](https://ant.design/)
3. 一套高质量可靠的 React Hooks库:[ahooks](https://ahooks.js.org/)

## 贡献者
感谢所有已经做出贡献的人!
### 后端
<a href="https://github.com/i4de/iThings/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=i4de/iThings" />
</a>

### 前端
<a href="https://github.com/i4de/iThings-admin-react/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=i4de/iThings-admin-react" />
</a>

## 社区
- 微信: `godLei6` (需备注“来自github”)
- 官网:[https://ithings.pages.dev/](https://ithings.pages.dev/)
- 微信二维码:
- <img src="./doc/assets/微信二维码.jpg">
## 收藏
<img src="https://starchart.cc/i4de/ithings.svg">
