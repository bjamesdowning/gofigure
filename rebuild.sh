#! /bin/bash
#rebuild env for testing
if [ "ls | grep gofigure" ]; then
   rm gofigure
   echo Old binary removed
fi

GOOS=linux go build
echo New binary build

if [ "docker-compose ps | grep gofigure_gofigure" ]; then
   docker-compose down
fi

sleep 5

docker-compose up -d --build
