# Estágio de build
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/api/main.go

# Estágio final (imagem menor)
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/internal/infra/security/keys ./internal/infra/security/keys

RUN mkdir -p uploads

EXPOSE 8080

CMD ["./server"]
