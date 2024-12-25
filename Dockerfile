FROM --platform=linux/amd64 golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api /app/

EXPOSE 8080
