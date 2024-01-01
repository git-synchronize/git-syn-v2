/*
 * Copyright (C) 2019-2022 Lucas Ramage <lucas.ramage@infinite-omicron.com>
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

#include "syn.h"

int main(int argc, char **argv)
{
    bool install_extension = false,
         pre_push = false,
	 remove_extension = false;

    static int debug, verbose;

    char cwd[PATH_MAX];

    char *git_repo_dir;

    int extension_installed, opt;

    int ret = EXIT_FAILURE;

    static struct option long_options[] = {
        { "debug", no_argument, &debug, 1 },
        { "help", no_argument, 0, 'h' },
        { "verbose", no_argument, &verbose, 1 },
        { "version", no_argument, 0, 'v' },
        { 0, 0, 0, 0 }
    };

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
            if (long_options[option_index].flag != 0)
                break;
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
        } else if (strcmp(argv[optind], "pre-push") == 0) {
            pre_push = true;
        }
    }

    if (debug) {
        for (int i = 1; i < argc; i++) {
            printf("argv[%d]: %s\n", i, argv[i]);
        }

        printf("argc: %d\n", argc);
    }

    if (verbose) {
        printf("Not yet implemented\n");
    }

    sds git_config_env = sdsnew(getenv("GIT_CONFIG"));

    if (debug) {
        printf("GIT_CONFIG: %s\n", git_config_env);
    }
    // First check environment variable, then current firevtory
    if (git_config_env[0] == '\0') {
        if (debug) {
            printf("git_config_env is empty\n");
        }

        if (getcwd(cwd, sizeof(cwd)) == NULL) {
            if (debug) {
                printf("getcwd failed\n");
            }
            ret = EXIT_FAILURE;
        } else {

            if (debug) {
                printf("cwd: %s\n", cwd);
            }

            git_repo_dir = cwd;
        }
    } else {
        if (debug) {
            printf("git_config_env: %s\n", git_config_env);
        }

        git_repo_dir = git_config_env;
    }

    if (debug) {
        printf("git_repo_dir: %s\n", git_repo_dir);
    }

    git_libgit2_init();

    if (git_repository_open_ext
        (NULL, git_repo_dir, GIT_REPOSITORY_OPEN_NO_SEARCH, NULL) == 0) {

        if (install_extension) {
            extension_installed = init_repo(git_repo_dir);
            if (extension_installed == EXIT_SUCCESS) {
                printf("Git SYN initialized.\n");
            }
        } else if (remove_extension) {
            printf("Not yet implemented\n");
        } else if (pre_push) {
            printf("Not yet implemented\n");
        }

        ret = EXIT_SUCCESS;
    } else {
        if (debug) {
            printf("Cannot open git_repo_dir\n");
        } else if (remove_extension) {
            printf("Not yet implemented\n");
        }
    }

    git_libgit2_shutdown();

    exit(ret);
}
