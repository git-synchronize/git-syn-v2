FROM alpine:latest

RUN apk update && apk add \
      check gcc make

WORKDIR /usr/src/git-syn
COPY . /usr/src/git-syn
