FROM alpine:latest

RUN apk update && apk add \
      check \
      gcc \
      indent \
      libgit2-dev \
      make \
      musl-dev

WORKDIR /usr/src/git-syn
COPY . /usr/src/git-syn
