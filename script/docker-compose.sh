#!/bin/sh
docker-compose rm -f -v
docker-compose pull
docker-compose up --force-recreate
