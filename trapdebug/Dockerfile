FROM alpine:latest
MAINTAINER Prawn
USER root

ENV LISTEN_PORT=162

EXPOSE 162/udp

COPY trapdebug /usr/local/bin/trapdebug

CMD exec /usr/local/bin/trapdebug -listenport=$LISTEN_PORT

# docker build -t "prawn/snmp-trapdebug" .
