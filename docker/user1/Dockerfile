FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o user1 ./cmd/user1

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user1 .
COPY --from=builder /app/internal/templates /app/internal/templates
COPY --from=builder /app/static /app/static
COPY --from=builder /app/config_redis.json .

EXPOSE 8082

CMD ["./user1"]