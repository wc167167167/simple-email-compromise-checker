#!/bin/bash
set -e

docker build . -t simplechecker:test
docker build -f nginx.Dockerfile . -t ngingx:test
docker pull amazon/dynamodb-local

docker run -p 8990:8000 -d amazon/dynamodb-local

docker run -p 8989:8989 -d simplechecker:test
docker run -p 8988:8989 -d simplechecker:test
docker run -p 8987:8989 -d simplechecker:test

docker run -p 80:80 -d ngingx:test
