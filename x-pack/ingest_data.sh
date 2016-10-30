#!/bin/sh

set -e

# Add Kibana Dashboard
curl -XPUT http://elasticsearch:9200/kibana-int/dashboard/dashboard-name -d@nginx_kibana.json

# Use Logstash to ingest data into Elasticsearch
cat nginx_logs | /logstash-entrypoint.sh -f nginx_logstash.conf
