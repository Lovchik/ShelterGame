FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download
COPY . .
RUN go build -o main ShelterGame/cmd/app

######## Start a new stage #######
FROM alpine:latest

ENV TZ=Europe/Minsk

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
