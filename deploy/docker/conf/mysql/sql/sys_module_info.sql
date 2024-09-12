create table sys_module_info
(
    id           bigint auto_increment
        primary key,
    code         varchar(50) charset utf8mb3              not null,
    type         bigint      default 1                    not null,
    order_num    bigint      default 1                    not null,
    name         varchar(50) charset utf8mb3              not null,
    path         varchar(64) charset utf8mb3              not null,
    url          varchar(200) charset utf8mb3             not null,
    icon         varchar(64) charset utf8mb3              not null,
    body         varchar(1024) charset utf8mb3            null,
    hide_in_menu bigint      default 2                    not null,
    `desc`       varchar(100) charset utf8mb3             not null,
    created_time datetime(3) default current_timestamp(3) not null,
    updated_time datetime(3) default current_timestamp(3) not null,
    deleted_time bigint      default 0                    null,
    created_by   bigint                                   null,
    deleted_by   bigint                                   null,
    updated_by   bigint                                   null,
    sub_type     bigint      default 1                    not null,
    `order`      bigint      default 1                    not null,
    tag          bigint      default 1                    not null,
    constraint code
        unique (code, deleted_time)
)
    collate = utf8mb4_bin;

create index idx_sys_module_info_created_time
    on sys_module_info (created_time);

INSERT INTO iThings.sys_module_info (id, code, type, order_num, name, path, url, icon, body, hide_in_menu, `desc`, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by, sub_type, `order`, tag) VALUES (15, 'systemManage', 1, 1, '系统管理', 'system', '', 'icon-menu-xitong', '{}', 2, '', '2024-02-18 22:13:31.502', '2024-09-03 23:21:05.519', 0, 1740358057038188544, 0, 1740358057038188544, 3, 2, 1);
INSERT INTO iThings.sys_module_info (id, code, type, order_num, name, path, url, icon, body, hide_in_menu, `desc`, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by, sub_type, `order`, tag) VALUES (18, 'things', 1, 1, '物联网', 'things', '/app/things', 'icon-menu-yingyong2', '{"microAppUrl":"/app/things","microAppName":"物联网","microAppBaseroute":"things"}', 2, '', '2024-03-01 21:09:12.557', '2024-04-01 15:26:45.063', 0, 1740358057038188544, 0, 1740358057038188544, 1, 1, 1);
INSERT INTO iThings.sys_module_info (id, code, type, order_num, name, path, url, icon, body, hide_in_menu, `desc`, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by, sub_type, `order`, tag) VALUES (24, 'myThings', 1, 1, '我的物联', 'myThings', '/app/things', 'icon-menu-haoyou', '{"microAppUrl":"/app/things","microAppName":"我的物联","microAppBaseroute":"myThings"}', 2, '', '2024-06-27 20:49:40.913', '2024-09-02 21:13:47.497', 0, 1740358057038188544, 0, 1740358057038188544, 1, 99, 1);
