#!/bin/bash
set -x  # 开启命令显示模式

# 获取当前Git仓库的本地分支名称
current_branch=$(git rev-parse --abbrev-ref HEAD)

# 输出当前分支信息（可选，用于确认）
echo "当前仓库所在分支: $current_branch"

# 执行go get命令，拉取指定分支(dev-1.6)的包
echo "开始执行: go get gitee.com/unitedrhino/things@$current_branch"
#go get gitee.com/unitedrhino/things@${current_branch}
go get gitee.com/unitedrhino/core@saas3
# 检查命令执行结果
if [ $? -eq 0 ]; then
    echo "go get 执行成功"
else
    echo "go get 执行失败，请检查网络或包路径是否正确"
    exit 1
fi