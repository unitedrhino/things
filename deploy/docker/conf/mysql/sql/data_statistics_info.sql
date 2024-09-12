use iThings;

create table data_statistics_info
(
    id                bigint auto_increment
        primary key,
    is_filter_tenant  bigint        default 1                    not null,
    is_filter_project bigint        default 1                    not null,
    is_filter_area    bigint        default 1                    not null,
    code              varchar(120)                               not null,
    type              varchar(120)                               not null,
    `table`           varchar(120)  default ''                   null,
    columns           varchar(512)  default ''                   null,
    arg_columns       longtext      default '{}'                 not null
        check (json_valid(`arg_columns`)),
    omits             varchar(120)  default ''                   null,
    is_to_hump        bigint        default 1                    not null,
    `sql`             varchar(2000) default ''                   null,
    order_by          varchar(120)  default 'created_time desc'  null,
    filter            longtext      default '{}'                 not null
        check (json_valid(`filter`)),
    created_time      datetime(3)   default current_timestamp(3) not null,
    updated_time      datetime(3)   default current_timestamp(3) not null,
    created_by        bigint                                     null,
    deleted_by        bigint                                     null,
    updated_by        bigint                                     null,
    deleted_time      bigint        default 0                    null,
    is_soft_delete    bigint        default 1                    not null,
    filter_slot_code  varchar(120)  default ''                   null,
    filter_roles      varchar(120)  default ''                   null,
    constraint `key`
        unique (code, deleted_time)
)
    collate = utf8mb4_bin;

create index idx_data_statistics_info_created_time
    on data_statistics_info (created_time);

INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (1, 1, 1, 1, 'dmDevicePower', 'table', 'data_dm_device_power', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{"endDate":{"sql":"?\\u003c=date","valNum":1,"type":"date"},"startDate":{"sql":"?\\u003e=date","valNum":1,"type":"date"}}', '2024-03-31 23:49:25.734', '2024-03-31 23:49:25.734', 0, 0, 0, 0, 1, '', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (2, 1, 1, 2, 'dmDeviceCount', 'table', 'dm_device_info', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{"areas": {"sql": "area_id in ?","valNum": 1,"type": "array"}}', '2024-04-18 10:39:28.304', '2024-04-18 10:39:28.304', null, null, null, 0, 1, '', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (3, 2, 2, 2, 'sysOpsWorkOrder', 'table', 'sys_ops_work_order', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{}', '2024-04-18 10:39:28.304', '2024-04-18 10:39:28.304', null, null, null, 0, 1, '', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (4, 2, 2, 2, 'sysUserAreaApply', 'table', 'sys_user_area_apply', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{}', '2024-04-18 10:39:28.304', '2024-04-18 10:39:28.304', null, null, null, 0, 1, '', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (6, 1, 2, 2, 'dmDeviceCountDistributor', 'table', 'dm_device_info', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{"areas": {"sql": "area_id in ?","valNum": 1,"type": "array"}}', '2024-04-18 10:39:28.304', '2024-04-18 10:39:28.304', null, null, null, 0, 1, 'distributor', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (7, 2, 2, 2, 'dmProductCount', 'table', 'dm_product_info', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{}', '2024-04-18 10:39:28.304', '2024-04-18 10:39:28.304', null, null, null, 0, 1, '', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (8, 1, 2, 2, 'saleDistributorCount', 'table', 'sale_distributor_info', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{}', '2024-04-18 10:39:28.304', '2024-04-18 10:39:28.304', null, null, null, 0, 1, 'distributor', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (9, 1, 2, 2, 'saleDistributorApplyCount', 'table', 'sale_distributor_apply', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{}', '2024-04-18 10:39:28.304', '2024-04-18 10:39:28.304', null, null, null, 0, 1, 'distributor', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (10, 1, 2, 2, 'saleDistributorWaterCount', 'table', 'sale_distributor_water', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{}', '2024-04-18 10:39:28.000', '2024-04-18 10:39:28.000', null, null, null, 0, 2, 'distributor', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (11, 1, 2, 2, 'sysUserInfo', 'table', 'sys_user_info', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{}', '2024-04-18 10:39:28.304', '2024-04-18 10:39:28.304', null, null, null, 0, 1, '', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (12, 1, 2, 2, 'saleOrderInfoCount', 'table', 'sale_order_info', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{}', '2024-04-18 10:39:28.000', '2024-04-18 10:39:28.000', null, null, null, 0, 2, 'distributor', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (13, 2, 2, 2, 'dmDeviceMsgCount', 'table', 'dm_device_msg_count', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{}', '2024-09-02 21:08:19.432', '2024-09-02 21:08:19.432', 1740358057038188544, 0, 0, 0, 2, '', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (14, 1, 2, 2, 'sysLoginLog', 'table', 'sys_login_log', '', '{}', 'created_time,updated_time', 2, '', 'created_time desc', '{}', '2024-09-02 21:08:59.004', '2024-09-02 21:08:59.004', 1740358057038188544, 0, 0, 0, 1, '', '');
INSERT INTO iThings.data_statistics_info (id, is_filter_tenant, is_filter_project, is_filter_area, code, type, `table`, columns, arg_columns, omits, is_to_hump, `sql`, order_by, filter, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_soft_delete, filter_slot_code, filter_roles) VALUES (15, 1, 2, 2, 'sysOperLog', 'table', 'sys_oper_log', '', '{}', 'created_time,updated_time', 1, '', 'created_time desc', '{}', '2024-09-02 21:09:35.066', '2024-09-02 21:09:35.066', 1740358057038188544, 0, 0, 0, 2, '', '');
