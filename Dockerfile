FROM registry.cn-qingdao.aliyuncs.com/unitedrhino/golang:1.21.13-alpine3.20 as go-builder
WORKDIR /unitedrhino/
COPY ./ ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN cd ./service/apisvr && go mod tidy && go build -tags no_k8s -ldflags="-s -w" .


FROM registry.cn-qingdao.aliyuncs.com/unitedrhino/alpine:3.20
LABEL homepage="https://gitee.com/unitedrhino"
ENV TZ Asia/Shanghai
RUN apk add tzdata

WORKDIR /unitedrhino/
COPY --from=go-builder /unitedrhino/service/apisvr/apisvr ./apisvr
#COPY --from=go-builder /unitedrhino/deploy/conf/things/etc/ ./etc

ENTRYPOINT ["./apisvr"]
