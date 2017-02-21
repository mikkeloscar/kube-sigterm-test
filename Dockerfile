FROM alpine:latest

ADD build/linux/sigterm-test /

ENTRYPOINT ["/sigterm-test"]
