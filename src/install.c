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

#include <git2.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

#include "git-syn.h"

int install_hooks(const char *git_hook_dir)
{
    int ret = EXIT_FAILURE;

    ret = copy_file(PRE_PUSH_HOOK, git_hook_dir);

    return ret;
}

int init_repo()
{
    int ret = EXIT_FAILURE;

    //ret = install_hooks();
    if (ret == EXIT_SUCCESS) {
        printf("Updated git hooks.\n");
        ret = EXIT_SUCCESS;
    }

    return ret;
}
