FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/api .

COPY --from=builder /app/frontend ./frontend

COPY --from=builder /app/.env .env

CMD ["./api"]
