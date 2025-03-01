#!/bin/bash
set -e

echo "host replication replicator 127.0.0.1/32 md5" >> /var/lib/postgresql/data/pg_hba.conf

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    create user replicator with replication password 'secret';
    alter system set listen_addresses to '*';
    alter system set wal_level to logical;
    alter system set port to 2345;
    create database sourcedb;
    \c sourcedb
    create table tb(id integer not null primary key, name varchar(64));
    grant select on all tables in schema public to replicator;
    create publication publication_one for all tables;
EOSQL
