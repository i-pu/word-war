FROM envoyproxy/envoy:latest

ARG ENVOY_FILE=./envoy.yaml

COPY ${ENVOY_FILE} /etc/envoy/envoy.yaml

CMD /usr/local/bin/envoy -c /etc/envoy/envoy.yaml

EXPOSE 50051
