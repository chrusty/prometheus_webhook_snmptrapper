FROM alpine:latest

RUN apk update
RUN apk add net-snmp

COPY PROMETHEUS-TRAPPER-MIB.txt /usr/share/snmp/mibs/

CMD ["/usr/sbin/snmptrapd", "-f", "-Lo", "-m", "PROMETHEUS-TRAPPER-MIB", "-M", "/usr/share/snmp/mibs"]

# docker build -t "prawn/snmptrapd" .
# docker run --net=prometheus_default --ip=172.15.0.24 --name=snmptrapd -d prawn/snmptrapd
