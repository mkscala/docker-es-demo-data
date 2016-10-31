FROM gliderlabs/alpine:3.4

MAINTAINER blacktop, https://github.com/blacktop

ENV STACK 5.0.0

ENV GOLANG_VERSION 1.7.3
ENV GOLANG_SRC_URL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz
ENV GOLANG_SRC_SHA256 79430a0027a09b0b3ad57e214c4c1acfdd7af290961dd08d322818895af1ef44
# https://golang.org/issue/14851
COPY no-pic.patch /

COPY nginx_data /nginx_data
COPY main.go /go/src/main.go

ADD https://github.com/lalyos/docker-upx/releases/download/v3.91/upx /bin/upx

RUN apk-install -t .build-deps bash wget ca-certificates go git mercurial build-base musl-dev openssl \
  && set -x \
  && echo "Upgrade Golang to $GOLANG_VERSION..." \
  && export GOROOT_BOOTSTRAP="$(go env GOROOT)" \
  && wget -q "$GOLANG_SRC_URL" -O golang.tar.gz \
  && echo "$GOLANG_SRC_SHA256  golang.tar.gz" | sha256sum -c - \
  && tar -C /usr/local -xzf golang.tar.gz \
  && cd /usr/local/go/src \
  && patch -p2 -i /no-pic.patch \
  && ./make.bash \
  && rm -rf /*.patch \
  && echo "Build [import-dashboard]" \
  && export GOPATH=/go \
  && go version \
  && git clone -b v$STACK https://github.com/elastic/beats.git $GOPATH/src/github.com/elastic/beats \
  && cd $GOPATH/src/github.com/elastic/beats/libbeat/dashboards \
  && go build -ldflags "-X main.beat=" -o import_dashboards \
  && mv import_dashboards /bin/import_dashboards \
  # && echo "Building ingest Go binary..." \
  # && cd /go/src/ \
  # && export GOPATH=/go \
  # && go version \
  # && go get -d \
  # && go build -ldflags="-s -w" -o /bin/ingest \
  # && chmod +x /bin/upx \
  # && upx --brute /bin/ingest \
  && rm -rf /go /usr/local/go /usr/lib/go /tmp/* \
  && apk del --purge .build-deps

ENTRYPOINT ingest
