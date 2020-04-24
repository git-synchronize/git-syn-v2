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

#include <getopt.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "git-syn.h"

int main(int argc, char **argv)
{
    bool install_extension = false,
         remove_extension = false,
         monitor_repository = false;

    int opt;

    static struct option long_options[] = {
        {"help", no_argument, 0, 'h'},
        {"version", no_argument, 0, 0},
        {0, 0, 0, 0}
    };

    while (true) {
        int option_index = 0;
        opt = getopt_long(argc, argv, "h0", long_options, &option_index);

        if (opt == -1)
            break;

        switch (opt) {
        case 'h':
            print_usage(EXIT_SUCCESS);
            break;
        case 0:
            if (strcmp(long_options[option_index].name, "help") == 0)
                print_usage(EXIT_SUCCESS);
            if (strcmp(long_options[option_index].name, "version") == 0)
                print_version(EXIT_SUCCESS);
            break;
        case '?':
            print_usage(EXIT_FAILURE);
        }
    }

    for (; optind < argc; optind++) {
        if (strcmp(argv[optind], "install") == 0) {
            install_extension = true;
        } else if (strcmp(argv[optind], "uninstall") == 0) {
            remove_extension = true;
        } else if (strcmp(argv[optind], "monitor") == 0) {
            monitor_repository = true;
        }
    }

    (void) install_extension;
    (void) remove_extension;
    (void) monitor_repository;

    parse_config();

    exit(EXIT_SUCCESS);
}
