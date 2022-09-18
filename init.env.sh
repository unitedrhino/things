#!/bin/bash
#DATA_DIR=eval echo '$'$env_name

echo "well come to go-things,we need init docker with docker-compose first"

function init_docker(){
 echo "init docker"
 mkdir -p /etc/docker
 cp -rf ./deploy/docker/* /etc/docker/
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

function get_ip() {
    ip add | awk -F ':' '$1~"^[1-9]" {print $2}'  |grep -v "lo" >eth.list
    while :
    do
      eths=`cat eth.list |xargs`
      read -p "请输入正确的网卡名如('$eths'):" e
      if [ -z "$e" ]
      then
          echo "网卡名不能为空"
      continue
      fi

      if [ $e = "lo" ]; then
    	read -p "lo ip 可能导致emq不可用，是否重新输入网卡名[y/n]:" b
            if [ $b = "y" -o $b = "Y" ]; then
    	    continue
            fi
      fi

      if grep -qw "$e" eth.list
      then
      break
      else
          echo "请重新输入正确的网卡名"
      continue
      fi
      done
    ip=''
    ipad() {
     ip add show dev $e |grep -w 'inet'|awk '{print $2}'|awk -F '/' '{print $1}' >$e.txt
     n=$(wc -l $e.txt|awk '{print $1}')
     if [ $n -eq 0 ]
      then
      echo "无IP地址"
      else
      echo "IP地址是：`cat $e.txt`"
      ip=$(cat $e.txt)
      #echo "ip: $ip"
      #replace all of the ip 127.0.0.1 in the file emqx_auth_http.conf to the real ip of service
      sed -i "s#127.0.0.1#$ip#g" ./deploy/conf/emqx/etc/plugins/emqx_auth_http.conf
      fi
    }
    ipad "$e"
}

get_ip
init_module
echo "now build and start i-Things needs mirror image"
echo
echo "docker-compose -f docker-compose.env.yml up -d" >> /etc/rc.local
sleep 1
echo "start docker compose "
docker-compose -f docker-compose.env.yml up -d

sleep 10 #这里必须等待足够长时间，等容器中mysql正常启动才能执行后续导入脚本命令
init_mysql_db_table

# 初始化tdengine的表
curl -u root:taosdata -d 'create database if not exists iThings;' 127.0.0.1:6041/rest/sql
