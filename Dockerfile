FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .
COPY tmp ./tmp

CMD ["./main"]
