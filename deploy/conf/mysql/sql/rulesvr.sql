SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

create database if not EXISTS iThings;
use iThings;

CREATE TABLE if not exists `rule_flow_info`
(
    `id`         bigint     auto_increment  NOT NULL COMMENT 'id',
    `name`    varchar(128)           DEFAULT '' COMMENT '流的名称',
    `password`    char(32)     NOT NULL DEFAULT '' COMMENT '登录密码',
    `desc`       varchar(512)          DEFAULT '' COMMENT '描述',
    `isDisabled`  tinyint(1) DEFAULT 2      not null COMMENT '是否禁用 1:是 2:否',
    `createdTime` datetime     not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime     NULL     DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime              DEFAULT NULL COMMENT '删除时间，默认为空，表示未删除，非空表示已删除',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `unique_name` (`name`) USING BTREE,
    KEY `user_deletedTime` (`deletedTime`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  ROW_FORMAT = COMPACT COMMENT ='用户登录信息表';
