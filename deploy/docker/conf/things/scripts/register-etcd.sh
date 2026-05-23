#!/bin/sh
# 在 etcd 中注册 dm.rpc 和 dg.rpc 端点，维持租约保活

ETCD_ENDPOINT="${ETCD_ENDPOINT:-etcd:2379}"
LEASE_TTL=30
N=0

# 获取本容器网段内 things 的 IP
THINGS_IP=$(getent hosts things | awk '{print $1}')

while true; do
  N=$((N + 1))
  LEASE=$(etcdctl --endpoints="$ETCD_ENDPOINT" lease grant $LEASE_TTL | awk '{print $2}')
  if [ -n "$LEASE" ]; then
    etcdctl --endpoints="$ETCD_ENDPOINT" put "dm.rpc/${$}-${N}" "${THINGS_IP}:9081" --lease="$LEASE"
    etcdctl --endpoints="$ETCD_ENDPOINT" put "dg.rpc/${$}-${N}" "${THINGS_IP}:6166" --lease="$LEASE"
    etcdctl --endpoints="$ETCD_ENDPOINT" lease keep-alive "$LEASE"
  fi
  sleep 5
done
