FROM golang:alpine as build

RUN apk add ca-certificates git gcc musl-dev

WORKDIR /opt

COPY go.mod go.sum ./
RUN  go mod download

COPY cmd/client cmd/client
COPY internal internal
COPY config config
COPY config.toml /srv/

RUN cd /opt/cmd/client && go build && mv client /srv/client

FROM alpine:latest

COPY --from=build /srv /srv

WORKDIR /srv
CMD /srv/client
