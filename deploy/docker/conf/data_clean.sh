#!/bin/bash
rm -rf oss/*
rm -rf pgsql/data/*
rm -rf pgsql/log/*
rm -rf redis/data/*
rm -rf taos/dnode/data/*
rm -rf taos/dnode/log/*
rm -rf minio/data/*
rm -rf minio/config/*
rm -rf mysql/data/*
echo "清理完成"
