# Build stage
FROM golang:alpine AS builder
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Runtime stage
FROM alpine:latest
COPY --from=builder /go/bin/goose /usr/local/bin/goose
ENTRYPOINT ["goose"]