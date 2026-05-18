# Build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /build
RUN git clone https://github.com/sqlc-dev/sqlc.git .
RUN go build -o sqlc ./cmd/sqlc

# Runtime stage
FROM alpine:latest
COPY --from=builder /build/sqlc /usr/local/bin/sqlc
WORKDIR /src
ENTRYPOINT ["sqlc"]
