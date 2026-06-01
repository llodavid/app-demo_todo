#!/bin/bash
# ~/apps/mariadb/mariadb-run.sh

# make port 9002 available in firewall
#
# run the mariadb container
# sql scripts in ./resources/ are used for database migration
docker run --name mariadb --rm -d -v ~/apps/mariadb/db:/var/lib/mysql:Z  -v ~/apps/mariadb/resources:/docker-entrypoint-initdb.d -p 9002:3306 -e MARIADB_USER=myuser -e MARIADB_PASSWORD=mypw -e MARIADB_ROOT_PASSWORD=myrootpw -e MARIADB_DATABASE=demo_todo mariadb:12.2 
