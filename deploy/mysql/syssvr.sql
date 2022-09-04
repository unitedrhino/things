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
    `backgroundUrl`  varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' comment '后台地址',
    `createdTime`     datetime not NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updatedTime` 	  datetime NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '更新时间',
    `deletedTime` 	  datetime DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `nameIndex` (`name`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='菜单管理表';

