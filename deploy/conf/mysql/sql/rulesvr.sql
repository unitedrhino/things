SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

create database if not EXISTS iThings;
use iThings;

CREATE TABLE if not exists `rule_scene_info`
(
    `id`         bigint     auto_increment  NOT NULL COMMENT 'id',
    `name`    varchar(128)           DEFAULT '' COMMENT '场景名称',
    `triggerType`    varchar(24)  not null     COMMENT '触发器类型 device: 设备触发 timer: 定时触发 manual:手动触发',
    `trigger`    json       COMMENT '触发器内容-根据触发器类型改变',
    `when`  json       COMMENT '触发条件',
    `then`  json      COMMENT '满足条件时执行的动作',
    `desc`       varchar(512)          DEFAULT '' COMMENT '描述',
    `state`           tinyint(1) NOT NULL DEFAULT 2 COMMENT '状态（1启用 2禁用）',
    `createdTime` datetime     not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime     NULL     DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime              DEFAULT NULL COMMENT '删除时间，默认为空，表示未删除，非空表示已删除',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `unique_name` (`name`) USING BTREE,
    KEY `triggerType` (`triggerType`) USING BTREE,
    KEY `state` (`state`) USING BTREE,
    KEY `deletedTime` (`deletedTime`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  ROW_FORMAT = COMPACT COMMENT ='规则引擎-场景联动信息表';


CREATE TABLE if not exists `rule_alarm_scene`
(
    `id`          bigint auto_increment comment 'id编号',
    `alarmID`      bigint      null comment '告警配置ID',
    `sceneID`      int      null comment '场景ID',
    `createdTime` datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime NULL     DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime          DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `alarmIDSceneIDIndex` (`alarmID`, `sceneID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPACT COMMENT ='告警配置与场景关联表';


CREATE TABLE if not exists `rule_alarm_info`
(
    `id`              bigint auto_increment comment '编号',
    `name`            varchar(100) NOT NULL DEFAULT '' comment '告警配置名称',
    `desc`            varchar(100) NOT NULL DEFAULT '' comment '告警配置说明',
    `type`            tinyint(1) NOT NULL COMMENT '告警配置类型（1产品 2设备 3其它）',
    `level`           tinyint(1) NOT NULL COMMENT '告警配置级别（1提醒 2一般 3严重 4紧急 5超紧急）',
    `state`           tinyint(1) NOT NULL DEFAULT 2 COMMENT '告警配置状态（1启用 2禁用）',
    `dealState`       tinyint(1) NOT NULL DEFAULT 1 COMMENT '告警记录状态（1告警中 2已处理）',
    `lastAlarm`       datetime              DEFAULT NULL COMMENT '最新告警时间',
    `createdTime`     datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` 	  datetime NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` 	  datetime DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `alarm_config_name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='告警配置信息表';

CREATE TABLE if not exists `rule_alarm_log`
(
    `id`              bigint auto_increment comment '编号',
    `alarmID`         bigint NOT NULL comment '告警记录ID',
    `serial`          varchar(1024) NOT NULL DEFAULT '' comment '告警流水',
    `sceneName`       varchar(100) NOT NULL DEFAULT '' comment '场景名称',
    `sceneID`         int      null comment '场景ID',
    `desc`            varchar(1024) NOT NULL DEFAULT '' comment '告警说明',
    `createdTime`     datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '告警时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `alarm_time` (`createdTime`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='告警日志表';


CREATE TABLE if not exists `rule_alarm_deal_record`
(
    `id`              bigint auto_increment comment '编号',
    `alarmID`         bigint NOT NULL comment '告警配置ID',
    `result`          varchar(1024) NOT NULL DEFAULT '' comment '告警处理结果',
    `type`            tinyint(1) NOT NULL COMMENT '告警处理类型（1人工 2其它）',
    `alarmTime`       datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '告警时间',
    `createdTime`     datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '告警处理时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `alarm_deal_time` (`createdTime`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='告警处理记录表';