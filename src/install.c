/*
 * Copyright (C) 2019-2022 Lucas Ramage <ramage.lucas@protonmail.com>
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
#include <sds/sds.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

#include "syn.h"


#define PRE_PUSH_HOOK "pre-push.sh"
#define PRE_PUSH_HOOK_FILE HOOK_DIR "/" PRE_PUSH_HOOK

int init_repo(char *git_repository_path)
{
    int ret = EXIT_FAILURE;
    sds pre_push_hook = sdsempty();

    pre_push_hook = sdscat(pre_push_hook, git_repo_dir);
    pre_push_hook = sdscat(pre_push_hook, "/.git/hooks/");
    pre_push_hook = sdscat(pre_push_hook, PRE_PUSH_HOOK);

    printf("pre_push_hook: %s\n", pre_push_hook);
    printf("PRE_PUSH_HOOK_FILE: %s\n", PRE_PUSH_HOOK_FILE);

    ret = copy_file(PRE_PUSH_HOOK_FILE, pre_push_hook);

    if (ret == EXIT_SUCCESS) {
        printf("Updated git hooks.\n");
        ret = EXIT_SUCCESS;
    }

    return ret;
}
