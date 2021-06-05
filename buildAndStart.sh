#!/bin/bash
set -e

docker build . -t simplechecker:test
docker pull amazon/dynamodb-local

docker run -p 8990:8000 amazon/dynamodb-local &
docker run -p 8989:8989 simplechecker:test &
