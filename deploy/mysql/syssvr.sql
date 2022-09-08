create database if not EXISTS things_sys;
use things_sys;
CREATE TABLE if not exists `user_info`
(
    `uid`         bigint                                                  NOT NULL COMMENT '用户id',
    `userName`    varchar(20)                                                      DEFAULT NULL COMMENT '登录用户名',
    `password`    char(32) CHARACTER SET utf8 COLLATE utf8_general_ci     NOT NULL DEFAULT '' COMMENT '登录密码',
    `email`       varchar(255)                                                     DEFAULT NULL COMMENT '邮箱',
    `phone`       varchar(20)                                                      DEFAULT NULL COMMENT '手机号',
    `wechat`      varchar(20)                                                      DEFAULT NULL COMMENT '微信union id',
    `lastIP`      varchar(40) CHARACTER SET utf8 COLLATE utf8_general_ci  NOT NULL DEFAULT '' COMMENT '最后登录ip',
    `regIP`       varchar(40) CHARACTER SET utf8 COLLATE utf8_general_ci  NOT NULL DEFAULT '' COMMENT '注册ip',
    `nickName`    varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci  NOT NULL DEFAULT '' COMMENT '用户的昵称',
    `sex`         tinyint(1)                                              NOT NULL DEFAULT '3' COMMENT '用户的性别，值为1时是男性，值为2时是女性，其他值为未知',
    `city`        varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci  NOT NULL DEFAULT '' COMMENT '用户所在城市',
    `country`     varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci  NOT NULL DEFAULT '' COMMENT '用户所在国家',
    `province`    varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci  NOT NULL DEFAULT '' COMMENT '用户所在省份',
    `language`    varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci  NOT NULL DEFAULT '' COMMENT '用户的语言，简体中文为zh_CN',
    `headImgUrl`  varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户头像',
    `role`        bigint                                                  not null COMMENT '用户角色',
    `createdTime` datetime                                                not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime                                                NULL     DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime                                                         DEFAULT NULL COMMENT '删除时间，默认为空，表示未删除，非空表示已删除',
    PRIMARY KEY (`uid`) USING BTREE,
    UNIQUE KEY `user_username` (`userName`) USING BTREE,
    UNIQUE KEY `user_phone` (`phone`) USING BTREE,
    UNIQUE KEY `user_email` (`email`) USING BTREE,
    UNIQUE KEY `user_wechat` (`wechat`) USING BTREE,
    KEY `user_deletedTime` (`deletedTime`) USING BTREE
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8
    ROW_FORMAT = COMPACT COMMENT ='用户登录信息表';

# 新增root用户
INSERT INTO `user_info`(`uid`, `userName`, `password`, `email`, `phone`, `wechat`, `lastIP`, `regIP`, `role`,`nickName`,
                        `sex`, `city`, `country`, `province`,`language`, `headImgUrl`,`deletedTime`)
VALUES (1740358057038188544, 'administrator', '4f0fded4a38abe7a3ea32f898bb82298', '163', '13911110000', 'wechat', '0.0.0.0', '0.0.0.0', 1,'liangjuan',
         1, 'shenzhen', 'Ut', 'guangdong', 'eiusmod', 'http',NULL);

CREATE TABLE if not exists `role_info`
(
    `id`          bigint auto_increment comment 'id编号',
    `name`        varchar(100)  CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '角色名称',
    `remark`      varchar(100)  CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '备注',
    `createdTime` datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime DEFAULT NULL,
    `status`      int default 1 null comment '状态  1:启用,2:禁用',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `nameIndex` (`name`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='角色管理表';

CREATE TABLE if not exists `role_menu`
(
    `id`           bigint auto_increment comment 'id编号',
    `roleID`       int null comment '角色ID',
    `menuID`       int null comment '菜单ID',
    `createdTime`  datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime`  datetime NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime`  datetime DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `roleIDMenuIDIndex` (`roleID`, `menuID`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='角色菜单关联表';

CREATE TABLE if not exists `menu_info`
(
    `id`              bigint auto_increment comment '编号',
    `parentID`        int null comment '父菜单ID，一级菜单为1',
    `type`            int null comment '类型   1：目录   2：菜单   3：按钮',
    `order`           int null comment '左侧table排序序号',
    `name`            varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '菜单名称',
    `path`            varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '系统的path',
    `component`       varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '页面',
    `icon`            varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '图标',
    `redirect`        varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '路由重定向',
    `backgroundUrl`   varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '后台地址',
    `hideInMenu`      int(11) not null default 2 comment '是否隐藏菜单 1-是 2-否';
    `createdTime`     datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` 	  datetime NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` 	  datetime DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `nameIndex` (`name`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='菜单管理表';

insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (2,"设备管理","icon_data_01","/deviceMangers",1,1,"./deviceMangers/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (3,"系统管理","icon_system","/systemManagers",1,2,"./systemManagers/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (4,"运维监控","icon_system","/operationsMonitorings",1,3,"./operationsMonitorings/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (5,"规则引擎","icon_system","/ruleEngines",1,4,"./ruleEngines/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (6,"产品","icon_system","/deviceMangers/products/index",2,1,"./deviceMangers/products/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (7,"产品详情-${id}","icon_system","/deviceMangers/products/details/:id",2,2,"./deviceMangers/products/details/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (8,"设备","icon_system","/deviceMangers/devices/index",2,3,"./deviceMangers/devices/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (9,"设备详情-${id}","icon_system","/deviceMangers/devices/details/:id",2,3,"./deviceMangers/devices/details/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (10,"用户管理","icon_system","/systemManagers/users/index",3,1,"./systemManagers/users/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (11,"角色管理","icon_system","/systemManagers/roles/index",3,2,"./systemManagers/roles/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (12,"菜单列表","icon_system","/systemManagers/menus/index",3,3,"./systemManagers/menus/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (13,"固件升级","icon_system","/operationsMonitorings/firmwareUpgrades/index",4,1,"./operationsMonitorings/firmwareUpgrades/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (14,"告警记录","icon_system","/operationsMonitorings/alarmRecords/index",4,2,"./operationsMonitorings/alarmRecords/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (15,"资源管理","icon_system","/operationsMonitorings/resourceManagements/index",4,3,"./operationsMonitorings/resourceManagements/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (16,"远程配置","icon_system","/operationsMonitorings/remoteConfigurations/index",4,4,"./operationsMonitorings/remoteConfigurations/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (17,"告警中心","icon_system","/operationsMonitorings/alarmCenters/index",4,5,"./operationsMonitorings/alarmCenters/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (18,"在线调试","icon_system","/operationsMonitorings/onlineDebugs/index",4,6,"./operationsMonitorings/onlineDebugs/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (19,"消息规则","icon_system","/ruleEngines/messageRules/index",5,1,"./ruleEngines/messageRules/index.tsx",2,2);
insert into menu_info (id,name,icon,`path`,parentID,`order`,component,`type`,`hideInMenu`) values (20,"规则日志","icon_system","/ruleEngines/ruleLogs/index",5,2,"./ruleEngines/ruleLogs/index.tsx",2,2);
