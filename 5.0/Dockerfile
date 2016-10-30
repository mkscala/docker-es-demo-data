FROM blacktop/logstash

MAINTAINER blacktop, https://github.com/blacktop

COPY nginx_ElasticStack_Example /nginx_ElasticStack_Example

WORKDIR /nginx_ElasticStack_Example

RUN cat nginx_logs | logstash -f nginx_logstash.conf

ENTRYPOINT ["/logstash-entrypoint.sh"]
CMD ["-e", ""]
