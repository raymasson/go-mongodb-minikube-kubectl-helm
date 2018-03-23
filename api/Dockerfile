FROM alpine:3.6

LABEL maintainer="raymasson"

RUN apk add --no-cache ca-certificates

ADD person-api /

ENTRYPOINT ["/person-api"]

EXPOSE 8000