## 协议组件 

docker-compose-protocol.yml是给开源使用的联犀第三方协议组件

注意需要先启动run-xx 下面的服务这里才能启动

启动命令: ` docker compose -f docker-compose-protocol.yml up -d`

## 监控组件

docker-compose-monitor.yml是给开源使用的监控系统

注意需要先启动run-xx 下面的服务这里才能启动
1. 首先打开 docker-compose.yml文件中core和things volumes 中的文件映射注释
2. 打开 `conf/core/etc/apiDocker.yaml` 和 `conf/things/etc/apiDocker.yaml` 中的 Telemetry 和 Prometheus 配置
3. 启动 ` docker compose -f docker-compose-monitor.yml up -d`