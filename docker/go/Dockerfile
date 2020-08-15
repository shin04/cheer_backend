FROM golang:1.14.2-alpine

ENV GO111MODULE=on

RUN apk update && apk add git bash

RUN apk add build-base

RUN go get github.com/pilu/fresh && \
    go get -u github.com/labstack/echo/... && \
    go get github.com/lib/pq && \
    go get -v github.com/rubenv/sql-migrate/...
