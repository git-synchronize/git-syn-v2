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
#include <getopt.h>
#include <string.h>

#include "git-syn.h"

int main(int argc, char **argv)
{
    int c;

    while (1) {
        int option_index = 0;
        static struct option long_options[] = {
            { "help", no_argument, 0, 0 },
            { "version", no_argument, 0, 0 },
            { 0, 0, 0, 0 }
        };

        c = getopt_long(argc, argv, "0", long_options, &option_index);

        if (c == -1)
            break;

        switch (c) {
        case 0:
            if (strcmp(long_options[option_index].name, "help") == 0)
                print_usage();
            if (strcmp(long_options[option_index].name, "version") == 0)
                print_version();
            break;
        default:
            exit(EXIT_FAILURE);
        }
    }

    exit(EXIT_SUCCESS);
}
