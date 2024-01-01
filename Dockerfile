FROM alpine:latest as dev

RUN apk update && \
    apk add --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing \
      check-dev \
      gcc \
      git \
      go \
      indent \
      less \
      libgit2-dev \
      libsds-dev \
      make \
      musl-dev \
      pandoc \
      vim

WORKDIR /usr/src/git-syn
COPY . /usr/src/git-syn

RUN git config --global --add safe.directory /usr/src/git-syn

#RUN mkdir -p /usr/local/share/man/man1 && \
#    make && \
#    make install

