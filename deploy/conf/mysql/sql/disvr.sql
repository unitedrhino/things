-- 设备交互模块SQL
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

create database if not EXISTS iThings;
use iThings;



CREATE TABLE if not exists `di_device_shadow`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `productID`   char(11)     NOT NULL COMMENT '产品id',
    `deviceName`  varchar(100) NOT NULL COMMENT '设备名称',
    `dataID`  varchar(100) NOT NULL COMMENT '属性id',
    `value`  varchar(100) NOT NULL COMMENT '属性值',
    PRIMARY KEY (`id`),
    UNIQUE KEY `data` (`productID`,`deviceName`,`dataID`) USING BTREE,
    `createdTime` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updatedTime` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deletedTime` datetime              DEFAULT NULL,
    ) ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = COMPACT COMMENT ='设备影子表';