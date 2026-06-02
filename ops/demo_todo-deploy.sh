#!/bin/bash

# chmod +x demo_todo-*.sh
#
# copy tar and script files from apps to destination directory
rm -rf ~/apps/demo_todo
mkdir -p ~/apps/demo_todo
mv ~/apps/demo_todo-* ~/apps/demo_todo
mv ~/apps/demo_todo.tar ~/apps/demo_todo
mv ~/apps/demo_todo/demo_todo-env ~/apps/demo_todo/.env
#
# load the docker image and check
docker load -i ~/apps/demo_todo/demo_todo.tar
docker image ls | grep demo_todo
