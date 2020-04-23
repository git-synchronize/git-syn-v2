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
    printf("Usage:\n\t%s [option] ... [command] ...\n%s\n", PROGRAM_NAME, PROGRAM_DESCRIPTION);
    exit(status);
}

void print_version(int status)
{
    printf("%s version %s\n", PROGRAM_NAME, VERSION);
    exit(status);
}
