FROM golang:1.19-alpine3.16 as go-builder
ARG GOPROXY=goproxy.cn
ENV GOPROXY=https://${GOPROXY},direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache make
WORKDIR /ithings/
COPY ./ ./
RUN cd ./src/apisvr && go build .

FROM node:19-alpine3.16 as web-builder
WORKDIR /ithings/
COPY ./assets ./assets
RUN cd assets && yarn
RUN cd assets && yarn build

FROM alpine:3.16
LABEL homepage="https://github.com/i4de/ithings"
ENV TZ Asia/Shanghai
RUN apk add tzdata

WORKDIR /ithings/
COPY --from=go-builder /ithings/src/apisvr/apisvr ./apisvr
COPY --from=go-builder /ithings/src/apisvr/etc ./etc
COPY --from=web-builder /ithings/assets/dist/ ./dist/front/iThingsCore

ENTRYPOINT ["./apisvr"]
