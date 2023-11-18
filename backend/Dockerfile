FROM golang:1.21-bullseye AS builder

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

