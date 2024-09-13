#!/bin/bash
sleep 10 # 等待tdengine容器启动
apk add curl
curl --location --request POST 'http://172.20.0.1:6041/rest/sql/ithings' \
--header 'Authorization: Basic cm9vdDp0YW9zZGF0YQ==' \
--header 'Content-Type: text/plain' \
--data-raw 'create database if not exists ithings;'