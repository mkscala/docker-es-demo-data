#!/bin/sh

set -e

# Ingest data into Elasticsearch
ingest

# Add Kibana Dashboard
if [ -z ${ES_USERNAME} ]; then
  import_dashboards -es http://elasticsearch:9200 -dir /nginx_data/nginx-dashboard;
else
  import_dashboards -es http://elasticsearch:9200 -user $ES_USERNAME -pass $ES_PASSWORD -dir /nginx_data/nginx-dashboard;
fi
