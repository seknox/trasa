#!/bin/bash
npm run build && \
rm -rdf /etc/trasa/build
mv build/ /etc/trasa/build/