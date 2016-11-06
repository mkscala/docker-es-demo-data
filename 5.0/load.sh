#!/bin/sh

set -e

# Ingest data into Elasticsearch
ingest

# Add Kibana Dashboard
import_dashboards -es http://elasticsearch:9200 -dir /nginx_data/nginx-dashboard
