#!/bin/bash -x
function runsvr(){
	echo "run "$1
	killall $1
  cd src/$1/
  go build
  nohup ./$1 &
  cd ../..
}
runsvr dmsvr
runsvr usersvr
runsvr dcsvr
runsvr webapi
