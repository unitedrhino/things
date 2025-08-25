# -*- coding:utf-8 -*-
.PHONY: update

start:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@docker compose up -d


update:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@docker compose pull core things
	@docker compose down core
	@docker compose down things
	@docker compose up -d

reset:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@docker stop mariadb
	@docker rm mariadb
	@docker stop tdengine
	@docker rm tdengine
	@docker compose pull
	@docker compose down
	@docker compose up -d

update-core:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	docker compose pull core
	rm -rf  ../conf/core/dist2
	mv ../conf/core/dist ../conf/core/dist2
	mkdir ../conf/core/dist
	docker compose down core
	docker compose up -d core
	rm -rf ../conf/core/dist2/app/core
	rm -rf ../conf/core/dist2/app/lowcode
	rm -rf ../conf/core/dist2/app/my-things
	rm -rf ../conf/core/dist2/app/things
	rm -rf ../conf/core/dist2/app/static
	cp -r ../conf/core/dist2/app/* ../conf/core/dist/app/
