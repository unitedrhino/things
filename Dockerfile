FROM golang:1.19-alpine3.16 as go-builder
WORKDIR /ithings/
COPY ./go.mod ./go.mod
RUN go mod download
COPY ./ ./
RUN cd ./src/apisvr && go mod tidy && go build .

FROM node:19 as web-builder
WORKDIR /ithings/
COPY ./assets/package.json ./assets/package.json
RUN cd assets && yarn install --no-lockfile
COPY ./assets ./assets
RUN cd assets && yarn build

FROM alpine:3.16
LABEL homepage="https://github.com/i-Things/iThings"
ENV TZ Asia/Shanghai
RUN apk add tzdata

WORKDIR /ithings/
COPY --from=go-builder /ithings/src/apisvr/apisvr ./apisvr
COPY --from=go-builder /ithings/src/apisvr/etc/ ./etc
COPY --from=web-builder /ithings/assets/dist/ ./dist/front/iThingsCore

ENTRYPOINT ["./apisvr"]
