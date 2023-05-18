-- 设备管理模块SQL
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

create database if not EXISTS iThings;
use iThings;

CREATE TABLE if not exists `dm_product_info`
(
    `productID`    char(11)     NOT NULL COMMENT '产品id',
    `productName`  varchar(100) NOT NULL COMMENT '产品名称',
    `productImg`   varchar(200) DEFAULT '' COMMENT '产品图片',
    `productType`  tinyint(1)            DEFAULT '1' COMMENT '产品状态:1:开发中,2:审核中,3:已发布',
    `authMode`     tinyint(1)            DEFAULT '1' COMMENT '认证方式:1:账密认证,2:秘钥认证',
    `deviceType`   tinyint(1)            DEFAULT '1' COMMENT '设备类型:1:设备,2:网关,3:子设备',
    `categoryID`   int(10)               DEFAULT '1' COMMENT '产品品类',
    `netType`      tinyint(1)            DEFAULT '1' COMMENT '通讯方式:1:其他,2:wi-fi,3:2G/3G/4G,4:5G,5:BLE,6:LoRaWAN',
    `dataProto`    tinyint(1)            DEFAULT '1' COMMENT '数据协议:1:自定义,2:数据模板',
    `autoRegister` tinyint(1)            DEFAULT '1' COMMENT '动态注册:1:关闭,2:打开,3:打开并自动创建设备',
    `secret`       varchar(50)           DEFAULT '' COMMENT '动态注册产品秘钥',
    `desc`  varchar(200)          DEFAULT '' COMMENT '描述',
    `createdTime`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updatedTime`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deletedTime`  datetime              DEFAULT NULL,
    `devStatus`    varchar(20)  NOT NULL DEFAULT '' COMMENT '产品状态',
    `tags`        json not null comment '产品标签',
    PRIMARY KEY (`productID`),
    UNIQUE KEY `productName` (`productName`) USING BTREE,
    KEY `deviceType` (`deviceType`) USING BTREE,
    KEY `createdTime` (`createdTime`) USING BTREE
    ) ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = COMPACT COMMENT ='产品信息表';


CREATE TABLE if not exists `dm_product_custom`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `productID`   char(11) NOT NULL COMMENT '产品id',
    `scriptLang` tinyint(1)  default 1 comment '脚本语言类型 1:JavaScript 2:lua 3:python',
    `customTopic` json  null comment '自定义topic数组',
    `transformScript` text   null comment '协议转换脚本',
    `loginAuthScript` text   null comment '登录认证脚本',
    `createdTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updatedTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deletedTime` datetime          DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `productID` (`productID`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 4
  DEFAULT CHARSET = utf8mb4
  ROW_FORMAT = COMPACT COMMENT ='产品自定义协议表';


CREATE TABLE if not exists `dm_product_schema`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `productID`   char(11) NOT NULL COMMENT '产品id',
    `tag` tinyint(1)  default 1 comment '物模型标签 1:自定义 2:可选 3:必选  必选不可删除',
    `type` tinyint(1)  default 1 comment '物模型类型 1:property属性 2:event事件 3:action行为',
    `identifier`      varchar(100) not null COMMENT '标识符',
    `name`      varchar(100) not null  COMMENT '功能名称',
    `desc` varchar(200) default '' COMMENT '描述',
    `required` tinyint(1)  default 2 comment '是否必须,1是 2否',
    `affordance` json not null  COMMENT '各类型的自定义功能定义',
    `createdTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updatedTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deletedTime` datetime          DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `productID_identifier` (`productID`, `identifier`),
    KEY `productID_type` (`productID`, `type`)
    ) ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = COMPACT COMMENT ='产品物模型表';

