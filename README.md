# 物联网云平台go-things

#### 介绍


依照腾讯云物联网平台打造,无缝迁移,基于go-zero框架,目前支持mqtt协议,支持数据模板检验,日志记录,数据流转,实时数据反馈,低代码数据流转,用户及设备互联互通.使用了kafka,redis,mysql,MongoDB,etcd开源框架.  
git地址:[https://gitee.com/godLei6/things](https://gitee.com/godLei6/things)  
我的微信:17052709767  
欢迎大家的加入  
![微信二维码](https://gitee.com/godLei6/things/raw/master/doc/assets/%E5%BE%AE%E4%BF%A1%E4%BA%8C%E7%BB%B4%E7%A0%81.jpg)



#### 软件架构

软件架构说明  
 ![go-things架构图.jpg](https://gitee.com/godLei6/things/raw/master/doc/assets/go-things%E6%9E%B6%E6%9E%84%E5%9B%BE.jpg)  
设备接入流程图  
![设备连接流程图.jpg](https://gitee.com/godLei6/things/raw/master/doc/assets/%E8%AE%BE%E5%A4%87%E8%BF%9E%E6%8E%A5%E6%B5%81%E7%A8%8B%E5%9B%BE.jpg)


#### 目录说明

- doc:该项目的文档都放在这里
- shared:所有该项目及其他项目所公用的代码都放在这里
- src:存放了所有服务的源码

#### 文档
介绍及说明文档:[https://www.yuque.com/gothings/umcf39](https://www.yuque.com/gothings/umcf39)
​

#### 安装教程
##### 环境依赖安装
在go-things中依赖mongodb,mysql,redis,etcd,kafka,zookeeper
* 在初始目录中提供了docker-compose文件,如果安装好了docker及docker-compose可以直接
docker-compose up 即可更新
* 如果都没有安装则sudo ./init.sh即会安装docker及docker-compose及第三方依赖
* 然后将db中的sql导入mysql中即可

##### 服务运行
1. 进入src目录进入对应的服务
2. 修改etc目录下的配置文件将对应的依赖改为本地的ip地址
3. 直接go build即可享受

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


