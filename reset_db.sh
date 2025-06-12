#!/bin/sh

DB_STRING=$(jq '.db_url' -r ~/.gatorconfig.json);

cd /home/badger/boot.dev/go-blog-agg/sql/schema;

goose postgres $DB_STRING down;
# goose postgres $DB_STRING down -v;

goose postgres $DB_STRING up;
# goose postgres $DB_STRING up -v;
