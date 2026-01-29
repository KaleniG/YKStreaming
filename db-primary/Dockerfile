FROM postgres:16-alpine

LABEL maintainer="KaleniG <kalenig.01@gmail.com>"
LABEL version="1.0"
LABEL description="PostgreSQL primary database for YKStreaming service"

COPY ./config/postgresql.conf /etc/postgresql/postgresql.conf
COPY ./config/pg_hba.conf /etc/postgresql/pg_hba.conf

EXPOSE 5432

CMD ["postgres", \
  "-c", "config_file=/etc/postgresql/postgresql.conf", \
  "-c", "hba_file=/etc/postgresql/pg_hba.conf"]