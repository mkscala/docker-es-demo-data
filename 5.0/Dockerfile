FROM gliderlabs/alpine:3.4

MAINTAINER blacktop, https://github.com/blacktop

ENV STACK 5.0.0

# ENV GOLANG_VERSION 1.7.3
# ENV GOLANG_SRC_URL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz
# ENV GOLANG_SRC_SHA256 79430a0027a09b0b3ad57e214c4c1acfdd7af290961dd08d322818895af1ef44
# https://golang.org/issue/14851
# COPY no-pic.patch /

COPY nginx_data /nginx_data
COPY main.go /go/src/main.go

# ADD https://github.com/lalyos/docker-upx/releases/download/v3.91/upx /bin/upx

RUN apk-install -t .build-deps \
                    ca-certificates \
                    build-base \
                    mercurial \
                    musl-dev \
                    openssl \
                    wget \
                    bash \
                    git \
                    go \
  && set -x \
  && echo "Build [import-dashboard]" \
  && export GOPATH=/go \
  && go version \
  && git clone -b v$STACK https://github.com/elastic/beats.git $GOPATH/src/github.com/elastic/beats \
  && cd $GOPATH/src/github.com/elastic/beats/libbeat/dashboards \
  && go build -ldflags "-s -w -X main.beat=" -o import_dashboards \
  && mv import_dashboards /bin/import_dashboards \
  # && echo "Building ingest Go binary..." \
  # && cd /go/src/ \
  # && export GOPATH=/go \
  # && go version \
  # && go get -d \
  # && go build -ldflags="-s -w" -o /bin/ingest \
  # && chmod +x /bin/upx \
  # && upx --brute /bin/import_dashboards \
  && rm -rf /go /usr/local/go /usr/lib/go /tmp/* \
  && apk del --purge .build-deps

ENTRYPOINT import_dashboards
