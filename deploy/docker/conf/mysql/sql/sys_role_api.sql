use iThings;
create table sys_role_api
(
    id           bigint auto_increment
        primary key,
    tenant_code  varchar(50)                              not null,
    role_id      bigint                                   not null,
    app_code     varchar(50)                              not null,
    module_code  varchar(50)                              not null,
    api_id       bigint                                   not null,
    created_time datetime(3) default current_timestamp(3) not null,
    updated_time datetime(3) default current_timestamp(3) not null,
    deleted_time bigint unsigned                          null,
    created_by   bigint                                   null,
    deleted_by   bigint                                   null,
    updated_by   bigint                                   null,
    constraint ri_mi
        unique (tenant_code, role_id, app_code, module_code, api_id, deleted_time)
)
    charset = utf8mb3;

create index idx_sys_role_api_created_time
    on sys_role_api (created_time);

create index idx_sys_tenant_role_api_created_time
    on sys_role_api (created_time);

INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (228, 'default', 1, 'core', 'SystemManage', 1, '2024-01-06 21:34:54.455', '2024-01-06 21:34:54.455', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (229, 'default', 1, 'core', 'SystemManage', 3, '2024-01-06 21:34:54.455', '2024-01-06 21:34:54.455', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (230, 'default', 1, 'core', 'SystemManage', 4, '2024-01-06 21:34:54.455', '2024-01-06 21:34:54.455', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (231, 'default', 1, 'core', 'SystemManage', 5, '2024-01-06 21:34:54.455', '2024-01-06 21:34:54.455', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (232, 'default', 1, 'core', 'SystemManage', 7, '2024-01-06 21:34:54.455', '2024-01-06 21:34:54.455', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (233, 'default', 1, 'core', 'SystemManage', 8, '2024-01-06 21:34:54.455', '2024-01-06 21:34:54.455', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (234, 'default', 1, 'core', 'SystemManage', 9, '2024-01-06 21:34:54.455', '2024-01-06 21:34:54.455', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (235, 'default', 1, 'core', 'SystemManage', 10, '2024-01-06 21:34:54.455', '2024-01-06 21:34:54.455', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (236, 'default', 1, 'core', 'SystemManage', 11, '2024-01-06 21:34:54.455', '2024-01-06 21:34:54.455', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (237, 'default', 1, 'core', 'SystemManage', 12, '2024-01-06 21:34:54.455', '2024-01-06 21:34:54.455', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (238, 'default', 1, 'test-app-code', 'SystemManage', 1, '2024-01-06 21:49:00.033', '2024-01-06 21:49:00.033', 0, null, null, null);
INSERT INTO sys_role_api (id, tenant_code, role_id, app_code, module_code, api_id, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (239, 'default', 1, 'test-app-code', 'DeviceManage', 2, '2024-01-06 21:49:04.119', '2024-01-06 21:49:04.119', 0, null, null, null);
