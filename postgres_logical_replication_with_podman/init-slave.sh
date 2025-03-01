#!/bin/bash
set -e

sleep 5s

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    create database destinationdb;
    \c destinationdb
    create table tb(id integer not null primary key, name varchar(64));
    create subscription slavesub connection 'user=replicator password=secret host=127.0.0.1 port=2345 dbname=sourcedb' \
    publication publication_one;
EOSQL
