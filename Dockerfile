FROM golang:1.8-alpine
COPY ./redgo /usr/bin/redgo
RUN chmod 777 /usr/bin/redgo && \
    mkdir -p /etc/redgo && \
    apk update && apk add dumb-init

COPY config.yaml /etc/redgo/config.yaml

EXPOSE 6379

ENTRYPOINT ["dumb-init", "--"]
CMD ["redgo", "-configure", "/etc/redgo/config.yaml"]

