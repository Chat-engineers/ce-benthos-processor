FROM jeffail/benthos:4

WORKDIR /home/benthos

COPY "resources" "resources"
COPY "streams" "streams"
COPY "benthos.yaml" "benthos.yaml"

ENTRYPOINT [ "/benthos", "-c", "benthos.yaml" ]
CMD ["-r", "resources/*.yaml", "streams", "streams/*.yaml" ]
