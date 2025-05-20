FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8000
ENTRYPOINT ["./main"]