# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o program.exe ./cmd

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/program.exe .

EXPOSE 8080
CMD ["./program.exe"]
