# -*- coding:utf-8 -*-
.PHONY: update

start:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@if [ ! -d ../conf/core/dist/app ]; then \
		echo "初始化 dist 文件（从 core 镜像提取）..."; \
		CID=$$(docker create docker.unitedrhino.com/ee/core:latest); \
		docker cp $$CID:/unitedrhino/dist/. ../conf/core/dist/; \
		docker rm $$CID > /dev/null; \
	fi
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
	docker compose down core
	rm -rf ../conf/core/dist/app
	CID=$$(docker create docker.unitedrhino.com/ee/core:latest); \
	docker cp $$CID:/unitedrhino/dist/. ../conf/core/dist/; \
	docker rm $$CID > /dev/null
	docker compose up -d core
