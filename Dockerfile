FROM golang:alpine as go-builder
WORKDIR /ithings/
COPY ./ ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN cd ./service/apisvr && go mod tidy && go build -tags no_k8s -ldflags="-s -w" .


FROM alpine:latest
LABEL homepage="https://github.com/i-Things/iThings"
ENV TZ Asia/Shanghai
RUN apk add tzdata

WORKDIR /ithings/
COPY --from=go-builder /ithings/service/apisvr/apisvr ./apisvr
#COPY --from=go-builder /ithings/deploy/conf/things/etc/ ./etc

ENTRYPOINT ["./apisvr"]
