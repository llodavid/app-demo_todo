#!/bin/bash

# chmod +x mariadb-*.sh
#
# copy tar and script files to destination directory
rm -rf ~/apps/mariadb
mkdir -p ~/apps/mariadb
mv ~/mariadb-* ~/apps/mariadb
#
cd ~/apps/mariadb
mkdir -p ./db
mkdir -p ./resources
mv -r ./mariadb-resources ./resources
rm -rf ./mariadb-resources