CREATE TABLE if not exists `dm_device_info` (
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `productID`   char(11)     NOT NULL COMMENT '产品id',
    `deviceName`  varchar(100) NOT NULL COMMENT '设备名称',
    `deviceAlias` varchar(100) NOT NULL COMMENT '设备别名',
    `secret`      varchar(50)  NOT NULL DEFAULT '' COMMENT '设备秘钥',
    `cert`        varchar(512) NOT NULL DEFAULT '' COMMENT '设备证书',
    `imei`        varchar(15)  NOT NULL DEFAULT '' COMMENT 'IMEI号信息',
    `mac`         varchar(17)  NOT NULL DEFAULT '' COMMENT 'MAC号信息',
    `version`     varchar(64)  NOT NULL DEFAULT '' COMMENT '固件版本',
    `hardInfo`    varchar(64)  NOT NULL DEFAULT '' COMMENT '模组硬件型号',
    `softInfo`    varchar(64)  NOT NULL DEFAULT '' COMMENT '模组软件版本',
    `position`    point        NOT NULL COMMENT '设备的位置(默认百度坐标系BD09)',
    `address`     varchar(512) NOT NULL DEFAULT '' COMMENT '所在地址',
    `tags`        json         NOT NULL COMMENT '设备标签',
    `isOnline`    tinyint(1)   NOT NULL DEFAULT 2 COMMENT '是否在线,1是2否',
    `firstLogin`  datetime              DEFAULT NULL COMMENT '激活时间',
    `lastLogin`   datetime              DEFAULT NULL COMMENT '最后上线时间',
    `logLevel`    tinyint(1)   NOT NULL DEFAULT '1' COMMENT '日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试',
    `createdTime` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updatedTime` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deletedTime` datetime              DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `productID_deviceName`(`productID`, `deviceName`),
    KEY `createdTime`(`createdTime`) USING BTREE
)
    ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARSET = utf8mb4 COMMENT ='设备信息表';


-- # CREATE TABLE if not exists `category_detail`
-- # (
-- #     `id`           bigint  NOT NULL,
-- #     `categoryKey`  varchar(20)      NOT NULL COMMENT '产品种类英文key',
-- #     `categoryName` varchar(100)     NOT NULL COMMENT '产品种类名字',
-- #     `parentID`     int(10) DEFAULT 0 COMMENT '父类id',
-- #     `isLeaf`       tinyint default 0 comment '是否是叶子节点',
-- #     `listOrder`    int(10) DEFAULT '' COMMENT '排序',
-- #     `template`     text COMMENT '物模型模板',
-- #     `createdTime`  datetime         NOT NULL,
-- #     PRIMARY KEY (`id`)
-- # ) ENGINE = InnoDB
-- #   DEFAULT CHARSET = utf8mb4 COMMENT ='产品品类详情';


CREATE TABLE if not exists `dm_product_firmware`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `productID`   char(11)     NOT NULL COMMENT '产品id',
    `version`     varchar(64)           DEFAULT '' COMMENT '固件版本',
    `createdTime` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updatedTime` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deletedTime` datetime              DEFAULT NULL,
    `name`        varchar(64)           DEFAULT '' COMMENT '固件名称',
    `desc` varchar(200)          DEFAULT '' COMMENT '描述',
    `size`        bigint       NOT NULL COMMENT '固件大小',
    `dir`         varchar(128) NOT NULL COMMENT '固件标识,拿来下载文件',
    PRIMARY KEY (`id`),
    UNIQUE KEY `deviceVersion` (`productID`, `version`)
    ) ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = COMPACT COMMENT ='产品固件信息表';



CREATE TABLE if not exists `dm_group_info`
(
    `groupID`     bigint COMMENT '分组ID',
    `parentID`    bigint NOT NULL DEFAULT 0 COMMENT '父组ID 0-根组',
    `groupName`   VARCHAR(100) NOT NULL COMMENT '分组名称',
    `desc` VARCHAR(200) DEFAULT '' COMMENT '描述',
    `tags`        json not null comment '分组标签',
    `createdTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY `groupID` (`groupID`),
    UNIQUE KEY `groupNameIndex` (`groupName`)
    ) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COMMENT = '设备分组信息表';


CREATE TABLE if not exists `dm_group_device`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `groupID`     bigint NOT NULL COMMENT '分组ID',
    `productID`   char(11) NOT NULL COMMENT '产品id',
    `deviceName`  varchar(100) NOT NULL COMMENT '设备名称',
    `createdTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `groupID_productID_deviceName` (`groupID`,`productID`,`deviceName`)
    ) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COMMENT = '分组与设备关系表';


CREATE TABLE if not exists `dm_gateway_device`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `gatewayProductID`   char(11) NOT NULL COMMENT '网关产品id',
    `gatewayDeviceName`  varchar(100) NOT NULL COMMENT '网关设备名称',
    `productID`   char(11) NOT NULL COMMENT '子设备产品id',
    `deviceName`  varchar(100) NOT NULL COMMENT '子设备名称',
    `createdTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `gatewayProductID_gatewayDeviceName_productID_deviceName` (`gatewayProductID`,`gatewayDeviceName`,`productID`,`deviceName`)
    ) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COMMENT = '网关与子设备关系表';

CREATE TABLE if not exists `dm_product_remote_config`
(
    `id`          bigint   NOT NULL AUTO_INCREMENT,
    `productID`   char(11) NOT NULL COMMENT '产品id',
    `content`     json not null comment '配置内容',
    `createdTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `productID` (`productID`)
    ) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COMMENT = '产品远程配置表';
