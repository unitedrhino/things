SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

create database if not EXISTS iThings;
use iThings;

CREATE TABLE if not exists `rule_scene_info`
(
    `id`         bigint     auto_increment  NOT NULL COMMENT 'id',
    `name`    varchar(128)           DEFAULT '' COMMENT '场景名称',
    `trigger`    json       COMMENT '触发器',
    `when`  json       COMMENT '触发条件',
    `then`  json      COMMENT '满足条件时执行的动作',
    `desc`       varchar(512)          DEFAULT '' COMMENT '描述',
    `createdTime` datetime     not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime     NULL     DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime              DEFAULT NULL COMMENT '删除时间，默认为空，表示未删除，非空表示已删除',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `unique_name` (`name`) USING BTREE,
    KEY `deletedTime` (`deletedTime`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  ROW_FORMAT = COMPACT COMMENT ='规则引擎-场景联动信息表';
