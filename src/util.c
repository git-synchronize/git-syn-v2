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

int copy_file(const char *src, const char *dest)
{
    char c;
    int ret = EXIT_FAILURE;
    FILE *src_file, *dest_file;

    src_file = fopen(src, "r");

    if (src_file == NULL) {
        ret = EXIT_FAILURE;
    }

    dest_file = fopen(dest, "w");

    if (dest_file == NULL) {
        fclose(src_file);
        ret = EXIT_FAILURE;
    }

    while ((c = fgetc(src_file)) != EOF) {
        fputc(c, dest_file);
    }

    fclose(src_file);
    fclose(dest_file);

    ret = EXIT_SUCCESS;

    return ret;
}
