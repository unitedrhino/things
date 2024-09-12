use iThings;

create table sys_tenant_app_module
(
    tenant_code  varchar(50)                              not null,
    id           bigint auto_increment
        primary key,
    app_code     varchar(50)                              not null,
    module_code  varchar(50)                              not null,
    created_time datetime(3) default current_timestamp(3) not null,
    updated_time datetime(3) default current_timestamp(3) not null,
    deleted_time bigint      default 0                    null,
    created_by   bigint                                   null,
    deleted_by   bigint                                   null,
    updated_by   bigint                                   null,
    constraint tc_ac
        unique (tenant_code, app_code, module_code, deleted_time)
)
    charset = utf8mb3;

create index idx_sys_tenant_app_module_created_time
    on sys_tenant_app_module (created_time);

INSERT INTO iThings.sys_tenant_app_module (tenant_code, id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES ('default', 163, 'core', 'systemManage', '2024-03-31 23:53:25.905', '2024-03-31 23:53:25.905', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_tenant_app_module (tenant_code, id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES ('default', 173, 'core', 'things', '2024-04-13 15:59:35.068', '2024-04-13 15:59:35.068', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_tenant_app_module (tenant_code, id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES ('default', 188, 'client-mini-wx', 'things', '2024-05-30 15:40:59.042', '2024-05-30 15:40:59.042', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_tenant_app_module (tenant_code, id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES ('default', 191, 'client-mini-wx', 'systemManage', '2024-05-30 15:40:59.050', '2024-05-30 15:40:59.050', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_tenant_app_module (tenant_code, id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES ('default', 198, 'core', 'myThings', '2024-06-27 20:55:14.201', '2024-06-27 20:55:14.201', 0, 1740358057038188544, 0, 0);
