#!/bin/sh
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

function init_mysql_db_table(){
usersvr=$(cat $CURDIR/db/usersvr.sql)
sudo docker exec -it mysql-docker bash -c "mysql -uroot -ppassword -e  '$usersvr'"

}

type docker >/dev/null 2>&1 || init_docker;
type docker-compose >/dev/null 2>&1 || init_docker_compose;
echo "docker with docker-compose init success"
echo "now buid and start go-things needs mirror image"
echo "docker-compose -f $CURDIR/docker-compose.yml up -d" >> /etc/rc.d/rc.local
docker-compose up -d
#init_mysql_db_table