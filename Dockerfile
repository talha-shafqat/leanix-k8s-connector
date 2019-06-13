FROM alpine:3.9.4

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN update-ca-certificates

# This would be nicer as `nobody:nobody` but distroless has no such entries.
USER 65535:65535

ADD bin/leanix-k8s-connector /leanix-k8s-connector

ENTRYPOINT ["/leanix-k8s-connector"]
