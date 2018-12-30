FROM golang:alpine as builder

WORKDIR /go/src/github.com/bblenard/C2

COPY ./ /go/src/github.com/bblenard/C2/

ENV GOCACHE=off

RUN    set -x \
    && apk update \
    && apk add \
               gcc \
               git \
               libc-dev \
    && apk upgrade \
    && rm -f /var/cache/apk/* \
    && go get . \
       || : \
    && go build -ldflags="-w -s" -o /tmp/server server.go \
    && go build -ldflags="-w -s" -o /tmp/cli cli.go


FROM alpine:3.8

EXPOSE 8888/tcp

RUN    set -x \
    && apk update \
    && apk upgrade \
    && rm -f /var/cache/apk/*

WORKDIR /app
VOLUME /app

COPY --from=builder /tmp/server /usr/local/bin/server
COPY --from=builder /tmp/cli /usr/local/bin/cli

CMD ["/usr/local/bin/server"]

