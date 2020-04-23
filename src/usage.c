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

#include "git-syn.h"

void print_usage()
{
    printf("Usage:\n\tgit-syn [option] ... [command]\n");
}

void print_version()
{
    printf("git-syn version %s\n", VERSION);
}
