use iThings;

create table sys_app_info
(
    id                  bigint auto_increment
        primary key,
    code                varchar(50)                               not null,
    name                varchar(100)                              not null,
    `desc`              varchar(100)                              not null,
    base_url            varchar(100)                              null,
    logo_url            varchar(100)                              null,
    created_time        datetime(3)  default current_timestamp(3) not null,
    updated_time        datetime(3)  default current_timestamp(3) not null,
    deleted_time        bigint       default 0                    null,
    type                varchar(100) default 'web'                not null,
    created_by          bigint                                    null,
    deleted_by          bigint                                    null,
    updated_by          bigint                                    null,
    sub_type            varchar(100) default 'wx'                 not null,
    mini_wx_app_id      varchar(50)  default ''                   null,
    mini_wx_app_key     varchar(50)  default ''                   null,
    mini_wx_app_secret  varchar(200) default ''                   null,
    mini_wx_mini_app_id varchar(50)  default ''                   null,
    constraint code
        unique (code, deleted_time),
    constraint name
        unique (name, deleted_time)
)
    charset = utf8mb3;

create index idx_sys_app_info_created_time
    on sys_app_info (created_time);

INSERT INTO iThings.sys_app_info (id, code, name, `desc`, base_url, logo_url, created_time, updated_time, deleted_time, type, created_by, deleted_by, updated_by, sub_type, mini_wx_app_id, mini_wx_app_key, mini_wx_app_secret, mini_wx_mini_app_id) VALUES (1, 'core', '管理后台', ' ', '', '', '2023-12-29 19:31:30.838', '2024-09-07 13:56:00.747', 0, 'web', 0, 0, 1740358057038188544, 'web', '', '', '', '');
INSERT INTO iThings.sys_app_info (id, code, name, `desc`, base_url, logo_url, created_time, updated_time, deleted_time, type, created_by, deleted_by, updated_by, sub_type, mini_wx_app_id, mini_wx_app_key, mini_wx_app_secret, mini_wx_mini_app_id) VALUES (12, 'client-mini-wx', 'c端微信小程序', 'c端微信小程序', 'c端微信小程序', 'c端微信小程序', '2024-05-07 18:47:10.679', '2024-07-26 18:39:54.428', 0, 'mini', 1740358057038188544, 0, 1740358057038188544, 'wx', '', '', '', '');
