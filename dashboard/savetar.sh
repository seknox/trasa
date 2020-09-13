#!/bin/bash
registry=$1
version=`date +%F`
tag=$1:latest
npm run build && sudo docker build -t $tag . &&  sudo docker save -o dash-$version.tar $tag
