PROGRAM_NAME = git-syn

PREFIX = /usr/local
BIN_DIR = ${PREFIX}/bin
MAN_DIR = ${PREFIX}/share/man
SHARE_DIR = ${PREFIX}/share/${PROGRAM_NAME}

MAN_SRC = ${DOC_DIR}/man/${PROGRAM_NAME}.1.md

all: ${PROGRAM_NAME}

${PROGRAM_NAME}:
	go build -buildvcs=false .

check:
	@CFLAGS="${CFLAGS} -Wextra -Werror -Wno-sign-compare -fsyntax-only" make
	@echo "checking for errors... OK"

man:
	pandoc -s -t man ${MAN_SRC} -o ./${PROGRAM_NAME}.1

clean: 
	@rm -f ${PROGRAM_NAME} ${PROGRAM_NAME}.1

