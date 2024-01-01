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
SHARE_DIR = ${PREFIX}/share/${PROGRAM_NAME}

TARGET_DIR = .
BUILD_DIR = ${TARGET_DIR}
SRC_DIR = ${TARGET_DIR}/src
DOC_DIR = ${TARGET_DIR}/doc
TEST_DIR = ${TARGET_DIR}/test

MAN_SRC = ${DOC_DIR}/man/${PROGRAM_NAME}.1.md

SRC = ${SRC_DIR}/${PROGRAM_NAME}.c
HOOK_SRC = ${SRC_DIR}/hook/pre-push.sh

SRC = ${SRC_DIR}/main.c
SRC += ${SRC_DIR}/install.c
SRC += ${SRC_DIR}/usage.c
SRC += ${SRC_DIR}/util.c

CFLAGS += -fPIE -fno-stack-protector -Wall -Wextra -O2
LDFLAGS += -lgit2 -lsds

all: ${PROGRAM_NAME} test man hook

${PROGRAM_NAME}:
	${CC} ${CFLAGS} -I${SRC_DIR} -o $@ ${SRC} ${LDFLAGS}

check:
	@CFLAGS="${CFLAGS} -Wextra -Werror -Wno-sign-compare -fsyntax-only" make
	@echo "checking for errors... OK"

test: check_main check_install check_usage check_util

check_main:
	${CC} ${TEST_DIR}/$@.c -lcheck -o ${TEST_DIR}/$@
	${TEST_DIR}/$@

check_install:
	${CC} ${TEST_DIR}/$@.c -lcheck -o ${TEST_DIR}/$@
	${TEST_DIR}/$@

check_usage:
	${CC} ${TEST_DIR}/$@.c -lcheck -o ${TEST_DIR}/$@
	${TEST_DIR}/$@

check_util:
	${CC} ${TEST_DIR}/$@.c -lcheck -o ${TEST_DIR}/$@
	${TEST_DIR}/$@

debug:
	@LDFLAGS="${LDFLAGS} -ggdb" make ${PROGRAM_NAME}

reformat:
	@VERSION_CONTROL=none ${INDENT} ${SRC}

install: all
	${INSTALL_PROGRAM} ${BUILD_DIR}/${PROGRAM_NAME} ${BIN_DIR}
	${INSTALL_DATA} ${BUILD_DIR}/${PROGRAM_NAME}.1 ${MAN_DIR}/man1

man:
	${PANDOC} -s -t man ${MAN_SRC} -o ${BUILD_DIR}/${PROGRAM_NAME}.1

clean: 
	@rm -f ${PROGRAM_NAME} ${PROGRAM_NAME}.1

