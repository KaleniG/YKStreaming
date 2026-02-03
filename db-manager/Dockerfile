FROM ghcr.io/pgdogdev/pgdog:8d3ea65


COPY ./config/pgdog.toml /etc/pgdog/
COPY ./config/users.toml /etc/pgdog/
EXPOSE 6433

CMD ["sh", "-c", "sleep 10 && pgdog --config /etc/pgdog/pgdog.toml --users /etc/pgdog/users.toml"]
