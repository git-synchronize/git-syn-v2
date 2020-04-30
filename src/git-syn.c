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
#include <git2.h>
#include <limits.h>
#include <sds/sds.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "git-syn.h"

int main(int argc, char **argv)
{
    bool install_extension = false,
        remove_extension = false, monitor_repository = false;

    char cwd[PATH_MAX];

    int extension_installed, opt;

    static struct option long_options[] = {
        {"help", no_argument, 0, 'h'},
        {"version", no_argument, 0, 'v'},
        {0, 0, 0, 0}
    };

    if (getcwd(cwd, sizeof(cwd)) == NULL) {
        exit(EXIT_FAILURE);
    }


    sds git_config_env = sdsnew(getenv("GIT_CONFIG"));

    while (true) {
        int option_index = 0;
        opt = getopt_long(argc, argv, "hv0", long_options, &option_index);

        if (opt == -1)
            break;

        switch (opt) {
        case 'h':
            print_usage(EXIT_SUCCESS);
            break;
        case 'v':
            print_version(EXIT_SUCCESS);
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

    git_libgit2_init();

    (void) git_config_env;

    parse_config();

    if (install_extension) {
        extension_installed = init_repo();
        if (extension_installed == EXIT_SUCCESS) {
            printf("Git SYN initialized.\n");
        }
    } else if (remove_extension) {
        printf("Not yet implemented\n");
    } else if (monitor_repository) {
        printf("Not yet implemented\n");
    }

    git_libgit2_shutdown();

    exit(EXIT_SUCCESS);
}
