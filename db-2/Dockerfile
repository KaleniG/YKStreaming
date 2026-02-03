FROM debian:12.13

RUN apt-get update \
  && apt-get install -y patroni postgresql postgresql-client \
  && rm -rf /var/lib/apt/lists/*

USER postgres

RUN mkdir -p /var/lib/postgresql/data
RUN chown -R postgres:postgres /var/lib/postgresql/data
RUN chmod 700 /var/lib/postgresql/data 

COPY ./config /etc/patroni/

ENV PATH="/usr/lib/postgresql/15/bin:$PATH"

EXPOSE 5432 8008

CMD ["patroni", "/etc/patroni/patroni.yml"]