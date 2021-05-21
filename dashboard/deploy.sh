#!/bin/bash
registry=$1
version=`date +%F`
tag=$1:v$version
npm run build && sudo docker build -t $tag . && sudo docker push $tag
