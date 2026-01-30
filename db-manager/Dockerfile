FROM ghcr.io/pgdogdev/pgdog:8d3ea65

EXPOSE 6433

COPY ./config/pgdog.toml /etc/pgdog/
COPY ./config/users.toml /etc/pgdog/

CMD ["pgdog", "--config", "/etc/pgdog/pgdog.toml", "--users", "/etc/pgdog/users.toml"]
