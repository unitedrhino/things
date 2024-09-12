create table sys_app_module
(
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
        unique (app_code, module_code, deleted_time)
)
    charset = utf8mb3;

create index idx_sys_app_module_created_time
    on sys_app_module (created_time);

INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (5, 'test1234', 'SystemManage', '2023-12-31 15:10:21.071', '2023-12-31 15:10:21.071', 0, null, null, null);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (6, 'test1234', 'DeviceManage', '2023-12-31 15:10:21.071', '2023-12-31 15:10:21.071', 0, null, null, null);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (7, 'test111', 'SystemManage', '2024-01-02 22:32:07.187', '2024-01-02 22:32:07.187', 0, null, null, null);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (8, 'test111', 'DeviceManage', '2024-01-02 22:32:07.187', '2024-01-02 22:32:07.187', 0, null, null, null);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (9, 'test111', 'outmanager', '2024-01-02 22:32:07.187', '2024-01-02 22:32:07.187', 0, null, null, null);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (19, 'test-app-code', 'SystemManage', '2024-01-04 22:59:08.295', '2024-01-04 22:59:08.295', 0, null, null, null);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (20, 'test-app-code', 'DeviceManage', '2024-01-04 22:59:08.295', '2024-01-04 22:59:08.295', 0, null, null, null);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (21, 'test-app-code', 'outmanager', '2024-01-04 22:59:08.295', '2024-01-04 22:59:08.295', 0, null, null, null);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (25, 'app-code-2', 'SystemManage', '2024-01-06 14:22:19.651', '2024-01-06 14:22:19.651', 0, null, null, null);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (26, 'app-code-2', 'outmanager', '2024-01-06 14:22:19.651', '2024-01-06 14:22:19.651', 0, null, null, null);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (94, 'web-app-test', 'iframe-ithings', '2024-02-04 21:18:45.599', '2024-02-04 21:18:45.599', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (110, 'apptest', 'systemManage', '2024-02-21 22:14:48.486', '2024-02-21 22:14:48.486', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (111, 'apptest', 'antdv', '2024-02-21 22:14:48.486', '2024-02-21 22:14:48.486', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (112, 'apptest', 'tenantManage', '2024-02-21 22:14:48.486', '2024-02-21 22:14:48.486', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (116, 'IntelligentLighting', 'systemManage', '2024-03-01 18:23:28.233', '2024-03-01 18:23:28.233', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (117, 'IntelligentLighting', 'IntelligentLightingWeb', '2024-03-01 18:23:28.233', '2024-03-01 18:23:28.233', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (118, 'IntelligentLighting', 'Intelligent-lighting-mini', '2024-03-01 18:23:28.233', '2024-03-01 18:23:28.233', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (128, 'all', 'userManage', '2024-04-01 11:53:57.554', '2024-04-01 11:53:57.554', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (135, 'all', 'supperManage', '2024-04-01 15:06:13.860', '2024-04-01 15:06:13.860', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (143, 'all', 'operationManagement', '2024-04-01 15:42:57.115', '2024-04-01 15:42:57.115', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (184, 'all', 'platformManage', '2024-04-06 13:43:56.110', '2024-04-06 13:43:56.110', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (192, 'all', 'sale', '2024-06-03 19:14:14.041', '2024-06-03 19:14:14.041', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (204, 'all', 'myThings', '2024-06-27 20:49:40.914', '2024-06-27 20:49:40.914', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (210, 'all', 'iotOverview', '2024-09-02 20:39:21.936', '2024-09-02 20:39:21.936', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (217, 'swxjj-mini-wx', 'things', '2024-09-07 13:29:02.051', '2024-09-07 13:29:02.051', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (218, 'core', 'things', '2024-09-12 14:09:33.447', '2024-09-12 14:09:33.447', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (219, 'core', 'systemManage', '2024-09-12 14:09:33.447', '2024-09-12 14:09:33.447', 0, 1740358057038188544, 0, 0);
INSERT INTO iThings.sys_app_module (id, app_code, module_code, created_time, updated_time, deleted_time, created_by, deleted_by, updated_by) VALUES (220, 'core', 'myThings', '2024-09-12 14:09:33.447', '2024-09-12 14:09:33.447', 0, 1740358057038188544, 0, 0);
