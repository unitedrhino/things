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
4. 打开 `http://localhost:17000` 默认账号:root 密码: root.2020
注:如果是win或本地运行,则需要修改下面两个文件中的地址为物理机地址才可以访问,如果新增其他服务,同样在这里新增个文件即可
* deploy/docker/conf/nightingale/etc-categraf/input.prometheus/core.toml
* deploy/docker/conf/nightingale/etc-categraf/input.prometheus/things.toml