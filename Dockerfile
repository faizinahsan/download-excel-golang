# syntax=docker/dockerfile:1

# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app ./app
COPY Book1.xlsx ./Book1.xlsx
EXPOSE 3000
CMD ["./app"]

