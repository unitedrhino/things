FROM docker.unitedrhino.com/unitedrhino/golang:1.23.4-alpine3.21 as go-builder
WORKDIR /unitedrhino/
COPY ./ ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
ENV GOPRIVATE=*.gitee.com,gitee.com/*
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add  git
RUN go mod download
RUN cd ./service/apisvr  && go build -tags no_k8s -ldflags="-s -w" -o thingsvr .

FROM docker.unitedrhino.com/unitedrhino/alpine:3.20
LABEL homepage="https://gitee.com/unitedrhino"
ENV TZ Asia/Shanghai
RUN apk add tzdata

WORKDIR /unitedrhino/
COPY --from=go-builder /unitedrhino/service/apisvr/thingsvr ./thingsvr
COPY --from=go-builder /unitedrhino/service/apisvr/etc ./etc

ENTRYPOINT ["./thingsvr"]
