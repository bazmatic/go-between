#!/bin/sh
PGPASSWORD=docker psql -h localhost -p 5432 -U docker -d between --file=./init.sql