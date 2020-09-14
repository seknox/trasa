#!/bin/bash
npm run build && \
rm -rdf /var/trasa/dashboard
mv build/ /var/trasa/dashboard/
