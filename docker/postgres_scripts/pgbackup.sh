#!/bin/bash

set -e

# dump databases
# shellcheck disable=SC2006
for DATABASE in `psql -At -U postgres -c "select datname from pg_database where not datistemplate order by datname;" postgres`
do
  echo "Plain backup of $DATABASE"
  # shellcheck disable=SC2046
  pg_dump -U postgres -Fc "$DATABASE" > /opt/backups/"$DATABASE".$(date -d "today" +"%Y-%m-%d-%H-%M").dump
done

# delete files older than 7 days
find /opt/backups -mtime +7 -type f -delete
