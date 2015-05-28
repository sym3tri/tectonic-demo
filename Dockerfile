FROM scratch

ADD bin/demo /opt/demo
ADD static /opt/static
WORKDIR /opt
CMD ["/opt/demo"]
