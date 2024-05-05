FROM golang:1.21.6-alpine AS builder
WORKDIR /app
ENV GOPATH=/
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY ./storage/migrations ./migrations

FROM alpine:latest as migrations
WORKDIR /app
COPY --from=builder /bin/goose /app/migrations /app/
ENTRYPOINT [ "/app/goose" ]
CMD ["-dir", "/app", "up"]