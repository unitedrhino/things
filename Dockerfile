FROM golang:1.21.13-alpine3.20 as go-builder
WORKDIR /ithings/
COPY ./ ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN cd ./service/apisvr && go mod tidy && go build -ldflags="-s -w" .


FROM alpine:3.20
LABEL homepage="https://github.com/i-Things/iThings"
ENV TZ Asia/Shanghai
RUN apk add tzdata

WORKDIR /ithings/
COPY --from=go-builder /ithings/service/apisvr/apisvr ./apisvr
COPY --from=go-builder /ithings/deploy/conf/ithings/etc/ ./etc

ENTRYPOINT ["./apisvr"]
