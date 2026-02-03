FROM alpine:3.23.2

LABEL maintainer="KaleniG <kalenig.01@gmail.com>"
LABEL version="1.0"
LABEL description="Golang streaming backend for YKStreaming service"

RUN apk add --no-cache go ffmpeg

COPY . .

EXPOSE 9090

CMD go run ./cmd/api/main.go