FROM alpine:3.6

ADD tmp/_output/bin/starter-pack-operator /usr/local/bin/starter-pack-operator

RUN adduser -D starter-pack-operator
USER starter-pack-operator
