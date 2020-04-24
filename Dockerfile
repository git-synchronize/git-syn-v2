FROM alpine:latest

RUN apk update && \
    apk add --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing \
      check \
      gcc \
      indent \
      libgit2-dev \
      make \
      musl-dev \
      pandoc

WORKDIR /usr/src/git-syn
COPY . /usr/src/git-syn

RUN make && \
    make install
