FROM alpine:3.23.2

LABEL maintainer="KaleniG <kalenig.01@gmail.com>"
LABEL version="1.0"
LABEL description="NGINX/Vite streaming frontend for YKStreaming service"

WORKDIR /var/www/servers/portal

RUN apk add --no-cache nginx nodejs npm

RUN adduser -D -u 33 -G www-data www-data
RUN mkdir -p /var/www/tmp/client_body \
  && chown -R www-data:www-data /var/www/tmp

COPY ./config/nginx.conf /etc/nginx/nginx.conf
RUN rm -f /etc/nginx/http.d/default.conf

COPY . .
RUN npm install

EXPOSE 80 5173

CMD ["sh", "-c", "nginx -g 'daemon off;' & npm run dev"]
