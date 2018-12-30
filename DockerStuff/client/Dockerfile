FROM golang:1.10 as builder

WORKDIR /go/src/github.com/bblenard/C2

COPY ./ /go/src/github.com/bblenard/C2/

ENV GOCACHE=off

RUN    set -x \
    && apt-get update \
    && apt-get upgrade -y \
    && rm -rf /var/cache/apt/* \
    && go get . \
       || : \
    && go build -ldflags="-w -s" -o /tmp/client client.go


FROM centos:7

COPY --from=builder /tmp/client /usr/local/bin/virus
COPY ./DockerStuff/client/files/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN    set -x \
    && yum update -y \
    && rm -fr /var/cache/yum/* \
    && chmod 0755 /usr/local/bin/entrypoint.sh

WORKDIR /root

CMD ["/usr/local/bin/entrypoint.sh"]

