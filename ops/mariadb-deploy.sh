#!/bin/bash

# chmod +x mariadb-*.sh
#
# copy tar and script files to destination directory
rm -rf ~/apps/mariadb
mkdir -p ~/apps/mariadb
mv ~/apps/mariadb-* ~/apps/mariadb
mv ~/apps/mariadb/mariadb-env ~/apps/mariadb/.env
mv ~/apps/goose.exe ~/apps/mariadb
#
cd ~/apps/mariadb
mkdir -p ./db
mkdir -p ./resources
mv -r ./mariadb-resources ./resources
rm -rf ./mariadb-resources
mkdir -p ./migrations
mv -r ./mariadb-migrations ./migrations
rm -rf ./mariadb-migrations
#