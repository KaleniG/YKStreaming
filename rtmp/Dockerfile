FROM alpine:3.23.2

LABEL maintainer="KaleniG <kalenig.01@gmail.com>"
LABEL version="1.0"
LABEL description="NGINX RTMP streaming server for YKStreaming service"

WORKDIR /var/www/servers/rtmp

RUN apk add --no-cache nginx nginx-mod-rtmp

COPY nginx.conf /etc/nginx/nginx.conf

RUN adduser -D -u 33 -G www-data www-data \
  && mkdir -p /var/www/streams /var/www/vods \
  && chown -R www-data:www-data /var/www/streams /var/www/vods

EXPOSE 1935 8080

CMD ["nginx", "-g", "daemon off;"]