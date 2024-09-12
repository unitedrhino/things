create table sys_tenant_app
(
    id                     bigint auto_increment
        primary key,
    tenant_code            varchar(50)                               not null,
    app_code               varchar(50)                               not null,
    created_time           datetime(3)  default current_timestamp(3) not null,
    updated_time           datetime(3)  default current_timestamp(3) not null,
    deleted_time           bigint       default 0                    null,
    created_by             bigint                                    null,
    deleted_by             bigint                                    null,
    updated_by             bigint                                    null,
    mini_ding_app_id       varchar(50)  default ''                   null,
    mini_ding_app_key      varchar(50)  default ''                   null,
    mini_ding_app_secret   varchar(200) default ''                   null,
    mini_wx_app_id         varchar(50)  default ''                   null,
    mini_wx_app_key        varchar(50)  default ''                   null,
    mini_wx_app_secret     varchar(200) default ''                   null,
    mini_ding_mini_app_id  varchar(50)  default ''                   null,
    mini_wx_mini_app_id    varchar(50)  default ''                   null,
    official_wxapp_id      varchar(50)  default ''                   null,
    official_wxmini_app_id varchar(50)  default ''                   null,
    official_wxapp_key     varchar(50)  default ''                   null,
    official_wxapp_secret  varchar(200) default ''                   null,
    login_types            longtext collate utf8mb4_bin              null
        check (json_valid(`login_types`)),
    ding_mini_app_id       varchar(50)  default ''                   null,
    ding_mini_mini_app_id  varchar(50)  default ''                   null,
    ding_mini_app_key      varchar(50)  default ''                   null,
    ding_mini_app_secret   varchar(200) default ''                   null,
    wx_mini_app_id         varchar(50)  default ''                   null,
    wx_mini_app_key        varchar(50)  default ''                   null,
    wx_mini_app_secret     varchar(200) default ''                   null,
    wx_open_app_id         varchar(50)  default ''                   null,
    wx_open_mini_app_id    varchar(50)  default ''                   null,
    wx_open_app_key        varchar(50)  default ''                   null,
    wx_open_app_secret     varchar(200) default ''                   null,
    wx_mini_mini_app_id    varchar(50)  default ''                   null,
    is_auto_register       bigint       default 1                    null,
    android_version        varchar(64)                               null,
    android_file_path      varchar(256)                              null,
    android_version_desc   varchar(100)                              null,
    constraint tc_ac
        unique (tenant_code, app_code, deleted_time)
)
    charset = utf8mb3;

create index idx_sys_tenant_app_created_time
    on sys_tenant_app (created_time);

INSERT INTO iThings.sys_tenant_app (id, tenant_code, app_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by, mini_ding_app_id, mini_ding_app_key, mini_ding_app_secret, mini_wx_app_id, mini_wx_app_key, mini_wx_app_secret, mini_ding_mini_app_id, mini_wx_mini_app_id, official_wxapp_id, official_wxmini_app_id, official_wxapp_key, official_wxapp_secret, login_types, ding_mini_app_id, ding_mini_mini_app_id, ding_mini_app_key, ding_mini_app_secret, wx_mini_app_id, wx_mini_app_key, wx_mini_app_secret, wx_open_app_id, wx_open_mini_app_id, wx_open_app_key, wx_open_app_secret, wx_mini_mini_app_id, is_auto_register, android_version, android_file_path, android_version_desc) VALUES (65, 'default', 'core', '2024-03-02 12:52:34.037', '2024-08-24 15:21:30.322', 0, 1740358057038188544, 0, 1740358057038188544, '', '', '', 'xxx', '', 'xxx', '', '', '', '', '', '', '["phone","pwd","wxOpen"]', '', '', '', '', '', '', '', 'xxx', '', '', 'xxx', '', 1, null, null, null);
INSERT INTO iThings.sys_tenant_app (id, tenant_code, app_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by, mini_ding_app_id, mini_ding_app_key, mini_ding_app_secret, mini_wx_app_id, mini_wx_app_key, mini_wx_app_secret, mini_ding_mini_app_id, mini_wx_mini_app_id, official_wxapp_id, official_wxmini_app_id, official_wxapp_key, official_wxapp_secret, login_types, ding_mini_app_id, ding_mini_mini_app_id, ding_mini_app_key, ding_mini_app_secret, wx_mini_app_id, wx_mini_app_key, wx_mini_app_secret, wx_open_app_id, wx_open_mini_app_id, wx_open_app_key, wx_open_app_secret, wx_mini_mini_app_id, is_auto_register, android_version, android_file_path, android_version_desc) VALUES (73, 'default', 'client-mini-wx', '2024-05-30 15:40:59.033', '2024-08-24 15:21:04.135', 0, 1740358057038188544, 0, 1740358057038188544, '', '', '', 'xxx', '', 'xxx', '', '', '', '', '', '', '["phone","email","wxMiniP","pwd"]', '', '', '', '', 'xxx', '', 'xxx', '', '', '', '', '', 1, null, null, null);
