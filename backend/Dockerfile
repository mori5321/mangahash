## dev
FROM golang:1.21-bullseye AS dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest


## builder
FROM golang:1.21-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd
COPY ./queries ./queries
COPY ./migrations ./migrations
COPY ./Makefile ./Makefile

RUN mkdir -p /bin
RUN go build -o /bin/main cmd/main.go


## runner
FROM golang:1.21-bullseye AS runner

WORKDIR /app
COPY --from=builder /bin/main /bin/main

EXPOSE 9090
CMD ["/bin/main"]

