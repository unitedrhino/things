# -*- coding:utf-8 -*-
.PHONY: build

build:build.clean mod cp.etc build.api build.dg  build.dm  build.ud

runall:   run.dm run.dg run.ud run.api

buildone:build.clean mod cp.etc build.api  moduleupdate build.core build.front

#仅编译后端
buildback: build.clean mod cp.etc build.api

packback: buildback toremote



toremote:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>tormote cmd<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@rsync -r -v ./cmd/* root@47.94.112.109:/root/run


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
	@go mod tidy


cp.etc:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>copying etc<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@mkdir -p ./cmd/things/etc/
	@mkdir -p ./cmd/things/dist/
	@cp -rf ./service/apisvr/etc/* ./cmd/things/etc/


build.api:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build -ldflags="-s -w"  -tags no_k8s -o ./cmd/things/thingsvr ./service/apisvr


build.dg:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build -ldflags="-s -w"  -tags no_k8s  -o ./cmd/things/dgsvr ./service/dgsvr

build.dm:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build -ldflags="-s -w"  -tags no_k8s  -o ./cmd/things/dmsvr ./service/dmsvr


build.ud:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build -ldflags="-s -w"  -tags no_k8s  -o ./cmd/things/udsvr ./service/udsvr


run.api:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd/things && nohup ./apisvr &  cd ..


run.dg:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd/things && nohup ./dgsvr &  cd ..

run.dm:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd/things && nohup ./dmsvr &  cd ..


run.ud:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd/things && nohup ./udsvr &  cd ..

