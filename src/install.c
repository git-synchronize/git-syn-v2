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

void install_hooks()
{
    printf("Updated git hooks.\n");
    exit(EXIT_SUCCESS);
}

void init_repo()
{
    printf("Git SYN initialized.\n");
    exit(EXIT_SUCCESS);
}
