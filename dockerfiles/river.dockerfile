# Build stage
FROM golang:alpine AS builder
RUN go install github.com/riverqueue/river/cmd/river@latest

# Runtime stage
FROM alpine:latest
COPY --from=builder /go/bin/river /usr/local/bin/river
ENTRYPOINT ["river"]
