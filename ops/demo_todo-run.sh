#!/bin/bash
# ~/apps/demo_todo/demo_todo-run.sh

# make port 9003 available in firewall
#
# run the demo_todo container
docker run  --name demo_todo --rm -p 9003:80 --env-file ~/apps/demo_todo/.env demo_todo:latest
