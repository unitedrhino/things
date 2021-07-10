# 物联网平台

#### 介绍
依照腾讯云物联网平台打造,无缝迁移,基于go-zero框架,目前支持mqtt协议,支持数据模板检验,日志记录,数据流转,实时数据反馈,低代码数据流转,用户及设备互联互通.使用了kafka,redis,mysql,MongoDB,etcd开源框架
#### 软件架构
软件架构说明
![avatar](./doc/assets/架构图.jpg)

#### 目录说明
* doc:该项目的文档都放在这里
* shared:所有该项目及其他项目所公用的代码都放在这里
* src:存放了所有服务的源码
#### 服务介绍
所有服务都在src目录下,基于gozero进行开发
* webapi:http网关服务
* usersvr:用户服务,提供注册登录用户管理等接口
* dmsvr: 提供设备的管理,登入登出,日志记录等综合的服务

#### 安装教程
git clone 下来后进入src目录进入对应的服务直接go build即可


#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


