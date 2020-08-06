#!/usr/bin/env sh

cd ../build/test && docker-compose build && docker-compose up &

../build/test/wait-for-it.sh 127.0.0.1:443 -- echo "server is up"

#cd ../test
jest
