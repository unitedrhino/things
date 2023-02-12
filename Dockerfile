FROM golang:1.19-alpine3.16 as go-builder
ARG GOPROXY=goproxy.cn
ENV GOPROXY=https://${GOPROXY},direct
WORKDIR /ithings/
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download
COPY ./ ./
RUN cd ./src/apisvr && go build .

FROM node:19 as web-builder
WORKDIR /ithings/
COPY ./assets/package.json ./assets/package.json
COPY ./assets/yarn.lock ./assets/yarn.lock
RUN cd assets && yarn install
COPY ./assets ./assets
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
