#!/bin/bash
docker build . --build-arg HOSTNAME=localhost:8080 -t gowog_local|| exit
docker stop gowog
docker rm gowog
docker run --privileged -d --name gowog -p 8080:8080 gowog_local server -prod client/build/
