FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.* .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/main.go

FROM alpine:3.20

COPY --from=builder /app/main /app/falakhan

CMD ["/app/falakhan", "-config", "/app/config.yml"]
