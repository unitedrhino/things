#!/bin/bash

function init_mysql_db_table(){
 for (( i=0; i<300; i++)); do
   check_result=$(docker ps |grep mariadb)
   if [ -n "$check_result" ];then
       docker exec -it mariadb /bin/bash -c 'mysql -uroot -ppassword < /mysql/dmsvr.sql'
       docker exec -it mariadb /bin/bash -c 'mysql -uroot -ppassword < /mysql/syssvr.sql'
       sleep 5
       echo "has install mysql"
       break
   else
       echo "not install mysql, please make sure docker mysql is running"
       sleep 2
       continue
   fi
  done
}

# 初始化数据
init_mysql_db_table
# 初始化tdengine的表
curl -u root:taosdata -d 'create database if not exists iThings;' 127.0.0.1:6041/rest/sql
