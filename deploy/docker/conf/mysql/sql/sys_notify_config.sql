use iThings;

create table sys_notify_config
(
    id            bigint auto_increment
        primary key,
    `group`       varchar(50)                                               not null,
    code          varchar(50)                                               not null,
    name          varchar(50)                                               not null,
    support_types longtext collate utf8mb4_bin default '[]'                 not null
        check (json_valid(`support_types`)),
    `desc`        varchar(100)                                              not null,
    is_record     bigint                                                    null,
    params        longtext collate utf8mb4_bin default '{}'                 not null
        check (json_valid(`params`)),
    created_time  datetime(3)                  default current_timestamp(3) not null,
    updated_time  datetime(3)                  default current_timestamp(3) not null,
    created_by    bigint                                                    null,
    deleted_by    bigint                                                    null,
    updated_by    bigint                                                    null,
    deleted_time  bigint                       default 0                    null,
    enable_types  longtext collate utf8mb4_bin default '[]'                 not null
        check (json_valid(`enable_types`)),
    constraint ri_mi
        unique (code, deleted_time)
);

create index idx_sys_notify_config_created_time
    on sys_notify_config (created_time);

INSERT INTO iThings.sys_notify_config (id, `group`, code, name, support_types, `desc`, is_record, params, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, enable_types) VALUES (1, '验证码', 'sysUserRegisterCaptcha', '用户注册验证码', '["sms","email"]', '', 2, '{"code":"验证码","expr":"过期时间"}', '2024-06-17 13:35:58.172', '2024-09-12 14:12:13.056', 0, 0, 1740358057038188544, 0, '[]');
INSERT INTO iThings.sys_notify_config (id, `group`, code, name, support_types, `desc`, is_record, params, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, enable_types) VALUES (2, '验证码', 'sysUserLoginCaptcha', '用户登录验证码', '["sms","email"]', '', 2, '{"code":"验证码","expr":"过期时间"}', '2024-06-17 13:35:58.172', '2024-09-12 14:12:14.088', 0, 0, 1740358057038188544, 0, '[]');
INSERT INTO iThings.sys_notify_config (id, `group`, code, name, support_types, `desc`, is_record, params, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, enable_types) VALUES (3, '场景联动通知', 'ruleScene', '场景联动通知', '["sms","email","dingWebhook","wxEWebHook","wxMini","dingTalk","dingMini"]', '', 1, '{"body":"通知的内容"}', '2024-06-17 13:35:58.172', '2024-09-12 14:35:25.023', 0, 0, 1740358057038188544, 0, '[]');
INSERT INTO iThings.sys_notify_config (id, `group`, code, name, support_types, `desc`, is_record, params, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, enable_types) VALUES (4, '设备', 'ruleDeviceAlarm', '设备告警通知', '["sms","email","dingWebhook"]', '', 1, '{"productID":"产品ID(若为设备触发)","deviceName":"触发设备ID(若为设备触发)","sceneName":"触发场景名称","deviceAlias":"设备名称(若为设备触发)"}', '2024-06-17 13:35:58.172', '2024-06-17 13:35:58.172', 0, 0, 0, 0, '[]');
INSERT INTO iThings.sys_notify_config (id, `group`, code, name, support_types, `desc`, is_record, params, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, enable_types) VALUES (9, '系统公告', 'sysAnnouncement', '系统公告', '["sms","email","wxMini"]', '', 1, '{"body":"内容","title":"标题"}', '2024-09-04 21:49:18.487', '2024-09-04 21:52:27.055', 1740358057038188544, 0, 1740358057038188544, 0, '[]');
INSERT INTO iThings.sys_notify_config (id, `group`, code, name, support_types, `desc`, is_record, params, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, enable_types) VALUES (11, '验证码', 'sysUserChangePwdCaptcha', '用户修改密码', '["sms","email"]', '', 2, '{"code":"验证码","expr":"过期时间"}', '2024-10-15 17:38:37.035', '2024-10-15 21:39:51.210', 1740358057038188544, 0, 1740358057038188544, 0, '["email"]');