create database things;
use things;

CREATE TABLE `device_info` (
  `deviceID` bigint unsigned NOT NULL COMMENT '设备id',
  `productID` bigint NOT NULL COMMENT '产品id',
  `deviceName` varchar(100) NOT NULL COMMENT '设备名称',
  `secret` varchar(50) DEFAULT '' COMMENT '设备秘钥',
  `firstLogin` datetime DEFAULT NULL COMMENT '激活时间',
  `lastLogin` datetime DEFAULT NULL COMMENT '最后上线时间',
  `createdTime` datetime NOT NULL,
  `updatedTime` datetime DEFAULT NULL,
  `deletedTime` datetime DEFAULT NULL,
  PRIMARY KEY (`deviceID`),
  UNIQUE KEY `deviceName` (`productID`,`deviceName`),
  KEY `idx_product_id` (`productID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备信息表';


CREATE TABLE `product_info` (
  `productID` bigint unsigned NOT NULL COMMENT '产品id',
  `productName` varchar(100) NOT NULL COMMENT '产品名称',
  `authMode` int unsigned DEFAULT '0' COMMENT '认证方式:0:账密认证,1:秘钥认证',
  `deviceType` int unsigned DEFAULT '0' COMMENT '设备类型:0:设备,1:网关,2:子设备',
  `categoryID` int unsigned DEFAULT '0' COMMENT '产品品类',
  `netType` int unsigned DEFAULT '0' COMMENT '通讯方式:0:其他,1:wi-fi,2:2G/3G/4G,3:5G,4:BLE,5:LoRaWAN',
  `dataProto` int unsigned DEFAULT '0' COMMENT '数据协议:0:自定义,1:数据模板',
  `autoRegister` int unsigned DEFAULT '0' COMMENT '动态注册:0:关闭,1:打开,2:打开并自动创建设备',
  `secret` varchar(50) DEFAULT '' COMMENT '动态注册产品秘钥',
  `description` varchar(200) DEFAULT '' COMMENT '描述',
  `createdTime` datetime NOT NULL,
  `updatedTime` datetime DEFAULT NULL,
  `deletedTime` datetime DEFAULT NULL,
  PRIMARY KEY (`productID`),
  UNIQUE KEY `productName` (`productName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='产品信息表';