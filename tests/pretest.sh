#!/usr/bin/env sh

cd ../build/test && \
docker-compose build && \
docker-compose up --force-recreate & \
../build/test/wait-for-it.sh 127.0.0.1:443 -- echo "server is up"
jest
cd ../build/test && docker-compose down
