FROM blacktop/logstash

MAINTAINER blacktop, https://github.com/blacktop

COPY nginx_ElasticStack_Example /nginx_ElasticStack_Example

WORKDIR /nginx_ElasticStack_Example

ENV LS_SETTINGS_DIR /etc/logstash

CMD cat nginx_logs | /logstash-entrypoint.sh -f nginx_logstash.conf
