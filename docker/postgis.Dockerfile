FROM postgres:11

ENV POSTGISV 2.5

# add addressing dictionary
#RUN mkdir -p /opt/apps
#COPY ./pgsql-address-dictionary.zip /opt/apps/pgsql-address-dictionary.zip

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  postgresql-$PG_MAJOR-postgis-$POSTGISV \
  postgresql-$PG_MAJOR-postgis-$POSTGISV-scripts \
  postgresql-$PG_MAJOR-pgrouting \
  postgresql-$PG_MAJOR-pgrouting-scripts \
  postgresql-server-dev-$PG_MAJOR \
  unzip \
  make \
#  && cd /opt/apps \
#  && unzip pgsql-address-dictionary.zip \
#  && cd pgsql-addressing-dictionary-master \
#  && make install \
  && apt-get purge -y --auto-remove postgresql-server-dev-$PG_MAJOR make unzip

# add bakcup job
RUN mkdir -p /opt/backups
COPY postgres_scripts/pgbackup.sh /opt/pgbackup.sh
RUN chmod +x /opt/pgbackup.sh

# add init script
RUN mkdir -p /docker-entrypoint-initdb.d
COPY postgres_scripts/initdb-postgis.sh /docker-entrypoint-initdb.d/postgis.sh

# create volume for backups
VOLUME ["/opt/backups"]
