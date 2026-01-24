FROM alpine:3
RUN apk add --no-cache git
COPY dist/semrel /usr/bin/semrel
ENTRYPOINT ["/usr/bin/semrel"]
