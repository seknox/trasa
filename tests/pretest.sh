#!/usr/bin/env sh

cd ../build/test && docker-compose build && docker-compose up &

# BACK_PID=$!
# wait $BACK_PID

#../build/test/wait-for-it.sh 127.0.0.1:443 -- echo "=========================== SERVER IS UP ===================================="




jest --runInBand

cd ../build/test && docker-compose down
