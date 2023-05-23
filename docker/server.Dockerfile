FROM golang:alpine as build

RUN apk add ca-certificates git gcc musl-dev

WORKDIR /opt

COPY go.mod go.sum ./
RUN  go mod download

COPY cmd/server cmd/server
COPY internal internal
COPY config config
COPY config.toml /srv/

RUN cd /opt/cmd/server && go build && mv server /srv/server


FROM alpine:latest

COPY --from=build /srv /srv

WORKDIR /srv
CMD /srv/server