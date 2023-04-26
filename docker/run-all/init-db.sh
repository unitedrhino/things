#!/bin/bash
sleep 10 # 等待tdengine容器启动
taos -s 'create database if not exists iThings;'  -h 172.19.0.6 -P 6030 -u root -ptaosdata