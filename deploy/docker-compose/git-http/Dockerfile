# Stage 1: build git-syn itself. Build context is the repo root (see
# docker-compose.yml), since git-syn's go.mod lives there.
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /git-syn .

# Stage 2: nginx, reverse-proxying to `git-syn serve`. Debian-based, not
# nginx:alpine -- Alpine's git package doesn't ship git-http-backend, only
# client-side HTTP tools, and git-syn serve needs the real one.
FROM nginx:1.27

RUN apt-get update && \
    apt-get install -y --no-install-recommends git gosu && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /var/lib/git

COPY --from=build /git-syn /usr/local/bin/git-syn
COPY deploy/docker-compose/git-http/nginx.conf /etc/nginx/nginx.conf
COPY deploy/docker-compose/git-http/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]
