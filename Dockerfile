FROM scratch

ADD bin/demo /opt/demo
CMD ["/opt/demo"]
