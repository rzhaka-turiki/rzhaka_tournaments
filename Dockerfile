FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /rzhaka_tournaments ./cmd/rzhaka_tournaments

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /rzhaka_tournaments .
EXPOSE 8080
CMD ["./rzhaka_tournaments"]