FROM thinkingdata/ta-logbus-v2:latest

RUN mkdir -p /ta/logbus/conf

COPY build/build_thinkingdata_http_linux /ta/logbus/thinkingdata_http

COPY shell/run/* /opt/run/
COPY shell/run_all /opt/bin/

ENTRYPOINT ["/opt/bin/run_all"]
