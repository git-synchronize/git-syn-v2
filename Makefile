CC = gcc
INDENT = indent
INSTALL = /usr/bin/install
INSTALL_PROGRAM = ${INSTALL}
INSTALL_DATA = ${INSTALL} -m 644

PREFIX = /usr/local
BIN_DIR = ${PREFIX}/bin

SRC_DIR = ./src

SRC = ${SRC_DIR}/git-syn.c
SRC += ${SRC_DIR}/usage.c
SRC += ${SRC_DIR}/config.c

CFLAGS += -fPIE -fno-stack-protector -Wall -Wextra -O2
LDFLAGS += -lgit2

all: git-syn

git-syn:
	${CC} ${CFLAGS} -I${SRC_DIR} -o $@ ${SRC} ${LDFLAGS}

check:
	@CFLAGS="${CFLAGS} -Wextra -Werror -Wno-sign-compare -fsyntax-only" make
	@echo "checking for errors... OK"

debug:
	@LDFLAGS="${LDFLAGS} -ggdb" make

reformat:
	@VERSION_CONTROL=none ${INDENT} ${SRC}

install: git-syn
	${INSTALL_PROGRAM} git-syn ${BIN_DIR}

clean: 
	@rm -f git-syn
