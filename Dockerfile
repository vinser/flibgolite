# syntax=docker/dockerfile:1

# build dp
FROM --platform=$BUILDPLATFORM scratch
# FROM gcr.io/distroless/static-debian12
# FROM --platform=$BUILDPLATFORM alpine

ARG TARGETOS TARGETARCH TARGETVARIANT
COPY ./bin/flibgolite-$TARGETOS-$TARGETARCH${TARGETVARIANT:+-${TARGETVARIANT#v}} /flibgolite/flibgolite

# expose ports
EXPOSE 8085

WORKDIR /flibgolite
# run command
ENTRYPOINT ["/flibgolite/flibgolite"]
