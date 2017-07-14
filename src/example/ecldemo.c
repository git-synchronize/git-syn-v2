/*
 *     Example of a C program embedding ECL with callbacks to C functions.
 *         Compiled via: gcc ecldemo.c -lecl
 *          */
#include <stdio.h>
#include <stdlib.h>
#include "ecl/ecl.h"

#define DEFUN(name,fun,args) \
	cl_def_c_function(c_string_to_object(name), \
			(cl_objectfn_fixed)fun, \
			args)

cl_object foo() {
	return ecl_make_integer(42);
}

cl_object bar(cl_object a, cl_object b) {
	int aval = fix(a);
	int bval = fix(b);
	return ecl_make_integer(aval + bval);
}

/*
 *     Assumes the string is a valid call.
 *     */
cl_object ecl_call(char *call) {
	return cl_safe_eval(c_string_to_object(call), Cnil, Cnil);
}

void init() {
	cl_boot(1, (char **)&"");

	atexit(cl_shutdown);

	/* 
	 *     Uncomment these lines to place your code into a separate package,
	 *         They may then be called like (my-code:foo)
	 *          */
	//  ecl_call("(make-package :my-code)");
	//  //  ecl_call("(in-package my-code)");
	//
	//      DEFUN("foo", foo, 0);
	//          DEFUN("bar", bar, 2);
	//
	//          //  ecl_call("(export foo)");
	//          //  ecl_call("(export bar)");
	//          //  ecl_call("(in-package common-lisp-user)");
}

int main() {
	init();
	cl_object exit_obj = c_string_to_object(":EXIT");
	cl_object result = Cnil;

	while (cl_equal(exit_obj, result) == Cnil) {
		printf("\n> ");
		cl_object form = ecl_call("(read)");
		result = cl_safe_eval(form, Cnil, Cnil);
		cl_print(1, result);
	}
	putchar('\n');

	return 0;
}
