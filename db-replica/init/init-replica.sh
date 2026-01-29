#!/bin/bash
set -e

echo "Waiting for primary..."
until pg_isready -h db-primary -p 5432; do
  sleep 2
done

echo "Running base backup..."
rm -rf /var/lib/postgresql/data/*

PGPASSWORD=SofSvi_37 pg_basebackup \
  -h db-primary \
  -U replicator \
  -D /var/lib/postgresql/data \
  -Fp -Xs -P -R

echo "Replica ready"
