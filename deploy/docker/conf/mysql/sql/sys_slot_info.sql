use iThings;

create table sys_slot_info
(
    id           bigint auto_increment
        primary key,
    code         varchar(100)                                              not null,
    sub_code     varchar(100)                                              not null,
    slot_code    varchar(100)                                              not null,
    method       varchar(50)                                               not null,
    uri          varchar(100)                                              not null,
    hosts        longtext collate utf8mb4_bin default '[]'                 not null
        check (json_valid(`hosts`)),
    body         varchar(100)                                              not null,
    handler      longtext collate utf8mb4_bin default '{}'                 not null
        check (json_valid(`handler`)),
    auth_type    varchar(100)                                              not null,
    `desc`       varchar(500)                                              null,
    created_time datetime(3)                  default current_timestamp(3) not null,
    created_by   bigint                                                    null,
    updated_time datetime(3)                  default current_timestamp(3) not null,
    updated_by   bigint                                                    null,
    deleted_time bigint                       default 0                    null,
    deleted_by   bigint                                                    null,
    constraint code_slot
        unique (code, sub_code, slot_code)
);

create index idx_sys_slot_info_created_time
    on sys_slot_info (created_time);

create index idx_sys_slot_info_deleted_time
    on sys_slot_info (deleted_time);

INSERT INTO iThings.sys_slot_info (id, code, sub_code, slot_code, method, uri, hosts, body, handler, auth_type, `desc`, created_time, created_by, updated_time, updated_by, deleted_time, deleted_by) VALUES (1, 'areaInfo', 'create', 'ithings', 'POST', '/api/v1/things/slot/area/create', '["http://things:7788"]', '{"projectID":"{{.ProjectID}}","areaID":"{{.AreaID}}","parentAreaID":"{{.ParentAreaID}}"}', '{}', 'core', '', '2024-04-19 17:48:48.262', 0, '2024-04-19 17:48:48.262', 0, 0, 0);
INSERT INTO iThings.sys_slot_info (id, code, sub_code, slot_code, method, uri, hosts, body, handler, auth_type, `desc`, created_time, created_by, updated_time, updated_by, deleted_time, deleted_by) VALUES (2, 'areaInfo', 'delete', 'ithings', 'POST', '/api/v1/things/slot/area/delete', '["http://things:7788"]', '{"projectID":"{{.ProjectID}}","areaID":"{{.AreaID}}","parentAreaID":"{{.ParentAreaID}}"}', '{}', 'core', '', '2024-04-19 17:48:48.262', 0, '2024-04-19 17:48:48.262', 0, 0, 0);
INSERT INTO iThings.sys_slot_info (id, code, sub_code, slot_code, method, uri, hosts, body, handler, auth_type, `desc`, created_time, created_by, updated_time, updated_by, deleted_time, deleted_by) VALUES (3, 'userSubscribe', 'devicePropertyReport', 'ithings', 'POST', '/api/v1/things/slot/user/subscribe', '["http://things:7788"]', '', '{}', 'core', '', '2024-04-19 17:48:48.262', 0, '2024-04-19 17:48:48.262', 0, 0, 0);
INSERT INTO iThings.sys_slot_info (id, code, sub_code, slot_code, method, uri, hosts, body, handler, auth_type, `desc`, created_time, created_by, updated_time, updated_by, deleted_time, deleted_by) VALUES (4, 'userSubscribe', 'deviceConn', 'ithings', 'POST', '/api/v1/things/slot/user/subscribe', '["http://things:7788"]', '', '{}', 'core', '', '2024-04-19 17:48:48.262', 0, '2024-04-19 17:48:48.262', 0, 0, 0);
INSERT INTO iThings.sys_slot_info (id, code, sub_code, slot_code, method, uri, hosts, body, handler, auth_type, `desc`, created_time, created_by, updated_time, updated_by, deleted_time, deleted_by) VALUES (219, 'userSubscribe', 'deviceOtaReport', 'ithings', 'POST', '/api/v1/things/slot/user/subscribe', '["http://things:7788"]', '', '{}', 'core', '', '2024-04-19 17:48:48.000', 0, '2024-04-19 17:48:48.000', 0, 0, 0);
