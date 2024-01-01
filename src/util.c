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

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <sys/stat.h>

#include "syn.h"

#define BUF_SIZE 65536          //2^16

int copy_file(const char *src, const char *dest)
{
    char *file_buffer = calloc(BUF_SIZE, 1);
    FILE *src_file = NULL, *dest_file = NULL;
    int ret = EXIT_FAILURE;
    mode_t mode = 0700;
    size_t n;

    if ((src_file = fopen(src, "rb")) && (dest_file = fopen(dest, "rb"))) {
        while ((n = fread(file_buffer, 1, BUF_SIZE, src_file))
               && fwrite(file_buffer, 1, n, dest_file));
    }
    free(file_buffer);
    if (src_file) {
        fclose(src_file);
    }

    if (dest_file) {
        fclose(dest_file);
    }

    chmod(dest, mode);

    ret = EXIT_SUCCESS;

    return ret;
}
