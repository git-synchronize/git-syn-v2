#include <yaml.h>

yaml_parser_t parser;
yaml_event_t event;

int done = 0;

/* Create the Parser object. */
yaml_parser_initialize(&parser);

/* Set a string input. */
char *input = "...";
size_t length = strlen(input);

yaml_parser_set_input_string(&parser, input, length);

/* Set a file input. */
FILE *input = fopen("...", "rb");

yaml_parser_set_input_file(&parser, input);

/* Set a generic reader. */
void *ext = ...;
int read_handler(void *ext, char *buffer, int size, int *length) {
	/* ... */
	*buffer = ...;
	*length = ...;
	/* ... */
	return error ? 0 : 1;
}

yaml_parser_set_input(&parser, read_handler, ext);

/* Read the event sequence. */
while (!done) {

	/* Get the next event. */
	if (!yaml_parser_parse(&parser, &event))
		goto error;

	/*
	 *       ...
	 *             Process the event.
	 *                   ...
	 *                       */

	/* Are we finished? */
	done = (event.type == YAML_STREAM_END_EVENT);

	/* The application is responsible for destroying the event object. */
	yaml_event_delete(&event);

}

/* Destroy the Parser object. */
yaml_parser_delete(&parser);

return 1;

/* On error. */
error:

/* Destroy the Parser object. */
yaml_parser_delete(&parser);

return 0;
