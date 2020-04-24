/*
 * Copyright (C) 2019-2020 Lucas Ramage <ramage.lucas@protonmail.com>
 *
 * SPDX-License-Identifier: GPL-2.0-or-later
 *
 * This check template is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 */

#include <stdio.h>
#include <stdlib.h>

#include "git-syn.h"

void print_usage(int status)
{
    printf("Usage: %s [option] ... [command] ...\n"
           "\n%s\n"
           "\n  -h, --help       display this help and exit"
           "\n  -v, --version    output version information and exit"
           "\n  install          install extension to repository"
           "\n  uninstall        remove extension from repository"
           "\n  monitor          view or enable/disable repository sync"
           "\n\nGit SYN online help: <https://gitlab.com/oxr463/git-syn>"
           "\nFull documentation <https://oxr463.gitlab.io/git-syn>"
           "\nor available locally via: man git-syn\n",
           PROGRAM_NAME, PROGRAM_DESCRIPTION);

    exit(status);
}

void print_version(int status)
{
    printf("%s version %s\n", PROGRAM_NAME, VERSION);
    exit(status);
}
