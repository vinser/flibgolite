FROM golang:1.18.3 as builder

# build flibgolite
COPY . /flibgolite
WORKDIR /flibgolite
RUN go build ./cmd/flibgolite

FROM alpine:3.16.0
COPY --from=builder /flibgolite/flibgolite /flibgolite/flibgolite

# run flibgolite with musl instead of glibc
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# configure service and directories
RUN mkdir -p /flibgolite/config /var/flibgolite
COPY ops/docker-config.yml /flibgolite/config/config.yml

# expose ports
EXPOSE 8085

# probes
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8085/opds || exit 1

# run command
ENTRYPOINT ["/flibgolite/flibgolite"]
