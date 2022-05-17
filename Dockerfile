FROM jeffail/benthos:4

COPY "config.yaml" "/config.yaml"
CMD [ "-c", "/config.yaml" ]