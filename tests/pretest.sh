#!/usr/bin/env sh

cd ../build/test && docker-compose build && docker-compose up --force-recreate -d && cd ../../tests


max_iterations=30
wait_seconds=1
http_endpoint="https://app.trasa/idp/login"

iterations=0
while true
do
	((iterations++))
	echo "Attempt $iterations"
	sleep $wait_seconds

	http_code=$(curl  --insecure -XPOST -s -o /tmp/result.txt -w '%{http_code}' "$http_endpoint";)

	if [ "$http_code" -eq 200 ]; then
		echo "Server Up"
		break
	fi

	if [ "$iterations" -ge "$max_iterations" ]; then
		echo "Loop Timeout"
		exit 1
	fi
done

jest
cd ../build/test && docker-compose down
