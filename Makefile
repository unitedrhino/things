# -*- coding:utf-8 -*-
.PHOmakeNY: build

build:build.clean mod cp.etc build.api build.dg build.dm build.sys build.rule build.timedjob build.timedscheduler

runall:  run.timedjob run.timedscheduler run.sys run.dm run.dg run.rule run.api

buildone:build.clean mod cp.etc build.api


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
	@mkdir -p ./cmd/etc/
	@cp -rf ./src/apisvr/etc/* ./cmd/etc/


build.api:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build -o ./cmd/apisvr ./src/apisvr

build.dg:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/dgsvr ./src/dgsvr

build.dm:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/dmsvr ./src/dmsvr

build.sys:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/syssvr ./src/syssvr

build.rule:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/rulesvr ./src/rulesvr

build.timedjob:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/timedjobsvr ./src/timed/timedjobsvr

build.timedscheduler:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/timedschedulersvr ./src/timed/timedschedulersvr


run.api:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd && nohup ./apisvr &  cd ..

run.dg:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd && nohup ./dgsvr &  cd ..

run.dm:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd && nohup ./dmsvr &  cd ..

run.sys:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd && nohup ./syssvr &  cd ..

run.rule:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd && nohup ./rulesvr &  cd ..

run.timedjob:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cd cmd && nohup ./timedjobsvr &  cd ..

run.timedscheduler:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cd cmd && nohup ./timedschedulersvr &  cd ..
