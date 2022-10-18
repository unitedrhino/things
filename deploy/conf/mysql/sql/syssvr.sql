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

INSERT into  role_info (name) values ('admin');

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

INSERT into role_menu (roleID,menuID) values (1,2);
INSERT into role_menu (roleID,menuID) values (1,3);
INSERT into role_menu (roleID,menuID) values (1,4);
INSERT into role_menu (roleID,menuID) values (1,5);
INSERT into role_menu (roleID,menuID) values (1,6);
INSERT into role_menu (roleID,menuID) values (1,7);
INSERT into role_menu (roleID,menuID) values (1,8);
INSERT into role_menu (roleID,menuID) values (1,9);
INSERT into role_menu (roleID,menuID) values (1,10);
INSERT into role_menu (roleID,menuID) values (1,11);
INSERT into role_menu (roleID,menuID) values (1,12);
INSERT into role_menu (roleID,menuID) values (1,13);
INSERT into role_menu (roleID,menuID) values (1,14);
INSERT into role_menu (roleID,menuID) values (1,15);
INSERT into role_menu (roleID,menuID) values (1,16);
INSERT into role_menu (roleID,menuID) values (1,17);
INSERT into role_menu (roleID,menuID) values (1,18);
INSERT into role_menu (roleID,menuID) values (1,19);
INSERT into role_menu (roleID,menuID) values (1,20);

CREATE TABLE if not exists `menu_info`
(
    `id`              bigint auto_increment comment '编号',
    `parentID`        int not null default 1 comment '父菜单ID，一级菜单为1',
    `type`            int not null default 1 comment '类型   1：目录   2：菜单   3：按钮',
    `order`           int not null default 1 comment '左侧table排序序号',
    `name`            varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '菜单名称',
    `path`            varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '系统的path',
    `component`       varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '页面',
    `icon`            varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '图标',
    `redirect`        varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '路由重定向',
    `backgroundUrl`   varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '后台地址',
    `hideInMenu`      int(11) not null default 2 comment '是否隐藏菜单 1-是 2-否',
    `createdTime`     datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` 	  datetime NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` 	  datetime DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `nameIndex` (`name`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='菜单管理表';

INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(2, 1, 2, 1, '设备管理', '/deviceMangers', './deviceMangers/index.tsx', 'icon_data_01', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 15:38:54', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(3, 1, 2, 2, '系统管理', '/systemManagers', './systemManagers/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 15:38:54', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(4, 1, 2, 3, '运维监控', '/operationsMonitorings', './operationsMonitorings/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 15:38:54', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(5, 1, 0, 4, '规则引擎', '/ruleEngines', './ruleEngines/index.tsx', 'icon_system', '', '', 1, '2022-09-24 15:38:54', '2022-09-24 17:01:55', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(6, 2, 2, 1, '产品', '/deviceMangers/product/index', './deviceMangers/product/index', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 16:16:58', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(7, 2, 0, 2, '产品详情', '/deviceMangers/product/detail/:id', './deviceMangers/product/detail/index', 'icon_system', '', '', 1, '2022-09-24 15:38:54', '2022-09-24 16:59:45', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(8, 2, 0, 3, '设备', '/deviceMangers/device/index', './deviceMangers/device/index', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 17:03:49', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(9, 2, 0, 3, '设备详情', '/deviceMangers/device/detail/:id/:name', './deviceMangers/device/detail/index', 'icon_system', '', '', 1, '2022-09-24 15:38:54', '2022-09-24 17:03:44', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(10, 3, 0, 1, '用户管理', '/systemMangers/user/index', './systemMangers/user/index', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 16:59:01', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(11, 3, 2, 2, '角色管理', '/systemMangers/role/index', './systemMangers/role/index', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 16:15:37', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(12, 3, 2, 3, '菜单列表', '/systemMangers/menu/index', './systemMangers/menu/index', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 16:15:52', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(13, 4, 0, 1, '固件升级', '/operationsMonitorings/firmwareUpgrades/index', './operationsMonitorings/firmwareUpgrades/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 16:52:10', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(14, 4, 2, 2, '告警记录', '/operationsMonitorings/alarmRecords/index', './operationsMonitorings/alarmRecords/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 15:38:54', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(15, 4, 2, 3, '资源管理', '/operationsMonitorings/resourceManagements/index', './operationsMonitorings/resourceManagements/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 15:38:54', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(16, 4, 2, 4, '远程配置', '/operationsMonitorings/remoteConfigurations/index', './operationsMonitorings/remoteConfigurations/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 15:38:54', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(17, 4, 2, 5, '告警中心', '/operationsMonitorings/alarmCenters/index', './operationsMonitorings/alarmCenters/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 15:38:54', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(18, 4, 2, 6, '在线调试', '/operationsMonitorings/onlineDebugs/index', './operationsMonitorings/onlineDebugs/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 15:38:54', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(19, 5, 2, 1, '消息规则', '/ruleEngines/messageRules/index', './ruleEngines/messageRules/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 15:38:54', NULL);
INSERT INTO things_sys.menu_info
(id, parentID, `type`, `order`, name, `path`, component, icon, redirect, backgroundUrl, hideInMenu, createdTime, updatedTime, deletedTime)
VALUES(20, 5, 2, 2, '规则日志', '/ruleEngines/ruleLogs/index', './ruleEngines/ruleLogs/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54', '2022-09-24 15:38:54', NULL);



