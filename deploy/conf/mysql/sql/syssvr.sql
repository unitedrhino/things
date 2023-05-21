-- 系统管理模块SQL
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

create database if not EXISTS iThings;
use iThings;

CREATE TABLE if not exists `sys_user_info`
(
    `uid`         bigint       NOT NULL COMMENT '用户id',
    `userName`    varchar(20)           DEFAULT NULL COMMENT '登录用户名',
    `password`    char(32)     NOT NULL DEFAULT '' COMMENT '登录密码',
    `email`       varchar(255)   CHARACTER SET utf8 COLLATE utf8_general_ci       DEFAULT NULL COMMENT '邮箱',
    `phone`       varchar(20)           DEFAULT NULL COMMENT '手机号',
    `wechat`      varchar(20)           DEFAULT NULL COMMENT '微信union id',
    `lastIP`      varchar(40)  NOT NULL DEFAULT '' COMMENT '最后登录ip',
    `regIP`       varchar(40)  NOT NULL DEFAULT '' COMMENT '注册ip',
    `nickName`    varchar(60)  NOT NULL DEFAULT '' COMMENT '用户的昵称',
    `sex`         tinyint(1)   NOT NULL DEFAULT '3' COMMENT '用户的性别，值为1时是男性，值为2时是女性，其他值为未知',
    `city`        varchar(50)  NOT NULL DEFAULT '' COMMENT '用户所在城市',
    `country`     varchar(50)  NOT NULL DEFAULT '' COMMENT '用户所在国家',
    `province`    varchar(50)  NOT NULL DEFAULT '' COMMENT '用户所在省份',
    `language`    varchar(50)  NOT NULL DEFAULT '' COMMENT '用户的语言，简体中文为zh_CN',
    `headImgUrl`  varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户头像',
    `role`        bigint       not null COMMENT '用户角色',
    `isAllData`   tinyint(1)   UNSIGNED NOT NULL default 2 COMMENT '是否所有数据权限（1是，2否）',
    `createdTime` datetime     not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime     NULL     DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime              DEFAULT NULL COMMENT '删除时间，默认为空，表示未删除，非空表示已删除',
    PRIMARY KEY (`uid`) USING BTREE,
    UNIQUE KEY `user_username` (`userName`) USING BTREE,
    UNIQUE KEY `user_phone` (`phone`) USING BTREE,
    UNIQUE KEY `user_email` (`email`) USING BTREE,
    UNIQUE KEY `user_wechat` (`wechat`) USING BTREE,
    KEY `user_deletedTime` (`deletedTime`) USING BTREE
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = COMPACT COMMENT ='用户登录信息表';

# 新增root用户
INSERT IGNORE INTO `sys_user_info`(`uid`, `userName`, `password`, `email`, `phone`, `wechat`, `lastIP`, `regIP`, `role`,
                        `nickName`,
                        `sex`, `city`, `country`, `province`, `language`, `headImgUrl`, `deletedTime`)
VALUES (1740358057038188544, 'administrator', '4f0fded4a38abe7a3ea32f898bb82298', '163', '13911110000', 'wechat',
        '0.0.0.0', '0.0.0.0', 1, 'liangjuan',
        1, 'shenzhen', 'Ut', 'guangdong', 'eiusmod', 'http', NULL);

CREATE TABLE if not exists `sys_role_info`
(
    `id`          bigint auto_increment comment 'id编号',
    `name`        varchar(100) NOT NULL DEFAULT '' COMMENT '角色名称',
    `remark`      varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
    `createdTime` datetime     not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime     NULL     DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime              DEFAULT NULL,
    `status`      tinyint(1)                   default 1 null comment '状态  1:启用,2:禁用',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `nameIndex` (`name`) USING BTREE
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = COMPACT COMMENT ='角色管理表';

INSERT IGNORE INTO sys_role_info (id, name) values (1, 'admin');

CREATE TABLE if not exists `sys_role_menu`
(
    `id`          bigint auto_increment comment 'id编号',
    `roleID`      int      null comment '角色ID',
    `menuID`      int      null comment '菜单ID',
    `createdTime` datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` datetime NULL     DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` datetime          DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `roleIDMenuIDIndex` (`roleID`, `menuID`) USING BTREE
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = COMPACT COMMENT ='角色菜单关联表';

INSERT IGNORE into `sys_role_menu` (`roleID`,`menuID`)values(1,2),(1,3),(1,4),(1,5),(1,6),(1,7),(1,8),(1,9),(1,10),(1,11),(1,12),(1,13),(1,14),(1,15),(1,16),(1,17),(1,18),(1,19),(1,20),(1,21),(1,22),(1,23),(1,24),(1,25),(1,26),(1,27),(1,28),(1,29),(1,30),(1,31),(1,32),(1,33),(1,34),(1,35),(1,36),(1,37),(1,38),(1,39),(1,40),(1,41),(1,42),(1,43),(1,44),(1,45),(1,46),(1,47),(1,48),(1,49),(1,50),(1,51),(1,52),(1,53),(1,54),(1,55),(1,56),(1,57),(1,58),(1,59),(1,60),(1,61),(1,62),(1,63),(1,64),(1,65),(1,66),(1,67),(1,68),(1,69),(1,70),(1,71),(1,72),(1,73),(1,74),(1,75),(1,76),(1,77),(1,78),(1,79),(1,80),(1,81),(1,82),(1,83),(1,84),(1,85),(1,86),(1,87),(1,88),(1,89),(1,90),(1,91),(1,92),(1,93),(1,94),(1,95),(1,96),(1,97),(1,98),(1,99),(1,100);


CREATE TABLE if not exists `sys_menu_info`
(
    `id`            bigint auto_increment comment '编号',
    `parentID`      int          not null default 1 comment '父菜单ID，一级菜单为1',
    `type`          int          not null default 1 comment '类型   1：目录   2：菜单   3：按钮',
    `order`         int          not null default 1 comment '左侧table排序序号',
    `name`          varchar(50)  NOT NULL DEFAULT '' comment '菜单名称',
    `path`          varchar(64)  NOT NULL DEFAULT '' comment '系统的path',
    `component`     varchar(64)  NOT NULL DEFAULT '' comment '页面',
    `icon`          varchar(64)  NOT NULL DEFAULT '' comment '图标',
    `redirect`      varchar(64)  NOT NULL DEFAULT '' comment '路由重定向',
    `backgroundUrl` varchar(128) NOT NULL DEFAULT '' comment '后台地址',
    `hideInMenu`    int(11)      not null default 2 comment '是否隐藏菜单 1-是 2-否',
    `createdTime`   datetime     not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime`   datetime     NULL     DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime`   datetime              DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `nameIndex` (`name`) USING BTREE
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = COMPACT COMMENT ='菜单管理表';

INSERT IGNORE INTO `sys_menu_info`
VALUES (2, 1, 0, 2, '设备管理', '/deviceMangers', './deviceMangers/index.tsx', 'icon_data_01', '', '', 2,
        '2022-09-24 15:38:54', '2022-10-18 10:31:11', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (3, 1, 0, 9, '系统管理', '/systemManagers', './systemManagers/index.tsx', 'icon_system', '', '', 2,
        '2022-09-24 15:38:54', '2022-10-17 20:42:42', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (4, 1, 0, 4, '运维监控', '/operationsMonitorings', './operationsMonitorings/index.tsx', 'icon_hvac', '', '', 2,
        '2022-09-24 15:38:54', '2022-10-18 21:42:34', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (6, 2, 2, 1, '产品', '/deviceMangers/product/index', './deviceMangers/product/index', 'icon_system', '', '', 2,
        '2022-09-24 15:38:54', '2022-09-24 16:16:58', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (7, 2, 0, 1, '产品详情', '/deviceMangers/product/detail/:id', './deviceMangers/product/detail/index',
        'icon_system', '', '', 1, '2022-09-24 15:38:54', '2022-10-13 23:02:39', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (8, 2, 0, 2, '设备', '/deviceMangers/device/index', './deviceMangers/device/index', 'icon_system', '', '', 2,
        '2022-09-24 15:38:54', '2022-10-13 23:02:46', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (9, 2, 0, 2, '设备详情', '/deviceMangers/device/detail/:id/:name', './deviceMangers/device/detail/index',
        'icon_system', '', '', 1, '2022-09-24 15:38:54', '2022-10-13 23:02:51', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (10, 3, 0, 1, '用户管理', '/systemMangers/user/index', './systemMangers/user/index', 'icon_system', '', '', 2,
        '2022-09-24 15:38:54', '2022-09-24 16:59:01', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (11, 3, 2, 2, '角色管理', '/systemMangers/role/index', './systemMangers/role/index', 'icon_system', '', '', 2,
        '2022-09-24 15:38:54', '2022-09-24 16:15:37', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (12, 3, 2, 3, '菜单列表', '/systemMangers/menu/index', './systemMangers/menu/index', 'icon_system', '', '', 2,
        '2022-09-24 15:38:54', '2022-09-24 16:15:52', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (13, 4, 0, 1, '固件升级', '/operationsMonitorings/firmwareUpgrades/index',
        './operationsMonitorings/firmwareUpgrades/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54',
        '2022-10-17 20:47:13', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (15, 4, 0, 3, '资源管理', '/operationsMonitorings/resourceManagements/index',
        './operationsMonitorings/resourceManagements/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54',
        '2022-10-17 20:45:12', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (16, 4, 0, 4, '远程配置', '/operationsMonitorings/remoteConfiguration/index',
        './operationsMonitorings/remoteConfiguration/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54',
        '2022-10-17 20:45:19', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (18, 4, 2, 6, '在线调试', '/operationsMonitorings/onlineDebugs/index',
        './operationsMonitorings/onlineDebugs/index.tsx', 'icon_system', '', '', 2, '2022-09-24 15:38:54',
        '2022-09-24 15:38:54', NULL);

INSERT IGNORE INTO `sys_menu_info`
VALUES (23, 2, 0, 3, '分组', '/deviceMangers/group/index', './deviceMangers/group/index.tsx', 'icon_system', '', '', 2,
        '2022-10-13 23:04:01', '2022-10-13 23:04:01', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (24, 2, 0, 3, '分组详情', '/deviceMangers/group/detail/:id', './deviceMangers/group/detail/index.tsx',
        'icon_system', '', '', 1, '2022-10-13 23:04:44', '2022-10-13 23:06:45', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (25, 4, 0, 7, '日志服务', '/operationsMonitorings/logService/index',
        './operationsMonitorings/logService/index.tsx', 'icon_system', '', '', 2, '2022-10-16 23:04:36',
        '2022-10-16 23:04:36', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (35, 1, 1, 1, '首页', '/home',
        './home/index.tsx', 'icon_dosing', '', '', 2, '2022-10-16 23:04:36',
        '2022-10-16 23:04:36', NULL);
INSERT IGNORE INTO `sys_menu_info`
VALUES (38, 3, 1, 5, '日志管理', '/systemManagers/log/index',
        './systemManagers/log/index.tsx', 'icon_system', '', '', 2, '2023-02-07 11:28:11', '2023-02-18 11:09:33', null);
INSERT IGNORE INTO `sys_menu_info`
VALUES (39, 38, 1, 1, '操作日志', '/systemMangers/log/operationLog/index',
        './systemMangers/log/operationLog/index.tsx', 'icon_dosing', '', '', 2, '2023-02-07 11:31:07', '2023-02-18 10:45:50', null);
INSERT IGNORE INTO `sys_menu_info`
VALUES (41, 38, 1, 2, '登录日志', '/systemMangers/log/loginLog/index',
        './systemMangers/log/loginLog/index', 'icon_heat', '', '', 2, '2023-02-16 23:03:15', '2023-02-18 10:45:55', null);
INSERT IGNORE INTO `sys_menu_info`
VALUES (42, 3, 1, 4, '接口管理', '/systemMangers/api/index',
        './systemMangers/api/index', 'icon_system', '', '', 2, '2023-02-18 11:08:56', '2023-02-18 11:09:27', null);
INSERT IGNORE INTO `sys_menu_info`
    VALUES (43, 1, 1, 5, '告警管理', '/alarmMangers', './alarmMangers/index', 'icon_ap', '', '', 2, '2023-02-22 16:42:54', '2023-04-09 12:18:22', null);
INSERT IGNORE INTO `sys_menu_info`
    VALUES (44, 43, 1, 1, '告警配置', '/alarmMangers/alarmConfiguration/index', './alarmMangers/alarmConfiguration/index', 'icon_ap', '', '', 2, '2023-02-22 16:43:48', '2023-02-22 16:45:25', null);
INSERT IGNORE INTO `sys_menu_info`
    VALUES (53, 43, 1, 5, '新增告警配置', '/alarmMangers/alarmConfiguration/save', './alarmMangers/alarmConfiguration/addAlarmConfig/index', 'icon_ap', '', '', 1, '2023-03-26 15:11:15', '2023-04-09 11:13:26', null);
INSERT IGNORE INTO `sys_menu_info`
    VALUES (54, 43, 1, 5, '告警日志', '/alarmMangers/alarmConfiguration/log/detail/:id/:level', './alarmMangers/alarmLog/index', 'icon_ap', '', '', 1, '2023-04-09 17:25:32', '2023-04-15 18:05:14', null);
INSERT IGNORE INTO `sys_menu_info`
    VALUES (62, 43, 1, 5, '告警记录', '/alarmMangers/alarmConfiguration/log', './alarmMangers/alarmRecord/index', 'icon_ap', '', '', 2, '2023-04-09 18:47:12', '2023-04-09 18:47:12', null);
INSERT IGNORE INTO `sys_menu_info`
    VALUES (50, 1, 1, 5, '规则引擎', '/ruleEngine', './ruleEngine/index.tsx', 'icon_dosing', '', '', 2, '2023-05-31 23:07:08', '2023-06-06 16:03:32', null);
INSERT IGNORE INTO `sys_menu_info`
    VALUES (51, 50, 1, 1, '场景联动', '/ruleEngine/scene/index', './ruleEngine/scene/index.tsx', 'icon_device', '', '', 2, '2023-05-31 23:07:34', '2023-05-31 23:08:15', null);

DROP TABLE IF EXISTS `sys_login_log`;
CREATE TABLE `sys_login_log` (
                                 `id` bigint auto_increment COMMENT '编号',
                                 `uid`         bigint       NOT NULL COMMENT '用户id',
                                 `userName` varchar(50) DEFAULT '' COMMENT '登录账号',
                                 `ipAddr` varchar(50) DEFAULT '' COMMENT '登录IP地址',
                                 `loginLocation` varchar(100) DEFAULT '' COMMENT '登录地点',
                                 `browser` varchar(50) DEFAULT '' COMMENT '浏览器类型',
                                 `os` varchar(50) DEFAULT '' COMMENT '操作系统',
                                 `code` int(11) NOT NULL DEFAULT 200 COMMENT '登录状态（200成功 其它失败）',
                                 `msg` varchar(255) DEFAULT '' COMMENT '提示消息',
                                 `createdTime`     datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
                                 PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT  COMMENT='登录日志管理';

DROP TABLE IF EXISTS `sys_oper_log`;
CREATE TABLE `sys_oper_log` (
    `id` bigint auto_increment COMMENT '编号',
    `operUid`         bigint       NOT NULL COMMENT '用户id',
    `operUserName` varchar(50) DEFAULT '' COMMENT '操作人员名称',
    `operName` varchar(50) DEFAULT '' COMMENT '操作名称',
    `businessType` int(11) NOT NULL COMMENT '业务类型（1新增 2修改 3删除 4查询 5其它）',
    `uri` varchar(100) DEFAULT '' COMMENT '请求地址',
    `operIpAddr` varchar(50) DEFAULT '' COMMENT '主机地址',
    `operLocation` varchar(255) DEFAULT '' COMMENT '操作地点',
    `req` text COMMENT '请求参数',
    `resp` text COMMENT '返回参数',
    `code` int(11) NOT NULL DEFAULT 200 COMMENT '返回状态（200成功 其它失败）',
    `msg` varchar(255) DEFAULT '' COMMENT '提示消息',
    `createdTime`     datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT  COMMENT='操作日志管理';

CREATE TABLE if not exists `sys_api_info`
(
    `id`              bigint auto_increment comment '编号',
    `route`           varchar(100) NOT NULL DEFAULT '' comment '路由',
    `method`          int(11) NOT NULL comment '请求方式（1 GET 2 POST 3 HEAD 4 OPTIONS 5 PUT 6 DELETE 7 TRACE 8 CONNECT 9 其它）',
    `name`            varchar(100) NOT NULL DEFAULT '' comment '请求名称',
    `businessType`    int(11) NOT NULL COMMENT '业务类型（1新增 2修改 3删除 4查询 5其它）',
    `group`           varchar(100) NOT NULL DEFAULT '' comment '接口组',
    `desc`            varchar(100) NOT NULL DEFAULT '' comment '备注',
    `createdTime`     datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` 	  datetime NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` 	  datetime DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `routeIndex` (`route`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='接口管理';

INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/info/update',2,'更新产品',2,'','产品管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/info/create',2,'新增产品',1,'','产品管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/info/read',2,'获取产品详情',4,'','产品管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/info/delete',2,'删除产品',3,'','产品管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/info/index',2,'获取产品列表',4,'','产品管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/custom/read',2,'获取产品自定义信息',4,'','产品自定义信息');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/custom/update',2,'更新产品自定义信息',2,'','产品自定义信息');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/schema/index',2,'获取产品物模型列表',4,'','物模型');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/schema/tsl-import',2,'导入物模型tsl',1,'','物模型');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/schema/tsl-read',2,'获取产品物模型tsl',4,'','物模型');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/schema/create',2,'新增物模型功能',1,'','物模型');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/schema/update',2,'更新物模型功能',2,'','物模型');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/schema/delete',2,'删除物模型功能',3,'','物模型');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/remote-config/create',2,'创建配置',1,'','产品远程配置');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/remote-config/index',2,'获取配置列表',4,'','产品远程配置');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/remote-config/push-all',2,'推送配置',5,'','产品远程配置');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/product/remote-config/lastest-read',2,'获取最新配置',4,'','产品远程配置');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/group/info/create',2,'创建分组',1,'','设备分组');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/group/info/index',2,'获取分组列表',4,'','设备分组');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/group/info/read',2,'获取分组详情信息',4,'','设备分组');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/group/info/update',2,'更新分组信息',2,'','设备分组');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/group/info/delete',2,'删除分组',3,'','设备分组');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/group/device/index',2,'获取分组设备列表',4,'','设备分组');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/group/device/multi-create',2,'添加分组设备(支持批量)',1,'','设备分组');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/group/device/multi-delete',2,'删除分组设备(支持批量)',3,'','设备分组');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/info/index',2,'获取设备列表',4,'','设备管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/info/read',2,'获取设备详情',4,'','设备管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/info/create',2,'新增设备',1,'','设备管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/info/delete',2,'删除设备',3,'','设备管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/info/update',2,'更新设备',2,'','设备管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/info/count',2,'设备统计详情',4,'','设备管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/info/multi-import',2,'批量导入设备',1,'','设备管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/auth/login',2,'设备登录认证',5,'','设备鉴权');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/auth/root-check',2,'鉴定mqtt账号root权限',5,'','设备鉴权');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/auth/access',2,'设备操作认证',5,'','设备鉴权');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/msg/property-log/index',2,'获取单个id属性历史记录',4,'','设备消息');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/msg/sdk-log/index',2,'获取设备本地日志',4,'','设备消息');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/msg/hub-log/index',2,'获取云端诊断日志',4,'','设备消息');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/msg/property-latest/index',2,'获取最新属性',4,'','设备消息');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/msg/event-log/index',2,'获取物模型事件历史记录',4,'','设备消息');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/interact/send-action',2,'同步调用设备行为',5,'','设备交互');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/interact/send-property',2,'同步调用设备属性',5,'','设备交互');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/interact/get-property-reply',2,'请求设备获取设备最新属性',4,'','设备交互');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/interact/send-msg',2,'发送消息给设备',5,'','设备交互');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/gateway/multi-create',2,'批量添加网关子设备',1,'','网关子设备管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/gateway/multi-delete',2,'批量解绑网关子设备',3,'','网关子设备管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/device/gateway/index',2,'获取子设备列表',4,'','网关子设备管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/log/login/index',2,'获取登录日志列表',4,'','日志管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/log/oper/index',2,'获取操作日志列表',4,'','日志管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/role/create',2,'添加角色',1,'','角色管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/role/index',2,'获取角色列表',4,'','角色管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/role/update',2,'更新角色',2,'','角色管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/role/delete',2,'删除角色',3,'','角色管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/role/role-menu/update',2,'更新角色对应菜单列表',2,'','角色管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/menu/create',2,'添加菜单',1,'','菜单管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/menu/index',2,'获取菜单列表',4,'','菜单管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/menu/update',2,'更新菜单',2,'','菜单管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/menu/delete',2,'删除菜单',3,'','菜单管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/user/create',2,'创建用户信息',1,'','用户管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/user/captcha',2,'获取验证码',5,'','用户管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/user/login',2,'登录',5,'','用户管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/user/delete',2,'删除用户',3,'','用户管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/user/read',2,'获取用户信息',4,'','用户管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/user/update',2,'更新用户基本数据',2,'','用户管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/user/index',2,'获取用户信息列表',4,'','用户管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/user/resource-read',2,'获取用户资源',4,'','用户管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/common/config',2,'获取系统配置',4,'','系统配置');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/api/create',2,'添加接口',1,'','接口管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/api/index',2,'获取接口列表',4,'','接口管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/api/update',2,'更新接口',2,'','接口管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/api/delete',2,'删除接口',3,'','接口管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/auth/api/index',2,'获取API权限列表',4,'','权限管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/system/auth/api/multiUpdate',2,'更新API权限',2,'','权限管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/scene/info/read',2,'获取场景信息',4,'','场景联动');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/scene/info/index', 2,'获取场景列表',4,'','场景联动');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/scene/info/create',2,'创建场景信息',1,'','场景联动');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/scene/info/update',2,'更新场景信息',2,'','场景联动');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/scene/info/delete',2,'删除场景信息',3,'','场景联动');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/flow/info/index',2,'获取流列表',4,'','流');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/flow/info/create', 2,'创建流',1,'','流');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/flow/info/update', 2,'修改流',2,'','流');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/flow/info/delete', 2,'删除流',3,'','流');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/info/create',2,'新增告警',1,'','告警管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/info/update',2,'更新告警',2,'','告警管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/info/delete',2,'删除告警',3,'','告警管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/info/index', 2,'获取告警信息列表',4,'','告警管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/info/read',2,'获取告警详情',4,'','告警管理');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/scene/delete',2,'删除告警和场景的关联',3,'','场景联动');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/log/index',2,'获取告警流水日志记录列表',4,'','告警日志');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/record/index',2,'获取告警记录列表',4,'','告警记录');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/deal-record/create',2,'新增告警处理记录',1,'','处理记录');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/deal-record/index',2,'获取告警处理记录列表',4,'','处理记录');
INSERT IGNORE INTO sys_api_info (route, `method`, name, businessType, `desc`, `group`) VALUES('/api/v1/things/rule/alarm/scene/multi-update',2,'更新告警和场景的关联',2,'','场景联动');

