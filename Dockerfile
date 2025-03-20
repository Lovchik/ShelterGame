FROM golang:1.21 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y gcc libc-dev sqlite3 libsqlite3-dev

COPY . .

ENV CGO_ENABLED=1
RUN go build -o main ShelterGame/cmd/app

######## Новый этап ########
FROM debian:bookworm-slim

ENV TZ=Europe/Minsk
WORKDIR /app

COPY --from=builder /app/main .
COPY .env .
COPY base.db .
COPY sample .

RUN apt-get update && apt-get install -y ca-certificates tzdata sqlite3 && rm -rf /var/lib/apt/lists/*

CMD ["./main"]