# Stage 1: build the CGI bridge (see cgi-bridge/main.go for why it exists).
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY cgi-bridge/ .
RUN CGO_ENABLED=0 go build -o /cgi-bridge .

# Stage 2: nginx + git-http-backend, proxied through the bridge above.
# Debian-based, not nginx:alpine -- Alpine's git package doesn't ship
# git-http-backend, only client-side HTTP tools.
FROM nginx:1.27

RUN apt-get update && \
    apt-get install -y --no-install-recommends git gosu && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /var/lib/git

COPY --from=build /cgi-bridge /usr/local/bin/cgi-bridge
COPY nginx.conf /etc/nginx/nginx.conf
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]
