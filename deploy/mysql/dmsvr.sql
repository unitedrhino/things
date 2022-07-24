create database if not EXISTS things_dm;
use things_dm;

CREATE TABLE if not exists `product_info`
(
    `productID`    varchar(20)  NOT NULL COMMENT '产品id',
    `productName`  varchar(100) NOT NULL COMMENT '产品名称',
    `productType`  int(10) unsigned      DEFAULT '0' COMMENT '产品状态:0:开发中,1:审核中,2:已发布',
    `authMode`     int(10) unsigned      DEFAULT '0' COMMENT '认证方式:0:账密认证,1:秘钥认证',
    `deviceType`   int(10) unsigned      DEFAULT '0' COMMENT '设备类型:0:设备,1:网关,2:子设备',
    `categoryID`   int(10) unsigned      DEFAULT '0' COMMENT '产品品类',
    `netType`      int(10) unsigned      DEFAULT '0' COMMENT '通讯方式:0:其他,1:wi-fi,2:2G/3G/4G,3:5G,4:BLE,5:LoRaWAN',
    `dataProto`    int(10) unsigned      DEFAULT '0' COMMENT '数据协议:0:自定义,1:数据模板',
    `autoRegister` int(10) unsigned      DEFAULT '0' COMMENT '动态注册:0:关闭,1:打开,2:打开并自动创建设备',
    `secret`       varchar(50)           DEFAULT '' COMMENT '动态注册产品秘钥',
    `description`  varchar(200)          DEFAULT '' COMMENT '描述',
    `createdTime`  datetime     NOT NULL,
    `updatedTime`  datetime              DEFAULT NULL,
    `deletedTime`  datetime              DEFAULT NULL,
    `devStatus`    varchar(20)  NOT NULL DEFAULT '' COMMENT '产品状态',
    PRIMARY KEY (`productID`),
    UNIQUE KEY `productName` (`productName`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 4
  DEFAULT CHARSET = utf8mb4
  ROW_FORMAT = COMPACT COMMENT ='产品信息表';

CREATE TABLE if not exists `product_template`
(
    `productID`   varchar(20) NOT NULL COMMENT '产品id',
    `template`    text        not null COMMENT '物模型模板',
    `createdTime` datetime    NOT NULL,
    `updatedTime` datetime DEFAULT NULL,
    `deletedTime` datetime DEFAULT NULL,
    PRIMARY KEY (`productID`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 4
  DEFAULT CHARSET = utf8mb4
  ROW_FORMAT = COMPACT COMMENT ='产品物模型表';

CREATE TABLE if not exists `device_info`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `productID`   varchar(20)      NOT NULL COMMENT '产品id',
    `deviceName`  varchar(100)     NOT NULL COMMENT '设备名称',
    `secret`      varchar(50)  DEFAULT '' COMMENT '设备秘钥',
    `firstLogin`  datetime     DEFAULT NULL COMMENT '激活时间',
    `lastLogin`   datetime     DEFAULT NULL COMMENT '最后上线时间',
    `isOnline`      tinyint  default 0 comment '是否在线,0离线1在线',
    `createdTime` datetime         NOT NULL,
    `updatedTime` datetime     DEFAULT NULL,
    `deletedTime` datetime     DEFAULT NULL,
    `version`     varchar(64)  DEFAULT '' COMMENT '固件版本',
    `logLevel`    int(10)      DEFAULT '1' COMMENT '日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试',
    `cert`        varchar(512) DEFAULT '' COMMENT '设备证书',
    PRIMARY KEY (`id`),
    UNIQUE KEY `deviceName` (`productID`, `deviceName`),
    KEY `device_productID` (`productID`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 3
  DEFAULT CHARSET = utf8mb4 COMMENT ='设备信息表';




# CREATE TABLE if not exists `category_detail`
# (
#     `id`           bigint unsigned NOT NULL,
#     `categoryKey`  varchar(20)      NOT NULL COMMENT '产品种类英文key',
#     `categoryName` varchar(100)     NOT NULL COMMENT '产品种类名字',
#     `parentID`     int(10) DEFAULT 0 COMMENT '父类id',
#     `isLeaf`       tinyint default 0 comment '是否是叶子节点',
#     `listOrder`    int(10) DEFAULT '' COMMENT '排序',
#     `template`     text COMMENT '物模型模板',
#     `createdTime`  datetime         NOT NULL,
#     PRIMARY KEY (`id`)
# ) ENGINE = InnoDB
#   DEFAULT CHARSET = utf8mb4 COMMENT ='产品品类详情';



CREATE TABLE if not exists `product_firmware`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `productID`   varchar(20)      NOT NULL COMMENT '产品id',
    `version`     varchar(64)  DEFAULT '' COMMENT '固件版本',
    `createdTime` datetime         NOT NULL,
    `updatedTime` datetime     DEFAULT NULL,
    `deletedTime` datetime     DEFAULT NULL,
    `name`     varchar(64)  DEFAULT '' COMMENT '固件名称',
    `description`  varchar(200)          DEFAULT '' COMMENT '描述',
    `size`          bigint unsigned NOT NULL COMMENT '固件大小',
    `dir`     varchar(128)  NOT NULL COMMENT '固件标识,拿来下载文件',
    PRIMARY KEY (`id`),
    UNIQUE KEY `deviceVersion` (`productID`, `version`)
    ) ENGINE = InnoDB
    AUTO_INCREMENT = 4
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = COMPACT COMMENT ='产品固件信息表';