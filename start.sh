#!/bin/bash

docker build -t blog .
docker stop blog
docker rm blog
docker run -v /var/lib/blog:/go/src/github.com/avesanen/blog/db -p 5001:5000 -d --name blog blog
