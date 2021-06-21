FROM fedora:34

RUN dnf install -y golang xen-devel-4.14.1-7.fc34 yajl-devel
COPY ./ /tmp/xenlight_exporter/
RUN cd /tmp/xenlight_exporter/ \
    && go build -o /usr/local/bin/xen_exporter *.go \
    && cd / \
    && rm -rf /tmp/xenlight_exporter

ENTRYPOINT ["xen_exporter"]
