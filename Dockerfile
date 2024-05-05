FROM golang:1.21.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/app ./cmd/main.go

FROM alpine:latest as app
RUN  apk update --no-cache && \
apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /bin/app /app/app
ENTRYPOINT [ "/app/app" ]
CMD ["-conf=/app/config.yaml"]
EXPOSE 80