FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache curl && \
    curl -o minio https://dl.min.io/server/minio/release/linux-amd64/minio && \
    chmod +x minio

EXPOSE 9000
EXPOSE 9001

CMD ["./minio", "server", "/data", "--console-address", ":9001"]
