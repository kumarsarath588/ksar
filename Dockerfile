## Builder stage
FROM golang:1.17 as builder

COPY . /app

WORKDIR /app

RUN go mod download

ENV CGO_ENABLED=0

RUN GOOS=linux go build -o tabsquare main.go

## Container build stage
FROM alpine:3.14.3

RUN apk --no-cache add curl && \
addgroup -S admin && \
adduser -S admin -G admin

ENV APP_DB_HOST=mysql \
    APP_DB_PORT=3306 \
    APP_DB_NAME=tabsquare

USER admin

WORKDIR /app

COPY --from=builder /app/tabsquare .
COPY scripts/entrypoint.sh .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --retries=5 --start-period=15s CMD curl -f http://localhost:8080/health || exit 1

CMD ["/app/tabsquare"]

ENTRYPOINT [ "/app/entrypoint.sh" ]