CC = gcc
INDENT = indent
INSTALL = /usr/bin/install
INSTALL_PROGRAM = ${INSTALL}
INSTALL_DATA = ${INSTALL} -m 644
PANDOC = pandoc

PROGRAM_NAME = git-syn

PREFIX = /usr/local
BIN_DIR = ${PREFIX}/bin
MAN_DIR = ${PREFIX}/share/man

DOC_DIR = ./doc

MAN_SRC = ${DOC_DIR}/man/${PROGRAM_NAME}.1.md

SRC_DIR = ./src

SRC = ${SRC_DIR}/${PROGRAM_NAME}.c
SRC += ${SRC_DIR}/config.c
SRC += ${SRC_DIR}/usage.c

CFLAGS += -fPIE -fno-stack-protector -Wall -Wextra -O2
LDFLAGS += -lgit2

all: ${PROGRAM_NAME} man

${PROGRAM_NAME}:
	${CC} ${CFLAGS} -I${SRC_DIR} -o $@ ${SRC} ${LDFLAGS}

check:
	@CFLAGS="${CFLAGS} -Wextra -Werror -Wno-sign-compare -fsyntax-only" make
	@echo "checking for errors... OK"

debug:
	@LDFLAGS="${LDFLAGS} -ggdb" make

reformat:
	@VERSION_CONTROL=none ${INDENT} ${SRC}

install: ${PROGRAM_NAME} man
	${INSTALL_PROGRAM} ${PROGRAM_NAME} ${BIN_DIR}
	@mkdir -p ${MAN_DIR}
	${INSTALL_DATA} ${PROGRAM_NAME}.1 ${MAN_DIR}/man1

man:
	${PANDOC} -s -t man ${MAN_SRC} -o ${PROGRAM_NAME}.1

clean: 
	@rm -f ${PROGRAM_NAME} ${PROGRAM_NAME}.1

