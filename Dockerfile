
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o cron-job ./cmd/cron-job/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/cron-job .

EXPOSE 8000

CMD ["./api"]
