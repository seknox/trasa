#!/bin/sh
echo $TRASA_HOSTNAME > build/Constants.json
serve -s build