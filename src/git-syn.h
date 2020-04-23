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

#ifndef GIT_SYN_H
#define GIT_SYN_H

#ifndef VERSION
#define VERSION "0.0.1"
#endif

#define PROGRAM_NAME "git-syn"
#define PROGRAM_DESCRIPTION "Event-driven git remote repository syncing."

void parse_config();

void print_usage(int status);

void print_version(int status);

#endif /* GIT_SYN_H */
