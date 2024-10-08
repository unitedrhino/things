# -*- coding:utf-8 -*-
.PHOmakeNY: update

start:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@docker compose up -d


update:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@docker compose pull
	@docker compose down
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
