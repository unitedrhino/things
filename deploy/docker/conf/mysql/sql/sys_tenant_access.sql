create table sys_tenant_access
(
    id           bigint auto_increment
        primary key,
    tenant_code  varchar(50)                              not null,
    access_code  varchar(50)                              not null,
    created_time datetime(3) default current_timestamp(3) not null,
    updated_time datetime(3) default current_timestamp(3) not null,
    created_by   bigint                                   null,
    deleted_by   bigint                                   null,
    updated_by   bigint                                   null,
    deleted_time bigint      default 0                    null,
    constraint tenant_scope
        unique (tenant_code, access_code, deleted_time)
);

create index idx_sys_tenant_access_created_time
    on sys_tenant_access (created_time);

