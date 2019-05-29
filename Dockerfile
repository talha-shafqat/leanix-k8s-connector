FROM gcr.io/distroless/static

ADD bin/leanix-k8s-connector /leanix-k8s-connector

# This would be nicer as `nobody:nobody` but distroless has no such entries.
USER 65535:65535

ENTRYPOINT ["/leanix-k8s-connector"]
