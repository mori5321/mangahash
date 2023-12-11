## dev
FROM golang:1.21-bullseye AS dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

## builder
FROM golang:1.21-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd
COPY ./Makefile ./Makefile

RUN mkdir -p /bin
RUN go build -o /bin/main cmd/main.go


## runner
FROM golang:1.21-bullseye AS runner

WORKDIR /app
COPY --from=builder /bin/main /bin/main

EXPOSE 5050
CMD ["/bin/main"]

