FROM plugins/base:multiarch

LABEL maintainer="Gregory Boddin" \
  org.label-schema.name="Drone Github Status" \
  org.label-schema.vendor="Gregory Boddin" \
  org.label-schema.schema-version="1.0"

ADD release/linux/amd64/github-status /bin/
ENTRYPOINT ["/bin/github-status"]
