FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o rzhaka-tournaments \
    ./cmd/server


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/rzhaka-tournaments .

EXPOSE 8080

CMD ["./rzhaka-tournaments"]