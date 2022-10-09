# 物联网云平台ithings

## 介绍


iThings是一个基于golang开发的轻量级云原生微服务物联网平台.  
  
定位于:
* 云原生高性能 -- 使用golang编写,使用gozero微服务框架高性能的同时保证服务的稳定性
* 易拓展 -- 原生支持微服务部署,避免只支持集群模式后期难以拓展的尴尬
* 易部署 -- 一键安装所有依赖,一键运行iThings

## 架构
下图是 iThings 平台的整体架构:  
<img src="./doc/assets/iThings架构图.png">

#### 目录说明

- doc:该项目的文档都放在这里
- shared:所有该项目及其他项目所公用的代码都放在这里
- src:存放了所有服务的源码


## 特征

### 已完成
1. 物模型管理及日志记录
2. 设备本地日志
3. 设备云端调试日志
4. 在线设备调试
5. 产品管理
6. 设备管理及认证
7. 微服务和单体模式
8. 独立管理平台

### 待实现
1. 固件升级
2. 规则引擎
3. 网关及子设备
4. 大屏
5. 设备配置
6. 设备影子


## 安装与运行
在ithings中依赖tdengine,mysql,redis,etcd,nats,emqx
* `sudo ./init.sh`即会安装docker及docker-compose及第三方依赖及初始化数据库脚本(一定是root权限,不然可能会有问题)
* 然后 `./run.sh` 即可运行iThings所有服务

## 文档

- 开发文档: [https://ithings.pages.dev/](https://ithings.pages.dev/)
- 用户文档: [https://ithings.pages.dev/](https://ithings.pages.dev/)


## 贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request

## 社区
- 官网:[https://ithings.pages.dev/](https://ithings.pages.dev/)
- 微信: `godLei6` (需备注“来自github”)
<img src="https://ithings.pages.dev/assets/img/things/%E5%BE%AE%E4%BF%A1%E4%BA%8C%E7%BB%B4%E7%A0%812.jpg">

## 收藏
<img src="https://starchart.cc/i4de/ithings.svg">
