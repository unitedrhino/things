#!/bin/bash
CURDIR="`pwd`"/"`dirname $0`"
#echo $CURDIR
echo "well come to go-things,we need init docker with docker-compose first"

function init_docker(){
 echo "init docker"
 curl -sSL https://get.daocloud.io/docker | sh
 sudo systemctl start docker
 docker run hello-world
}

function init_docker_compose(){
 echo "init docker_compose"
 curl -L https://get.daocloud.io/docker/compose/releases/download/1.12.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
 chmod +x /usr/local/bin/docker-compose
 docker-compose version
}

function init_conf_path(){
  #预创建配置所需文件夹
  thingsPath="/opt/things"
  confPath="/opt/things/conf"
  emqxPath="/opt/things/conf/emqx"
  mysqlPath="/opt/things/conf/mysql"

  if [ ! -d "$thingsPAth" ]; then
    mkdir "$thingsPath"
  fi
  sleep 1
  if [ ! -d "$confPath" ]; then
    mkdir "$confPath"
  fi
  sleep 1
  if [ ! -d "$emqxPath" ]; then
    mkdir "$emqxPath"
  fi
  sleep 1
  if [ ! -d "$mysqlPath" ]; then
    mkdir "$mysqlPath"
  fi
  sleep 1

  #将emqx和mysql所在工程内的配置拷贝到物理机目标位置
  cp conf/emqx/* /opt/things/conf/emqx/
  cp conf/mysql/* /opt/things/conf/mysql/
}

function init_mysql_db_table(){
 for (( i=0; i<300; i++)); do
   check_result=`docker ps |grep mysql`
   if [ -n $("$check_result") ];then
       docker exec -it mysql-docker /bin/bash -c 'mysql -uroot -ppassword < admin.sql'
       docker exec -it mysql-docker /bin/bash -c 'mysql -uroot -ppassword < dcsvr.sql'
       docker exec -it mysql-docker /bin/bash -c 'mysql -uroot -ppassword < dmsvr.sql'
       docker exec -it mysql-docker /bin/bash -c 'mysql -uroot -ppassword < usersvr.sql'
       sleep 5
       echo "has install mysql"
       break
   else
       echo "not install mysql, please make sure docker mysql is running"
       sleep 2
   fi
  done
}

type docker >/dev/null 2>&1 || init_docker;
type docker-compose >/dev/null 2>&1 || init_docker_compose;
echo "docker with docker-compose init success"
echo "now buid and start go-things needs mirror image"
echo "docker-compose -f $CURDIR/docker-compose.yml up -d" >> /etc/rc.d/rc.local

init_conf_path
sleep 1
echo "start docker compose "
docker-compose up -d
sleep 10 #这里必须等待足够长时间，等容器中mysql正常启动才能执行后续导入脚本命令
init_mysql_db_table