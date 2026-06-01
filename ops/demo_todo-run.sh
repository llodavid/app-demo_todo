#!/bin/bash
# ~/apps/demo_todo/demo_todo-run.sh

# make port 9003 available in firewall
#
# run the demo_todo container
docker run  --name demo_todo --rm -p 9003:8080 -e OCI_PORT=9003 -e DB_PORT=9002 -e DB_HOST=home-testappserver.robertthecoder.org -e DB_USER=myuser -e DB_PASSWORD=mypw -e DB_NAME=demo_todo -e LOG_LEVEL=debug demo_todo:latest
