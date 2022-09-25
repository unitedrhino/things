FROM golang:1.19-alpine3.16 as go-builder
ARG SVR_NAME
ARG GOPROXY=goproxy.cn
ENV GOPROXY=https://${GOPROXY},direct

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache make

WORKDIR /ithings/

COPY src/ src/
COPY shared/ shared/
COPY go.mod go.sum Makefile ./

RUN go mod download -x
RUN make build.${SVR_NAME}

FROM alpine:3.16
ARG SVR_NAME
LABEL name="${SVR_NAME}svr"
LABEL author="wwhai"
LABEL email="cnwwhai@gmail.com"
LABEL homepage="https://github.com/i4de/ithings"

WORKDIR /ithings/
COPY --from=go-builder /ithings/cmd/${SVR_NAME}svr ./bin/svr
COPY --from=go-builder /ithings/src/${SVR_NAME}svr/etc ./etc

EXPOSE 2580
ENTRYPOINT ["./bin/svr"]
