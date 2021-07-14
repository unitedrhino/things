create database user;
use user;
CREATE TABLE `user_core` (
  `uid` bigint unsigned NOT NULL COMMENT '用户id',
  `userName` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '登录用户名',
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '登录密码',
  `email` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `phone` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `wechat` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '微信union id',
  `lastIP` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '最后登录ip',
  `regIP` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '注册ip',
  `createdTime` datetime DEFAULT NULL,
  `updatedTime` datetime DEFAULT NULL,
  `deletedTime` datetime DEFAULT NULL,
  `status` int DEFAULT '0' COMMENT '用户状态:0为未注册状态',
  PRIMARY KEY (`uid`) USING BTREE,
  KEY `user_username` (`userName`) USING BTREE,
  KEY `user_phone` (`phone`) USING BTREE,
  KEY `user_email` (`email`) USING BTREE,
  KEY `user_wechat` (`wechat`) USING BTREE,
  KEY `user_deletedTime` (`deletedTime`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='用户登录核心信息表';

CREATE TABLE `user_info` (
  `uid` bigint unsigned NOT NULL COMMENT '用户id',
  `userName` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `nickName` varchar(100) NOT NULL DEFAULT '' COMMENT '用户的昵称',
  `inviterUid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '邀请人用户id',
  `inviterId` varbinary(64) NOT NULL DEFAULT '' COMMENT '邀请码',
  `sex` int NOT NULL DEFAULT '0' COMMENT '用户的性别，值为1时是男性，值为2时是女性，值为0时是未知',
  `city` varchar(20) NOT NULL DEFAULT '' COMMENT '用户所在城市',
  `country` varchar(20) NOT NULL DEFAULT '' COMMENT '用户所在国家',
  `province` varchar(20) NOT NULL DEFAULT '' COMMENT '用户所在省份',
  `language` varchar(20) NOT NULL DEFAULT '' COMMENT '用户的语言，简体中文为zh_CN',
  `headImgUrl` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户头像',
  `createdTime` datetime DEFAULT NULL,
  `updatedTime` datetime DEFAULT NULL,
  `deletedTime` datetime DEFAULT NULL,
  PRIMARY KEY (`uid`) USING BTREE,
  KEY `user_inviterUid` (`inviterUid`) USING BTREE,
  KEY `user_deletedTime` (`deletedTime`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='用户详细信息表'