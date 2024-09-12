use iThings;
create table dm_product_category
(
    id           int auto_increment
        primary key,
    name         varchar(100)                             not null,
    `desc`       varchar(200)                             null,
    head_img     varchar(200)                             null,
    parent_id    bigint                                   not null,
    id_path      varchar(100)                             not null,
    created_time datetime(3) default current_timestamp(3) not null,
    updated_time datetime(3) default current_timestamp(3) not null,
    created_by   bigint                                   null,
    deleted_by   bigint                                   null,
    updated_by   bigint                                   null,
    deleted_time bigint      default 0                    null,
    is_leaf      bigint      default 1                    not null,
    device_count bigint      default 0                    null,
    constraint pn
        unique (name, deleted_time)
)
    collate = utf8mb4_bin;

create index idx_dm_product_category_created_time
    on dm_product_category (created_time);

INSERT INTO dm_product_category (id, name, `desc`, head_img, parent_id, id_path, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_leaf, device_count) VALUES (3, '照明设备', '', '', 1, '3-', '2024-03-31 23:49:25.933', '2024-09-11 20:32:37.075', 0, 0, 0, 0, 1, 0);
INSERT INTO dm_product_category (id, name, `desc`, head_img, parent_id, id_path, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_leaf, device_count) VALUES (4, '空调设备', '', '', 1, '4-', '2024-03-31 23:49:25.933', '2024-09-11 20:32:37.084', 0, 0, 0, 0, 1, 0);
INSERT INTO dm_product_category (id, name, `desc`, head_img, parent_id, id_path, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_leaf, device_count) VALUES (5, '风扇设备', '', '', 1, '5-', '2024-03-31 23:49:25.933', '2024-09-11 20:32:37.092', 0, 0, 0, 0, 1, 0);
INSERT INTO dm_product_category (id, name, `desc`, head_img, parent_id, id_path, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_leaf, device_count) VALUES (6, '传感器设备', '', '', 1, '6-', '2024-03-31 23:49:25.933', '2024-09-11 20:32:37.105', 0, 0, 0, 0, 1, 0);
INSERT INTO dm_product_category (id, name, `desc`, head_img, parent_id, id_path, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_leaf, device_count) VALUES (7, '4G直连', '', '', 1, '7-', '2024-04-14 15:08:27.464', '2024-09-11 20:32:37.119', 1740358057038188544, 0, 0, 0, 1, 0);
INSERT INTO dm_product_category (id, name, `desc`, head_img, parent_id, id_path, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_leaf, device_count) VALUES (8, '4G网关', '', '', 1, '8-', '2024-04-14 15:08:37.429', '2024-09-11 20:32:37.127', 1740358057038188544, 0, 0, 0, 1, 879);
INSERT INTO dm_product_category (id, name, `desc`, head_img, parent_id, id_path, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_leaf, device_count) VALUES (9, '安防监控', '', '', 1, '9-', '2024-05-19 22:47:32.900', '2024-09-11 20:32:37.130', 1740358057038188544, 0, 0, 0, 1, 0);
INSERT INTO dm_product_category (id, name, `desc`, head_img, parent_id, id_path, created_time, updated_time, created_by, deleted_by, updated_by, deleted_time, is_leaf, device_count) VALUES (10, 'wifi网关', '', '', 1, '10-', '2024-05-19 23:08:40.679', '2024-09-11 20:32:37.132', 1740358057038188544, 0, 0, 0, 1, 8);
