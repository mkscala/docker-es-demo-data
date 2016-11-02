#!/bin/sh

set -e

# Ingest data into Elasticsearch
/ingest

# Add Kibana Dashboard
/import_dashboards -file /nginx_data/nginx-dashboard.zip
