#!/bin/bash
ip add | awk -F ':' '$1~"^[1-9]" {print $2}' >eth.list
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
 ip add show dev "$e" |grep -w 'inet'|awk '{print $2}'|awk -F '/' '{print $1}' >$e.txt
 n=$(wc -l "$e".txt|awk '{print $1}')
 if [ $n -eq 0 ]
  then
  echo "无IP地址"
  else
  echo "IP地址是：`cat $e.txt`
  ip=$(cat $e.txt)
  #echo "ip: $ip"
  #replace all of the ip 127.0.0.1 in the file emqx_auth_http.conf to the real ip of service
  sed -i "s#127.0.0.1#$ip#g" /opt/things/conf/emqx/etc/plugins/emqx_auth_http.conf
  fi
}
ipad $e
