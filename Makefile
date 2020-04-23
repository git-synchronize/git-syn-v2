CC = gcc
INDENT = indent
INSTALL = /usr/bin/install
INSTALL_PROGRAM = ${INSTALL}
INSTALL_DATA = ${INSTALL} -m 644

SRC_DIR = ./src

PREFIX = /usr/local
BIN_DIR = ${PREFIX}/bin

GIT_SYN_SOURCES = ${SRC_DIR}/main.c
GIT_SYN_SOURCES += ${SRC_DIR}/usage.c

all: git-syn

git-syn:
	${CC} ${CFLAGS} ${LDFLAGS} -fPIE -Wall -o git-syn ${GIT_SYN_SOURCES}

check:
	@CFLAGS="${CFLAGS} -Wextra -Werror -Wno-sign-compare -fsyntax-only" make
	@echo "checking for errors... OK"

debug:
	@LDFLAGS="${LDFLAGS} -ggdb" make

reformat:
	@VERSION_CONTROL=none $(INDENT) ${SRC_DIR}/*.c

install: git-syn
	${INSTALL_PROGRAM} git-syn ${BIN_DIR}

clean: 
	@rm -f git-syn

