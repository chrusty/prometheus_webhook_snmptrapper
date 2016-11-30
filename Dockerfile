FROM alpine:latest
MAINTAINER Prawn
USER root

RUN apk update && apk add curl net-snmp net-snmp-tools

ENV SNMP_COMMUNITY="public"
ENV SNMP_RETRIES=1
ENV SNMP_TRAP_ADDRESS="localhost:162"
ENV WEBHOOK_ADDRESS="0.0.0.0:9099"
ENV SNMP_TRAP_BINARY="/usr/bin/snmptrap"

EXPOSE 9099

COPY prometheus_webhook_snmptrapper.linux-amd64 /usr/local/bin/prometheus_webhook_snmptrapper
COPY sample-alert.json /

CMD exec /usr/local/bin/prometheus_webhook_snmptrapper -snmpcommunity=$SNMP_COMMUNITY -snmpretries=$SNMP_RETRIES -snmptrapaddress=$SNMP_TRAP_ADDRESS -webhookaddress=$WEBHOOK_ADDRESS -snmptrapbinary=$SNMP_TRAP_BINARY

# docker build -t "camelot/prometheus-webhook" .
