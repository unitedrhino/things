#!/bin/bash
set -x
# 检查是否至少有一个参数
if [ $# -eq 0 ]; then
    echo "请输入tag标签。"
    exit 1
fi

# 获取第一个参数
tag="$1"
git tag $tag
git push origin $tag
#git push github $tag
git push gitee $tag
git checkout master
git pull origin master
git push  gitee master
#git push  github master
git checkout dev
