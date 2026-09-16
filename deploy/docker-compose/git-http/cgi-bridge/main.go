// cgi-bridge runs git-http-backend as a plain CGI program and exposes it
// over HTTP for nginx to reverse-proxy to. This avoids the FastCGI/CGI
// bridges (fcgiwrap, uwsgi's cgi plugin) that nginx otherwise needs to talk
// to a CGI script: nginx speaks ordinary HTTP/1.1 to this program, and
// net/http/cgi handles the CGI protocol (stdin/stdout, environment
// variables) to git-http-backend itself.
package main

import (
	"log"
	"net/http"
	"net/http/cgi"
	"os"
)

func main() {
	projectRoot := os.Getenv("GIT_PROJECT_ROOT")
	if projectRoot == "" {
		projectRoot = "/var/lib/git"
	}

	handler := &cgi.Handler{
		Path: "/usr/lib/git-core/git-http-backend",
		Env: []string{
			"GIT_PROJECT_ROOT=" + projectRoot,
			"GIT_HTTP_EXPORT_ALL=",
		},
	}

	addr := "127.0.0.1:8080"
	log.Printf("cgi-bridge listening on %s, serving git-http-backend for %s", addr, projectRoot)
	log.Fatal(http.ListenAndServe(addr, handler))
}
