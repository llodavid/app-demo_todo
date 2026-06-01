#!/bin/bash

# chmod +x demo_todo-*.sh
#
# copy tar and script files to destination directory
rm -rf ~/apps/demo_todo
mkdir -p ~/apps/demo_todo
mv ~/demo_todo-* ~/apps/demo_todo
#
# load the docker image and check
docker load -i ~/apps/demo_todo/demo_todo.tar
docker image ls | grep demo_todo
