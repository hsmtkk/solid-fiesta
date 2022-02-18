#!/bin/sh
docker-compose --file docker-compose/docker-compose.yml rm -f -v
docker-compose --file docker-compose/docker-compose.yml pull
docker-compose --file docker-compose/docker-compose.yml up --force-recreate
