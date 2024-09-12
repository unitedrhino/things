use iThings;
create table sys_role_app
(
    id           bigint auto_increment
        primary key,
    tenant_code  varchar(50)                              not null,
    role_id      bigint                                   not null,
    app_code     varchar(50)                              not null,
    created_time datetime(3) default current_timestamp(3) not null,
    updated_time datetime(3) default current_timestamp(3) not null,
    deleted_time bigint      default 0                    null,
    created_by   bigint                                   null,
    deleted_by   bigint                                   null,
    updated_by   bigint                                   null,
    constraint tc_ac
        unique (tenant_code, role_id, app_code, deleted_time)
)
    charset = utf8mb3;

create index idx_sys_role_app_created_time
    on sys_role_app (created_time);

create index idx_sys_tenant_role_app_created_time
    on sys_role_app (created_time);

INSERT INTO sys_role_app (id, tenant_code, role_id, app_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (84, 'default', 11, 'core', '2024-03-04 21:46:08.510', '2024-03-04 21:46:08.510', 0, 1740358057038188544, 0, 0);
INSERT INTO sys_role_app (id, tenant_code, role_id, app_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (85, 'default', 12, 'core', '2024-03-04 21:47:06.231', '2024-03-04 21:47:06.231', 0, 1740358057038188544, 0, 0);
INSERT INTO sys_role_app (id, tenant_code, role_id, app_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (87, 'default', 1, 'core', '2024-03-31 23:58:41.987', '2024-03-31 23:58:41.987', 0, 1740358057038188544, 0, 0);
INSERT INTO sys_role_app (id, tenant_code, role_id, app_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (89, 'default', 2, 'core', '2024-07-14 11:52:46.653', '2024-07-14 11:52:46.653', 0, 1740358057038188544, 0, 0);
INSERT INTO sys_role_app (id, tenant_code, role_id, app_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (91, 'default', 55, 'client-mini-wx', '2024-08-24 15:23:14.167', '2024-08-24 15:23:14.167', 0, 1740358057038188544, 0, 0);
INSERT INTO sys_role_app (id, tenant_code, role_id, app_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (92, 'default', 55, 'core', '2024-08-24 15:23:14.167', '2024-08-24 15:23:14.167', 0, 1740358057038188544, 0, 0);
INSERT INTO sys_role_app (id, tenant_code, role_id, app_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (93, 'default', 54, 'core', '2024-08-28 22:01:52.073', '2024-08-28 22:01:52.073', 0, 1740358057038188544, 0, 0);
INSERT INTO sys_role_app (id, tenant_code, role_id, app_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (94, 'default', 54, 'client-mini-wx', '2024-08-28 22:01:52.073', '2024-08-28 22:01:52.073', 0, 1740358057038188544, 0, 0);
