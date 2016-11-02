FROM blacktop/logstash

MAINTAINER blacktop, https://github.com/blacktop

RUN apk-install curl

ENV LS_SETTINGS_DIR /etc/logstash
COPY config/logstash/logstash.yml /etc/logstash/logstash.yml

COPY nginx_data /nginx_data
COPY ingest_data.sh /

WORKDIR /nginx_data

# Use Logstash to ingest data into Elasticsearch and Add Kibana Dashboard
ENTRYPOINT /ingest_data.sh
