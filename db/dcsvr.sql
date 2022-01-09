create database dc;
use dc;

CREATE TABLE if not exists `group_info`
(
    `groupID`     bigint(20) unsigned NOT NULL COMMENT '组id',
    `name`        varchar(128)        NOT NULL COMMENT '组名',
    `uid`         bigint(20) unsigned NOT NULL COMMENT '管理员用户id',
    `createdTime` datetime            NOT NULL,
    `updatedTime` datetime DEFAULT NULL,
    `deletedTime` datetime DEFAULT NULL,
    PRIMARY KEY (`groupID`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='组信息表';

CREATE TABLE if not exists `group_member`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `groupID`     bigint(20) unsigned NOT NULL COMMENT '组id',
    `memberID`    varchar(50)         NOT NULL COMMENT '成员id',
    `memberType`  int(10) unsigned    NOT NULL COMMENT '成员类型:1:设备 2:用户',
    `createdTime` datetime            NOT NULL,
    `updatedTime` datetime DEFAULT NULL,
    `deletedTime` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `groupId_memberId` (`groupID`, `memberID`, `memberType`),
    KEY `memberId` (`memberID`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COMMENT='组成员信息表';