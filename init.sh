#!/bin/bash
CURDIR="$(pwd)"/"$(dirname "$0")"
#echo $CURDIR
echo "well come to go-things,we need init docker with docker-compose first"

function init_docker(){
 echo "init docker"
 mkdir -p /etc/docker
 cp -rf "$confPath"/docker/* /etc/docker/
 curl -sSL https://get.daocloud.io/docker | sh
 sudo systemctl start docker
 docker run hello-world
 echo "init docker end"
}

function init_curl() {
  echo "init curl"
  apt  install curl
  echo "init curl end"
}

function init_docker_compose(){
 echo "init docker_compose"
 curl -L https://get.daocloud.io/docker/compose/releases/download/1.12.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
 chmod 751 /usr/local/bin/docker-compose
 docker-compose version
 echo "init docker_compose end"
}

function init_module() {
    type curl >/dev/null 2>&1 || init_curl;
    type docker >/dev/null 2>&1 || init_docker;
    type docker-compose >/dev/null 2>&1 || init_docker_compose;
}

function init_conf_path(){
  #预创建配置所需文件夹
  thingsPath="/opt/things"
  confPath="$thingsPath/conf"
  if [ ! -d "$confPath" ]; then
    mkdir -p "$confPath"
  fi
  #将docker映射的所在工程内的配置拷贝到物理机目标位置
  cp -rf ./deploy/* $confPath
}

function init_mysql_db_table(){
 for (( i=0; i<300; i++)); do
   check_result=$(docker ps |grep mysql)
   if [ -n "$check_result" ];then
       docker exec -it mysql-docker /bin/bash -c 'mysql -uroot -ppassword < admin.sql'
       docker exec -it mysql-docker /bin/bash -c 'mysql -uroot -ppassword < dmsvr.sql'
       docker exec -it mysql-docker /bin/bash -c 'mysql -uroot -ppassword < usersvr.sql'
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
init_conf_path
chmod 751 ./deploy/shell/getip.sh
sudo ./deploy/shell/getip.sh
init_module
echo "now build and start i-Things needs mirror image"
echo "docker-compose -f $CURDIR/docker-compose.yml up -d" >> /etc/rc.local
sleep 1
echo "start docker compose "
docker-compose up -d
sleep 10 #这里必须等待足够长时间，等容器中mysql正常启动才能执行后续导入脚本命令
init_mysql_db_table
# 初始化tdengine的表
curl -u root:taosdata -d 'create database if not exists iThings;' 127.0.0.1:6041/rest/sql
