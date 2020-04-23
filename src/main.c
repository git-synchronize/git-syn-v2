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

#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "syn.h"

int main(int argc, char *argv[])
{
    int opt;

    while ((opt = getopt(argc, argv, "h")) != -1) {
        switch (opt) {
        case 'h':
            print_usage();
            break;
        case '?':
            print_usage();
            exit(EXIT_FAILURE);
        }
    }

    for (; optind < argc; optind++) {
        if (strcmp(argv[optind], "install") == 0) {
            printf("not yet implemented\n");
        } else if (strcmp(argv[optind], "monitor") == 0) {
            printf("not yet implemented\n");
        }
    }

    exit(EXIT_SUCCESS);
}
