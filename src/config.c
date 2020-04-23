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

#include "git-syn.h"

void parse_config()
{
    git_buf *path = NULL;
    int error;

    git_libgit2_init();

    error = git_config_find_global(path);
    error = git_config_find_xdg(path);
    error = git_config_find_system(path);

    git_config *cfg = NULL;

    error = git_config_open_default(&cfg);

    (void) error;
}
