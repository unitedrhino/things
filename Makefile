# -*- coding:utf-8 -*-
.PHOmakeNY: build

build:build.clean mod cp.etc build.api build.dg build.view build.dm build.sys build.ud build.timedjob build.timedscheduler

runall:  run.timedjob run.timedscheduler run.sys run.dm run.dg run.ud run.api run.view

buildone:build.clean mod cp.etc build.api build.view moduleupdate build.core

moduleupdate:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@git submodule update --init --recursive
	@git submodule foreach git checkout master
	@git submodule foreach git pull

build.core:
	@mkdir -p ./cmd/core
	@cd module/core && make packone
	@cp -rf module/core/cmd/* ./cmd/core

toremote:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>tormote cmd<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@rsync -r -v ./cmd/* root@120.79.205.165:/root/git/iThings

killall:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>killing all<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@killall  apisvr &
	@killall  dgsvr &
	@killall  syssvr &
	@killall  dmsvr &
	@killall  timedjob &
	@killall  timedscheduler &

build.clean:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>clean cmd<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@rm -rf ./cmd/*

mod:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>downloading $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go mod download
	@go mod tidy


cp.etc:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>copying etc<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@mkdir -p ./cmd/things/etc/
	@mkdir -p ./cmd/things/dist/
	@cp -rf ./service/apisvr/etc/* ./cmd/things/etc/
	@cp -rf ./service/viewsvr/etc/* ./cmd/things/etc/


build.api:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build -o ./cmd/things/apisvr ./service/apisvr

build.view:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build -o ./cmd/things/viewsvr ./service/viewsvr

build.dg:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/things/dgsvr ./service/dgsvr

build.dm:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/things/dmsvr ./service/dmsvr

build.sys:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/things/syssvr ./service/syssvr

build.ud:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/things/udsvr ./service/udsvr

build.timedjob:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/things/timedjobsvr ./service/timed/timedjobsvr

build.timedscheduler:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/things/timedschedulersvr ./service/timed/timedschedulersvr


run.api:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd/things && nohup ./apisvr &  cd ..

run.view:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd/things && nohup ./viewsvr &  cd ..


run.dg:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd/things && nohup ./dgsvr &  cd ..

run.dm:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd/things && nohup ./dmsvr &  cd ..

run.sys:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd/things && nohup ./syssvr &  cd ..

run.ud:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd/things && nohup ./udsvr &  cd ..

run.timedjob:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cd cmd/things && nohup ./timedjobsvr &  cd ..

run.timedscheduler:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cd cmd/things && nohup ./timedschedulersvr &  cd ..
