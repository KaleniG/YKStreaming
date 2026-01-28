FROM alpine:3.23.2

LABEL maintainer="KaleniG <kalenig.01@gmail.com>"
LABEL version="1.0"
LABEL description="NGINX/Vite streaming backend for YKStreaming service"

WORKDIR /var/www/servers/api

RUN apk add --no-cache go ffmpeg

COPY . .
#RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

#RUN export PATH="$PATH:$(go env GOPATH)/bin" \
#  sqlc generate

EXPOSE 9090

CMD sh -c "go run ./cmd/api/main.go"