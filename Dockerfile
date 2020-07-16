# Any copyright is dedicated to the Public Domain.
# https://creativecommons.org/publicdomain/zero/1.0/

# == BUILD STAGE ==
FROM golang:alpine as build

RUN apk add git

RUN go get -u github.com/caddyserver/xcaddy/cmd/xcaddy

WORKDIR /src

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux \
        xcaddy build \
        --output ./caddy \
        --with github.com/lucaslorentz/caddy-docker-proxy@v2.1.0
        --with github.com/HeavenVolkoff/caddy-authelia

# == RUNTIME STAGE ==
FROM gcr.io/distroless/static

COPY --from=build /src/caddy /caddy

ENTRYPOINT ["/caddy"]