CREATE TABLE if not exists `sys_api_auth` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT comment '编号',
    `p_type` varchar(255) NOT NULL DEFAULT '' comment '策略类型，即策略的分类，例如"p"表示主体（provider）访问资源（resource）的许可权，"g"表示主体（provider）之间的关系访问控制',
    `v0` varchar(255) NOT NULL DEFAULT '' comment '策略中的第一个参数，通常用于表示资源的归属范围（即限制访问的对象），例如资源所属的机构、部门、业务线、地域等',
    `v1` varchar(255) NOT NULL DEFAULT '' comment '策略中的第二个参数，通常用于表示主体（provider），即需要访问资源的用户或者服务',
    `v2` varchar(255) NOT NULL DEFAULT '' comment '策略中的第三个参数，通常用于表示资源（resource），即需要进行访问的对象',
    `v3` varchar(255) NOT NULL DEFAULT '' comment '策略中的第四个参数，通常用于表示访问操作（permission），例如 “read”, “write”, “execute” 等',
    `v4` varchar(255) NOT NULL DEFAULT '' comment '策略中的第五个参数，通常用于表示资源的类型（object type），例如表示是文件或者数据库表等',
    `v5` varchar(255) NOT NULL DEFAULT '' comment '策略中的第六个参数，通常用于表示扩展信息，例如 IP 地址、端口号等',
    PRIMARY KEY (`id`),
    UNIQUE KEY `roleId_path_index` (`v0`, `v1`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='api权限管理';

INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/info/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/info/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/info/read',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/info/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/info/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/schema/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/schema/tsl-import',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/schema/tsl-read',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/schema/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/schema/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/schema/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/remote-config/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/remote-config/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/remote-config/push-all',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/remote-config/lastest-read',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/custom/read',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/product/custom/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/group/info/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/group/info/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/group/info/read',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/group/info/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/group/info/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/group/device/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/group/device/multi-create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/group/device/multi-delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/info/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/info/read',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/info/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/info/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/info/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/info/count',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/auth/login',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/auth/root-check',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/auth/access',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/msg/property-log/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/msg/sdk-log/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/msg/hub-log/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/msg/property-latest/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/msg/event-log/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/interact/send-action',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/interact/send-property',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/interact/send-msg',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/gateway/multi-create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/gateway/multi-delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/device/gateway/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/log/login/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/log/oper/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/role/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/role/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/role/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/role/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/role/role-menu/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/menu/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/menu/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/menu/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/menu/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/user/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/user/captcha',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/user/login',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/user/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/user/read',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/user/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/user/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/user/resource-read',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/common/config',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/api/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/api/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/api/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/api/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/auth/api/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/system/auth/api/multiUpdate',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/scene/info/read',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/scene/info/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/scene/info/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/scene/info/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/scene/info/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/flow/info/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/flow/info/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/flow/info/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/flow/info/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/info/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/info/update',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/info/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/info/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/info/read',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/scene/delete',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/log/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/record/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/deal-record/create',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/deal-record/index',2,'','','');
INSERT IGNORE INTO sys_api_auth (p_type, v0, v1, v2, v3, v4, v5) VALUES('p','1','/api/v1/things/rule/alarm/scene/multi-update',2,'','','');
