#!/bin/bash

# 检查是否传入日志文件路径
if [ $# -ne 1 ]; then
    echo "用法: $0 <日志文件路径>"
    exit 1
fi

# 从命令行参数获取日志文件路径
LOG_FILE="$1"

# 检查文件是否存在
if [ ! -f "$LOG_FILE" ]; then
    echo "日志文件 $LOG_FILE 不存在。"
    exit 1
fi

# 获取文件大小（单位：字节）
FILE_SIZE=$(stat -c %s "$LOG_FILE")

# 2GB 对应的字节数
MAX_SIZE=$((2 * 1024 * 1024 * 1024))

# 200MB 对应的字节数
KEEP_SIZE=$((200 * 1024 * 1024))

# 若文件大小超过 2GB
if [ $FILE_SIZE -gt $MAX_SIZE ]; then
    echo "日志文件 $LOG_FILE 大小超过 2GB，将保留最后 200MB 的数据。"
    # 创建临时文件
    TEMP_FILE="${LOG_FILE}.tmp"
    # 提取最后 200MB 的数据到临时文件
    tail -c $KEEP_SIZE "$LOG_FILE" > "$TEMP_FILE"
    # 移动临时文件覆盖原文件
    mv "$TEMP_FILE" "$LOG_FILE"
    echo "已保留日志文件 $LOG_FILE 最后 200MB 的数据。"
else
    echo "日志文件 $LOG_FILE 大小未超过 2GB，无需处理。"
fi
