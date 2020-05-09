FROM alpine:latest

RUN apk update && \
    apk add --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing \
      check \
      gcc \
      git \
      indent \
      libgit2-dev \
      make \
      musl-dev \
      pandoc

WORKDIR /usr/src/sds
ENV LIBSDS_SONAME libsds.so.2.0.0
RUN git clone https://github.com/antirez/sds.git . && \
    echo "typedef int sdsvoid;" >> sdsalloc.h && \
    gcc -fPIC -fstack-protector -std=c99 -pedantic -Wall \
        -Werror -shared -o "${LIBSDS_SONAME}" \
        -Wl,-soname="${LIBSDS_SONAME}" \
        sds.c sds.h sdsalloc.h && \
    mkdir -p /usr/local/lib && \
    cp -a "${LIBSDS_SONAME}" /usr/local/lib && \
    ln -s "/usr/local/lib/${LIBSDS_SONAME}" \
        /usr/local/lib/libsds.so && \
    mkdir -p /usr/local/include/sds && \
    cp -a sds.h /usr/local/include/sds && \
    ldconfig /etc/ld.so.conf.d

WORKDIR /usr/src/git-syn
COPY . /usr/src/git-syn

RUN make && \
    make install

