FROM postgres:16-alpine

LABEL maintainer="KaleniG <kalenig.01@gmail.com>"
LABEL version="1.0"
LABEL description="PostgreSQL database for YKStreaming service"

# Copy only a single SQL file into the init directory
COPY schema.sql /docker-entrypoint-initdb.d/01-schema.sql
