FROM golang:1.21.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/update_urls ./cmd/cron_lobs/update_urls/main.go

FROM alpine:latest as update_urls
RUN  apk update --no-cache && \
apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /bin/update_urls /app/update_urls
ENTRYPOINT [ "/app/update_urls" ]
CMD ["-conf=/app/config.yaml"]